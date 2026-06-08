package handler

import (
	"context"
	"corporate-translator-api/internal/model"
	"corporate-translator-api/internal/repository"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func okTranslator() *fakeTranslator {
	return &fakeTranslator{purify: func(_ context.Context, _, _ string, _ int, _ string) (string, error) {
		return "สุภาพแล้ว", nil
	}}
}

// ── GetStatus ─────────────────────────────────────────────────────────────────

func TestGetStatus_LoggedInUser(t *testing.T) {
	us := &fakeUserSvc{findByID: func(_ context.Context, _ string) (*model.User, error) {
		return &model.User{ID: "u1", Credit: 9, Username: "bob"}, nil
	}}
	h := NewGuestHandler(&fakeGuestSvc{}, us, okTranslator())
	app := testApp(func(a *fiber.App) { a.Get("/s", h.GetStatus) })

	req := httptest.NewRequest("GET", "/s", nil)
	req.Header.Set("X-User", "u1")
	resp, _ := app.Test(req)
	if resp.StatusCode != 200 {
		t.Fatalf("status = %d", resp.StatusCode)
	}
}

func TestGetStatus_LoggedInUser_Error(t *testing.T) {
	us := &fakeUserSvc{findByID: func(_ context.Context, _ string) (*model.User, error) {
		return nil, errors.New("db")
	}}
	h := NewGuestHandler(&fakeGuestSvc{}, us, okTranslator())
	app := testApp(func(a *fiber.App) { a.Get("/s", h.GetStatus) })

	req := httptest.NewRequest("GET", "/s", nil)
	req.Header.Set("X-User", "u1")
	resp, _ := app.Test(req)
	if resp.StatusCode != 500 {
		t.Errorf("status = %d, want 500", resp.StatusCode)
	}
}

func TestGetStatus_Guest(t *testing.T) {
	gs := &fakeGuestSvc{getStatus: func(_ context.Context, _ string) (*model.Guest, error) {
		return &model.Guest{Credit: 6}, nil
	}}
	h := NewGuestHandler(gs, &fakeUserSvc{}, okTranslator())
	app := testApp(func(a *fiber.App) { a.Get("/s", h.GetStatus) })

	req := httptest.NewRequest("GET", "/s", nil)
	req.Header.Set("X-Guest", "g1")
	resp, _ := app.Test(req)
	if resp.StatusCode != 200 {
		t.Errorf("status = %d", resp.StatusCode)
	}
}

func TestGetStatus_Guest_Error(t *testing.T) {
	gs := &fakeGuestSvc{getStatus: func(_ context.Context, _ string) (*model.Guest, error) {
		return nil, errors.New("redis")
	}}
	h := NewGuestHandler(gs, &fakeUserSvc{}, okTranslator())
	app := testApp(func(a *fiber.App) { a.Get("/s", h.GetStatus) })

	req := httptest.NewRequest("GET", "/s", nil)
	req.Header.Set("X-Guest", "g1")
	resp, _ := app.Test(req)
	if resp.StatusCode != 500 {
		t.Errorf("status = %d", resp.StatusCode)
	}
}

func TestGetStatus_NoSession(t *testing.T) {
	h := NewGuestHandler(&fakeGuestSvc{}, &fakeUserSvc{}, okTranslator())
	app := testApp(func(a *fiber.App) { a.Get("/s", h.GetStatus) })
	resp, _ := app.Test(httptest.NewRequest("GET", "/s", nil))
	if resp.StatusCode != 401 {
		t.Errorf("status = %d, want 401", resp.StatusCode)
	}
}

// ── Translate ─────────────────────────────────────────────────────────────────

func doTranslate(t *testing.T, h *GuestHandler, body, userHdr, guestHdr string) int {
	t.Helper()
	app := testApp(func(a *fiber.App) { a.Post("/t", h.Translate) })
	req := httptest.NewRequest("POST", "/t", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if userHdr != "" {
		req.Header.Set("X-User", userHdr)
	}
	if guestHdr != "" {
		req.Header.Set("X-Guest", guestHdr)
	}
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	return resp.StatusCode
}

func TestTranslate_InvalidBody(t *testing.T) {
	h := NewGuestHandler(&fakeGuestSvc{}, &fakeUserSvc{}, okTranslator())
	if code := doTranslate(t, h, "{not json", "", "g1"); code != 400 {
		t.Errorf("status = %d, want 400", code)
	}
}

func TestTranslate_EmptyText(t *testing.T) {
	h := NewGuestHandler(&fakeGuestSvc{}, &fakeUserSvc{}, okTranslator())
	if code := doTranslate(t, h, `{"text":""}`, "", "g1"); code != 400 {
		t.Errorf("status = %d, want 400", code)
	}
}

func TestTranslate_User_Success(t *testing.T) {
	us := &fakeUserSvc{deduct: func(_ context.Context, _ string, _ float64) error { return nil }}
	h := NewGuestHandler(&fakeGuestSvc{}, us, okTranslator())
	if code := doTranslate(t, h, `{"text":"hi","level":4,"target":"boss"}`, "u1", ""); code != 200 {
		t.Errorf("status = %d, want 200", code)
	}
}

func TestTranslate_User_InsufficientCredit(t *testing.T) {
	us := &fakeUserSvc{deduct: func(_ context.Context, _ string, _ float64) error {
		return errors.New("insufficient")
	}}
	h := NewGuestHandler(&fakeGuestSvc{}, us, okTranslator())
	if code := doTranslate(t, h, `{"text":"hi"}`, "u1", ""); code != 402 {
		t.Errorf("status = %d, want 402", code)
	}
}

func TestTranslate_User_TranslationFails(t *testing.T) {
	us := &fakeUserSvc{deduct: func(_ context.Context, _ string, _ float64) error { return nil }}
	tr := &fakeTranslator{purify: func(_ context.Context, _, _ string, _ int, _ string) (string, error) {
		return "", errors.New("ai down")
	}}
	h := NewGuestHandler(&fakeGuestSvc{}, us, tr)
	if code := doTranslate(t, h, `{"text":"hi"}`, "u1", ""); code != 500 {
		t.Errorf("status = %d, want 500", code)
	}
}

func TestTranslate_Guest_Success(t *testing.T) {
	gs := &fakeGuestSvc{useCredit: func(_ context.Context, _ string, fn func() error) error {
		return fn()
	}}
	h := NewGuestHandler(gs, &fakeUserSvc{}, okTranslator())
	if code := doTranslate(t, h, `{"text":"hi"}`, "", "g1"); code != 200 {
		t.Errorf("status = %d, want 200", code)
	}
}

func TestTranslate_Guest_Insufficient(t *testing.T) {
	gs := &fakeGuestSvc{useCredit: func(_ context.Context, _ string, _ func() error) error {
		return repository.ErrInsufficientCredit
	}}
	h := NewGuestHandler(gs, &fakeUserSvc{}, okTranslator())
	if code := doTranslate(t, h, `{"text":"hi"}`, "", "g1"); code != 402 {
		t.Errorf("status = %d, want 402", code)
	}
}

func TestTranslate_Guest_Expired(t *testing.T) {
	gs := &fakeGuestSvc{useCredit: func(_ context.Context, _ string, _ func() error) error {
		return repository.ErrGuestNotFound
	}}
	h := NewGuestHandler(gs, &fakeUserSvc{}, okTranslator())
	if code := doTranslate(t, h, `{"text":"hi"}`, "", "g1"); code != 401 {
		t.Errorf("status = %d, want 401", code)
	}
}

func TestTranslate_Guest_OtherError(t *testing.T) {
	gs := &fakeGuestSvc{useCredit: func(_ context.Context, _ string, _ func() error) error {
		return errors.New("boom")
	}}
	h := NewGuestHandler(gs, &fakeUserSvc{}, okTranslator())
	if code := doTranslate(t, h, `{"text":"hi"}`, "", "g1"); code != 500 {
		t.Errorf("status = %d, want 500", code)
	}
}

func TestTranslate_NoSession(t *testing.T) {
	h := NewGuestHandler(&fakeGuestSvc{}, &fakeUserSvc{}, okTranslator())
	if code := doTranslate(t, h, `{"text":"hi"}`, "", ""); code != 401 {
		t.Errorf("status = %d, want 401", code)
	}
}
