package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kyungseok-lee/fiber-gorm-starter/internal/db"
	"github.com/kyungseok-lee/fiber-gorm-starter/pkg/resp"
	"gorm.io/gorm"
)

// HealthHandler 헬스 체크 핸들러 / Health check handler
type HealthHandler struct {
	db *gorm.DB
}

// NewHealthHandler 새 헬스 체크 핸들러 생성 / Create new health check handler
func NewHealthHandler(database *gorm.DB) *HealthHandler {
	return &HealthHandler{
		db: database,
	}
}

// HealthResponse 헬스 체크 응답 구조체 / Health check response structure
type HealthResponse struct {
	Status  string            `json:"status"`
	Service string            `json:"service"`
	Version string            `json:"version"`
	Checks  map[string]string `json:"checks,omitempty"`
}

// Health 기본 헬스 체크 / Basic health check
// @Summary Health check
// @Description Get service health status
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func (h *HealthHandler) Health(c *fiber.Ctx) error {
	return resp.Success(c, HealthResponse{
		Status:  "ok",
		Service: "spindle",
		Version: "1.0.0",
	})
}

// Ready 준비 상태 체크 (의존성 포함) / Readiness check (including dependencies)
// @Summary Readiness check
// @Description Get service readiness status including dependencies
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse
// @Failure 503 {object} resp.ErrorResponse
// @Router /ready [get]
func (h *HealthHandler) Ready(c *fiber.Ctx) error {
	checks := make(map[string]string)

	// 데이터베이스 연결 상태 확인 / Check database connection status
	if err := db.HealthCheck(h.db); err != nil {
		checks["database"] = "fail"
		return resp.Error(c, fiber.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "Service not ready", checks)
	}
	checks["database"] = "ok"

	// 향후 추가 의존성 체크 / Future additional dependency checks
	// - Redis 연결 체크
	// - 외부 API 연결 체크
	// - 파일 시스템 체크
	// if err := checkRedis(); err != nil {
	//     checks["redis"] = "fail"
	//     return resp.Error(c, fiber.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "Service not ready", checks)
	// }
	// checks["redis"] = "ok"

	return resp.Success(c, HealthResponse{
		Status:  "ready",
		Service: "spindle",
		Version: "1.0.0",
		Checks:  checks,
	})
}