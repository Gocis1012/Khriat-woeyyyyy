package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	GuestCookieName = "guest_id"
	GuestCookieTTL  = 30 * 24 * time.Hour
)

func GuestSession(appEnv string) fiber.Handler {
	secureCookie := appEnv == "production"
	// Cross-origin requests (Vercel → Render) require SameSite=None; Secure.
	// SameSite=Lax blocks cookies on cross-origin POST, causing 401s.
	sameSite := "Lax"
	if secureCookie {
		sameSite = "None"
	}

	return func(c *fiber.Ctx) error {
		guestID := c.Cookies(GuestCookieName)

		if guestID == "" {
			guestID = uuid.New().String()

			c.Cookie(&fiber.Cookie{
				Name:     GuestCookieName,
				Value:    guestID,
				Expires:  time.Now().Add(GuestCookieTTL),
				HTTPOnly: true,
				Secure:   secureCookie,
				SameSite: sameSite,
			})
		}

		c.Locals("guest_id", guestID)

		return c.Next()
	}
}
