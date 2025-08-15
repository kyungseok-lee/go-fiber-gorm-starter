package middleware

// Swagger route handler using fiber-swagger

import (
	fiber "github.com/gofiber/fiber/v2"
	fiberswagger "github.com/swaggo/fiber-swagger"
)

// Swagger API 문서 미들웨어 / Swagger API documentation middleware
func Swagger() fiber.Handler {
	return fiberswagger.WrapHandler
}
