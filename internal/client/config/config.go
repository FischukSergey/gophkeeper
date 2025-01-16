package config

import (
	"flag"
	"os"

	"github.com/spf13/viper"
)

// Config структура для конфигурации.
type Config struct {
	ServerAddress string `yaml:"server_address"`
}

var (
	// Флаг для адреса сервера.
	FlagServerClientAddress string
)

func NewConfigServerClient() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./internal/client/config")

	var cfg Config
	if err := viper.ReadInConfig(); err != nil {
		cfg.ServerAddress = "localhost:8080"
		//return nil, fmt.Errorf("failed to read config: %w", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		cfg.ServerAddress = "localhost:8080"
		//return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	flag.StringVar(&FlagServerClientAddress, "server_address", cfg.ServerAddress, "server client address")
	flag.Parse()

	if envServerClientAddress := os.Getenv("SERVER_CLIENT_ADDRESS"); envServerClientAddress != "" {
		cfg.ServerAddress = envServerClientAddress
	} else {
		cfg.ServerAddress = FlagServerClientAddress
	}

	return &cfg, nil
}
