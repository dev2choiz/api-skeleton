package config

import (
	"github.com/caarlos0/env/v11"
)

// Config represents application environment configuration loaded from environment variables.
type Config struct {
	APIPort          string `env:"API_PORT,required"`
	JWTSecret        string `env:"JWT_SECRET,required"`
	PostgresUser     string `env:"POSTGRES_USER,required"`
	PostgresPassword string `env:"POSTGRES_PASSWORD,required"`
	PostgresDatabase string `env:"POSTGRES_DATABASE,required"`
	PostgresHost     string `env:"POSTGRES_HOST,required"`
	PostgresPort     int    `env:"POSTGRES_PORT,required"`
	RedisHost        string `env:"REDIS_HOST,required"`
	RedisPort        int    `env:"REDIS_PORT,required"`
}

func Load() (*Config, error) {
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
