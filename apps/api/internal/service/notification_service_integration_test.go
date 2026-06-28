//go:build integration

package service_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"lotto-journal/api/internal/models"
	"lotto-journal/api/internal/repository"
	"lotto-journal/api/internal/service"
)

func TestNotificationService_SendDrawNotifications_Winner(t *testing.T) {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "postgres://postgres:yourpassword@localhost:5432/lotto_journal?sslmode=disable"
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("Skipping notification integration test; cannot connect to DB: %v", err)
		return
	}

	tx := db.Begin()
	defer tx.Rollback()

	// Setup mock LINE server
	var receivedPushes []messaging_api.PushMessageRequest
	var mu sync.Mutex
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && r.URL.Path == "/v2/bot/message/push" {
			body, _ := io.ReadAll(r.Body)
			var pushReq messaging_api.PushMessageRequest
			_ = json.Unmarshal(body, &pushReq)

			mu.Lock()
			receivedPushes = append(receivedPushes, pushReq)
			mu.Unlock()

			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("{}"))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer mockServer.Close()

	// Initialize bot client targeting the mock server
	bot, err := messaging_api.NewMessagingApiAPI("dummy-token", messaging_api.WithEndpoint(mockServer.URL))
	if err != nil {
		t.Fatalf("Failed to initialize bot: %v", err)
	}

	// Seed test data
	user := models.User{
		ID:         uuid.New(),
		LineUserID: "U-winner-1",
		Status:     "active",
	}
	if err := tx.Create(&user).Error; err != nil {
		t.Fatalf("seed user failed: %v", err)
	}

	drawDate, _ := time.Parse("2006-01-02", "2026-08-01")
	drawRepo := repository.NewDrawRepository(tx)
	draw, err := drawRepo.FindOrCreate(drawDate)
	if err != nil {
		t.Fatalf("seed draw failed: %v", err)
	}
	// Make sure it is marked verified in the test DB context
	if err := tx.Model(draw).Update("is_verified", true).Error; err != nil {
		t.Fatalf("failed to update is_verified: %v", err)
	}

	// Define tickets and winnings
	ticket1 := models.Ticket{
		ID:        uuid.New(),
		OwnerID:   user.ID,
		DrawID:    draw.ID,
		Type:      "N3",
		Number:    "077",
		Quantity:  2,
		IsChecked: true,
	}
	if err := tx.Create(&ticket1).Error; err != nil {
		t.Fatalf("seed ticket failed: %v", err)
	}

	drawResult1 := models.DrawResult{
		ID:            uuid.New(),
		DrawID:        draw.ID,
		PrizeCategory: "n3_straight_three",
		WinningNumber: "077",
		PrizeAmount:   8000,
	}
	if err := tx.Create(&drawResult1).Error; err != nil {
		t.Fatalf("seed drawResult failed: %v", err)
	}

	// N3 Special jackpot result
	drawResultSpecial := models.DrawResult{
		ID:            uuid.New(),
		DrawID:        draw.ID,
		PrizeCategory: "n3_special",
		WinningNumber: "077000001685",
		PrizeAmount:   700000,
	}
	if err := tx.Create(&drawResultSpecial).Error; err != nil {
		t.Fatalf("seed drawResultSpecial failed: %v", err)
	}

	winning1 := models.UserWinning{
		ID:           uuid.New(),
		TicketID:     ticket1.ID,
		DrawResultID: drawResult1.ID,
		PrizeMoney:   16000, // 8000 * 2
		UserID:       user.ID,
	}
	if err := tx.Create(&winning1).Error; err != nil {
		t.Fatalf("seed winning failed: %v", err)
	}

	// Initialize repositories and service
	ticketRepo := repository.NewTicketRepository(tx)
	winningRepo := repository.NewUserWinningRepository(tx)
	drawResultRepo := repository.NewDrawResultRepository(tx)
	notificationSvc := service.NewNotificationService(tx, bot, ticketRepo, winningRepo, drawResultRepo)

	// Execute SendDrawNotifications
	err = notificationSvc.SendDrawNotifications(context.Background(), draw.ID, "2026-08-01")
	if err != nil {
		t.Fatalf("SendDrawNotifications failed: %v", err)
	}

	// Assertions on pushed messages
	mu.Lock()
	defer mu.Unlock()
	if len(receivedPushes) != 1 {
		t.Fatalf("Expected 1 push message, got %d", len(receivedPushes))
	}
	pushReq := receivedPushes[0]
	if pushReq.To != "U-winner-1" {
		t.Errorf("Expected recipient to be U-winner-1, got %s", pushReq.To)
	}
	if len(pushReq.Messages) != 1 {
		t.Fatalf("Expected 1 message, got %d", len(pushReq.Messages))
	}
	txtMsg, ok := pushReq.Messages[0].(messaging_api.TextMessage)
	if !ok {
		t.Fatalf("Expected TextMessage, got %T", pushReq.Messages[0])
	}

	// Check content of the message
	text := txtMsg.Text
	if !contains(text, "🏆 ยินดีด้วย!") {
		t.Errorf("Expected message to contain congrats, got: %s", text)
	}
	if !contains(text, "01/08/2026") {
		t.Errorf("Expected message to contain formatted date, got: %s", text)
	}
	if !contains(text, "077") {
		t.Errorf("Expected message to contain ticket number, got: %s", text)
	}
	if !contains(text, "3 ตัวตรง x2 ใบ") {
		t.Errorf("Expected message to contain prize category with quantity, got: %s", text)
	}
	if !contains(text, "16,000 บาท") {
		t.Errorf("Expected message to contain formatted prize amount, got: %s", text)
	}
	// Check Jackpot message inclusion
	if !contains(text, "ℹ️ ลุ้นรางวัลพิเศษสามตัวท้าย (N3 Jackpot)") {
		t.Errorf("Expected message to contain N3 Jackpot alert, got: %s", text)
	}
	if !contains(text, "077000001685") {
		t.Errorf("Expected message to contain jackpot winning number, got: %s", text)
	}
	if !contains(text, "700,000 บาท") {
		t.Errorf("Expected message to contain formatted jackpot amount, got: %s", text)
	}

	// Check if audit log was written
	var logs []models.NotificationLog
	if err := tx.Where("user_id = ?", user.ID).Find(&logs).Error; err != nil {
		t.Fatalf("failed to query logs: %v", err)
	}
	if len(logs) != 1 {
		t.Fatalf("Expected 1 notification log in DB, got %d", len(logs))
	}
	logRecord := logs[0]
	if logRecord.Status != "success" {
		t.Errorf("Expected status to be success, got %s", logRecord.Status)
	}
	if logRecord.NotificationType != "draw_result" {
		t.Errorf("Expected type to be draw_result, got %s", logRecord.NotificationType)
	}
	if logRecord.DrawID == nil || *logRecord.DrawID != draw.ID {
		t.Errorf("Expected draw ID to be seed draw, got %v", logRecord.DrawID)
	}
}

