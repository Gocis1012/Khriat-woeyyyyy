package handler

import (
	"context"
	"net/http"

	"corporate-translator-api/internal/service"

	"github.com/gofiber/fiber/v2"
)

type TranslatorHandler struct {
	svc *service.TranslationService
}

func NewTranslatorHandler(svc *service.TranslationService) *TranslatorHandler {
	return &TranslatorHandler{svc: svc}
}

type translateRequest struct {
	Text string `json:"text"`
}

type translateResponse struct {
	Result string `json:"result"`
}

func (h *TranslatorHandler) HandleTranslate(c *fiber.Ctx) error {
	var req translateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request payload"})
	}
	if req.Text == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "text is required"})
	}

	// call the translation service with a context derived from the request
	ctx := context.Background()
	result, err := h.svc.PurifyText(ctx, req.Text)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(translateResponse{Result: result})
}
