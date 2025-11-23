package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Auth struct {
		JwtSignToken string `env:"JWT_SIGN_TOKEN"`
	}
	Limiter struct {
		MaxReq int `env:"LIMITER_MAX_REQ" envDefault:"8"`
		Window int `env:"LIMITER_WINDOW" envDefault:"60"`
	}
}

func Load() (*Config, error) {
	if err := godotenv.Load("../.env"); err != nil {
		log.Printf("no .env file found, using system envs")
	}

	cfg := Config{}

	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("config parsing error: %w", err)
	}

	return &cfg, nil
}
