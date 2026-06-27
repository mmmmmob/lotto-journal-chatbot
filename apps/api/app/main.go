package main

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v3"
	recoverer "github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/gofiber/fiber/v3/middleware/timeout"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"

	"lotto-journal/api/internal/client"
	"lotto-journal/api/internal/config"
	"lotto-journal/api/internal/database"
	"lotto-journal/api/internal/handler"
	"lotto-journal/api/internal/repository"
	"lotto-journal/api/internal/service"
	"lotto-journal/api/middlewares"

	_ "lotto-journal/api/docs"

	"github.com/gofiber/contrib/v3/swaggo"
)

// @title Lotto Journal API
// @version 0.2
// @description Backend API for the Lotto Journal LINE Bot and verification engine.
// @host localhost:3000
// @BasePath /
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

	// GLO Lottery client
	lotteryClient := client.NewLotteryClient("")

	// Repositories
	userRepo := repository.NewUserRepository(db)
	drawRepo := repository.NewDrawRepository(db)
	ticketRepo := repository.NewTicketRepository(db)
	webhookRepo := repository.NewWebhookEventRepository(db)
	drawResultRepo := repository.NewDrawResultRepository(db)
	winningRepo := repository.NewUserWinningRepository(db)

	// Services
	userSvc := service.NewUserService(userRepo)
	drawSvc := service.NewDrawService(drawRepo, lotteryClient)
	ticketSvc := service.NewTicketService(ticketRepo, drawRepo, drawSvc)
	resultSvc := service.NewResultService(db, lotteryClient, drawRepo, drawResultRepo, ticketRepo, winningRepo)

	// Start background cron scheduler
	scheduler := service.NewCronScheduler(drawSvc, resultSvc, cfg.CronSyncSchedule, cfg.CronVerifySchedule)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go scheduler.Start(ctx)

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

	if cfg.APP_ENV != "production" {
		log.Println("[main] Swagger UI available at /swagger/index.html")
		app.Get("/swagger/*", swaggo.HandlerDefault)
	}

	// timeout.New wraps the handler: if it does not return within 25 s, Fiber responds 408.
	// Applied only to /webhook — other routes (e.g. /health) have no artificial deadline.
	app.Post("/webhook", timeout.New(lineHandler.Handle, timeout.Config{
		Timeout: 25 * time.Second,
	}))

	// Start server
	log.Fatal(app.Listen(cfg.PORT))
}
