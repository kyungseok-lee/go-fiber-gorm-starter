// Package logger provides structured logging functionality using Zap
package logger

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Init 로거 초기화 / Initialize logger
func Init(level string, env string) (*zap.Logger, error) {
	var config zap.Config

	if env == "prod" {
		config = zap.NewProductionConfig()
		// 프로덕션에서는 JSON 포맷 사용 / Use JSON format in production
		config.Encoding = "json"
	} else {
		config = zap.NewDevelopmentConfig()
		// 개발환경에서는 콘솔 포맷 사용 / Use console format in development
		config.Encoding = "console"
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// 로그 레벨 설정 / Set log level
	switch strings.ToLower(level) {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	default:
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}

	// 타임스탬프 형식 설정 / Set timestamp format
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// stdout으로 출력 / Output to stdout
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	// 글로벌 로거 교체 / Replace global logger
	zap.ReplaceGlobals(logger)

	return logger, nil
}

// WithRequestID 요청 ID를 포함한 로거 생성 / Create logger with request ID
func WithRequestID(requestID string) *zap.Logger {
	return zap.L().With(zap.String("request_id", requestID))
}

// LogLevel 환경변수에서 로그 레벨 가져오기 / Get log level from environment
func LogLevel() string {
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = "info"
	}
	return level
}
