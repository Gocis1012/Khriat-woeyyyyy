package config

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

type Config struct {
	Port                   string
	DatabaseURL            string
	RedisURL               string
	FrontendOrigin         string
	JWTSecret              string
	GoogleClientID         string
	DeepSeekAPIKey         string
	AutoMigrate            bool
	AppEnv                 string
	OmiseSecretKey         string
	OmiseWebhookAllowedIPs string
	OmiseWebhookSecret     string
}

func Load() (Config, error) {
	_ = loadDotEnv(".env")
	_ = loadDotEnv("../.env")
	_ = loadDotEnv("../../.env")

	cfg := Config{
		Port:                   getEnv("PORT", "8080"),
		DatabaseURL:            os.Getenv("DATABASE_URL"),
		RedisURL:               getEnv("REDIS_URL", "localhost:6379"),
		FrontendOrigin:         getEnv("FRONTEND_ORIGIN", "http://localhost:3000"),
		JWTSecret:              os.Getenv("JWT_SECRET"),
		GoogleClientID:         os.Getenv("GOOGLE_CLIENT_ID"),
		DeepSeekAPIKey:         os.Getenv("DEEPSEEK_API_KEY"),
		AutoMigrate:            getEnv("AUTO_MIGRATE", "true") == "true",
		AppEnv:                 getEnv("APP_ENV", "development"),
		OmiseSecretKey:         os.Getenv("OMISE_SECRET_KEY"),
		OmiseWebhookAllowedIPs: os.Getenv("OMISE_WEBHOOK_ALLOWED_IPS"),
		OmiseWebhookSecret:     os.Getenv("OMISE_WEBHOOK_SECRET"),
	}

	if cfg.DatabaseURL == "" {
		return cfg, errors.New("DATABASE_URL is required")
	}
	if cfg.JWTSecret == "" {
		return cfg, errors.New("JWT_SECRET is required")
	}
	if cfg.DeepSeekAPIKey == "" {
		return cfg, errors.New("DEEPSEEK_API_KEY is required")
	}
	if cfg.OmiseSecretKey == "" {
		return cfg, errors.New("OMISE_SECRET_KEY is required")
	}
	if cfg.OmiseWebhookAllowedIPs == "" {
		return cfg, errors.New("OMISE_WEBHOOK_ALLOWED_IPS is required")
	}
	if cfg.OmiseWebhookSecret == "" {
		return cfg, errors.New("OMISE_WEBHOOK_SECRET is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func loadDotEnv(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
		if os.Getenv(key) == "" {
			_ = os.Setenv(key, value)
		}
	}
	return scanner.Err()
}
