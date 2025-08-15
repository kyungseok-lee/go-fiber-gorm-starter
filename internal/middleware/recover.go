package middleware

// Panic recover middleware using Fiber's recover

import (
	fiber "github.com/gofiber/fiber/v2"
	recovermw "github.com/gofiber/fiber/v2/middleware/recover"
)

func Recover() fiber.Handler {
	return recovermw.New()
}