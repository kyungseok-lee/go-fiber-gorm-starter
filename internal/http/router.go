package http

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/kyungseok-lee/fiber-gorm-starter/internal/config"
	"github.com/kyungseok-lee/fiber-gorm-starter/internal/domain/user"
	"github.com/kyungseok-lee/fiber-gorm-starter/internal/middleware"
	"github.com/kyungseok-lee/fiber-gorm-starter/internal/metrics"
	"gorm.io/gorm"
)

// Router HTTP 라우터 설정 / HTTP router configuration
type Router struct {
	app    *fiber.App
	cfg    *config.Config
	db     *gorm.DB
	userH  *user.Handler
}

// NewRouter 새 라우터 생성 / Create new router
func NewRouter(cfg *config.Config, db *gorm.DB) *Router {
	// Fiber 앱 설정 / Fiber app configuration
	app := fiber.New(fiber.Config{
		AppName:      "spindle API", // 브랜딩 이름 사용 / Use branding name
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
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
	// 보안 헤더 미들웨어 / Security headers middleware
	r.app.Use(helmet.New(helmet.Config{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "DENY",
		HSTSMaxAge:            31536000,
		HSTSExcludeSubdomains: false,
		HSTSPreloadEnabled:    true,
	}))

	// 패닉 복구 미들웨어 / Panic recovery middleware
	r.app.Use(recover.New())

	// 요청 ID 미들웨어 / Request ID middleware
	r.app.Use(middleware.RequestID())

	// 로깅 미들웨어 / Logging middleware
	r.app.Use(middleware.Logger())

	// CORS 미들웨어 / CORS middleware
	r.app.Use(middleware.CORS(r.cfg))

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
		r.app.Get("/docs/*", swagger.HandlerDefault)
	}

	// API v1 라우트 / API v1 routes
	r.setupV1Routes()

	// 프로파일링 라우트 (활성화된 경우) / Profiling routes (if enabled)
	if r.cfg.PProfEnabled {
		r.setupPProfRoutes()
	}
}

// setupHealthRoutes 헬스 체크 라우트 설정 / Setup health check routes
func (r *Router) setupHealthRoutes() {
	r.app.Get("/health", NewHealthHandler(r.db).Health)
	r.app.Get("/ready", NewHealthHandler(r.db).Ready)
}

// setupV1Routes API v1 라우트 설정 / Setup API v1 routes
func (r *Router) setupV1Routes() {
	v1 := r.app.Group("/v1")

	// User 라우트 / User routes
	users := v1.Group("/users")
	users.Get("/", r.userH.List)                    // GET /v1/users
	users.Get("/:id", r.userH.GetByID)              // GET /v1/users/:id
	users.Post("/", r.userH.Create)                 // POST /v1/users
	users.Put("/:id", r.userH.Update)               // PUT /v1/users/:id
	users.Delete("/:id", r.userH.Delete)            // DELETE /v1/users/:id

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

// GetApp Fiber 앱 반환 / Return Fiber app
func (r *Router) GetApp() *fiber.App {
	return r.app
}