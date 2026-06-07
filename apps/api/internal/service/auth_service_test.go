package service

import (
	"context"
	"errors"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/api/idtoken"
)

func TestAuthService_JWT_RoundTrip(t *testing.T) {
	svc := NewAuthService("secret", "client-id")

	token, err := svc.GenerateJWT("user-1", "a@b.com")
	if err != nil {
		t.Fatalf("GenerateJWT error: %v", err)
	}
	if token == "" {
		t.Fatal("empty token")
	}

	claims, err := svc.ParseJWT(token)
	if err != nil {
		t.Fatalf("ParseJWT error: %v", err)
	}
	if claims.UserID != "user-1" {
		t.Errorf("UserID = %q", claims.UserID)
	}
	if claims.Email != "a@b.com" {
		t.Errorf("Email = %q", claims.Email)
	}
}

func TestAuthService_ParseJWT_Invalid(t *testing.T) {
	svc := NewAuthService("secret", "client-id")
	if _, err := svc.ParseJWT("not-a-jwt"); err == nil {
		t.Error("expected error for garbage token")
	}
}

func TestAuthService_ParseJWT_WrongSecret(t *testing.T) {
	signer := NewAuthService("secret-A", "client-id")
	verifier := NewAuthService("secret-B", "client-id")

	token, _ := signer.GenerateJWT("u", "e")
	if _, err := verifier.ParseJWT(token); err == nil {
		t.Error("expected signature verification to fail with different secret")
	}
}

func TestAuthService_ParseJWT_WrongSigningMethod(t *testing.T) {
	svc := NewAuthService("secret", "client-id")
	// Token with "none" alg should be rejected by the HMAC check.
	tok := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{"userId": "x"})
	s, _ := tok.SignedString(jwt.UnsafeAllowNoneSignatureType)
	if _, err := svc.ParseJWT(s); err == nil {
		t.Error("expected rejection of non-HMAC token")
	}
}

func TestAuthService_ValidateGoogleToken_Success(t *testing.T) {
	svc := NewAuthService("secret", "client-id")
	svc.validate = func(_ context.Context, _, _ string) (*idtoken.Payload, error) {
		return &idtoken.Payload{Claims: map[string]interface{}{
			"email":   "user@gmail.com",
			"name":    "User Name",
			"picture": "http://pic",
			"sub":     "google-123",
		}}, nil
	}

	profile, err := svc.ValidateGoogleToken(context.Background(), "fake-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if profile.GoogleID != "google-123" || profile.Email != "user@gmail.com" {
		t.Errorf("profile = %+v", profile)
	}
	if profile.Username != "User Name" || profile.AvatarURL != "http://pic" {
		t.Errorf("profile = %+v", profile)
	}
}

func TestAuthService_ValidateGoogleToken_ValidatorError(t *testing.T) {
	svc := NewAuthService("secret", "client-id")
	svc.validate = func(_ context.Context, _, _ string) (*idtoken.Payload, error) {
		return nil, errors.New("invalid")
	}
	if _, err := svc.ValidateGoogleToken(context.Background(), "x"); err == nil {
		t.Error("expected error from validator")
	}
}

func TestAuthService_ValidateGoogleToken_MissingClaims(t *testing.T) {
	svc := NewAuthService("secret", "client-id")
	svc.validate = func(_ context.Context, _, _ string) (*idtoken.Payload, error) {
		return &idtoken.Payload{Claims: map[string]interface{}{"name": "No Email"}}, nil
	}
	if _, err := svc.ValidateGoogleToken(context.Background(), "x"); err == nil {
		t.Error("expected error when email/sub missing")
	}
}
