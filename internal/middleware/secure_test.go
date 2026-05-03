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

func TestSecureHeaders(t *testing.T) {
	header := testSecureHeaders(t, &config.Config{Env: "local"})

	assert.Equal(t, "nosniff", header.Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", header.Get("X-Frame-Options"))
	assert.Equal(t, "1; mode=block", header.Get("X-XSS-Protection"))
	assert.Equal(t, "no-referrer", header.Get("Referrer-Policy"))
	assert.Equal(t, "none", header.Get("X-Permitted-Cross-Domain-Policies"))
	assert.Equal(t, "same-site", header.Get("Cross-Origin-Resource-Policy"))
	assert.Empty(t, header.Get("Strict-Transport-Security"))
}

func TestSecureHeadersAddsHSTSInProduction(t *testing.T) {
	header := testSecureHeaders(t, &config.Config{Env: "prod"})

	assert.Equal(t, "max-age=31536000; includeSubDomains", header.Get("Strict-Transport-Security"))
}

func testSecureHeaders(t *testing.T, cfg *config.Config) http.Header {
	t.Helper()

	app := fiber.New()
	app.Use(SecureHeaders(cfg))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/", nil))

	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())
	return resp.Header
}
