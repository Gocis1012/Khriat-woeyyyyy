package routes

import (
	"corporate-translator-api/internal/handler"
	"corporate-translator-api/internal/middleware"
	"corporate-translator-api/internal/service"

	"github.com/gofiber/fiber/v2"
)

func Setup(
	app *fiber.App,
	userHandler *handler.UserHandler,
	authHandler *handler.AuthHandler,
	authService *service.AuthService,
) {
	api := app.Group("/api/v1")

	// Auth routes (public)
	auth := api.Group("/auth")
	auth.Post("/google", authHandler.GoogleLogin)

	// User routes (protected)
	user := api.Group("/user")
	user.Use(middleware.RequireAuth(authService))
	user.Get("/me", authHandler.GetMe)
}
