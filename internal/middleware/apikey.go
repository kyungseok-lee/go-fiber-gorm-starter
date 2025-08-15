// Package middleware provides HTTP middleware functions for the Fiber application
package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/kyungseok-lee/go-fiber-gorm-starter/internal/config"
	"github.com/kyungseok-lee/go-fiber-gorm-starter/pkg/resp"
)

// APIKey API 키 인증 미들웨어 / API key authentication middleware
func APIKey(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// API 키가 설정되지 않은 경우 스킵 / Skip if API key is not configured
		if cfg.APIKey == "" {
			return c.Next()
		}

		// Authorization 헤더에서 API 키 추출 / Extract API key from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return resp.Unauthorized(c, "Missing authorization header")
		}

		// Bearer 토큰 형식 확인 / Check Bearer token format
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return resp.Unauthorized(c, "Invalid authorization header format")
		}

		// API 키 추출 및 검증 / Extract and validate API key
		apiKey := strings.TrimPrefix(authHeader, "Bearer ")
		if apiKey != cfg.APIKey {
			return resp.Unauthorized(c, "Invalid API key")
		}

		return c.Next()
	}
}

// OptionalAPIKey 선택적 API 키 인증 미들웨어 / Optional API key authentication middleware
// API 키가 제공되면 검증하지만 필수는 아님 / Validates API key if provided but not required
func OptionalAPIKey(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// API 키가 설정되지 않은 경우 스킵 / Skip if API key is not configured
		if cfg.APIKey == "" {
			return c.Next()
		}

		// Authorization 헤더가 없으면 스킵 / Skip if no Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}

		// Bearer 토큰 형식 확인 / Check Bearer token format
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return resp.Unauthorized(c, "Invalid authorization header format")
		}

		// API 키 추출 및 검증 / Extract and validate API key
		apiKey := strings.TrimPrefix(authHeader, "Bearer ")
		if apiKey != cfg.APIKey {
			return resp.Unauthorized(c, "Invalid API key")
		}

		// 인증된 사용자 표시 / Mark as authenticated user
		c.Locals("authenticated", true)
		return c.Next()
	}
}
