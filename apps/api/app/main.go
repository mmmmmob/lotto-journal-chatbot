package main

import (
	"lotto-journal/api/internal/config"
	"lotto-journal/api/internal/database"
	"lotto-journal/api/middlewares"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// load config
	cfg := config.LoadConfig()

	// connect to database
	database.ConnectDatabase(cfg.DB_DSN)

	// create fiber instance
	app := fiber.New()

	// run middlewares
	app.Use(middlewares.Logging)

	// api routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello! Lotto Journal!")
	})
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// TODO(T-002): register LINE webhook handler here
	// app.Post("/webhook", lineHandler.Handle)

	// start server
	app.Listen(cfg.PORT)
}
