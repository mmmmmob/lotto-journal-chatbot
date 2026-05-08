package middlewares

import (
	"fmt"
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

	fmt.Printf("[%s] %s - %d - %s - req_id: %s\n",
		c.Method(), c.OriginalURL(), status, duration, reqID)

	return err
}
