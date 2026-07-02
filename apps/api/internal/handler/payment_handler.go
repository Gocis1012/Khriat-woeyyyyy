package handler

import (
	"corporate-translator-api/internal/repository"
	"corporate-translator-api/internal/service"
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

// Header names per docs.omise.co/api-webhooks#protecting-your-endpoints.
const (
	omiseSignatureHeader          = "Omise-Signature"
	omiseSignatureTimestampHeader = "Omise-Signature-Timestamp"
)

type PaymentHandler struct {
	paymentService paymentSvc
}

func NewPaymentHandler(paymentService paymentSvc) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}

type createChargeRequest struct {
	Amount float64 `json:"amount"`
}

// CreateCharge handles POST /payments/create (requires auth).
func (h *PaymentHandler) CreateCharge(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	var req createChargeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	payment, qrCodeURI, err := h.paymentService.CreateCharge(c.Context(), userID, req.Amount)
	if err != nil {
		if errors.Is(err, service.ErrInvalidAmount) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "amount must be between 20 and 5000 THB",
			})
		}
		slog.Error("CreateCharge failed", "error", err, "user_id", userID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create payment",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"paymentId": payment.ID,
		"status":    payment.Status,
		"amount":    payment.Amount,
		"currency":  payment.Currency,
		"qrCodeUri": qrCodeURI,
	})
}

// GetStatus handles GET /payments/:id/status (requires auth).
func (h *PaymentHandler) GetStatus(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	paymentID := c.Params("id")
	payment, err := h.paymentService.GetStatus(c.Context(), paymentID, userID)
	if err != nil {
		if errors.Is(err, repository.ErrPaymentNotFound) || errors.Is(err, service.ErrPaymentForbidden) {
			// Same response for "not found" and "not yours" — avoid leaking existence.
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "payment not found",
			})
		}
		slog.Error("GetStatus failed", "error", err, "payment_id", paymentID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get payment status",
		})
	}

	return c.JSON(fiber.Map{
		"paymentId": payment.ID,
		"status":    payment.Status,
		"amount":    payment.Amount,
	})
}

// HandleWebhook handles POST /webhooks/omise (public, no user auth).
// Request source is additionally restricted upstream by the
// OmiseIPAllowlist middleware (defense in depth). Uses the raw request body
// for HMAC signature verification — the body must not be parsed as JSON
// before the signature check.
func (h *PaymentHandler) HandleWebhook(c *fiber.Ctx) error {
	rawBody := c.Body()
	timestamp := c.Get(omiseSignatureTimestampHeader)
	signature := c.Get(omiseSignatureHeader)

	if !h.paymentService.VerifyWebhookSignature(rawBody, timestamp, signature) {
		slog.Warn("Omise webhook signature mismatch")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if err := h.paymentService.ProcessWebhook(c.Context(), rawBody); err != nil {
		slog.Error("Omise webhook processing failed", "error", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}
