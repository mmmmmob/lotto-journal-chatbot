package handler

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// HealthHandler handles GET /health.
type HealthHandler struct {
	db *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// Handle performs a liveness + DB readiness check.
//
// @Summary Perform health check
// @Description Checks if the API server is live and the database connection is healthy.
// @Tags System
// @Produce json
// @Success 200 {object} map[string]string "Successful health status"
// @Failure 503 {object} map[string]string "Database connection is unhealthy"
// @Router /health [get]
func (h *HealthHandler) Handle(c fiber.Ctx) error {
	sqlDB, err := h.db.DB()
	if err != nil {
		log.Printf("[health] failed to get underlying sql.DB: %v", err)
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status":     "degraded",
			"db":         "error",
			"request_id": c.Get(fiber.HeaderXRequestID),
		})
	}

	// Enforce a 5-second timeout on the database ping to prevent hanging
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		log.Printf("[health] database ping failed: %v", err)
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status":     "degraded",
			"db":         "unhealthy",
			"request_id": c.Get(fiber.HeaderXRequestID),
		})
	}

	return c.JSON(fiber.Map{
		"status": "ok",
		"db":     "ok",
	})
}