func TestNotificationService_SendDrawNotifications_NonWinner(t *testing.T) {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "postgres://postgres:yourpassword@localhost:5432/lotto_journal?sslmode=disable"
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("Skipping notification integration test; cannot connect to DB: %v", err)
		return
	}

	tx := db.Begin()
	defer tx.Rollback()

	// Setup mock LINE server
	var receivedPushes []messaging_api.PushMessageRequest
	var mu sync.Mutex
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && r.URL.Path == "/v2/bot/message/push" {
			body, _ := io.ReadAll(r.Body)
			var pushReq messaging_api.PushMessageRequest
			_ = json.Unmarshal(body, &pushReq)

			mu.Lock()
			receivedPushes = append(receivedPushes, pushReq)
			mu.Unlock()

			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("{}"))
			return
		}
	}))
	defer mockServer.Close()

	bot, err := messaging_api.NewMessagingApiAPI("dummy-token", messaging_api.WithEndpoint(mockServer.URL))
	if err != nil {
		t.Fatalf("Failed to initialize bot: %v", err)
	}

	user := models.User{
		ID:         uuid.New(),
		LineUserID: "U-loser-1",
		Status:     "active",
	}
	if err := tx.Create(&user).Error; err != nil {
		t.Fatalf("seed user failed: %v", err)
	}

	drawDate, _ := time.Parse("2006-01-02", "2026-08-02")
	drawRepo := repository.NewDrawRepository(tx)
	draw, err := drawRepo.FindOrCreate(drawDate)
	if err != nil {
		t.Fatalf("seed draw failed: %v", err)
	}
	if err := tx.Model(draw).Update("is_verified", true).Error; err != nil {
		t.Fatalf("failed to update is_verified: %v", err)
	}

	ticket := models.Ticket{
		ID:        uuid.New(),
		OwnerID:   user.ID,
		DrawID:    draw.ID,
		Type:      "L6",
		Number:    "999999",
		Quantity:  1,
		IsChecked: true,
	}
	if err := tx.Create(&ticket).Error; err != nil {
		t.Fatalf("seed ticket failed: %v", err)
	}

	ticketRepo := repository.NewTicketRepository(tx)
	winningRepo := repository.NewUserWinningRepository(tx)
	drawResultRepo := repository.NewDrawResultRepository(tx)
	notificationSvc := service.NewNotificationService(tx, bot, ticketRepo, winningRepo, drawResultRepo)

	err = notificationSvc.SendDrawNotifications(context.Background(), draw.ID, "2026-08-02")
	if err != nil {
		t.Fatalf("SendDrawNotifications failed: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if len(receivedPushes) != 1 {
		t.Fatalf("Expected 1 push message, got %d", len(receivedPushes))
	}
	txtMsg, ok := receivedPushes[0].Messages[0].(messaging_api.TextMessage)
	if !ok {
		t.Fatalf("Expected TextMessage, got %T", receivedPushes[0].Messages[0])
	}
	text := txtMsg.Text

	if !contains(text, "เสียใจด้วยครับ") {
		t.Errorf("Expected loser message to contain consolation, got: %s", text)
	}
	if !contains(text, "999999 (L6) x1 ใบ") {
		t.Errorf("Expected ticket details in loser message, got: %s", text)
	}

	// Check if audit log was written
	var logs []models.NotificationLog
	tx.Where("user_id = ?", user.ID).Find(&logs)
	if len(logs) != 1 {
		t.Fatalf("Expected 1 notification log in DB, got %d", len(logs))
	}
	if logs[0].Status != "success" {
		t.Errorf("Expected status success, got %s", logs[0].Status)
	}
}

