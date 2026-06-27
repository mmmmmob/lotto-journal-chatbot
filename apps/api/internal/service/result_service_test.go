package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"lotto-journal/api/internal/client"
	"lotto-journal/api/internal/models"
	"lotto-journal/api/internal/repository"
)

func TestVerifyDrawResults_Integration(t *testing.T) {
	// 1. Connect to local PostgreSQL test database using env vars
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		// Fallback to local default for development testing
		dsn = "postgres://postgres:yourpassword@localhost:5432/lotto_journal?sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("Skipping integration test; cannot connect to DB: %v", err)
		return
	}

	// 2. Start a transaction that we will roll back at the end
	tx := db.Begin()
	defer tx.Rollback()

	// 3. Create mock GLO API server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"statusMessage": "getLatestLottery - Success",
			"statusCode": 200,
			"status": true,
			"response": {
				"date": "2026-07-01",
				"data": {
					"first": { "price": "6,000,000.00", "number": [{"round": 1, "value": "123456"}] },
					"second": { "price": "200000.00", "number": [] },
					"third": { "price": "80000.00", "number": [] },
					"fourth": { "price": "40000.00", "number": [] },
					"fifth": { "price": "20000.00", "number": [] },
					"last2": { "price": "2000.00", "number": [{"round": 1, "value": "43"}] },
					"last3f": { "price": "4000.00", "number": [{"round": 1, "value": "267"}] },
					"last3b": { "price": "4000.00", "number": [{"round": 1, "value": "065"}] },
					"near1": { "price": "100000.00", "number": [] }
				},
				"n3": {
					"straight3": { "price": "8000.00", "number": [{"round": 1, "value": "077"}] },
					"shuffle3": { "price": "3000.00", "number": [{"round": 1, "value": "707"}] },
					"straight2": { "price": "500.00", "number": [{"round": 1, "value": "43"}] },
					"special": { "price": "700000.00", "number": [{"round": 1, "value": "077000001685"}] }
				}
			}
		}`))
	}))
	defer mockServer.Close()

	// 4. Initialize repositories and service inside the transaction
	lotteryClient := client.NewLotteryClient(mockServer.URL)
	drawRepo := repository.NewDrawRepository(tx)
	drawResultRepo := repository.NewDrawResultRepository(tx)
	ticketRepo := repository.NewTicketRepository(tx)
	winningRepo := repository.NewUserWinningRepository(tx)

	resultSvc := NewResultService(tx, lotteryClient, drawRepo, drawResultRepo, ticketRepo, winningRepo)

	// 5. Seed test user
	user := models.User{
		ID:         uuid.New(),
		LineUserID: "U" + uuid.New().String()[:10],
		Status:     "active",
	}
	if err := tx.Create(&user).Error; err != nil {
		t.Fatalf("failed to seed test user: %v", err)
	}

	// 6. Seed test draw
	drawDate, _ := time.Parse("2006-01-02", "2026-07-01")
	draw, err := drawRepo.FindOrCreate(drawDate)
	if err != nil {
		t.Fatalf("failed to seed test draw: %v", err)
	}

	// 7. Seed test tickets
	tickets := []*models.Ticket{
		{ID: uuid.New(), OwnerID: user.ID, DrawID: draw.ID, Type: "L6", Number: "123456", Quantity: 2, IsChecked: false}, // L6 winner first prize
		{ID: uuid.New(), OwnerID: user.ID, DrawID: draw.ID, Type: "L6", Number: "000043", Quantity: 1, IsChecked: false}, // L6 winner last 2
		{ID: uuid.New(), OwnerID: user.ID, DrawID: draw.ID, Type: "N3", Number: "077", Quantity: 1, IsChecked: false},    // N3 winner straight 3
		{ID: uuid.New(), OwnerID: user.ID, DrawID: draw.ID, Type: "N3", Number: "707", Quantity: 1, IsChecked: false},    // N3 winner shuffle 3
		{ID: uuid.New(), OwnerID: user.ID, DrawID: draw.ID, Type: "N3", Number: "943", Quantity: 1, IsChecked: false},    // N3 winner straight 2
		{ID: uuid.New(), OwnerID: user.ID, DrawID: draw.ID, Type: "L6", Number: "999999", Quantity: 1, IsChecked: false}, // L6 non-winner
	}

	for _, ticket := range tickets {
		if err := tx.Create(ticket).Error; err != nil {
			t.Fatalf("failed to seed ticket %s: %v", ticket.Number, err)
		}
	}

	// 8. Run Verification Service
	err = resultSvc.VerifyDrawResults(context.Background(), drawDate)
	if err != nil {
		t.Fatalf("VerifyDrawResults failed: %v", err)
	}

	// 9. Assertions
	// Check if draw is marked verified
	var updatedDraw models.Draw
	if err := tx.First(&updatedDraw, "id = ?", draw.ID).Error; err != nil {
		t.Fatalf("failed to retrieve draw: %v", err)
	}
	if !updatedDraw.IsVerified {
		t.Errorf("expected draw to be verified, but got false")
	}

	// Check if all tickets are marked as checked
	var checkedCount int64
	tx.Model(&models.Ticket{}).Where("draw_id = ? AND is_checked = true", draw.ID).Count(&checkedCount)
	if checkedCount != int64(len(tickets)) {
		t.Errorf("expected all %d tickets to be checked, got %d", len(tickets), checkedCount)
	}

	// Check user winnings populated
	var winnings []models.UserWinning
	if err := tx.Where("user_id = ?", user.ID).Find(&winnings).Error; err != nil {
		t.Fatalf("failed to query winnings: %v", err)
	}

	// Expected winnings count:
	// - L6 first prize: 123456 x2 -> 6,000,000 * 2 = 12,000,000
	// - L6 last 2: 000043 x1 -> 2,000 * 1 = 2,000
	// - N3 straight 3: 077 x1 -> 8,000 * 1 = 8,000
	// - N3 shuffle 3: 707 x1 -> 3,000 * 1 = 3,000
	// - N3 straight 2: 943 x1 -> 500 * 1 = 500
	// Total winning entries: 5
	if len(winnings) != 5 {
		t.Errorf("expected 5 winning entries, got %d: %+v", len(winnings), winnings)
	}

	expectedAmounts := map[uuid.UUID]int{
		tickets[0].ID: 12000000, // first prize
		tickets[1].ID: 2000,     // last 2
		tickets[2].ID: 8000,     // straight 3
		tickets[3].ID: 3000,     // shuffle 3
		tickets[4].ID: 500,      // straight 2
	}

	for _, w := range winnings {
		expected, exists := expectedAmounts[w.TicketID]
		if !exists {
			t.Errorf("unexpected winning record for ticket ID %s", w.TicketID)
			continue
		}
		if w.PrizeMoney != expected {
			t.Errorf("expected prize money %d for ticket %s, got %d", expected, w.TicketID, w.PrizeMoney)
		}
	}
}

func TestVerifyDrawResults_PendingState(t *testing.T) {
	// Create mock server returning a date other than requested
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"statusMessage": "getLatestLottery - Success",
			"statusCode": 200,
			"status": true,
			"response": {
				"date": "2026-06-16",
				"data": { "first": { "price": "6000000.00", "number": [] } }
			}
		}`))
	}))
	defer mockServer.Close()

	lotteryClient := client.NewLotteryClient(mockServer.URL)
	resultSvc := NewResultService(nil, lotteryClient, nil, nil, nil, nil)

	drawDate, _ := time.Parse("2006-01-02", "2026-07-01")
	err := resultSvc.VerifyDrawResults(context.Background(), drawDate)

	if err == nil {
		t.Fatalf("expected error due to results pending, but got nil")
	}

	expectedErr := "results pending: latest GLO draw is 2026-06-16, expected 2026-07-01"
	if err.Error() != expectedErr {
		t.Errorf("expected error message %q, got %q", expectedErr, err.Error())
	}
}

func TestParsePrizeAmount(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"6000000.00", 6000000},
		{"6,000,000.00", 6000000},
		{"200000", 200000},
		{"2,000", 2000},
		{"500.00", 500},
		{"0", 0},
		{"", 0},
		{"invalid", 0},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got := parsePrizeAmount(tc.input)
			if got != tc.expected {
				t.Errorf("parsePrizeAmount(%q) = %d; expected %d", tc.input, got, tc.expected)
			}
		})
	}
}
