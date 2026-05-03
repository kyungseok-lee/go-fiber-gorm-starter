package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/kyungseok-lee/go-fiber-gorm-starter/internal/config"
)

// CORS CORS 미들웨어 설정 / CORS middleware configuration
func CORS(cfg *config.Config) fiber.Handler {
	corsConfig := cors.Config{
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-Request-ID",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
	}

	// 프로덕션 환경에서는 특정 도메인만 허용 / Allow only specific domains in production
	if cfg.IsProd() {
		if cfg.CORSAllowedOrigins == "" {
			corsConfig.AllowOriginsFunc = func(_ string) bool {
				return false
			}
		} else {
			corsConfig.AllowOrigins = cfg.CORSAllowedOrigins
			corsConfig.AllowCredentials = cfg.CORSAllowCredentials
		}
	} else {
		// 개발환경에서는 모든 오리진 허용 / Allow all origins in development
		corsConfig.AllowOrigins = "*"
	}

	return cors.New(corsConfig)
}
