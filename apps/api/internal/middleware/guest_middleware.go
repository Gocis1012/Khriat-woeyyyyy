package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	GuestCookieName = "guest_id"
	GuestCookieTTL  = 30 * 24 * time.Hour // 30 วัน
)


func GuestSession() fiber.Handler {
	return func(c *fiber.Ctx) error {

		guestID := c.Cookies(GuestCookieName)

		if guestID == "" {

			guestID = uuid.New().String()

			c.Cookie(&fiber.Cookie{
				Name: GuestCookieName,
				Value: guestID,
				Expires: time.Now().Add(GuestCookieTTL),
				 HTTPOnly: true,  // JS อ่านไม่ได้ → ป้องกัน XSS
                Secure:   false, // true ใน production (HTTPS only)
                SameSite: "Lax", // ป้องกัน CSRF พื้นฐาน
			})
		}

		c.Locals("guest_id", guestID)

		return c.Next()
	}
}