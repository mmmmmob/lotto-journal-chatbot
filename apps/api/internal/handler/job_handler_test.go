package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

func TestJobHandler_Authorize(t *testing.T) {
	app := fiber.New()

	// Create job handler with empty services (nil) since we only want to test auth first
	h := NewJobHandler(nil, nil, "super-secret-token")

	app.Post("/jobs/sync-schedule", h.SyncSchedule)
	app.Post("/jobs/verify-results", h.VerifyResults)

	t.Run("missing authorization header", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/jobs/sync-schedule", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("invalid token prefix", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/jobs/sync-schedule", nil)
		req.Header.Set("Authorization", "Basic super-secret-token")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("incorrect token", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/jobs/sync-schedule", nil)
		req.Header.Set("Authorization", "Bearer wrong-token")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("empty configured secret blocks all", func(t *testing.T) {
		hEmpty := NewJobHandler(nil, nil, "")
		appEmpty := fiber.New()
		appEmpty.Post("/jobs/sync-schedule", hEmpty.SyncSchedule)

		req := httptest.NewRequest("POST", "/jobs/sync-schedule", nil)
		req.Header.Set("Authorization", "Bearer ")
		resp, err := appEmpty.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
