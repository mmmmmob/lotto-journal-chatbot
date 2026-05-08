package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v3"
	recoverer "github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/gofiber/fiber/v3/middleware/timeout"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"

	"lotto-journal/api/internal/config"
	"lotto-journal/api/internal/database"
	"lotto-journal/api/internal/handler"
	"lotto-journal/api/internal/repository"
	"lotto-journal/api/internal/service"
	"lotto-journal/api/middlewares"
)

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Connect to database
	database.ConnectDatabase(cfg.DB_DSN)
	db := database.DB

	// LINE bot client
	bot, err := messaging_api.NewMessagingApiAPI(cfg.LineChannelAccessToken)
	if err != nil {
		log.Fatalf("Failed to create LINE bot client: %v", err)
	}

	// Repositories
	userRepo := repository.NewUserRepository(db)
	drawRepo := repository.NewDrawRepository(db)
	ticketRepo := repository.NewTicketRepository(db)
	webhookRepo := repository.NewWebhookEventRepository(db)

	// Services
	userSvc := service.NewUserService(userRepo)
	ticketSvc := service.NewTicketService(ticketRepo, drawRepo)

	// Handlers
	healthHandler := handler.NewHealthHandler(db)
	lineHandler := handler.NewLineHandler(
		cfg.LineChannelSecret,
		bot,
		userSvc,
		ticketSvc,
		webhookRepo,
	)

	// Fiber app
	app := fiber.New()

	// Global middleware — order matters:
	//   1. recoverer : catch panics before anything else runs, return 500 + stack trace
	//   2. requestid : assign a trace ID so subsequent middleware can read it
	//   3. Logging   : log after the handler finishes (reads requestid + status code)
	app.Use(recoverer.New(recoverer.Config{EnableStackTrace: true}))
	app.Use(requestid.New())
	app.Use(middlewares.Logging)

	// Routes
	app.Get("/health", healthHandler.Handle)

	// timeout.New wraps the handler: if it does not return within 25 s, Fiber responds 408.
	// Applied only to /webhook — other routes (e.g. /health) have no artificial deadline.
	app.Post("/webhook", timeout.New(lineHandler.Handle, timeout.Config{
		Timeout: 25 * time.Second,
	}))

	// Start server
	log.Fatal(app.Listen(cfg.PORT))
}
