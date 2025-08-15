// Package db provides database connection and management functionality
package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/kyungseok-lee/go-fiber-gorm-starter/internal/config"
)

// Connect 데이터베이스 연결 / Connect to database
func Connect(cfg *config.Config) (*gorm.DB, error) {
	dialector, err := createDialector(cfg)
	if err != nil {
		return nil, err
	}

	gormConfig := createGormConfig(cfg)

	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := configureConnectionPool(db, cfg); err != nil {
		return nil, err
	}

	if err := verifyConnection(db); err != nil {
		return nil, err
	}

	logConnectionSuccess(cfg)
	return db, nil
}

// createDialector 드라이버에 따른 dialector 생성 / Create dialector based on driver
func createDialector(cfg *config.Config) (gorm.Dialector, error) {
	switch cfg.DBDriver {
	case "postgres":
		return postgres.Open(cfg.GetDBDSN()), nil
	case "mysql":
		return mysql.Open(cfg.GetDBDSN()), nil
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.DBDriver)
	}
}

// createGormConfig GORM 설정 생성 / Create GORM configuration
func createGormConfig(cfg *config.Config) *gorm.Config {
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
		NowFunc: func() time.Time {
			loc, _ := time.LoadLocation("Asia/Seoul")
			return time.Now().In(loc)
		},
	}

	if cfg.IsDev() {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	return gormConfig
}

// configureConnectionPool 커넥션 풀 설정 / Configure connection pool
func configureConnectionPool(db *gorm.DB, cfg *config.Config) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.DBMaxOpen)
	sqlDB.SetMaxIdleConns(cfg.DBMaxIdle)
	sqlDB.SetConnMaxLifetime(cfg.DBMaxLifetime)

	return nil
}

// verifyConnection 연결 확인 / Verify connection
func verifyConnection(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

// logConnectionSuccess 연결 성공 로그 / Log connection success
func logConnectionSuccess(cfg *config.Config) {
	zap.L().Info("Database connected successfully",
		zap.String("driver", cfg.DBDriver),
		zap.String("host", cfg.DBHost),
		zap.String("port", cfg.DBPort),
		zap.String("database", cfg.DBName),
	)
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
