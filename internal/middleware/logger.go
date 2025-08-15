package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const (
	httpStatusInternalServerError = 500
	httpStatusBadRequest          = 400
)

// Logger 요청 로깅 미들웨어 / Request logging middleware
func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// 요청 처리 / Process request
		err := c.Next()

		// 로깅 / Logging
		duration := time.Since(start)
		requestID := GetRequestID(c)

		// 로그 필드 구성 / Configure log fields
		fields := []zap.Field{
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.String("ip", c.IP()),
			zap.String("user_agent", c.Get("User-Agent")),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("latency", duration),
			zap.Int("bytes_in", len(c.Body())),
			zap.Int("bytes_out", len(c.Response().Body())),
		}

		// 요청 ID가 있으면 추가 / Add request ID if available
		if requestID != "" {
			fields = append(fields, zap.String("request_id", requestID))
		}

		// 에러가 있으면 에러 로그, 없으면 정보 로그 / Error log if error exists, otherwise info log
		if err != nil {
			fields = append(fields, zap.Error(err))
			zap.L().Error("Request completed with error", fields...)
		} else {
			// HTTP 상태 코드에 따른 로그 레벨 조정 / Adjust log level based on HTTP status code
			status := c.Response().StatusCode()
			switch {
			case status >= httpStatusInternalServerError:
				zap.L().Error("Request completed", fields...)
			case status >= httpStatusBadRequest:
				zap.L().Warn("Request completed", fields...)
			default:
				zap.L().Info("Request completed", fields...)
			}
		}

		return err
	}
}

// RequestLogger 간소화된 요청 로깅 미들웨어 / Simplified request logging middleware
func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		rid, _ := c.Locals(RequestIDContextKey).(string)
		method := c.Method()
		path := c.Path()
		err := c.Next()
		status := c.Response().StatusCode()
		latency := time.Since(start)

		// Avoid PII in logs; only standard fields
		zap.L().Info("http",
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", status),
			zap.Duration("latency", latency),
			zap.String("rid", rid),
		)
		return err
	}
}
