package database

import (
	"embed"
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*
var migrationFiles embed.FS

func RunMigrations(databaseURL string) error {
	// แปลง embed.FS ให้เป็น Source Driver ที่ golang-migrate เข้าใจ
	d, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		return fmt.Errorf("failed to create iofs source: %w", err)
	}

	// เริ่มต้นระบบ Migrate โดยส่ง Source และ Database URL เข้าไป
	m, err := migrate.NewWithSourceInstance("iofs", d, databaseURL)
	if err != nil {
		return fmt.Errorf("failed to initialize migrate instance: %w", err)
	}
	defer m.Close()

	// สั่งดัน Migration ขึ้นไปเวอร์ชันล่าสุด (Up)
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("Database schema is up to date. No migrations applied.")
			return nil
		}
		return fmt.Errorf("failed to run up migrations: %w", err)
	}

	log.Println("Database migrations applied successfully!")
	return nil
}
