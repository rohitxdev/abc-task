package config

import (
	"fmt"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
)

type Config struct {
	Env             string `validate:"required,oneof=development production"`
	Host            string `validate:"required,ip"`
	Port            string `validate:"required,number"`
	DatabaseURL     string `validate:"required,filepath"`
	ShutdownTimeout time.Duration
}

func Load() (*Config, error) {
	shutdownTimeout, err := time.ParseDuration(os.Getenv("SHUTDOWN_TIMEOUT"))
	if err != nil {
		return nil, fmt.Errorf("Failed to parse SHUTDOWN_TIMEOUT: %w", err)
	}

	cfg := Config{
		Env:             os.Getenv("ENV"),
		Host:            os.Getenv("HOST"),
		Port:            os.Getenv("PORT"),
		DatabaseURL:     os.Getenv("DATABASE_URL"),
		ShutdownTimeout: shutdownTimeout,
	}

	if err = validator.New().Struct(cfg); err != nil {
		return nil, fmt.Errorf("Failed to validate config: %w", err)
	}

	return &cfg, nil
}
