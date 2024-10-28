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
	Cfg            *config.Config // конфигурация
)

// InitConfig функция для инициализации конфигурации.
func InitConfig() {
	flag.StringVar(&FlagConfigPath, "config", "", "path to config file")
	flag.StringVar(&FlagDBPassword, "db_password", "", "database password")
	flag.Parse()

	if envConfigPath := os.Getenv("CONFIG_PATH"); envConfigPath != "" {
		FlagConfigPath = envConfigPath
	}

	if envDBPassword := os.Getenv("DB_PASSWORD"); envDBPassword != "" {
		FlagDBPassword = envDBPassword
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
