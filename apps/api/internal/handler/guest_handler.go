package handler

import (
	"corporate-translator-api/internal/repository"
	"corporate-translator-api/internal/service"
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

type GuestHandler struct {
	guestService     guestSvc
	userService      service.UserService
	translateService translatorSvc
}

type translateGuestRequest struct {
	Text   string `json:"text"`
	Level  int    `json:"level"`  // 1-5, default 3
	Target string `json:"target"` // boss | client | teacher | friend
	Lang   string `json:"lang"`   // "th" | "en", default "th"
}

func NewGuestHandler(
	guestService guestSvc,
	userService service.UserService,
	translateService translatorSvc,
) *GuestHandler {
	return &GuestHandler{
		guestService:     guestService,
		userService:      userService,
		translateService: translateService,
	}
}

func (h *GuestHandler) GetStatus(c *fiber.Ctx) error {
	// If logged-in user, return their credit
	if userID, ok := c.Locals("user_id").(string); ok && userID != "" {
		user, err := h.userService.FindByID(c.Context(), userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to get user status",
			})
		}
		return c.JSON(fiber.Map{
			"user_id":    user.ID,
			"credit":     user.Credit,
			"username":   user.Username,
			"created_at": user.CreatedAt,
			"logged_in":  true,
		})
	}

	// Guest flow
	guestID, ok := c.Locals("guest_id").(string)
	if !ok || guestID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "guest session not found",
		})
	}

	guest, err := h.guestService.GetStatus(c.Context(), guestID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get guest status",
		})
	}

	return c.JSON(fiber.Map{
		"guest_id":   guestID,
		"credit":     guest.Credit,
		"created_at": guest.CreatedAt,
		"logged_in":  false,
	})
}

func (h *GuestHandler) Translate(c *fiber.Ctx) error {
	var req translateGuestRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}
	if req.Text == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "text is required",
		})
	}
	if len([]rune(req.Text)) > 3000 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "text exceeds 3000 character limit",
		})
	}
	if req.Level == 0 {
		req.Level = 3 // default
	}

	// ── Logged-in user flow ──────────────────────
	if userID, ok := c.Locals("user_id").(string); ok && userID != "" {
		// Deduct credit from Postgres
		if err := h.userService.DeductCredit(c.Context(), userID, 1.0); err != nil {
			return c.Status(fiber.StatusPaymentRequired).JSON(fiber.Map{
				"error":   "insufficient_credit",
				"message": "เครดิตหมดแล้ว",
			})
		}

		result, err := h.translateService.PurifyText(c.Context(), req.Text, req.Target, req.Level, req.Lang)
		if err != nil {
			slog.Error("Translation failed for user", "error", err, "user_id", userID)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "translation failed",
			})
		}

		return c.JSON(fiber.Map{
			"result": result,
			"level":  req.Level,
			"target": req.Target,
		})
	}

	// ── Guest flow ───────────────────────────────
	guestID, ok := c.Locals("guest_id").(string)
	if !ok || guestID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "guest session not found",
		})
	}

	var result string
	err := h.guestService.UseCredit(c.Context(), guestID, func() error {
		var err error
		result, err = h.translateService.PurifyText(c.Context(), req.Text, req.Target, req.Level, req.Lang)
		return err
	})

	if err != nil {
		if errors.Is(err, repository.ErrInsufficientCredit) {
			return c.Status(fiber.StatusPaymentRequired).JSON(fiber.Map{
				"error":   "insufficient_credit",
				"message": "เครดิตหมดแล้ว กรุณา Login เพื่อรับเครดิตเพิ่ม",
			})
		}
		if errors.Is(err, repository.ErrGuestNotFound) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "guest session expired",
			})
		}

		slog.Error("Translation failed", "error", err, "guest_id", guestID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "translation failed",
		})
	}

	return c.JSON(fiber.Map{
		"result": result,
		"level":  req.Level,
		"target": req.Target,
	})
}
