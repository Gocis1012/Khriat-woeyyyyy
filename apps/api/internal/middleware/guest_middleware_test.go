package middleware

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestGuestSession_NewGuestGetsCookie(t *testing.T) {
	app := fiber.New()
	app.Use(GuestSession("development"))
	app.Get("/", func(c *fiber.Ctx) error {
		id, _ := c.Locals("guest_id").(string)
		return c.SendString(id)
	})

	resp, err := app.Test(httptest.NewRequest("GET", "/", nil))
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	body, _ := io.ReadAll(resp.Body)
	if len(body) == 0 {
		t.Error("expected a generated guest_id in body")
	}
	if sc := resp.Header.Get("Set-Cookie"); !strings.Contains(sc, "guest_id=") {
		t.Errorf("expected Set-Cookie with guest_id, got %q", sc)
	}
}

func TestGuestSession_ExistingCookiePreserved(t *testing.T) {
	app := fiber.New()
	app.Use(GuestSession("development"))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString(c.Locals("guest_id").(string))
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Cookie", "guest_id=existing-123")
	resp, _ := app.Test(req)
	body, _ := io.ReadAll(resp.Body)
	if string(body) != "existing-123" {
		t.Errorf("guest_id = %q, want existing-123", string(body))
	}
}

func TestGuestSession_ProductionSecureCookie(t *testing.T) {
	app := fiber.New()
	app.Use(GuestSession("production"))
	app.Get("/", func(c *fiber.Ctx) error { return c.SendString("ok") })

	resp, _ := app.Test(httptest.NewRequest("GET", "/", nil))
	if sc := resp.Header.Get("Set-Cookie"); !strings.Contains(strings.ToLower(sc), "secure") {
		t.Errorf("production cookie should be Secure, got %q", sc)
	}
}
