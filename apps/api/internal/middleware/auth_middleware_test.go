package middleware

import (
	"net/http/httptest"
	"testing"

	"corporate-translator-api/internal/service"

	"github.com/gofiber/fiber/v2"
)

func newAuth() (*service.AuthService, string) {
	svc := service.NewAuthService("secret", "client-id")
	token, _ := svc.GenerateJWT("user-1", "a@b.com")
	return svc, token
}

func TestRequireAuth_ValidToken(t *testing.T) {
	svc, token := newAuth()
	app := fiber.New()
	app.Use(RequireAuth(svc))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString(c.Locals("user_id").(string))
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ := app.Test(req)
	if resp.StatusCode != 200 {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
}

func TestRequireAuth_Rejections(t *testing.T) {
	svc, _ := newAuth()
	app := fiber.New()
	app.Use(RequireAuth(svc))
	app.Get("/", func(c *fiber.Ctx) error { return c.SendString("ok") })

	cases := []struct {
		name   string
		header string
	}{
		{"missing header", ""},
		{"wrong scheme", "Token abc"},
		{"invalid token", "Bearer garbage"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			if tc.header != "" {
				req.Header.Set("Authorization", tc.header)
			}
			resp, _ := app.Test(req)
			if resp.StatusCode != 401 {
				t.Errorf("status = %d, want 401", resp.StatusCode)
			}
		})
	}
}

func TestOptionalAuth_ValidToken(t *testing.T) {
	svc, token := newAuth()
	app := fiber.New()
	app.Use(OptionalAuth(svc))
	app.Get("/", func(c *fiber.Ctx) error {
		if id, ok := c.Locals("user_id").(string); ok {
			return c.SendString(id)
		}
		return c.SendString("guest")
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ := app.Test(req)
	if resp.StatusCode != 200 {
		t.Fatalf("status = %d", resp.StatusCode)
	}
}

func TestOptionalAuth_PassThroughAsGuest(t *testing.T) {
	svc, _ := newAuth()
	app := fiber.New()
	app.Use(OptionalAuth(svc))
	app.Get("/", func(c *fiber.Ctx) error {
		if _, ok := c.Locals("user_id").(string); ok {
			return c.SendString("user")
		}
		return c.SendString("guest")
	})

	// No header → guest
	resp, _ := app.Test(httptest.NewRequest("GET", "/", nil))
	if resp.StatusCode != 200 {
		t.Errorf("no-header status = %d", resp.StatusCode)
	}

	// Wrong scheme → guest (passes through)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Token abc")
	resp2, _ := app.Test(req)
	if resp2.StatusCode != 200 {
		t.Errorf("wrong-scheme status = %d", resp2.StatusCode)
	}

	// Invalid token → guest (passes through)
	req3 := httptest.NewRequest("GET", "/", nil)
	req3.Header.Set("Authorization", "Bearer garbage")
	resp3, _ := app.Test(req3)
	if resp3.StatusCode != 200 {
		t.Errorf("invalid-token status = %d", resp3.StatusCode)
	}
}
