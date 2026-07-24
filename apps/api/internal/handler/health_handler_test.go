package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestHealthHandler_Handle(t *testing.T) {
	t.Run("happy path - status ok", func(t *testing.T) {
		app := fiber.New()
		app.Use(requestid.New())

		mockDB, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		require.NoError(t, err)
		defer mockDB.Close()

		mock.ExpectPing()

		// Construct gorm.DB directly to bypass GORM dialector version check queries
		db := &gorm.DB{
			Config: &gorm.Config{
				ConnPool: mockDB,
			},
		}

		h := NewHealthHandler(db)
		app.Get("/health", h.Handle)

		req := httptest.NewRequest("GET", "/health", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		require.NotNil(t, resp)
		t.Cleanup(func() { resp.Body.Close() })

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body HealthResponse
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)
		assert.Equal(t, "ok", body.Status)
		assert.Equal(t, "ok", body.DB)
		assert.NotEmpty(t, body.RequestID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db.DB() failure path - status degraded", func(t *testing.T) {
		app := fiber.New()
		app.Use(requestid.New())

		// Force db.DB() to fail by leaving ConnPool as nil
		db := &gorm.DB{
			Config: &gorm.Config{
				ConnPool: nil,
			},
		}

		h := NewHealthHandler(db)
		app.Get("/health", h.Handle)

		req := httptest.NewRequest("GET", "/health", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		require.NotNil(t, resp)
		t.Cleanup(func() { resp.Body.Close() })

		assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)

		var body HealthResponse
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)
		assert.Equal(t, "degraded", body.Status)
		assert.Equal(t, "unhealthy", body.DB)
		assert.NotEmpty(t, body.RequestID)
	})

	t.Run("PingContext failure - status degraded", func(t *testing.T) {
		app := fiber.New()
		app.Use(requestid.New())

		mockDB, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		require.NoError(t, err)
		defer mockDB.Close()

		// Mock ping to return an error
		mock.ExpectPing().WillReturnError(errors.New("db connection failure"))

		db := &gorm.DB{
			Config: &gorm.Config{
				ConnPool: mockDB,
			},
		}

		h := NewHealthHandler(db)
		app.Get("/health", h.Handle)

		req := httptest.NewRequest("GET", "/health", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		require.NotNil(t, resp)
		t.Cleanup(func() { resp.Body.Close() })

		assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)

		var body HealthResponse
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)
		assert.Equal(t, "degraded", body.Status)
		assert.Equal(t, "unhealthy", body.DB)
		assert.NotEmpty(t, body.RequestID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("timeout/cancellation behavior - status degraded", func(t *testing.T) {
		app := fiber.New()
		app.Use(requestid.New())

		mockDB, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		require.NoError(t, err)
		defer mockDB.Close()

		// Mock ping to return context.Canceled (simulating timeout or cancellation)
		mock.ExpectPing().WillReturnError(context.Canceled)

		db := &gorm.DB{
			Config: &gorm.Config{
				ConnPool: mockDB,
			},
		}

		h := NewHealthHandler(db)
		app.Get("/health", h.Handle)

		req := httptest.NewRequest("GET", "/health", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		require.NotNil(t, resp)
		t.Cleanup(func() { resp.Body.Close() })

		assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)

		var body HealthResponse
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)
		assert.Equal(t, "degraded", body.Status)
		assert.Equal(t, "unhealthy", body.DB)
		assert.NotEmpty(t, body.RequestID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
