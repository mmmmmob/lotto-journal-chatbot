package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
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
	lineHandler := handler.NewLineHandler(
		cfg.LineChannelSecret,
		bot,
		userSvc,
		ticketSvc,
		webhookRepo,
	)

	// Fiber app
	app := fiber.New()
	app.Use(middlewares.Logging)

	// Routes
	app.Post("/webhook", lineHandler.Handle)

	// Start server
	app.Listen(cfg.PORT)
}
