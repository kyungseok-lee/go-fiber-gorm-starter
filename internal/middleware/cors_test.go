package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kyungseok-lee/go-fiber-gorm-starter/internal/config"
)

func TestCORSLocalAllowsAnyOrigin(t *testing.T) {
	header := testCORS(t, &config.Config{Env: "local"}, "https://app.example.com")

	assert.Equal(t, "*", header.Get(fiber.HeaderAccessControlAllowOrigin))
}

func TestCORSProdWithoutOriginsDeniesBrowserOrigins(t *testing.T) {
	header := testCORS(t, &config.Config{Env: "prod"}, "https://app.example.com")

	assert.Empty(t, header.Get(fiber.HeaderAccessControlAllowOrigin))
}

func TestCORSProdAllowsConfiguredOrigin(t *testing.T) {
	header := testCORS(t, &config.Config{
		Env:                  "prod",
		CORSAllowedOrigins:   "https://app.example.com",
		CORSAllowCredentials: true,
	}, "https://app.example.com")

	assert.Equal(t, "https://app.example.com", header.Get(fiber.HeaderAccessControlAllowOrigin))
	assert.Equal(t, "true", header.Get(fiber.HeaderAccessControlAllowCredentials))
}

func testCORS(t *testing.T, cfg *config.Config, origin string) http.Header {
	t.Helper()

	app := fiber.New()
	app.Use(CORS(cfg))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})

	req := httptest.NewRequest(fiber.MethodGet, "/", nil)
	req.Header.Set("Origin", origin)
	resp, err := app.Test(req)

	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())
	return resp.Header
}
