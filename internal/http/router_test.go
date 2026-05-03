package http

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kyungseok-lee/go-fiber-gorm-starter/internal/config"
)

func TestRouter_PProfRoutesEnabledInDevelopment(t *testing.T) {
	router := NewRouter(&config.Config{
		Env:          "local",
		PProfEnabled: true,
	}, nil)
	router.Setup()

	resp, err := router.GetApp().Test(httptest.NewRequest("GET", "/debug/pprof/", nil), 5000)

	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestRouter_PProfRoutesDisabledOutsideDevelopment(t *testing.T) {
	router := NewRouter(&config.Config{
		Env:          "prod",
		PProfEnabled: true,
	}, nil)
	router.Setup()

	resp, err := router.GetApp().Test(httptest.NewRequest("GET", "/debug/pprof/", nil), 5000)

	require.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)
}
