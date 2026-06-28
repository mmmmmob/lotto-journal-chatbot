package service

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"lotto-journal/api/internal/client"
	"lotto-journal/api/internal/models"
	"lotto-journal/api/internal/repository"
)

type ResultService struct {
	db              *gorm.DB
	client          *client.LotteryClient
	drawRepo        repository.DrawRepositoryInterface
	drawResultRepo  repository.DrawResultRepositoryInterface
	ticketRepo      repository.TicketRepositoryInterface
	winningRepo     repository.UserWinningRepositoryInterface
	notificationSvc NotificationServiceInterface
}

func NewResultService(
	db *gorm.DB,
	client *client.LotteryClient,
	drawRepo repository.DrawRepositoryInterface,
	drawResultRepo repository.DrawResultRepositoryInterface,
	ticketRepo repository.TicketRepositoryInterface,
	winningRepo repository.UserWinningRepositoryInterface,
	notificationSvc NotificationServiceInterface,
) *ResultService {
	return &ResultService{
		db:              db,
		client:          client,
		drawRepo:        drawRepo,
		drawResultRepo:  drawResultRepo,
		ticketRepo:      ticketRepo,
		winningRepo:     winningRepo,
		notificationSvc: notificationSvc,
	}
}

// VerifyDrawResults pulls GLO results for the given date, saves them,
// matches against unchecked user tickets, and records any winning tickets.
func (s *ResultService) VerifyDrawResults(ctx context.Context, drawDate time.Time) error {
	// 1. Fetch latest GLO results
	latest, err := s.client.FetchLatestResult(ctx)
	if err != nil {
		return fmt.Errorf("fetch latest results: %w", err)
	}

	// 2. Validate that the returned draw date matches what the scheduler requested
	expectedDateStr := drawDate.Format("2006-01-02")
	if latest.Response.Date != expectedDateStr {
		return fmt.Errorf("results pending: latest GLO draw is %s, expected %s", latest.Response.Date, expectedDateStr)
	}

	// 3. Find or create the draw record
	draw, err := s.drawRepo.FindOrCreate(drawDate)
	if err != nil {
		return fmt.Errorf("resolve draw record: %w", err)
	}

	// If already verified, nothing to do
	if draw.IsVerified {
		log.Printf("[result_service] draw %s is already verified. Skipping checking.", expectedDateStr)
		return nil
	}

	// 4. Parse winning results from response
	var drawResults []*models.DrawResult
	drawResults = append(drawResults, parsePrize(draw.ID, "l6_first", latest.Response.Data.First)...)
	drawResults = append(drawResults, parsePrize(draw.ID, "l6_second", latest.Response.Data.Second)...)
	drawResults = append(drawResults, parsePrize(draw.ID, "l6_third", latest.Response.Data.Third)...)
	drawResults = append(drawResults, parsePrize(draw.ID, "l6_fourth", latest.Response.Data.Fourth)...)
	drawResults = append(drawResults, parsePrize(draw.ID, "l6_fifth", latest.Response.Data.Fifth)...)
	drawResults = append(drawResults, parsePrize(draw.ID, "l6_last2", latest.Response.Data.Last2)...)
	drawResults = append(drawResults, parsePrize(draw.ID, "l6_last3f", latest.Response.Data.Last3f)...)
	drawResults = append(drawResults, parsePrize(draw.ID, "l6_last3b", latest.Response.Data.Last3b)...)
	drawResults = append(drawResults, parsePrize(draw.ID, "l6_near_first", latest.Response.Data.Near1)...)

	drawResults = append(drawResults, parsePrize(draw.ID, "n3_straight_three", latest.Response.N3.Straight3)...)
	drawResults = append(drawResults, parsePrize(draw.ID, "n3_shuffle", latest.Response.N3.Shuffle3)...)
	drawResults = append(drawResults, parsePrize(draw.ID, "n3_straight_two", latest.Response.N3.Straight2)...)
	drawResults = append(drawResults, parsePrize(draw.ID, "n3_special", latest.Response.N3.Special)...)

	// Map results by category for fast lookup
	prizesByCategory := make(map[string][]*models.DrawResult)
	for _, r := range drawResults {
		prizesByCategory[r.PrizeCategory] = append(prizesByCategory[r.PrizeCategory], r)
	}

	// 5. Run checking in a database transaction
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// Clean up any existing winnings, results, or checked statuses for this draw to ensure absolute idempotency
		if err := s.winningRepo.DeleteByDrawIDInTransaction(tx, draw.ID); err != nil {
			return fmt.Errorf("clear existing user winnings: %w", err)
		}

		if err := s.drawResultRepo.DeleteByDrawIDInTransaction(tx, draw.ID); err != nil {
			return fmt.Errorf("clear existing draw results: %w", err)
		}

		if err := s.ticketRepo.ResetCheckedStatusByDrawIDInTransaction(tx, draw.ID); err != nil {
			return fmt.Errorf("reset tickets checked status: %w", err)
		}

		// Save all parsed draw results (handles GORM UUID generation in-place)
		if err := s.drawResultRepo.CreateInBatchesInTransaction(tx, drawResults); err != nil {
			return fmt.Errorf("bulk insert draw results: %w", err)
		}

		// Query unchecked tickets
		uncheckedTickets, err := s.ticketRepo.FindUncheckedInTransaction(tx, draw.ID)
		if err != nil {
			return fmt.Errorf("retrieve unchecked tickets: %w", err)
		}

		var winnings []*models.UserWinning
		var processedTicketIDs []uuid.UUID

		for _, ticket := range uncheckedTickets {
			processedTicketIDs = append(processedTicketIDs, ticket.ID)

			if ticket.Type == "L6" {
				if len(ticket.Number) != 6 {
					continue
				}

				// Check 6-digit category matches
				for _, category := range []string{"l6_first", "l6_second", "l6_third", "l6_fourth", "l6_fifth", "l6_near_first"} {
					for _, prize := range prizesByCategory[category] {
						if ticket.Number == prize.WinningNumber {
							winnings = append(winnings, &models.UserWinning{
								TicketID:     ticket.ID,
								DrawResultID: prize.ID,
								PrizeMoney:   prize.PrizeAmount * ticket.Quantity,
								UserID:       ticket.OwnerID,
							})
						}
					}
				}

				// Check last 2 digits
				for _, prize := range prizesByCategory["l6_last2"] {
					if ticket.Number[4:] == prize.WinningNumber {
						winnings = append(winnings, &models.UserWinning{
							TicketID:     ticket.ID,
							DrawResultID: prize.ID,
							PrizeMoney:   prize.PrizeAmount * ticket.Quantity,
							UserID:       ticket.OwnerID,
						})
					}
				}

				// Check first 3 digits
				for _, prize := range prizesByCategory["l6_last3f"] {
					if ticket.Number[:3] == prize.WinningNumber {
						winnings = append(winnings, &models.UserWinning{
							TicketID:     ticket.ID,
							DrawResultID: prize.ID,
							PrizeMoney:   prize.PrizeAmount * ticket.Quantity,
							UserID:       ticket.OwnerID,
						})
					}
				}

				// Check last 3 digits
				for _, prize := range prizesByCategory["l6_last3b"] {
					if ticket.Number[3:] == prize.WinningNumber {
						winnings = append(winnings, &models.UserWinning{
							TicketID:     ticket.ID,
							DrawResultID: prize.ID,
							PrizeMoney:   prize.PrizeAmount * ticket.Quantity,
							UserID:       ticket.OwnerID,
						})
					}
				}

			} else if ticket.Type == "N3" {
				if len(ticket.Number) != 3 {
					continue
				}

				// Check straight 3
				for _, prize := range prizesByCategory["n3_straight_three"] {
					if ticket.Number == prize.WinningNumber {
						winnings = append(winnings, &models.UserWinning{
							TicketID:     ticket.ID,
							DrawResultID: prize.ID,
							PrizeMoney:   prize.PrizeAmount * ticket.Quantity,
							UserID:       ticket.OwnerID,
						})
					}
				}

				// Check shuffle 3
				for _, prize := range prizesByCategory["n3_shuffle"] {
					if ticket.Number == prize.WinningNumber {
						winnings = append(winnings, &models.UserWinning{
							TicketID:     ticket.ID,
							DrawResultID: prize.ID,
							PrizeMoney:   prize.PrizeAmount * ticket.Quantity,
							UserID:       ticket.OwnerID,
						})
					}
				}

				// Check straight 2
				for _, prize := range prizesByCategory["n3_straight_two"] {
					if ticket.Number[1:] == prize.WinningNumber {
						winnings = append(winnings, &models.UserWinning{
							TicketID:     ticket.ID,
							DrawResultID: prize.ID,
							PrizeMoney:   prize.PrizeAmount * ticket.Quantity,
							UserID:       ticket.OwnerID,
						})
					}
				}
			}
		}

		// Save winnings (if any)
		if len(winnings) > 0 {
			if err := s.winningRepo.CreateInBatchesInTransaction(tx, winnings); err != nil {
				return fmt.Errorf("bulk insert user winnings: %w", err)
			}
		}

		// Mark processed tickets as checked
		if len(processedTicketIDs) > 0 {
			if err := s.ticketRepo.MarkCheckedInTransaction(tx, processedTicketIDs); err != nil {
				return fmt.Errorf("mark tickets as checked: %w", err)
			}
		}

		// Mark draw verified
		if err := s.drawRepo.MarkVerifiedInTransaction(tx, draw.ID); err != nil {
			return fmt.Errorf("mark draw verified: %w", err)
		}

		log.Printf("[result_service] Draw check complete for %s. Checked %d tickets. Recorded %d winning entries.",
			expectedDateStr, len(uncheckedTickets), len(winnings))

		return nil
	})

	if err != nil {
		return fmt.Errorf("database transaction: %w", err)
	}

	// 6. Push notifications
	if err := s.notificationSvc.SendDrawNotifications(ctx, draw.ID, expectedDateStr); err != nil {
		log.Printf("[result_service] failed to send draw notifications: %v", err)
	}

	return nil
}

