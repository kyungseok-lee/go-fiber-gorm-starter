package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/kyungseok-lee/fiber-gorm-starter/internal/config"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect 데이터베이스 연결 / Connect to database
func Connect(cfg *config.Config) (*gorm.DB, error) {
	var dialector gorm.Dialector

	// 드라이버에 따른 dialector 선택 / Select dialector based on driver
	switch cfg.DBDriver {
	case "postgres":
		dialector = postgres.Open(cfg.GetDBDSN())
	case "mysql":
		dialector = mysql.Open(cfg.GetDBDSN())
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.DBDriver)
	}

	// GORM 설정 / GORM configuration
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn), // Warning 레벨 이상만 로깅 / Log only warning level and above
		NowFunc: func() time.Time {
			// 한국 시간대 설정 / Set Korean timezone
			loc, _ := time.LoadLocation("Asia/Seoul")
			return time.Now().In(loc)
		},
	}

	// 개발환경에서는 SQL 쿼리 로깅 활성화 / Enable SQL query logging in development
	if cfg.IsDev() {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	// 데이터베이스 연결 / Connect to database
	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 커넥션 풀 설정 / Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// 최대 열린 연결 수 / Maximum number of open connections
	sqlDB.SetMaxOpenConns(cfg.DBMaxOpen)
	// 최대 유휴 연결 수 / Maximum number of idle connections
	sqlDB.SetMaxIdleConns(cfg.DBMaxIdle)
	// 연결 최대 생존 시간 / Maximum lifetime of connections
	sqlDB.SetConnMaxLifetime(cfg.DBMaxLifetime)

	// 연결 확인 / Verify connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	zap.L().Info("Database connected successfully",
		zap.String("driver", cfg.DBDriver),
		zap.String("host", cfg.DBHost),
		zap.String("port", cfg.DBPort),
		zap.String("database", cfg.DBName),
	)

	// TODO: 향후 추가 확장 가능한 훅들
	// - Query 성능 모니터링 훅
	// - 슬로우 쿼리 로깅 훅
	// - 캐시 무효화 훅
	// - 감사 로그 훅

	return db, nil
}

// HealthCheck 데이터베이스 상태 확인 / Check database health
func HealthCheck(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// 1초 타임아웃으로 ping 실행 / Execute ping with 1 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	return sqlDB.PingContext(ctx)
}

// GetConnectionStats 연결 통계 정보 반환 / Return connection statistics
func GetConnectionStats(db *gorm.DB) (*sql.DBStats, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	stats := sqlDB.Stats()
	return &stats, nil
}