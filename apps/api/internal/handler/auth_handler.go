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
	authService  authSvc
	userService  service.UserService
	guestService guestSvc
}

func NewAuthHandler(auth authSvc, userSvc service.UserService, guest guestSvc) *AuthHandler {
	return &AuthHandler{
		authService:  auth,
		userService:  userSvc,
		guestService: guest,
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

		// New user — base 10 credits + any remaining guest credits
		const baseCredit = 10.0
		initialCredit := baseCredit

		guestID, _ := c.Locals("guest_id").(string)
		if guestID != "" {
			if guest, err := h.guestService.GetStatus(c.Context(), guestID); err == nil {
				initialCredit += guest.Credit
			}
		}

		user = &model.User{
			GoogleID:  profile.GoogleID,
			Email:     profile.Email,
			Username:  profile.Username,
			AvatarURL: &profile.AvatarURL,
			Credit:    initialCredit,
		}
		if err := h.userService.Insert(c.Context(), user); err != nil {
			slog.Error("Failed to create user", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to create user",
			})
		}

		// Clean up guest session (best-effort; 24h TTL covers failures)
		if guestID != "" {
			if err := h.guestService.DeleteSession(c.Context(), guestID); err != nil {
				slog.Warn("Failed to delete guest session after registration", "guest_id", guestID, "error", err)
			}
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
