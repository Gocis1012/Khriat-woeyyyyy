package service

import (
	"context"
	"corporate-translator-api/internal/model"
	"corporate-translator-api/internal/repository"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	minChargeAmountTHB = 20.0
	maxChargeAmountTHB = 5000.0

	// creditPerTHB is the conversion rate applied when a payment settles:
	// 1 THB paid = 1 credit granted.
	creditPerTHB = 1.0

	// webhookTimestampTolerance bounds how old/future-dated a webhook's
	// Omise-Signature-Timestamp may be, to prevent replay of a captured
	// valid signature.
	webhookTimestampTolerance = 5 * time.Minute
)

var (
	ErrInvalidAmount    = errors.New("invalid amount")
	ErrPaymentForbidden = errors.New("payment does not belong to user")
)

// omiseClient is the slice of *OmiseService this service depends on.
type omiseClient interface {
	CreatePromptPayCharge(ctx context.Context, amountTHB float64) (*OmiseQRCharge, error)
}

type PaymentService struct {
	repo          repository.PaymentRepository
	omise         omiseClient
	webhookSecret string // base64-encoded, from Omise dashboard "Webhook secret"
}

func NewPaymentService(repo repository.PaymentRepository, omise omiseClient, webhookSecret string) *PaymentService {
	return &PaymentService{repo: repo, omise: omise, webhookSecret: webhookSecret}
}

// CreateCharge creates an Omise PromptPay charge and a matching pending
// payment record. Returns the payment row and the QR code image URI.
func (s *PaymentService) CreateCharge(ctx context.Context, userID string, amountTHB float64) (*model.Payment, string, error) {
	if amountTHB < minChargeAmountTHB || amountTHB > maxChargeAmountTHB {
		return nil, "", ErrInvalidAmount
	}

	qr, err := s.omise.CreatePromptPayCharge(ctx, amountTHB)
	if err != nil {
		return nil, "", fmt.Errorf("PaymentService.CreateCharge omise: %w", err)
	}

	payment, err := s.repo.CreatePending(ctx, userID, amountTHB, "THB", "promptpay", qr.ChargeID)
	if err != nil {
		return nil, "", fmt.Errorf("PaymentService.CreateCharge repo: %w", err)
	}

	return payment, qr.QRCodeURI, nil
}

// GetStatus returns a payment, scoped to the requesting user.
func (s *PaymentService) GetStatus(ctx context.Context, paymentID, userID string) (*model.Payment, error) {
	payment, err := s.repo.FindByID(ctx, paymentID)
	if err != nil {
		return nil, err
	}
	if payment.UserID != userID {
		return nil, ErrPaymentForbidden
	}
	return payment, nil
}

// VerifyWebhookSignature validates the Omise-Signature-Timestamp /
// Omise-Signature headers per
// https://docs.omise.co/api-webhooks#protecting-your-endpoints.
//
// The signed payload is "<timestamp>.<raw_body>", HMAC-SHA256'd with the
// base64-decoded webhook secret and hex-encoded. During secret rotation,
// Omise-Signature may contain two comma-separated signatures — either
// matching is accepted. The timestamp is bounds-checked to reject replay of
// an old captured signature.
func (s *PaymentService) VerifyWebhookSignature(rawBody []byte, timestampHeader, signatureHeader string) bool {
	if timestampHeader == "" || signatureHeader == "" || s.webhookSecret == "" {
		return false
	}

	ts, err := strconv.ParseInt(timestampHeader, 10, 64)
	if err != nil {
		return false
	}
	if age := time.Since(time.Unix(ts, 0)); age > webhookTimestampTolerance || age < -webhookTimestampTolerance {
		return false
	}

	secret, err := base64.StdEncoding.DecodeString(s.webhookSecret)
	if err != nil {
		return false
	}

	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(timestampHeader + "." + string(rawBody)))
	expected := hex.EncodeToString(mac.Sum(nil))

	for _, candidate := range strings.Split(signatureHeader, ",") {
		if hmac.Equal([]byte(expected), []byte(strings.TrimSpace(candidate))) {
			return true
		}
	}
	return false
}

type omiseWebhookEvent struct {
	Data struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	} `json:"data"`
}

// ProcessWebhook parses an Omise event (the caller must have already
// restricted the request to Omise's allowlisted source IPs), checks
// idempotency via provider_tx_id, and applies the credit update
// transactionally when the charge succeeded.
func (s *PaymentService) ProcessWebhook(ctx context.Context, rawBody []byte) error {
	var event omiseWebhookEvent
	if err := json.Unmarshal(rawBody, &event); err != nil {
		return fmt.Errorf("ProcessWebhook decode: %w", err)
	}
	if event.Data.ID == "" {
		return fmt.Errorf("ProcessWebhook: missing charge id")
	}

	payment, err := s.repo.FindByProviderTxID(ctx, event.Data.ID)
	if err != nil {
		if errors.Is(err, repository.ErrPaymentNotFound) {
			// No local record for this charge — nothing to reconcile.
			return nil
		}
		return fmt.Errorf("ProcessWebhook lookup: %w", err)
	}

	// Idempotency: this charge was already processed.
	if payment.WebhookProcessedAt != nil {
		return nil
	}

	var newStatus string
	switch event.Data.Status {
	case "successful":
		newStatus = "success"
	case "failed", "expired":
		newStatus = "failed"
	default:
		// Intermediate status (e.g. "pending") — nothing to reconcile yet.
		return nil
	}

	creditDelta := 0.0
	if newStatus == "success" {
		creditDelta = payment.Amount * creditPerTHB
	}

	if err := s.repo.MarkProcessed(ctx, payment.ID, newStatus, creditDelta); err != nil {
		return fmt.Errorf("ProcessWebhook mark: %w", err)
	}
	return nil
}
