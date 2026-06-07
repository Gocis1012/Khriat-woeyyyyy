package service

import (
	"context"
	"corporate-translator-api/internal/model"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/api/idtoken"
)

// googleValidateFunc matches idtoken.Validate so it can be swapped in tests.
type googleValidateFunc func(ctx context.Context, idToken, audience string) (*idtoken.Payload, error)

type AuthService struct {
	jwtSecret      string
	googleClientID string
	validate       googleValidateFunc
}

type JWTClaims struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func NewAuthService(jwtSecret, googleClientID string) *AuthService {
	return &AuthService{
		jwtSecret:      jwtSecret,
		googleClientID: googleClientID,
		validate:       idtoken.Validate,
	}
}

// ValidateGoogleToken verifies a Google ID token and extracts profile info.
func (s *AuthService) ValidateGoogleToken(ctx context.Context, idToken string) (*model.GoogleProfile, error) {
	payload, err := s.validate(ctx, idToken, s.googleClientID)
	if err != nil {
		return nil, fmt.Errorf("invalid Google token: %w", err)
	}

	email, _ := payload.Claims["email"].(string)
	name, _ := payload.Claims["name"].(string)
	picture, _ := payload.Claims["picture"].(string)
	sub, _ := payload.Claims["sub"].(string)

	if email == "" || sub == "" {
		return nil, fmt.Errorf("Google token missing email or sub")
	}

	return &model.GoogleProfile{
		GoogleID:  sub,
		Email:     email,
		Username:  name,
		AvatarURL: picture,
	}, nil
}

// GenerateJWT creates a signed JWT with a 7-day expiry.
func (s *AuthService) GenerateJWT(userID, email string) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// ParseJWT validates a JWT string and returns the claims.
func (s *AuthService) ParseJWT(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}
