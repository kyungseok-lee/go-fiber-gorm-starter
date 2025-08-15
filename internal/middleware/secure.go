package middleware

// Basic security headers

import (
	fiber "github.com/gofiber/fiber/v2"
)

func SecureHeaders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-XSS-Protection", "1; mode=block")
		return c.Next()
	}
}