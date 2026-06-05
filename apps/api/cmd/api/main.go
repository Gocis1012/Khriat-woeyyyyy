package main

import (
	"context"
	"corporate-translator-api/internal/config"
	"corporate-translator-api/internal/database"
	"corporate-translator-api/internal/middleware"
	"corporate-translator-api/internal/repository"
	"corporate-translator-api/internal/repository/users"
	"corporate-translator-api/internal/routes"
	"time"

	"fmt"
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"

	"corporate-translator-api/internal/handler"
	"corporate-translator-api/internal/service"
)

func main() {
    wd, _ := os.Getwd()
    slog.Info("Current Working Directory", "path", wd)

    ctx := context.Background()
    env, err := config.Load()
    if err != nil {
        slog.Error("Failed to load config", "error", err)
        os.Exit(1)
    }
    slog.Info("Config loaded", "database_url", env.DatabaseURL)

    // ── Database ─────────────────────────────────────────
    db, err := database.Connect(ctx, env.DatabaseURL)
    if err != nil {
        slog.Error("Failed to connect to database", "error", err)
        os.Exit(1)
    }
    defer db.Close()
    slog.Info("Successfully connected to PostgreSQL")

    if env.AutoMigrate {
        slog.Info("Running database migrations...")
        if err := database.RunMigrations(env.DatabaseURL); err != nil {
            slog.Error("Failed to run migrations", "error", err)
            os.Exit(1)
        }
    }

    // ── Redis ─────────────────────────────────────────────
    redisClient, err := database.NewRedisClient(env.RedisURL)
    if err != nil {
        slog.Error("Failed to connect to Redis", "error", err)
        os.Exit(1)
    }
    slog.Info("Successfully connected to Redis")

    // ── Services ──────────────────────────────────────────
    userRepo     := users.NewPostgresRepository(db)
    userService  := service.NewUserService(userRepo)
    userHandler  := handler.NewUserHandler(userService)

    guestRepo    := repository.NewGuestRepository(redisClient)
    guestSvc     := service.NewGuestService(guestRepo)

    translationSvc, err := service.NewTranslationService(env.DeepSeekAPIKey)
    if err != nil {
        slog.Error("Failed to initialize translation service", "error", err)
        os.Exit(1)
    }

    guestHandler := handler.NewGuestHandler(guestSvc, translationSvc)

    // ── Fiber ─────────────────────────────────────────────
    app := fiber.New()

    // Middleware — ลำดับสำคัญ
    app.Use(cors.New(cors.Config{
        AllowOrigins:     env.FrontendOrigin,
        AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
        AllowHeaders:     "Content-Type,Authorization",
        AllowCredentials: true, // required for cookies with credentials:"include"
    }))
    app.Use(middleware.GuestSession(env.AppEnv))
    app.Use("/translate", limiter.New(limiter.Config{
        Max:        10,
        Expiration: 1 * time.Minute,
    }))

    // Routes
    app.Get("/health", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{"status": "ok"})
    })
    app.Get("/guest/status", guestHandler.GetStatus)
    app.Post("/translate",   guestHandler.Translate)
    routes.Setup(app, userHandler)

    // ── Start ─────────────────────────────────────────────
    slog.Info("Server listening", "port", env.Port)
    if err := app.Listen(fmt.Sprintf(":%s", env.Port)); err != nil {
        slog.Error("Server error", "error", err)
        os.Exit(1)
    }
}