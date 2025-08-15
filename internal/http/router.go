// Package http provides HTTP router configuration and setup for the application
package http

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/kyungseok-lee/go-fiber-gorm-starter/internal/config"
	"github.com/kyungseok-lee/go-fiber-gorm-starter/internal/domain/user"
	"github.com/kyungseok-lee/go-fiber-gorm-starter/internal/http/health"
	"github.com/kyungseok-lee/go-fiber-gorm-starter/internal/metrics"
	"github.com/kyungseok-lee/go-fiber-gorm-starter/internal/middleware"
	"github.com/kyungseok-lee/go-fiber-gorm-starter/pkg/resp"
)

const (
	readTimeoutSeconds  = 10
	writeTimeoutSeconds = 10
	idleTimeoutSeconds  = 120
)

// Router HTTP 라우터 설정 / HTTP router configuration
type Router struct {
	app   *fiber.App
	cfg   *config.Config
	db    *gorm.DB
	userH *user.Handler
}

// NewRouter 새 라우터 생성 / Create new router
func NewRouter(cfg *config.Config, db *gorm.DB) *Router {
	// Fiber 앱 설정 / Fiber app configuration
	app := fiber.New(fiber.Config{
		AppName:      "spindle API", // 브랜딩 이름 사용 / Use branding name
		ReadTimeout:  readTimeoutSeconds * time.Second,
		WriteTimeout: writeTimeoutSeconds * time.Second,
		IdleTimeout:  idleTimeoutSeconds * time.Second,
		ServerHeader: "spindle",
		// JSON 엔코더 최적화 옵션 (필요시 주석 해제) / JSON encoder optimization option (uncomment if needed)
		// JSONEncoder: json.Marshal,   // 기본 encoding/json 사용 / Use default encoding/json
		// JSONDecoder: json.Unmarshal, // goccy/go-json으로 교체 가능 / Can be replaced with goccy/go-json
	})

	// User 도메인 초기화 / Initialize User domain
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	return &Router{
		app:   app,
		cfg:   cfg,
		db:    db,
		userH: userHandler,
	}
}

// Setup 라우터 설정 / Setup router
func (r *Router) Setup() {
	// 패닉 복구 미들웨어 / Panic recovery middleware
	r.app.Use(middleware.Recover())

	// 보안 헤더 미들웨어 / Security headers middleware
	r.app.Use(middleware.SecureHeaders())

	// 요청 ID 미들웨어 / Request ID middleware
	r.app.Use(middleware.RequestID())

	// 로깅 미들웨어 / Logging middleware
	r.app.Use(middleware.RequestLogger())

	// CORS 미들웨어 / CORS middleware
	r.app.Use(middleware.CORS(r.cfg))

	// API 키 미들웨어 (설정된 경우) / API key middleware (if configured)
	if r.cfg.APIKey != "" {
		r.app.Use(middleware.APIKey(r.cfg))
	}

	// 메트릭 미들웨어 (활성화된 경우) / Metrics middleware (if enabled)
	if r.cfg.MetricsEnabled {
		prometheus := metrics.NewPrometheus()
		r.app.Use(prometheus.Middleware())
		prometheus.RegisterAt(r.app, "/metrics")
	}

	// Health 체크 라우트 / Health check routes
	r.setupHealthRoutes()

	// Swagger 문서 (개발환경에서만) / Swagger documentation (development only)
	if r.cfg.IsDev() {
		r.app.Get("/docs/*", middleware.Swagger())
	}

	// API v1 라우트 / API v1 routes
	r.setupV1Routes()

	// 프로파일링 라우트 (활성화된 경우) / Profiling routes (if enabled)
	if r.cfg.PProfEnabled {
		r.setupPProfRoutes()
	}

	// 404 핸들러 / 404 handler
	r.setup404Handler()
}

// setupHealthRoutes 헬스 체크 라우트 설정 / Setup health check routes
func (r *Router) setupHealthRoutes() {
	healthHandler := health.New(r.db)
	r.app.Get("/health", healthHandler.Health)
	r.app.Get("/ready", healthHandler.Ready)
}

// setupV1Routes API v1 라우트 설정 / Setup API v1 routes
func (r *Router) setupV1Routes() {
	v1 := r.app.Group("/v1")

	// User 라우트 / User routes
	users := v1.Group("/users")
	users.Get("/", r.userH.List)         // GET /v1/users
	users.Get("/:id", r.userH.GetByID)   // GET /v1/users/:id
	users.Post("/", r.userH.Create)      // POST /v1/users
	users.Put("/:id", r.userH.Update)    // PUT /v1/users/:id
	users.Delete("/:id", r.userH.Delete) // DELETE /v1/users/:id

	// 향후 확장 가능한 라우트들 / Future extensible routes
	// auth := v1.Group("/auth")
	// auth.Post("/login", authHandler.Login)
	// auth.Post("/logout", authHandler.Logout)
	// auth.Post("/refresh", authHandler.Refresh)

	// protected := v1.Group("/protected")
	// protected.Use(middleware.APIKey(r.cfg)) // API 키 인증 필요 / Requires API key authentication
	// protected.Get("/admin", adminHandler.Dashboard)
}

// setupPProfRoutes 프로파일링 라우트 설정 / Setup profiling routes
func (r *Router) setupPProfRoutes() {
	// pprof 라우트는 보안상 개발환경에서만 활성화하는 것을 권장 / Recommend enabling pprof routes only in development for security
	if r.cfg.IsDev() {
		// TODO: net/http/pprof 패키지 통합 / Integrate net/http/pprof package
		// r.app.Get("/debug/pprof/*", adaptor.HTTPHandler(http.DefaultServeMux))
	}
}

// setup404Handler 404 에러 핸들러 설정 / Setup 404 error handler
func (r *Router) setup404Handler() {
	r.app.Use(func(c *fiber.Ctx) error {
		return resp.NotFound(c, "route not found")
	})
}

// GetApp Fiber 앱 반환 / Return Fiber app
func (r *Router) GetApp() *fiber.App {
	return r.app
}
