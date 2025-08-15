// Package health provides health check handlers for the HTTP server
package health

// Health and readiness handlers

import (
	fiber "github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/kyungseok-lee/go-fiber-gorm-starter/internal/db"
	"github.com/kyungseok-lee/go-fiber-gorm-starter/pkg/resp"
)

// Handler 헬스 체크 핸들러 / Health check handler
type Handler struct{ db *gorm.DB }

// New 새로운 헬스 체크 핸들러 생성 / Create new health check handler
func New(db *gorm.DB) *Handler { return &Handler{db: db} }

// Response 헬스 체크 응답 구조체 / Health check response structure
type Response struct {
	Status  string            `json:"status"`
	Service string            `json:"service"`
	Version string            `json:"version"`
	Checks  map[string]string `json:"checks,omitempty"`
}

// Health returns static 200 OK.
// @Summary Health check
// @Description Get service health status
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} Response
// @Router /health [get]
func (h *Handler) Health(c *fiber.Ctx) error {
	return resp.Success(c, Response{
		Status:  "ok",
		Service: "fiber-gorm-starter",
		Version: "1.0.0",
	})
}

// Ready checks DB ping.
// @Summary Readiness check
// @Description Get service readiness status including dependencies
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} Response
// @Failure 503 {object} resp.ErrorResponse
// @Router /ready [get]
func (h *Handler) Ready(c *fiber.Ctx) error {
	checks := make(map[string]string)

	// 데이터베이스 연결 상태 확인 / Check database connection status
	if err := db.HealthCheck(h.db); err != nil {
		checks["database"] = "fail"
		return resp.Error(c, fiber.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "Service not ready", checks)
	}
	checks["database"] = "ok"

	return resp.Success(c, Response{
		Status:  "ready",
		Service: "fiber-gorm-starter",
		Version: "1.0.0",
		Checks:  checks,
	})
}
