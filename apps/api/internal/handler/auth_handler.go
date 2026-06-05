package handler

import (
	"corporate-translator-api/internal/model"
	"corporate-translator-api/internal/repository/users"
	"corporate-translator-api/internal/service"
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService *service.AuthService
	userService service.UserService
}

func NewAuthHandler(authSvc *service.AuthService, userSvc service.UserService) *AuthHandler {
	return &AuthHandler{
		authService: authSvc,
		userService: userSvc,
	}
}

// GoogleLogin handles POST /api/v1/auth/google
// Expects { "idToken": "..." } from Google One Tap.
func (h *AuthHandler) GoogleLogin(c *fiber.Ctx) error {
	var req model.GoogleAuthRequest
	if err := c.BodyParser(&req); err != nil || req.IDToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "idToken is required",
		})
	}

	// 1. Verify the Google ID token
	profile, err := h.authService.ValidateGoogleToken(c.Context(), req.IDToken)
	if err != nil {
		slog.Error("Google token validation failed", "error", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid Google token",
		})
	}

	// 2. Find or create user
	user, err := h.userService.FindByGoogleID(c.Context(), profile.GoogleID)
	if err != nil {
		if !errors.Is(err, users.ErrUserNotFound) {
			slog.Error("FindByGoogleID failed", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		}

		// New user — create with 10 credits (DB default)
		user = &model.User{
			GoogleID:  profile.GoogleID,
			Email:     profile.Email,
			Username:  profile.Username,
			AvatarURL: &profile.AvatarURL,
		}
		if err := h.userService.Insert(c.Context(), user); err != nil {
			slog.Error("Failed to create user", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to create user",
			})
		}
	}

	// 3. Generate JWT
	token, err := h.authService.GenerateJWT(user.ID, user.Email)
	if err != nil {
		slog.Error("Failed to generate JWT", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to generate token",
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.AuthResponse{
		Token: token,
		User:  *user,
	})
}

// GetMe handles GET /api/v1/user/me (requires auth middleware)
func (h *AuthHandler) GetMe(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	user, err := h.userService.FindByID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	return c.JSON(user)
}