// parsePrize is a helper mapping GLO prize structures to DrawResult database models.
func parsePrize(drawID uuid.UUID, category string, prize client.GLOPrize) []*models.DrawResult {
	amount := ParsePrizeAmount(prize.Price)
	var results []*models.DrawResult
	for _, num := range prize.Number {
		if num.Value == "" {
			continue
		}
		results = append(results, &models.DrawResult{
			DrawID:        drawID,
			PrizeCategory: category,
			WinningNumber: num.Value,
			PrizeAmount:   amount,
		})
	}
	return results
}

func ParsePrizeAmount(price string) int {
	price = strings.ReplaceAll(price, ",", "")
	var f float64
	_, err := fmt.Sscanf(price, "%f", &f)
	if err != nil {
		// Fallback to basic string parsing if format differs
		val, err := strconv.ParseFloat(price, 64)
		if err != nil {
			return 0
		}
		f = val
	}
	return int(f)
}

// VerifyLatestDrawResults fetches the latest draw results from GLO API and verifies that specific draw.
func (s *ResultService) VerifyLatestDrawResults(ctx context.Context) error {
	latest, err := s.client.FetchLatestResult(ctx)
	if err != nil {
		return fmt.Errorf("fetch latest results: %w", err)
	}

	drawDate, err := time.Parse("2006-01-02", latest.Response.Date)
	if err != nil {
		return fmt.Errorf("parse GLO draw date %s: %w", latest.Response.Date, err)
	}

	return s.VerifyDrawResults(ctx, drawDate)
}
