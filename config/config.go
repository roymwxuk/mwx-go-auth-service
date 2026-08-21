package config

import (
	"errors"
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string

	GoogleClientID     string
	GoogleClientSecret string
}

func Load() (*Config, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	port := os.Getenv("PORT")
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	if databaseURL == "" {
		return nil, errors.New("Missing config: DATABASE_URL is required")
	}

	if port == "" {
		port = "8080"
	}

	return &Config{
		DatabaseURL:        databaseURL,
		Port:               port,
		GoogleClientID:     googleClientID,
		GoogleClientSecret: googleClientSecret,
	}, nil
}
