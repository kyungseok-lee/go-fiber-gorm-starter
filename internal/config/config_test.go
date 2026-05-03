package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadParsesDurationValues(t *testing.T) {
	t.Setenv("DB_MAX_LIFETIME", "300s")

	cfg, err := Load()

	require.NoError(t, err)
	assert.Equal(t, 300*time.Second, cfg.DBMaxLifetime)
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
