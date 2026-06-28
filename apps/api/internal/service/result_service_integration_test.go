//go:build integration

package service_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"lotto-journal/api/internal/client"
	"lotto-journal/api/internal/mocks"
	"lotto-journal/api/internal/models"
	"lotto-journal/api/internal/repository"
	"lotto-journal/api/internal/service"
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
				"date": "2026-09-01",
				"data": {
					"first": { "price": "6,000,000.00", "number": [{"round": 1, "value": "123456"}] },
					"second": { "price": "200,000.00", "number": [{"round": 1, "value": "222222"}] },
					"third": { "price": "80,000.00", "number": [{"round": 1, "value": "333333"}] },
					"fourth": { "price": "40,000.00", "number": [{"round": 1, "value": "444444"}] },
					"fifth": { "price": "20,000.00", "number": [{"round": 1, "value": "555555"}] },
					"last2": { "price": "2,000.00", "number": [{"round": 1, "value": "43"}] },
					"last3f": { "price": "4,000.00", "number": [{"round": 1, "value": "267"}] },
					"last3b": { "price": "4,000.00", "number": [{"round": 1, "value": "065"}] },
					"near1": { "price": "100,000.00", "number": [{"round": 1, "value": "123455"}, {"round": 1, "value": "123457"}] }
				},
				"n3": {
					"straight3": { "price": "8,000.00", "number": [{"round": 1, "value": "077"}] },
					"shuffle3": { "price": "3,000.00", "number": [{"round": 1, "value": "707"}] },
					"straight2": { "price": "500.00", "number": [{"round": 1, "value": "43"}] },
					"special": { "price": "700,000.00", "number": [{"round": 1, "value": "077000001685"}] }
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

	mockNotifSvc := mocks.NewMockNotificationServiceInterface(t)
	mockNotifSvc.On("SendDrawNotifications", mock.Anything, mock.Anything, "2026-09-01").Return(nil)

	resultSvc := service.NewResultService(tx, lotteryClient, drawRepo, drawResultRepo, ticketRepo, winningRepo, mockNotifSvc)

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
	drawDate, _ := time.Parse("2006-01-02", "2026-09-01")
	draw, err := drawRepo.FindOrCreate(drawDate)
	if err != nil {
		t.Fatalf("failed to seed test draw: %v", err)
	}

	// 7. Seed test tickets
	tickets := []*models.Ticket{
		{ID: uuid.New(), OwnerID: user.ID, DrawID: draw.ID, Type: "L6", Number: "123456", Quantity: 2, IsChecked: false}, // L6 winner first prize
		{ID: uuid.New(), OwnerID: user.ID, DrawID: draw.ID, Type: "L6", Number: "222222", Quantity: 1, IsChecked: false}, // L6 winner second prize
		{ID: uuid.New(), OwnerID: user.ID, DrawID: draw.ID, Type: "L6", Number: "333333", Quantity: 1, IsChecked: false}, // L6 winner third prize
		{ID: uuid.New(), OwnerID: user.ID, DrawID: draw.ID, Type: "L6", Number: "444444", Quantity: 1, IsChecked: false}, // L6 winner fourth prize
		{ID: uuid.New(), OwnerID: user.ID, DrawID: draw.ID, Type: "L6", Number: "555555", Quantity: 1, IsChecked: false}, // L6 winner fifth prize
		{ID: uuid.New(), OwnerID: user.ID, DrawID: draw.ID, Type: "L6", Number: "123455", Quantity: 1, IsChecked: false}, // L6 winner near first
		{ID: uuid.New(), OwnerID: user.ID, DrawID: draw.ID, Type: "L6", Number: "000043", Quantity: 1, IsChecked: false}, // L6 winner last 2
		{ID: uuid.New(), OwnerID: user.ID, DrawID: draw.ID, Type: "L6", Number: "267000", Quantity: 1, IsChecked: false}, // L6 winner first 3 (front)
		{ID: uuid.New(), OwnerID: user.ID, DrawID: draw.ID, Type: "L6", Number: "000065", Quantity: 1, IsChecked: false}, // L6 winner last 3 (back)
		{ID: uuid.New(), OwnerID: user.ID, DrawID: draw.ID, Type: "L6", Number: "267043", Quantity: 1, IsChecked: false}, // L6 winner first 3 AND last 2
		{ID: uuid.New(), OwnerID: user.ID, DrawID: draw.ID, Type: "N3", Number: "077", Quantity: 1, IsChecked: false},    // N3 winner straight 3
		{ID: uuid.New(), OwnerID: user.ID, DrawID: draw.ID, Type: "N3", Number: "707", Quantity: 1, IsChecked: false},    // N3 winner shuffle 3
		{ID: uuid.New(), OwnerID: user.ID, DrawID: draw.ID, Type: "N3", Number: "943", Quantity: 1, IsChecked: false},    // N3 winner straight 2
		{ID: uuid.New(), OwnerID: user.ID, DrawID: draw.ID, Type: "L6", Number: "999999", Quantity: 1, IsChecked: false}, // L6 non-winner
		{ID: uuid.New(), OwnerID: user.ID, DrawID: draw.ID, Type: "L6", Number: "12345", Quantity: 1, IsChecked: false},  // L6 invalid length (ignored)
		{ID: uuid.New(), OwnerID: user.ID, DrawID: draw.ID, Type: "N3", Number: "12", Quantity: 1, IsChecked: false},     // N3 invalid length (ignored)
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
	// - L6 first prize: 123456 x2 -> 6,000,000 * 2 = 12,000,000 (1 entry)
	// - L6 second prize: 222222 x1 -> 200,000 * 1 = 200,000 (1 entry)
	// - L6 third prize: 333333 x1 -> 80,000 * 1 = 80,000 (1 entry)
	// - L6 fourth prize: 444444 x1 -> 40,000 * 1 = 40,000 (1 entry)
	// - L6 fifth prize: 555555 x1 -> 20,000 * 1 = 20,000 (1 entry)
	// - L6 near first: 123455 x1 -> 100,000 * 1 = 100,000 (1 entry)
	// - L6 last 2: 000043 x1 -> 2,000 * 1 = 2,000 (1 entry)
	// - L6 first 3: 267000 x1 -> 4,000 * 1 = 4,000 (1 entry)
	// - L6 last 3: 000065 x1 -> 4,000 * 1 = 4,000 (1 entry)
	// - L6 first 3 AND last 2: 267043 x1 -> wins first 3 (4,000) and last 2 (2,000) (2 entries)
	// - N3 straight 3: 077 x1 -> 8,000 * 1 = 8,000 (1 entry)
	// - N3 shuffle 3: 707 x1 -> 3,000 * 1 = 3,000 (1 entry)
	// - N3 straight 2: 943 x1 -> 500 * 1 = 500 (1 entry)
	// Total winning entries: 14
	if len(winnings) != 14 {
		t.Errorf("expected 14 winning entries, got %d: %+v", len(winnings), winnings)
	}

	expectedAmounts := map[uuid.UUID][]int{
		tickets[0].ID:  {12000000},
		tickets[1].ID:  {200000},
		tickets[2].ID:  {80000},
		tickets[3].ID:  {40000},
		tickets[4].ID:  {20000},
		tickets[5].ID:  {100000},
		tickets[6].ID:  {2000},
		tickets[7].ID:  {4000},
		tickets[8].ID:  {4000},
		tickets[9].ID:  {4000, 2000},
		tickets[10].ID: {8000},
		tickets[11].ID: {3000},
		tickets[12].ID: {500},
	}

	actualAmounts := make(map[uuid.UUID][]int)
	for _, w := range winnings {
		actualAmounts[w.TicketID] = append(actualAmounts[w.TicketID], w.PrizeMoney)
	}

	for ticketID, expectedSlice := range expectedAmounts {
		actualSlice := actualAmounts[ticketID]
		if len(actualSlice) != len(expectedSlice) {
			t.Errorf("ticket %s: expected %d winning entries, got %d (actual: %v, expected: %v)", ticketID, len(expectedSlice), len(actualSlice), actualSlice, expectedSlice)
			continue
		}
		tempActual := append([]int(nil), actualSlice...)
		for _, exp := range expectedSlice {
			found := false
			for i, act := range tempActual {
				if act == exp {
					tempActual = append(tempActual[:i], tempActual[i+1:]...)
					found = true
					break
				}
			}
			if !found {
				t.Errorf("ticket %s: expected prize amount %d not found in actual winnings %v", ticketID, exp, actualSlice)
			}
		}
	}
}

func TestVerifyDrawResults_PendingState(t *testing.T) {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "postgres://postgres:yourpassword@localhost:5432/lotto_journal?sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("Skipping integration test; cannot connect to DB: %v", err)
		return
	}

	tx := db.Begin()
	defer tx.Rollback()

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
	drawRepo := repository.NewDrawRepository(tx)
	drawResultRepo := repository.NewDrawResultRepository(tx)
	ticketRepo := repository.NewTicketRepository(tx)
	winningRepo := repository.NewUserWinningRepository(tx)
	mockNotifSvc := mocks.NewMockNotificationServiceInterface(t)
	resultSvc := service.NewResultService(tx, lotteryClient, drawRepo, drawResultRepo, ticketRepo, winningRepo, mockNotifSvc)

	drawDate, _ := time.Parse("2006-01-02", "2026-07-01")
	err = resultSvc.VerifyDrawResults(context.Background(), drawDate)

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
			got := service.ParsePrizeAmount(tc.input)
			if got != tc.expected {
				t.Errorf("ParsePrizeAmount(%q) = %d; expected %d", tc.input, got, tc.expected)
			}
		})
	}
}

