package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadParsesDurationValues(t *testing.T) {
	t.Setenv("ENV", "local")
	t.Setenv("DB_MAX_LIFETIME", "300s")

	cfg, err := Load()

	require.NoError(t, err)
	assert.Equal(t, 300*time.Second, cfg.DBMaxLifetime)
}

func TestLoadRejectsUnsafeProductionDefaults(t *testing.T) {
	testCases := []struct {
		name          string
		apiKey        string
		dbPass        string
		corsOrigins   string
		errorContains string
	}{
		{
			name:          "missing api key",
			dbPass:        "strong-db-pass",
			errorContains: "API_KEY",
		},
		{
			name:          "placeholder api key",
			apiKey:        "your-api-key-here",
			dbPass:        "strong-db-pass",
			errorContains: "API_KEY",
		},
		{
			name:          "default database password",
			apiKey:        "strong-api-key",
			dbPass:        "password",
			errorContains: "DB_PASS",
		},
		{
			name:          "wildcard cors origin",
			apiKey:        "strong-api-key",
			dbPass:        "strong-db-pass",
			corsOrigins:   "*",
			errorContains: "CORS_ALLOWED_ORIGINS",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("ENV", "prod")
			t.Setenv("API_KEY", tc.apiKey)
			t.Setenv("DB_PASS", tc.dbPass)
			t.Setenv("CORS_ALLOWED_ORIGINS", tc.corsOrigins)

			cfg, err := Load()

			assert.Nil(t, cfg)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.errorContains)
		})
	}
}

func TestLoadAcceptsSecureProductionConfig(t *testing.T) {
	t.Setenv("ENV", "prod")
	t.Setenv("API_KEY", "strong-api-key")
	t.Setenv("DB_PASS", "strong-db-pass")
	t.Setenv("CORS_ALLOWED_ORIGINS", "https://app.example.com")
	t.Setenv("CORS_ALLOW_CREDENTIALS", "true")

	cfg, err := Load()

	require.NoError(t, err)
	assert.Equal(t, "strong-api-key", cfg.APIKey)
	assert.Equal(t, "strong-db-pass", cfg.DBPass)
	assert.Equal(t, "https://app.example.com", cfg.CORSAllowedOrigins)
	assert.True(t, cfg.CORSAllowCredentials)
}

func TestGetDBDSNMySQL(t *testing.T) {
	cfg := &Config{
		DBDriver: "mysql",
		DBHost:   "localhost",
		DBPort:   "3306",
		DBUser:   "user",
		DBPass:   "password",
		DBName:   "fiber_gorm_starter",
	}

	dsn := cfg.GetDBDSN()

	assert.Contains(t, dsn, "user:password@tcp(localhost:3306)/fiber_gorm_starter?")
	assert.Contains(t, dsn, "charset=utf8mb4")
	assert.Contains(t, dsn, "parseTime=True")
	assert.Contains(t, dsn, "loc=Asia%2FSeoul")
}

func TestGetDBDSNPostgres(t *testing.T) {
	cfg := &Config{
		DBDriver:  "postgres",
		DBHost:    "postgres",
		DBPort:    "5432",
		DBUser:    "user",
		DBPass:    "password",
		DBName:    "fiber_gorm_starter",
		DBSSLMode: "disable",
	}

	dsn := cfg.GetDBDSN()

	want := "host=postgres port=5432 user=user password=password dbname=fiber_gorm_starter " +
		"sslmode=disable TimeZone=Asia/Seoul"
	assert.Equal(t, want, dsn)
}
