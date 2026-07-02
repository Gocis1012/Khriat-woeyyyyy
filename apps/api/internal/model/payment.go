package model

import "time"

type Payment struct {
	ID                 string     `json:"id"`
	UserID             string     `json:"userId"`
	Amount             float64    `json:"amount"`
	Currency           string     `json:"currency"`
	Method             string     `json:"method"`
	Status             string     `json:"status"`
	ProviderTxID       *string    `json:"providerTxId,omitempty"`
	WebhookProcessedAt *time.Time `json:"webhookProcessedAt,omitempty"`
	CreatedAt          time.Time  `json:"createdAt"`
}
