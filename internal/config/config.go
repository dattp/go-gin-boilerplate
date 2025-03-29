package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	Port     string `env:"PORT,required"`
	Env      string `env:"ENV,required"`
	LogLevel string `env:"LOG_LEVEL,required"`
	LogJSON  bool   `env:"LOG_JSON,required"`
	Redis    struct {
		Addr     string `env:"REDIS_ADDR,required"`
		Password string `env:"REDIS_PASSWORD"`
		DB       int    `env:"REDIS_DB" envDefault:"0"`
	}
	MongoDB struct {
		URI      string `env:"MONGODB_URI,required"`
		Database string `env:"MONGODB_DATABASE,required"`
	}
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	return cfg, nil
}
