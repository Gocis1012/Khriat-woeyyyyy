package config

import (
	"os"
	"path/filepath"
	"testing"
)

// setRequired sets the three required env vars for a successful Load().
func setRequired(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://localhost:5432/db")
	t.Setenv("JWT_SECRET", "test-secret")
	t.Setenv("DEEPSEEK_API_KEY", "sk-test")
}

func TestLoad_Success_Defaults(t *testing.T) {
	setRequired(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.DatabaseURL != "postgres://localhost:5432/db" {
		t.Errorf("DatabaseURL = %q", cfg.DatabaseURL)
	}
	if cfg.Port != "8080" {
		t.Errorf("default Port = %q, want 8080", cfg.Port)
	}
	if cfg.RedisURL != "localhost:6379" {
		t.Errorf("default RedisURL = %q", cfg.RedisURL)
	}
	if cfg.FrontendOrigin != "http://localhost:3000" {
		t.Errorf("default FrontendOrigin = %q", cfg.FrontendOrigin)
	}
	if cfg.AppEnv != "development" {
		t.Errorf("default AppEnv = %q", cfg.AppEnv)
	}
	if !cfg.AutoMigrate {
		t.Error("AutoMigrate should default to true")
	}
}

func TestLoad_Overrides(t *testing.T) {
	setRequired(t)
	t.Setenv("PORT", "9999")
	t.Setenv("REDIS_URL", "rediss://upstash:6380")
	t.Setenv("APP_ENV", "production")
	t.Setenv("AUTO_MIGRATE", "false")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Port != "9999" {
		t.Errorf("Port = %q", cfg.Port)
	}
	if cfg.RedisURL != "rediss://upstash:6380" {
		t.Errorf("RedisURL = %q", cfg.RedisURL)
	}
	if cfg.AppEnv != "production" {
		t.Errorf("AppEnv = %q", cfg.AppEnv)
	}
	if cfg.AutoMigrate {
		t.Error("AutoMigrate should be false")
	}
}

func TestLoad_MissingRequired(t *testing.T) {
	cases := []struct {
		name      string
		database  string
		jwt       string
		deepseek  string
		wantError bool
	}{
		{"all present", "db", "jwt", "ds", false},
		{"missing database", "", "jwt", "ds", true},
		{"missing jwt", "db", "", "ds", true},
		{"missing deepseek", "db", "jwt", "", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("DATABASE_URL", tc.database)
			t.Setenv("JWT_SECRET", tc.jwt)
			t.Setenv("DEEPSEEK_API_KEY", tc.deepseek)

			_, err := Load()
			if tc.wantError && err == nil {
				t.Error("expected error, got nil")
			}
			if !tc.wantError && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestGetEnv(t *testing.T) {
	t.Setenv("PRESENT_KEY", "value")
	if got := getEnv("PRESENT_KEY", "fallback"); got != "value" {
		t.Errorf("getEnv present = %q, want value", got)
	}
	if got := getEnv("ABSENT_KEY_XYZ", "fallback"); got != "fallback" {
		t.Errorf("getEnv absent = %q, want fallback", got)
	}
}

func TestLoadDotEnv(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	content := "# a comment\n\nFOO_TEST_VAR=hello\nBAR_TEST_VAR=\"quoted\"\nMALFORMED_LINE\n"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	// Ensure clean slate
	os.Unsetenv("FOO_TEST_VAR")
	os.Unsetenv("BAR_TEST_VAR")
	t.Cleanup(func() {
		os.Unsetenv("FOO_TEST_VAR")
		os.Unsetenv("BAR_TEST_VAR")
	})

	if err := loadDotEnv(path); err != nil {
		t.Fatalf("loadDotEnv error: %v", err)
	}
	if os.Getenv("FOO_TEST_VAR") != "hello" {
		t.Errorf("FOO_TEST_VAR = %q", os.Getenv("FOO_TEST_VAR"))
	}
	if os.Getenv("BAR_TEST_VAR") != "quoted" {
		t.Errorf("BAR_TEST_VAR = %q (quotes should be stripped)", os.Getenv("BAR_TEST_VAR"))
	}
}

func TestLoadDotEnv_MissingFile(t *testing.T) {
	if err := loadDotEnv("/nonexistent/path/.env"); err == nil {
		t.Error("expected error for missing file")
	}
}
