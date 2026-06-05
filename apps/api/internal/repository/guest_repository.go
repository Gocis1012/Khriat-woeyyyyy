package repository

import (
	"context"
	"corporate-translator-api/internal/model"
)

type GuestRepository interface {
	GetOrCreate(ctx context.Context, guestID string) (*model.Guest, error)
	DeductCredit(ctx context.Context, guestID string, amount float64) error
	RefundCredit(ctx context.Context, guestID string, amount float64) error
	Delete(ctx context.Context, guestID string) error
}

