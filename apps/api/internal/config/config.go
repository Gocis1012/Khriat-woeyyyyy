package config

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

type Config struct {
	Port           string
	DatabaseURL    string
	RedisURL       string
	FrontendOrigin string
	JWTSecret      string
	GoogleClientID string
	GoogleAPIKey   string
	OpenAIAPIKey   string
	OpenAIModel    string
	AutoMigrate    bool
}

func Load() (Config, error) {
	// Try common .env locations so running from different working dirs works
	_ = loadDotEnv(".env")
	_ = loadDotEnv("../.env")
	_ = loadDotEnv("../../.env")

	cfg := Config{
		Port: getEnv("PORT", "8080"),
		// ⭐️ เปลี่ยนจาก os.Getenv("DATABASE_URL") เป็นบรรทัดล่างนี้:
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		RedisURL:       getEnv("REDIS_URL", "localhost:6379"),
		FrontendOrigin: getEnv("FRONTEND_ORIGIN", "http://localhost:3000"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		GoogleClientID: os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleAPIKey:   os.Getenv("GOOGLE_API_KEY"),
		OpenAIAPIKey:   os.Getenv("OPENAI_API_KEY"),
		OpenAIModel:    getEnv("OPENAI_MODEL", "gpt-4o-mini"),
		AutoMigrate:    getEnv("AUTO_MIGRATE", "true") == "true",
	}

	if cfg.DatabaseURL == "" {
		return cfg, errors.New("DATABASE_URL is required")
	}
	if cfg.JWTSecret == "" {
		return cfg, errors.New("JWT_SECRET is required")
	}
	if cfg.GoogleClientID == "" {
		return cfg, errors.New("GOOGLE_CLIENT_ID is required")
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
