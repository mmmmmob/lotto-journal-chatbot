package middlewares

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

// Logging logs method, path, status code, duration, and request ID for every request.
// It expects the requestid middleware to run before it so that requestid.FromContext(c) is set.
func Logging(c fiber.Ctx) error {
	start := time.Now()
	err := c.Next()
	duration := time.Since(start)

	reqID := requestid.FromContext(c)
	status := c.Response().StatusCode()

	// Suppress Fly platform health-check noise while keeping user/manual /health logs.
	if strings.HasPrefix(c.OriginalURL(), "/health") && strings.EqualFold(c.Get("X-Health-Check"), "fly") {
		return err
	}

	fmt.Printf("[%s] %s - %d - %s - req_id: %s\n",
		c.Method(), c.OriginalURL(), status, duration, reqID)

	return err
}
