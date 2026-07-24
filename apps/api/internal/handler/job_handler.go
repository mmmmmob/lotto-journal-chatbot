package handler

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"log"
	"strings"

	"github.com/gofiber/fiber/v3"
)

type DrawServiceInterface interface {
	SyncDrawSchedule(ctx context.Context) error
}

type ResultServiceInterface interface {
	VerifyLatestDrawResults(ctx context.Context) error
}

type JobHandler struct {
	drawService    DrawServiceInterface
	resultService  ResultServiceInterface
	cronSecretHash [32]byte
	hasSecret      bool
}

func NewJobHandler(
	drawService DrawServiceInterface,
	resultService ResultServiceInterface,
	cronSecret string,
) *JobHandler {
	var secretHash [32]byte
	hasSecret := cronSecret != ""
	if hasSecret {
		secretHash = sha256.Sum256([]byte(cronSecret))
	}
	return &JobHandler{
		drawService:    drawService,
		resultService:  resultService,
		cronSecretHash: secretHash,
		hasSecret:      hasSecret,
	}
}

// authorize checks if the request contains the correct CRON_SECRET token using constant-time comparison on SHA-256 hashes
func (h *JobHandler) authorize(c fiber.Ctx) bool {
	if !h.hasSecret {
		return false
	}

	authHeader := c.Get("Authorization")
	const prefix = "Bearer "
	if len(authHeader) < len(prefix) || !strings.EqualFold(authHeader[:len(prefix)], prefix) {
		return false
	}

	token := strings.TrimSpace(authHeader[len(prefix):])
	if token == "" {
		return false
	}

	// Hash the incoming token to ensure the comparison operates on fixed-length (32-byte) inputs,
	// preventing leaking the length of the expected h.cronSecret through timing side channels.
	tokenHash := sha256.Sum256([]byte(token))

	return subtle.ConstantTimeCompare(tokenHash[:], h.cronSecretHash[:]) == 1
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

	if h.drawService == nil {
		log.Println("[job_handler] Error: drawService dependency is nil")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":      "Internal Server Error",
			"request_id": c.Get(fiber.HeaderXRequestID),
		})
	}

	log.Printf("[job_handler] Triggering draw schedule sync (request_id: %s)...", c.Get(fiber.HeaderXRequestID))
	if err := h.drawService.SyncDrawSchedule(c.Context()); err != nil {
		log.Printf("[job_handler] Draw schedule sync failed (request_id: %s): %v", c.Get(fiber.HeaderXRequestID), err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":      "Failed to sync draw schedule",
			"request_id": c.Get(fiber.HeaderXRequestID),
		})
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

	if h.resultService == nil {
		log.Println("[job_handler] Error: resultService dependency is nil")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":      "Internal Server Error",
			"request_id": c.Get(fiber.HeaderXRequestID),
		})
	}

	log.Printf("[job_handler] Triggering lottery results check (request_id: %s)...", c.Get(fiber.HeaderXRequestID))
	if err := h.resultService.VerifyLatestDrawResults(c.Context()); err != nil {
		log.Printf("[job_handler] Verification check failed (request_id: %s): %v", c.Get(fiber.HeaderXRequestID), err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":      "Failed to verify draw results",
			"request_id": c.Get(fiber.HeaderXRequestID),
		})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Result check trigger complete"})
}
