package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// структура конфигурации приложения
type Config struct {
	TokenTTL  time.Duration  `yaml:"token_ttl" env-required:"true"`
	SecretKey string         `yaml:"secret_key"`
	GRPC      GRPCConfig     `yaml:"grpc"`
	Postgres  PostgresConfig `yaml:"postgres"`
	Log       Log            `yaml:"log"`
	JWT       JWTConfig      `yaml:"jwt"`
}

// структура конфигурации grpc
type GRPCConfig struct {
	Port    string          `yaml:"port"`
	Timeout time.Duration   `yaml:"timeout"`
}

// структура конфигурации postgres
type PostgresConfig struct {
	DSN string `yaml:"dsn"`
}

// структура конфигурации логирования
type Log struct {
	Level string `yaml:"level"`
}

// структура конфигурации jwt
type JWTConfig struct {
	SecretKey string        `yaml:"secret_key"`
	ExpiresKey time.Duration `yaml:"expires_key"`
}

// функция для загрузки конфигурации
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

