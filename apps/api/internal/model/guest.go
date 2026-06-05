package model

import "time"

type Guest struct {
	ID        string    `json:"id"`
	Credit    float64   `json:"credit"`
	CreatedAt time.Time `json:"created_at"`
}
