package handler

import (
	"corporate-translator-api/internal/repository"
	"corporate-translator-api/internal/service"
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

type GuestHandler struct {
	guestService     *service.GuestService
	translateService *service.TranslationService
}

type translateGuestRequest struct {
	Text string `json:"text"`
}

func NewGuestHandler(guestService *service.GuestService, translateService *service.TranslationService) *GuestHandler {
	return &GuestHandler{
		guestService:     guestService,
		translateService: translateService,
	}
}

func (h *GuestHandler) GetStatus(c *fiber.Ctx) error {
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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"guest_id":   guestID,
		"credit":     guest.Credit,
		"created_at": guest.CreatedAt,
	})
}

func (h *GuestHandler) Translate(c *fiber.Ctx) error {
	guestID, ok := c.Locals("guest_id").(string)

	if !ok || guestID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "guest session not found",
		})
	}

	var req translateGuestRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	var result string

	err := h.guestService.UseCredit(c.Context(), guestID, func() error {
		var err error
		result, err = h.translateService.PurifyText(c.Context(), req.Text)
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

		// เปลี่ยนเป็นแบบนี้
		slog.Error("Translation failed", "error", err, "guest_id", guestID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "translation failed",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"result": result,
	})
}
