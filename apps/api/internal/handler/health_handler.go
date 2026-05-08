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
// Response shape:
//
//	200  {"status":"ok",       "db":"ok"}
//	503  {"status":"degraded", "db":"<error message>"}
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