func TestNotificationService_SendDrawNotifications_RetryFailure(t *testing.T) {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "postgres://postgres:yourpassword@localhost:5432/lotto_journal?sslmode=disable"
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("Skipping notification integration test; cannot connect to DB: %v", err)
		return
	}

	tx := db.Begin()
	defer tx.Rollback()

	// Setup mock LINE server returning 500
	var attempts int
	var mu sync.Mutex
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		attempts++
		mu.Unlock()
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"message": "Internal error"}`))
	}))
	defer mockServer.Close()

	bot, err := messaging_api.NewMessagingApiAPI("dummy-token", messaging_api.WithEndpoint(mockServer.URL))
	if err != nil {
		t.Fatalf("Failed to initialize bot: %v", err)
	}

	user := models.User{
		ID:         uuid.New(),
		LineUserID: "U-failed-1",
		Status:     "active",
	}
	if err := tx.Create(&user).Error; err != nil {
		t.Fatalf("seed user failed: %v", err)
	}

	drawDate, _ := time.Parse("2006-01-02", "2026-08-03")
	drawRepo := repository.NewDrawRepository(tx)
	draw, err := drawRepo.FindOrCreate(drawDate)
	if err != nil {
		t.Fatalf("seed draw failed: %v", err)
	}
	if err := tx.Model(draw).Update("is_verified", true).Error; err != nil {
		t.Fatalf("failed to update is_verified: %v", err)
	}

	ticket := models.Ticket{
		ID:        uuid.New(),
		OwnerID:   user.ID,
		DrawID:    draw.ID,
		Type:      "L6",
		Number:    "123456",
		Quantity:  1,
		IsChecked: true,
	}
	if err := tx.Create(&ticket).Error; err != nil {
		t.Fatalf("seed ticket failed: %v", err)
	}

	ticketRepo := repository.NewTicketRepository(tx)
	winningRepo := repository.NewUserWinningRepository(tx)
	drawResultRepo := repository.NewDrawResultRepository(tx)
	notificationSvc := service.NewNotificationService(tx, bot, ticketRepo, winningRepo, drawResultRepo)

	err = notificationSvc.SendDrawNotifications(context.Background(), draw.ID, "2026-08-03")
	if err != nil {
		t.Fatalf("SendDrawNotifications should not fail the overall process even if push fails: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if attempts != 3 {
		t.Errorf("Expected exactly 3 retry attempts, got %d", attempts)
	}

	// Check if audit log was written as failed
	var logs []models.NotificationLog
	tx.Where("user_id = ?", user.ID).Find(&logs)
	if len(logs) != 1 {
		t.Fatalf("Expected 1 notification log in DB, got %d", len(logs))
	}
	if logs[0].Status != "failed" {
		t.Errorf("Expected status failed, got %s", logs[0].Status)
	}
	if logs[0].ErrorMessage == nil || *logs[0].ErrorMessage == "" {
		t.Errorf("Expected error message to be set, got empty/nil")
	}
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
