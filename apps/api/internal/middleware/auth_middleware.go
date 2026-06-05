package middleware

import (
	"corporate-translator-api/internal/service"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// RequireAuth validates the JWT in the Authorization header.
// On success, it sets "user_id" and "user_email" in c.Locals.
func RequireAuth(authService *service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing authorization header",
			})
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid authorization format",
			})
		}

		claims, err := authService.ParseJWT(parts[1])
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid or expired token",
			})
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("user_email", claims.Email)

		return c.Next()
	}
}

// OptionalAuth tries to parse the JWT but does NOT block the request if absent.
// If a valid token is present, user_id is set in Locals.
// If no token or invalid token, the request proceeds as a guest.
func OptionalAuth(authService *service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Next()
		}

		claims, err := authService.ParseJWT(parts[1])
		if err != nil {
			return c.Next() // Invalid token → treat as guest
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("user_email", claims.Email)
		return c.Next()
	}
}
