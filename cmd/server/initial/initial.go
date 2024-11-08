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
	FlagConfigPath string         // путь к файлу конфигурации
	FlagDBPassword string         // пароль для подключения к базе данных
	FlagDBTest     bool           // флаг для тестовой базы данных
	Cfg            *config.Config // конфигурация
)

// InitConfig функция для инициализации конфигурации.
func InitConfig() {
	flag.StringVar(&FlagConfigPath, "config", "", "path to config file")
	flag.StringVar(&FlagDBPassword, "db_password", "", "database password")
	flag.BoolVar(&FlagDBTest, "db_test", false, "use test database")
	flag.Parse()

	if envConfigPath := os.Getenv("CONFIG_PATH"); envConfigPath != "" {
		FlagConfigPath = envConfigPath
	}

	if envDBPassword := os.Getenv("DB_PASSWORD"); envDBPassword != "" {
		FlagDBPassword = envDBPassword
	}
	if envDBTest, ok := os.LookupEnv("DB_TEST"); ok && envDBTest == "true" {
		FlagDBTest = true
	}

	Cfg = config.MustLoad(FlagConfigPath) // загрузка конфигурации	.yaml
}

// InitStorage функция для инициализации подключения к базе данных.
func InitStorage() (*dbstorage.Storage, error) {
	var dbConfig *pgconn.Config
	dbConfig, err := pgconn.ParseConfig(Cfg.Postgres.DSN)
	if err != nil {
		return nil, fmt.Errorf("error parsing database DSN: %w", err)
	}

	dbconn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		dbConfig.User, FlagDBPassword, dbConfig.Host, strconv.Itoa(int(dbConfig.Port)), dbConfig.Database)

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
func InitS3() (*s3.S3, error) {
	session, err := session.NewSession(&aws.Config{
		Region:   aws.String(Cfg.S3.Region),
		Endpoint: aws.String(Cfg.S3.Endpoint),
		Credentials: credentials.NewStaticCredentials(
			Cfg.S3.AccessKey,
			Cfg.S3.SecretKey,
			"",
		),
		S3ForcePathStyle: aws.Bool(Cfg.S3.ForcePathStyle),
	})
	if err != nil {
		return nil, fmt.Errorf("error creating session s3: %w", err)
	}

	return &s3.S3{S3Session: session}, nil
}
