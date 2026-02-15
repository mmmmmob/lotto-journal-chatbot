package middlewares

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Logging(c *fiber.Ctx) error {
	start := time.Now()
	err := c.Next()
	duration := time.Since(start)

	fmt.Printf("[%s] %s - duration: %s\n", c.Method(), c.OriginalURL(), duration)

	return err
}
