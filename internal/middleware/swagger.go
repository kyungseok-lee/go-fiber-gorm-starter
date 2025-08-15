package middleware

// Swagger route handler using fiber-swagger

import (
	fiber "github.com/gofiber/fiber/v2"
	fiberswagger "github.com/swaggo/fiber-swagger"
)

func Swagger() fiber.Handler {
	return fiberswagger.WrapHandler
}