func TestVerifyLatestDrawResults(t *testing.T) {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "postgres://postgres:yourpassword@localhost:5432/lotto_journal?sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("Skipping integration test; cannot connect to DB: %v", err)
		return
	}

	tx := db.Begin()
	defer tx.Rollback()

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"statusMessage": "getLatestLottery - Success",
			"statusCode": 200,
			"status": true,
			"response": {
				"date": "2026-06-16",
				"data": {
					"first": { "price": "6000000.00", "number": [] }
				}
			}
		}`))
	}))
	defer mockServer.Close()

	lotteryClient := client.NewLotteryClient(mockServer.URL)
	drawRepo := repository.NewDrawRepository(tx)
	drawResultRepo := repository.NewDrawResultRepository(tx)
	ticketRepo := repository.NewTicketRepository(tx)
	winningRepo := repository.NewUserWinningRepository(tx)
	mockNotifSvc := mocks.NewMockNotificationServiceInterface(t)
	mockNotifSvc.On("SendDrawNotifications", mock.Anything, mock.Anything, "2026-06-16").Return(nil)

	resultSvc := service.NewResultService(tx, lotteryClient, drawRepo, drawResultRepo, ticketRepo, winningRepo, mockNotifSvc)

	err = resultSvc.VerifyLatestDrawResults(context.Background())
	if err != nil {
		t.Fatalf("expected VerifyLatestDrawResults to succeed, got: %v", err)
	}

	var draw models.Draw
	if err := tx.First(&draw, "draw_date = ?", "2026-06-16").Error; err != nil {
		t.Fatalf("failed to find draw record: %v", err)
	}
	if !draw.IsVerified {
		t.Errorf("expected draw to be marked verified, but got false")
	}
}
