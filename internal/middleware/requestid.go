package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RequestIDHeaderKey 요청 ID 헤더 키 / Request ID header key
const RequestIDHeaderKey = "X-Request-ID"

// RequestIDContextKey 요청 ID 컨텍스트 키 / Request ID context key
const RequestIDContextKey = "request_id"

// RequestID 요청 ID 미들웨어 / Request ID middleware
func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 기존 요청 ID가 있으면 사용, 없으면 새로 생성 / Use existing request ID if available, otherwise generate new one
		requestID := c.Get(RequestIDHeaderKey)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// 응답 헤더에 요청 ID 설정 / Set request ID in response header
		c.Set(RequestIDHeaderKey, requestID)

		// 컨텍스트에 요청 ID 저장 / Store request ID in context
		c.Locals(RequestIDContextKey, requestID)

		return c.Next()
	}
}

// GetRequestID 컨텍스트에서 요청 ID 가져오기 / Get request ID from context
func GetRequestID(c *fiber.Ctx) string {
	if requestID, ok := c.Locals(RequestIDContextKey).(string); ok {
		return requestID
	}
	return ""
}