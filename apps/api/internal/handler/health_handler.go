package handler

import (
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
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status": "degraded",
			"db":     err.Error(),
		})
	}

	if err := sqlDB.Ping(); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status": "degraded",
			"db":     err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "ok",
		"db":     "ok",
	})
}
