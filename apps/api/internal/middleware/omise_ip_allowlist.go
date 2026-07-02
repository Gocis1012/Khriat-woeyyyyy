package middleware

import (
	"net"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// OmiseIPAllowlist restricts a route to Omise's published webhook source
// IPs. Omise does not sign webhook payloads (no HMAC/shared-secret scheme
// like Stripe or GitHub), so an IP allowlist is the verification mechanism
// it actually supports.
//
// allowedCSV is a comma-separated list of IPs and/or CIDR ranges, e.g.
// "203.0.113.4,198.51.100.0/24". Populate it from Omise's current webhook
// IP documentation/dashboard — the list can change, so keep it up to date.
func OmiseIPAllowlist(allowedCSV string) fiber.Handler {
	var nets []*net.IPNet
	var ips []net.IP

	for _, raw := range strings.Split(allowedCSV, ",") {
		entry := strings.TrimSpace(raw)
		if entry == "" {
			continue
		}
		if strings.Contains(entry, "/") {
			if _, ipNet, err := net.ParseCIDR(entry); err == nil {
				nets = append(nets, ipNet)
			}
			continue
		}
		if ip := net.ParseIP(entry); ip != nil {
			ips = append(ips, ip)
		}
	}

	return func(c *fiber.Ctx) error {
		parsed := net.ParseIP(requestIP(c))

		allowed := false
		if parsed != nil {
			for _, ip := range ips {
				if ip.Equal(parsed) {
					allowed = true
					break
				}
			}
			if !allowed {
				for _, n := range nets {
					if n.Contains(parsed) {
						allowed = true
						break
					}
				}
			}
		}

		if !allowed {
			return c.SendStatus(fiber.StatusForbidden)
		}
		return c.Next()
	}
}

// requestIP returns the first hop in X-Forwarded-For (the original client),
// since the app runs behind a reverse proxy (e.g. Render) that would
// otherwise mask the real source IP. Falls back to the direct connection IP
// when the header is absent (e.g. local development).
func requestIP(c *fiber.Ctx) string {
	if xff := c.Get(fiber.HeaderXForwardedFor); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	return c.IP()
}
