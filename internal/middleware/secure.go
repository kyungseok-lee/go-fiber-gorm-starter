package middleware

// Basic security headers

import (
	fiber "github.com/gofiber/fiber/v2"

	"github.com/kyungseok-lee/go-fiber-gorm-starter/internal/config"
)

// SecureHeaders 보안 헤더 미들웨어 / Security headers middleware
func SecureHeaders(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("Referrer-Policy", "no-referrer")
		c.Set("X-Permitted-Cross-Domain-Policies", "none")
		c.Set("Cross-Origin-Resource-Policy", "same-site")
		if cfg != nil && cfg.IsProd() {
			c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		return c.Next()
	}
}
