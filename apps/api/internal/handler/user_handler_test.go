package handler

import (
	"context"
	"corporate-translator-api/internal/model"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func doInsert(t *testing.T, h *UserHandler, body string) int {
	t.Helper()
	app := testApp(func(a *fiber.App) { a.Post("/u", h.Insert) })
	req := httptest.NewRequest("POST", "/u", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	return resp.StatusCode
}

func TestUserHandler_Insert_InvalidBody(t *testing.T) {
	h := NewUserHandler(&fakeUserSvc{
		insert: func(_ context.Context, _ *model.User) error { return nil },
	})
	if code := doInsert(t, h, "{bad"); code != 400 {
		t.Errorf("status = %d, want 400", code)
	}
}

func TestUserHandler_Insert_Success(t *testing.T) {
	h := NewUserHandler(&fakeUserSvc{
		insert: func(_ context.Context, _ *model.User) error { return nil },
	})
	if code := doInsert(t, h, `{"email":"a@b.com"}`); code != 201 {
		t.Errorf("status = %d, want 201", code)
	}
}

func TestUserHandler_Insert_ServiceError(t *testing.T) {
	h := NewUserHandler(&fakeUserSvc{
		insert: func(_ context.Context, _ *model.User) error { return errors.New("invalid") },
	})
	if code := doInsert(t, h, `{"email":""}`); code != 400 {
		t.Errorf("status = %d, want 400", code)
	}
}
