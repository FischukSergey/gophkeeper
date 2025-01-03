package initial

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/FischukSergey/gophkeeper/internal/config"
	"github.com/FischukSergey/gophkeeper/internal/storage/dbstorage"
	"github.com/FischukSergey/gophkeeper/internal/storage/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	s3bucket "github.com/aws/aws-sdk-go/service/s3"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	EnvLocal = "local"
	EnvProd  = "prod"
	EnvDev   = "dev"
)

// глобальные переменные для хранения флагов и конфигурации.
var (
	FlagConfigPath  string         // путь к файлу конфигурации
	FlagDBUser      string         // имя пользователя для подключения к базе данных
	FlagDBPassword  string         // пароль для подключения к базе данных
	FlagDBHost      string         // хост для подключения к базе данных
	FlagDBPort      string         // порт для подключения к базе данных
	FlagDBTest      bool           // флаг для тестовой базы данных
	FlagS3SecretKey string         // секретный ключ для доступа к S3
	FlagS3AccessKey string         // ключ для доступа к S3
	FlagS3Bucket    string         // бакет для доступа к S3
	Cfg             *config.Config // конфигурация
)

// InitConfig функция для инициализации конфигурации.
func InitConfig() {
	flag.StringVar(&FlagConfigPath, "config", "", "path to config file")
	flag.StringVar(&FlagDBUser, "db_user", "", "database user")
	flag.StringVar(&FlagDBPassword, "db_password", "", "database password")
	flag.StringVar(&FlagDBHost, "db_host", "", "database host")
	flag.StringVar(&FlagDBPort, "db_port", "", "database port")
	flag.BoolVar(&FlagDBTest, "db_test", false, "use test database")
	flag.StringVar(&FlagS3SecretKey, "s3_secret_key", "", "s3 secret key")
	flag.StringVar(&FlagS3AccessKey, "s3_access_key", "", "s3 access key")
	flag.StringVar(&FlagS3Bucket, "s3_bucket", "", "s3 bucket")
	flag.Parse()


	//получаем путь к файлу конфигурации
	if envConfigPath := os.Getenv("CONFIG_PATH"); envConfigPath != "" {
		FlagConfigPath = envConfigPath
	}

	//загружаем конфигурацию
	Cfg = config.MustLoad(FlagConfigPath) // загрузка конфигурации	.yaml

	//получаем оставшиеся параметры из переменных окружения
	if envDBUser := os.Getenv("DB_USER"); envDBUser != "" {
		FlagDBUser = envDBUser
	}
	if envDBPassword := os.Getenv("DB_PASSWORD"); envDBPassword != "" {
		FlagDBPassword = envDBPassword
	}
	if envDBTest, ok := os.LookupEnv("DB_TEST"); ok && envDBTest == "true" {
		FlagDBTest = true
	}
	if envDBHost := os.Getenv("DB_HOST"); envDBHost != "" {
		FlagDBHost = envDBHost
	}
	if envDBPort := os.Getenv("DB_PORT"); envDBPort != "" {
		FlagDBPort = envDBPort
	}

	if envS3SecretKey := os.Getenv("S3_SECRET_KEY"); envS3SecretKey != "" {
		FlagS3SecretKey = envS3SecretKey
	}
	if envS3AccessKey := os.Getenv("S3_ACCESS_KEY"); envS3AccessKey != "" {
		FlagS3AccessKey = envS3AccessKey
	} else if FlagS3AccessKey == "" {
		FlagS3AccessKey = Cfg.S3.AccessKey
	}
	if envS3Bucket := os.Getenv("S3_BUCKET"); envS3Bucket != "" {
		FlagS3Bucket = envS3Bucket
	} else if FlagS3Bucket == "" {
		FlagS3Bucket = Cfg.S3.Bucket
	}
}

// InitStorage функция для инициализации подключения к базе данных.
func InitStorage() (*dbstorage.Storage, error) {
	var dbConfig *pgconn.Config
	dbConfig, err := pgconn.ParseConfig(Cfg.Postgres.DSN)
	if err != nil {
		return nil, fmt.Errorf("error parsing database DSN: %w", err)
	}
	var host string
	if FlagDBHost != "" {
		host = FlagDBHost
	} else {
		host = dbConfig.Host
	}
	var port string
	if FlagDBPort != "" {
		port = FlagDBPort
	} else {
		port = strconv.Itoa(int(dbConfig.Port))
	}
	var user string
	if FlagDBUser != "" {
		user = FlagDBUser
	} else {
		user = dbConfig.User
	}
	dbconn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		user, FlagDBPassword, host, port, dbConfig.Database)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := pgxpool.New(ctx, dbconn)
	if err != nil {
		return nil, fmt.Errorf("%w, unable to create connection db:%s", err, dbConfig.Database)
	}
	err = db.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}
	return &dbstorage.Storage{DB: db}, nil
}

// InitLogger функция для инициализации логгера.
func InitLogger() *slog.Logger {
	var log *slog.Logger
	switch Cfg.Log.Level {
	case EnvLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case EnvDev:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case EnvProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}

// InitS3 функция для инициализации подключения к S3.
func InitS3() (*s3.S3Storage, error) {
	session, err := session.NewSession(&aws.Config{
		Region:   aws.String(Cfg.S3.Region),
		Endpoint: aws.String(Cfg.S3.Endpoint),
		Credentials: credentials.NewStaticCredentials(
			FlagS3AccessKey,
			FlagS3SecretKey,
			"",
		),
		S3ForcePathStyle: aws.Bool(Cfg.S3.ForcePathStyle),
	})
	if err != nil {
		return nil, fmt.Errorf("error creating session s3: %w", err)
	}

	//проверка на наличие бакета
	bucket := FlagS3Bucket
	svc := s3bucket.New(session)

	// Проверяем существование бакета
	_, err = svc.HeadBucketWithContext(context.Background(), &s3bucket.HeadBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		// Создаем новый приватный бакет
		_, err = svc.CreateBucketWithContext(context.Background(), &s3bucket.CreateBucketInput{
			Bucket: aws.String(bucket),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}
	return &s3.S3Storage{S3Session: session, BucketName: bucket}, nil
}
