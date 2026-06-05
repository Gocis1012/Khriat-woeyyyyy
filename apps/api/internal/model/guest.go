package model

import "time"

type Guest struct {
	ID        string    `json:"id"`
	Credit    int       `json:"credit"`
	CreatedAt time.Time `json:"created_at"`
}
