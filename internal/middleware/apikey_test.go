package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kyungseok-lee/go-fiber-gorm-starter/internal/config"
)

func TestAPIKeyAllowsValidBearerToken(t *testing.T) {
	app := fiber.New()
	app.Use(APIKey(&config.Config{APIKey: "expected-secret"}))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})

	req := httptest.NewRequest(fiber.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer expected-secret")
	resp, err := app.Test(req)

	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())
	assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
}

func TestAPIKeyRejectsMissingAndInvalidTokens(t *testing.T) {
	testCases := []struct {
		name          string
		authorization string
	}{
		{name: "missing header"},
		{name: "wrong scheme", authorization: "Basic expected-secret"},
		{name: "empty bearer", authorization: "Bearer "},
		{name: "wrong token", authorization: "Bearer wrong-secret"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			app := fiber.New()
			app.Use(APIKey(&config.Config{APIKey: "expected-secret"}))
			app.Get("/", func(c *fiber.Ctx) error {
				return c.SendStatus(fiber.StatusNoContent)
			})

			req := httptest.NewRequest(fiber.MethodGet, "/", nil)
			if tc.authorization != "" {
				req.Header.Set("Authorization", tc.authorization)
			}
			resp, err := app.Test(req)

			require.NoError(t, err)
			require.NoError(t, resp.Body.Close())
			assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
		})
	}
}

func TestOptionalAPIKeyAllowsAnonymousAndRejectsInvalidToken(t *testing.T) {
	app := fiber.New()
	app.Use(OptionalAPIKey(&config.Config{APIKey: "expected-secret"}))
	app.Get("/", func(c *fiber.Ctx) error {
		if c.Locals("authenticated") == true {
			return c.SendStatus(fiber.StatusCreated)
		}
		return c.SendStatus(fiber.StatusNoContent)
	})

	anonymousResp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/", nil))
	require.NoError(t, err)
	require.NoError(t, anonymousResp.Body.Close())
	assert.Equal(t, fiber.StatusNoContent, anonymousResp.StatusCode)

	validReq := httptest.NewRequest(fiber.MethodGet, "/", nil)
	validReq.Header.Set("Authorization", "Bearer expected-secret")
	validResp, err := app.Test(validReq)
	require.NoError(t, err)
	require.NoError(t, validResp.Body.Close())
	assert.Equal(t, fiber.StatusCreated, validResp.StatusCode)

	invalidReq := httptest.NewRequest(fiber.MethodGet, "/", nil)
	invalidReq.Header.Set("Authorization", "Bearer wrong-secret")
	invalidResp, err := app.Test(invalidReq)
	require.NoError(t, err)
	require.NoError(t, invalidResp.Body.Close())
	assert.Equal(t, fiber.StatusUnauthorized, invalidResp.StatusCode)
}
