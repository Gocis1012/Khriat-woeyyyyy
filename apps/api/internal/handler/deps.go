package handler

import (
	"context"

	"corporate-translator-api/internal/model"
)

// These consumer-defined interfaces describe the slice of each service the
// handlers actually use. The concrete *service.* types satisfy them, so
// main.go is unchanged, while tests can inject lightweight fakes.

type guestSvc interface {
	GetStatus(ctx context.Context, guestID string) (*model.Guest, error)
	UseCredit(ctx context.Context, guestID string, fn func() error) error
	DeleteSession(ctx context.Context, guestID string) error
}

type translatorSvc interface {
	PurifyText(ctx context.Context, text, target string, level int, lang string) (string, error)
}

type authSvc interface {
	ValidateGoogleToken(ctx context.Context, idToken string) (*model.GoogleProfile, error)
	GenerateJWT(userID, email string) (string, error)
}

type paymentSvc interface {
	CreateCharge(ctx context.Context, userID string, amountTHB float64) (*model.Payment, string, error)
	GetStatus(ctx context.Context, paymentID, userID string) (*model.Payment, error)
	VerifyWebhookSignature(rawBody []byte, timestampHeader, signatureHeader string) bool
	ProcessWebhook(ctx context.Context, rawBody []byte) error
}
