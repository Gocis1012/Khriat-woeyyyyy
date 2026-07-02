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
	paymentHandler *handler.PaymentHandler,
) {
	api := app.Group("/api/v1")

	// Auth routes (public)
	auth := api.Group("/auth")
	auth.Post("/google", authHandler.GoogleLogin)

	// User routes (protected)
	user := api.Group("/user")
	user.Use(middleware.RequireAuth(authService))
	user.Get("/me", authHandler.GetMe)

	// Payment routes (protected) — guests must log in before charging.
	payments := api.Group("/payments")
	payments.Use(middleware.RequireAuth(authService))
	payments.Post("/create", paymentHandler.CreateCharge)
	payments.Get("/:id/status", paymentHandler.GetStatus)
}
