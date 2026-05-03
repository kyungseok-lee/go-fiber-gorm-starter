// Package config provides configuration management for the application
package config

import (
	"fmt"
	"net/url"
	"time"

	"github.com/caarlos0/env/v11"
)

// Config 환경설정 구조체 / Application configuration structure
type Config struct {
	// Environment settings
	Env  string `env:"ENV" envDefault:"local"`
	Port string `env:"PORT" envDefault:"8080"`

	// Database settings
	DBDriver      string        `env:"DB_DRIVER" envDefault:"mysql"`
	DBHost        string        `env:"DB_HOST" envDefault:"localhost"`
	DBPort        string        `env:"DB_PORT" envDefault:"3306"`
	DBUser        string        `env:"DB_USER" envDefault:"user"`
	DBPass        string        `env:"DB_PASS" envDefault:"password"`
	DBName        string        `env:"DB_NAME" envDefault:"fiber_gorm_starter"`
	DBSSLMode     string        `env:"DB_SSL_MODE" envDefault:"disable"`
	DBMaxOpen     int           `env:"DB_MAX_OPEN" envDefault:"25"`
	DBMaxIdle     int           `env:"DB_MAX_IDLE" envDefault:"10"`
	DBMaxLifetime time.Duration `env:"DB_MAX_LIFETIME" envDefault:"300s"`

	// Security settings
	APIKey string `env:"API_KEY" envDefault:""`

	// Logging settings
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`

	// Metrics settings
	MetricsEnabled bool `env:"METRICS_ENABLED" envDefault:"true"`

	// Profiling settings
	PProfEnabled bool `env:"PPROF_ENABLED" envDefault:"false"`
}

// Load 환경변수에서 설정을 로드 / Load configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// IsDev 개발 환경인지 확인 / Check if running in development environment
func (c *Config) IsDev() bool {
	return c.Env == "dev" || c.Env == "local"
}

// IsProd 프로덕션 환경인지 확인 / Check if running in production environment
func (c *Config) IsProd() bool {
	return c.Env == "prod"
}

// GetDBDSN 데이터베이스 DSN 생성 / Generate database DSN
func (c *Config) GetDBDSN() string {
	switch c.DBDriver {
	case "postgres":
		return fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Seoul",
			c.DBHost,
			c.DBPort,
			c.DBUser,
			c.DBPass,
			c.DBName,
			c.DBSSLMode,
		)
	case "mysql":
		fallthrough
	default:
		query := url.Values{}
		query.Set("charset", "utf8mb4")
		query.Set("parseTime", "True")
		query.Set("loc", "Asia/Seoul")

		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBName, query.Encode())
	}
}
