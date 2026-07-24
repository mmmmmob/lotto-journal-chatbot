package handler

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v3"

	"lotto-journal/api/internal/service"
)

type JobHandler struct {
	drawService   *service.DrawService
	resultService *service.ResultService
	cronSecret    string
}

func NewJobHandler(
	drawService *service.DrawService,
	resultService *service.ResultService,
	cronSecret string,
) *JobHandler {
	return &JobHandler{
		drawService:   drawService,
		resultService: resultService,
		cronSecret:    cronSecret,
	}
}

// authorize checks if the request contains the correct CRON_SECRET token
func (h *JobHandler) authorize(c fiber.Ctx) bool {
	if h.cronSecret == "" {
		log.Println("[job_handler] Warning: CRON_SECRET is not configured. Request rejected.")
		return false
	}

	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return false
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return false
	}

	return parts[1] == h.cronSecret
}

// SyncSchedule triggers the daily draw schedule sync
//
// @Summary Sync draw schedule
// @Description Syncs the lottery draw dates from the GLO API for the current year. Secured with CRON_SECRET.
// @Tags Jobs
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} map[string]string "Draw schedule sync complete"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /jobs/sync-schedule [post]
func (h *JobHandler) SyncSchedule(c fiber.Ctx) error {
	if !h.authorize(c) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	log.Println("[job_handler] Triggering manually initiated draw schedule sync...")
	if err := h.drawService.SyncDrawSchedule(c.Context()); err != nil {
		log.Printf("[job_handler] Draw schedule sync failed: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Draw schedule sync complete"})
}

// VerifyResults triggers verification of the latest lottery results
//
// @Summary Verify lottery results
// @Description Fetches the latest lottery results from GLO API and processes unchecked user tickets. Secured with CRON_SECRET.
// @Tags Jobs
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} map[string]string "Result check trigger complete"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /jobs/verify-results [post]
func (h *JobHandler) VerifyResults(c fiber.Ctx) error {
	if !h.authorize(c) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	log.Println("[job_handler] Triggering manually initiated lottery results check...")
	if err := h.resultService.VerifyLatestDrawResults(c.Context()); err != nil {
		log.Printf("[job_handler] Verification check failed: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Result check trigger complete"})
}
