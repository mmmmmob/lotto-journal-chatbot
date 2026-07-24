package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubDrawService struct {
	syncCalled bool
	syncErr    error
}

func (s *stubDrawService) SyncDrawSchedule(ctx context.Context) error {
	s.syncCalled = true
	return s.syncErr
}

type stubResultService struct {
	verifyCalled bool
	verifyErr    error
}

func (s *stubResultService) VerifyLatestDrawResults(ctx context.Context) error {
	s.verifyCalled = true
	return s.verifyErr
}

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
		require.NotNil(t, resp)
		t.Cleanup(func() { resp.Body.Close() })
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("invalid token prefix", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/jobs/sync-schedule", nil)
		req.Header.Set("Authorization", "Basic super-secret-token")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		require.NotNil(t, resp)
		t.Cleanup(func() { resp.Body.Close() })
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("incorrect token", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/jobs/sync-schedule", nil)
		req.Header.Set("Authorization", "Bearer wrong-token")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		require.NotNil(t, resp)
		t.Cleanup(func() { resp.Body.Close() })
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
		require.NotNil(t, resp)
		t.Cleanup(func() { resp.Body.Close() })
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestJobHandler_Endpoints(t *testing.T) {
	t.Run("sync-schedule success happy path", func(t *testing.T) {
		app := fiber.New()
		drawSvc := &stubDrawService{}
		h := NewJobHandler(drawSvc, nil, "super-secret-token")
		app.Post("/jobs/sync-schedule", h.SyncSchedule)

		req := httptest.NewRequest("POST", "/jobs/sync-schedule", nil)
		req.Header.Set("Authorization", "Bearer super-secret-token")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		require.NotNil(t, resp)
		t.Cleanup(func() { resp.Body.Close() })
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.True(t, drawSvc.syncCalled)
	})

	t.Run("verify-results success happy path", func(t *testing.T) {
		app := fiber.New()
		resultSvc := &stubResultService{}
		h := NewJobHandler(nil, resultSvc, "super-secret-token")
		app.Post("/jobs/verify-results", h.VerifyResults)

		req := httptest.NewRequest("POST", "/jobs/verify-results", nil)
		req.Header.Set("Authorization", "Bearer super-secret-token")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		require.NotNil(t, resp)
		t.Cleanup(func() { resp.Body.Close() })
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.True(t, resultSvc.verifyCalled)
	})

	t.Run("sync-schedule service failure returns 500", func(t *testing.T) {
		app := fiber.New()
		drawSvc := &stubDrawService{syncErr: assert.AnError}
		h := NewJobHandler(drawSvc, nil, "super-secret-token")
		app.Post("/jobs/sync-schedule", h.SyncSchedule)

		req := httptest.NewRequest("POST", "/jobs/sync-schedule", nil)
		req.Header.Set("Authorization", "Bearer super-secret-token")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		require.NotNil(t, resp)
		t.Cleanup(func() { resp.Body.Close() })
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.True(t, drawSvc.syncCalled)
	})
	t.Run("verify-results service failure returns 500", func(t *testing.T) {
		app := fiber.New()
		resultSvc := &stubResultService{verifyErr: assert.AnError}
		h := NewJobHandler(nil, resultSvc, "super-secret-token")
		app.Post("/jobs/verify-results", h.VerifyResults)

		req := httptest.NewRequest("POST", "/jobs/verify-results", nil)
		req.Header.Set("Authorization", "Bearer super-secret-token")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		require.NotNil(t, resp)
		t.Cleanup(func() { resp.Body.Close() })
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.True(t, resultSvc.verifyCalled)
	})

	t.Run("sync-schedule nil service returns 500", func(t *testing.T) {
		app := fiber.New()
		h := NewJobHandler(nil, nil, "super-secret-token")
		app.Post("/jobs/sync-schedule", h.SyncSchedule)

		req := httptest.NewRequest("POST", "/jobs/sync-schedule", nil)
		req.Header.Set("Authorization", "Bearer super-secret-token")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		require.NotNil(t, resp)
		t.Cleanup(func() { resp.Body.Close() })
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("verify-results nil service returns 500", func(t *testing.T) {
		app := fiber.New()
		h := NewJobHandler(nil, nil, "super-secret-token")
		app.Post("/jobs/verify-results", h.VerifyResults)

		req := httptest.NewRequest("POST", "/jobs/verify-results", nil)
		req.Header.Set("Authorization", "Bearer super-secret-token")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		require.NotNil(t, resp)
		t.Cleanup(func() { resp.Body.Close() })
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
