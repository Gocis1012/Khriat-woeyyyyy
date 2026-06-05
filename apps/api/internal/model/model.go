package model

import "time"

type User struct {
	ID                string     `json:"id"`
	GoogleID          string     `json:"googleID"`
	Email             string     `json:"email"`
	Username          string     `json:"username"`
	AvatarURL         *string    `json:"avatarUrl"`
	Credit            float64    `json:"credit"`
	MemberType        string     `json:"memberType"`
	LastDailyCreditAt *time.Time `json:"lastDailyCreditAt"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
}

type GoogleProfile struct {
	GoogleID  string
	Email     string
	Username  string
	AvatarURL string
}

type GoogleAuthRequest struct {
	IDToken string `json:"idToken"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type CreditResponse struct {
	Credit  float64 `json:"credit"`
	Claimed bool    `json:"claimed,omitempty"`
}

type TranslateRequest struct {
	Input string `json:"input"`
	Tone  string `json:"tone"`
}

type TranslateResponse struct {
	ID             string  `json:"id"`
	TranslatedText string  `json:"translatedText"`
	Credit         float64 `json:"credit"`
	CreditUsed     float64 `json:"creditUsed"`
}

type TransWord struct {
	ID             string    `json:"id"`
	UserID         string    `json:"userId"`
	RawText        string    `json:"rawText"`
	TranslatedText *string   `json:"translatedText"`
	ToneMode       string    `json:"toneMode"`
	CreditUsed     float64   `json:"creditUsed"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
