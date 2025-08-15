package middleware

// Panic recover middleware using Fiber's recover

import (
	fiber "github.com/gofiber/fiber/v2"
	recovermw "github.com/gofiber/fiber/v2/middleware/recover"
)

// Recover 패닉 복구 미들웨어 / Panic recovery middleware
func Recover() fiber.Handler {
	return recovermw.New()
}
