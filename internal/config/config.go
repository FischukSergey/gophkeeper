package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config структура конфигурации приложения.
type Config struct {
	GRPC      GRPCConfig     `yaml:"grpc"`
	JWT       JWTConfig      `yaml:"jwt"`
	Log       Log            `yaml:"log"`
	Postgres  PostgresConfig `yaml:"postgres"`
	S3        S3Config       `yaml:"s3"`
	SecretKey string         `yaml:"secret_key"`
	TokenTTL  time.Duration  `yaml:"token_ttl" env-required:"true"`
}

// GRPCConfig структура конфигурации grpc.
type GRPCConfig struct {
	Port    string        `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

// PostgresConfig структура конфигурации postgres.
type PostgresConfig struct {
	DSN string `yaml:"dsn"`
}

// Log структура конфигурации логирования.
type Log struct {
	Level string `yaml:"level"`
}

// JWTConfig структура конфигурации jwt.
type JWTConfig struct {
	SecretKey  string        `yaml:"secret_key"`
	ExpiresKey time.Duration `yaml:"expires_key"`
}

// S3Config структура конфигурации s3.
type S3Config struct {
	Region         string `yaml:"region"`
	Bucket         string `yaml:"bucket"`
	Endpoint       string `yaml:"endpoint"`
	AccessKey      string `yaml:"access_key"`
	SecretKey      string `yaml:"secret_key"`
	ForcePathStyle bool   `yaml:"force_path_style"`
}

// MustLoad функция для загрузки конфигурации.
func MustLoad(path string) *Config {
	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file not found: " + path)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}
