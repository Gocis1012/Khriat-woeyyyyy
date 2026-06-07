package handler

import (
	"context"
	"corporate-translator-api/internal/model"
	"corporate-translator-api/internal/repository/users"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func okAuth() *fakeAuthSvc {
	return &fakeAuthSvc{
		validateGoogle: func(_ context.Context, _ string) (*model.GoogleProfile, error) {
			return &model.GoogleProfile{GoogleID: "g", Email: "e@x.com", Username: "n"}, nil
		},
		generateJWT: func(_, _ string) (string, error) { return "jwt-token", nil },
	}
}

func doLogin(t *testing.T, h *AuthHandler, body, guestHdr string) int {
	t.Helper()
	app := testApp(func(a *fiber.App) { a.Post("/login", h.GoogleLogin) })
	req := httptest.NewRequest("POST", "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if guestHdr != "" {
		req.Header.Set("X-Guest", guestHdr)
	}
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	return resp.StatusCode
}

func TestGoogleLogin_MissingToken(t *testing.T) {
	h := NewAuthHandler(okAuth(), &fakeUserSvc{}, &fakeGuestSvc{})
	if code := doLogin(t, h, `{}`, ""); code != 400 {
		t.Errorf("status = %d, want 400", code)
	}
}

func TestGoogleLogin_InvalidToken(t *testing.T) {
	auth := &fakeAuthSvc{
		validateGoogle: func(_ context.Context, _ string) (*model.GoogleProfile, error) {
			return nil, errors.New("bad token")
		},
	}
	h := NewAuthHandler(auth, &fakeUserSvc{}, &fakeGuestSvc{})
	if code := doLogin(t, h, `{"idToken":"x"}`, ""); code != 401 {
		t.Errorf("status = %d, want 401", code)
	}
}

func TestGoogleLogin_ExistingUser(t *testing.T) {
	us := &fakeUserSvc{findByGoogleID: func(_ context.Context, _ string) (*model.User, error) {
		return &model.User{ID: "u1", Email: "e@x.com", Credit: 9}, nil
	}}
	h := NewAuthHandler(okAuth(), us, &fakeGuestSvc{})
	if code := doLogin(t, h, `{"idToken":"x"}`, ""); code != 200 {
		t.Errorf("status = %d, want 200", code)
	}
}

func TestGoogleLogin_NewUser_WithGuestMerge(t *testing.T) {
	inserted := false
	deleted := false
	us := &fakeUserSvc{
		findByGoogleID: func(_ context.Context, _ string) (*model.User, error) {
			return nil, users.ErrUserNotFound
		},
		insert: func(_ context.Context, u *model.User) error {
			inserted = true
			u.ID = "new-user"
			return nil
		},
	}
	gs := &fakeGuestSvc{
		getStatus: func(_ context.Context, _ string) (*model.Guest, error) {
			return &model.Guest{Credit: 4}, nil
		},
		del: func(_ context.Context, _ string) error { deleted = true; return nil },
	}
	h := NewAuthHandler(okAuth(), us, gs)
	if code := doLogin(t, h, `{"idToken":"x"}`, "g1"); code != 200 {
		t.Errorf("status = %d, want 200", code)
	}
	if !inserted {
		t.Error("expected user insert")
	}
	if !deleted {
		t.Error("expected guest session deletion")
	}
}

func TestGoogleLogin_FindError(t *testing.T) {
	us := &fakeUserSvc{findByGoogleID: func(_ context.Context, _ string) (*model.User, error) {
		return nil, errors.New("db down")
	}}
	h := NewAuthHandler(okAuth(), us, &fakeGuestSvc{})
	if code := doLogin(t, h, `{"idToken":"x"}`, ""); code != 500 {
		t.Errorf("status = %d, want 500", code)
	}
}

func TestGoogleLogin_InsertError(t *testing.T) {
	us := &fakeUserSvc{
		findByGoogleID: func(_ context.Context, _ string) (*model.User, error) {
			return nil, users.ErrUserNotFound
		},
		insert: func(_ context.Context, _ *model.User) error { return errors.New("insert fail") },
	}
	h := NewAuthHandler(okAuth(), us, &fakeGuestSvc{})
	if code := doLogin(t, h, `{"idToken":"x"}`, ""); code != 500 {
		t.Errorf("status = %d, want 500", code)
	}
}

func TestGoogleLogin_JWTError(t *testing.T) {
	auth := &fakeAuthSvc{
		validateGoogle: func(_ context.Context, _ string) (*model.GoogleProfile, error) {
			return &model.GoogleProfile{GoogleID: "g", Email: "e@x.com"}, nil
		},
		generateJWT: func(_, _ string) (string, error) { return "", errors.New("sign fail") },
	}
	us := &fakeUserSvc{findByGoogleID: func(_ context.Context, _ string) (*model.User, error) {
		return &model.User{ID: "u1"}, nil
	}}
	h := NewAuthHandler(auth, us, &fakeGuestSvc{})
	if code := doLogin(t, h, `{"idToken":"x"}`, ""); code != 500 {
		t.Errorf("status = %d, want 500", code)
	}
}

func TestGetMe_Unauthorized(t *testing.T) {
	h := NewAuthHandler(okAuth(), &fakeUserSvc{}, &fakeGuestSvc{})
	app := testApp(func(a *fiber.App) { a.Get("/me", h.GetMe) })
	resp, _ := app.Test(httptest.NewRequest("GET", "/me", nil))
	if resp.StatusCode != 401 {
		t.Errorf("status = %d, want 401", resp.StatusCode)
	}
}

func TestGetMe_Success(t *testing.T) {
	us := &fakeUserSvc{findByID: func(_ context.Context, _ string) (*model.User, error) {
		return &model.User{ID: "u1"}, nil
	}}
	h := NewAuthHandler(okAuth(), us, &fakeGuestSvc{})
	app := testApp(func(a *fiber.App) { a.Get("/me", h.GetMe) })
	req := httptest.NewRequest("GET", "/me", nil)
	req.Header.Set("X-User", "u1")
	resp, _ := app.Test(req)
	if resp.StatusCode != 200 {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}
}

func TestGetMe_NotFound(t *testing.T) {
	us := &fakeUserSvc{findByID: func(_ context.Context, _ string) (*model.User, error) {
		return nil, errors.New("not found")
	}}
	h := NewAuthHandler(okAuth(), us, &fakeGuestSvc{})
	app := testApp(func(a *fiber.App) { a.Get("/me", h.GetMe) })
	req := httptest.NewRequest("GET", "/me", nil)
	req.Header.Set("X-User", "u1")
	resp, _ := app.Test(req)
	if resp.StatusCode != 404 {
		t.Errorf("status = %d, want 404", resp.StatusCode)
	}
}
