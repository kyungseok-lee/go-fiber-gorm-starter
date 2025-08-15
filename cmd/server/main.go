package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/kyungseok-lee/fiber-gorm-starter/internal/config"
	"github.com/kyungseok-lee/fiber-gorm-starter/internal/db"
	"github.com/kyungseok-lee/fiber-gorm-starter/internal/domain/user"
	"github.com/kyungseok-lee/fiber-gorm-starter/internal/http"
	"github.com/kyungseok-lee/fiber-gorm-starter/internal/logger"
	"go.uber.org/zap"
)

// @title           Spindle API
// @version         1.0
// @description     A production-ready REST API built with Go Fiber and GORM
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// 헬스 체크 모드 / Health check mode
	if len(os.Args) > 1 && os.Args[1] == "--health-check" {
		healthCheck()
		return
	}

	// .env 파일 로드 / Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// 설정 로드 / Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 로거 초기화 / Initialize logger
	logger, err := logger.Init(cfg.LogLevel, cfg.Env)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	zap.L().Info("Starting Spindle API Server",
		zap.String("env", cfg.Env),
		zap.String("port", cfg.Port),
		zap.String("db_driver", cfg.DBDriver),
	)

	// 데이터베이스 연결 / Connect to database
	database, err := db.Connect(cfg)
	if err != nil {
		zap.L().Fatal("Failed to connect to database", zap.Error(err))
	}

	sqlDB, err := database.DB()
	if err != nil {
		zap.L().Fatal("Failed to get underlying sql.DB", zap.Error(err))
	}
	defer sqlDB.Close()

	// Auto-migrate 테이블 / Auto-migrate tables
	if err := database.AutoMigrate(&user.User{}); err != nil {
		zap.L().Fatal("Failed to auto-migrate database", zap.Error(err))
	}

	// HTTP 라우터 설정 / Setup HTTP router
	router := http.NewRouter(cfg, database)
	router.Setup()

	app := router.GetApp()

	// Graceful shutdown 설정 / Setup graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// 서버 시작 / Start server
	go func() {
		if err := app.Listen(":" + cfg.Port); err != nil {
			zap.L().Fatal("Failed to start server", zap.Error(err))
		}
	}()

	zap.L().Info("Server started successfully",
		zap.String("address", ":"+cfg.Port),
		zap.String("env", cfg.Env),
		zap.String("docs", "http://localhost:"+cfg.Port+"/docs/index.html"),
	)

	// 종료 신호 대기 / Wait for shutdown signal
	<-c

	zap.L().Info("Shutting down server...")

	// 30초 타임아웃으로 graceful shutdown / Graceful shutdown with 30 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		zap.L().Error("Server forced to shutdown", zap.Error(err))
	}

	zap.L().Info("Server exited")
}

// healthCheck 컨테이너 헬스 체크 / Container health check
func healthCheck() {
	// 단순한 HTTP 요청으로 헬스 체크 / Simple HTTP request for health check
	// 실제 구현에서는 HTTP 클라이언트로 /health 엔드포인트 호출 / In actual implementation, call /health endpoint with HTTP client
	if err := godotenv.Load(); err != nil {
		// 헬스 체크에서는 에러 무시 / Ignore error in health check
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Println("Health check failed: config load error")
		os.Exit(1)
	}

	database, err := db.Connect(cfg)
	if err != nil {
		fmt.Println("Health check failed: database connection error")
		os.Exit(1)
	}

	if err := db.HealthCheck(database); err != nil {
		fmt.Println("Health check failed: database health check error")
		os.Exit(1)
	}

	fmt.Println("Health check passed")
	os.Exit(0)
}