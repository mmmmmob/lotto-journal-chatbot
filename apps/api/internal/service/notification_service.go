package service

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"gorm.io/gorm"

	"lotto-journal/api/internal/localization"
	"lotto-journal/api/internal/models"
	"lotto-journal/api/internal/repository"
)

type NotificationService struct {
	db             *gorm.DB
	bot            *messaging_api.MessagingApiAPI
	ticketRepo     repository.TicketRepositoryInterface
	winningRepo    repository.UserWinningRepositoryInterface
	drawResultRepo repository.DrawResultRepositoryInterface
}

func NewNotificationService(
	db *gorm.DB,
	bot *messaging_api.MessagingApiAPI,
	ticketRepo repository.TicketRepositoryInterface,
	winningRepo repository.UserWinningRepositoryInterface,
	drawResultRepo repository.DrawResultRepositoryInterface,
) *NotificationService {
	return &NotificationService{
		db:             db,
		bot:            bot,
		ticketRepo:     ticketRepo,
		winningRepo:    winningRepo,
		drawResultRepo: drawResultRepo,
	}
}

// SendDrawNotifications queries active tickets and winnings for the draw,
// groups them, and pushes formatted messages to each user.
func (s *NotificationService) SendDrawNotifications(ctx context.Context, drawID uuid.UUID, drawDateStr string) error {
	// 1. Fetch tickets with owners
	tickets, err := s.ticketRepo.FindDrawTicketsWithOwners(drawID)
	if err != nil {
		return fmt.Errorf("fetch tickets with owners: %w", err)
	}
	if len(tickets) == 0 {
		log.Printf("[notification_service] No active tickets found for draw %s. Skipping notifications.", drawDateStr)
		return nil
	}

	// 2. Fetch winnings
	winnings, err := s.winningRepo.FindDrawWinnings(drawID)
	if err != nil {
		return fmt.Errorf("fetch draw winnings: %w", err)
	}

	// 3. Fetch n3_special result (if any)
	var specialResult *models.DrawResult
	hasSpecial := false
	sResult, err := s.drawResultRepo.FindSpecialResultByDrawID(drawID)
	if err == nil {
		specialResult = sResult
		hasSpecial = true
	} else if err != gorm.ErrRecordNotFound {
		log.Printf("[notification_service] Error querying n3_special result: %v", err)
	}

	// 4. Group tickets by owner ID
	type UserInfo struct {
		ID         uuid.UUID
		LineUserID string
		Language   string
	}
	userTickets := make(map[UserInfo][]repository.DrawTicketWithOwner)
	for _, t := range tickets {
		u := UserInfo{ID: t.OwnerID, LineUserID: t.LineUserID, Language: t.Language}
		userTickets[u] = append(userTickets[u], t)
	}

	// 5. Group winnings by ticket ID
	ticketWinnings := make(map[uuid.UUID][]repository.DrawWinningDetail)
	for _, w := range winnings {
		ticketWinnings[w.TicketID] = append(ticketWinnings[w.TicketID], w)
	}

	// 6. Format draw date for display (e.g. 2026-07-01 -> 01/07/2026)
	formattedDrawDate := drawDateStr
	if t, err := time.Parse("2006-01-02", drawDateStr); err == nil {
		formattedDrawDate = t.Format("02/01/2006")
	}

	// 7. Iterate through users to build and send notifications
	for u, tkts := range userTickets {
		var winDetails []string
		var loseDetails []string
		totalPrize := 0
		hasN3Straight := false

		dict := localization.GetDictionary(u.Language)

		for _, t := range tkts {
			wins, isWinner := ticketWinnings[t.ID]
			if isWinner {
				for _, w := range wins {
					totalPrize += w.PrizeMoney
					detail := fmt.Sprintf(dict.DrawWinDetail,
						t.Number, t.Type, localization.GetPrizeName(w.PrizeCategory, u.Language), t.Quantity, formatMoney(w.PrizeMoney))
					winDetails = append(winDetails, detail)

					if w.PrizeCategory == "n3_straight_three" {
						hasN3Straight = true
					}
				}
			} else {
				detail := fmt.Sprintf(dict.DrawLoseDetail, t.Number, t.Type, t.Quantity)
				loseDetails = append(loseDetails, detail)
			}
		}

		var messageText string
		if len(winDetails) > 0 {
			// User is a winner!
			winDetailsText := strings.Join(winDetails, "\n")
			messageText = fmt.Sprintf(dict.DrawWinMessage, formattedDrawDate, winDetailsText, formatMoney(totalPrize))

			if hasN3Straight && hasSpecial {
				messageText += fmt.Sprintf(dict.DrawJackpotInfo, specialResult.WinningNumber, formatMoney(specialResult.PrizeAmount))
			}

			messageText += dict.DrawFootnote
		} else if len(loseDetails) > 0 {
			// User didn't win anything
			loseDetailsText := strings.Join(loseDetails, "\n")
			messageText = fmt.Sprintf(dict.DrawLoseMessage, formattedDrawDate, loseDetailsText)
		}

		if messageText != "" {
			err := s.pushMessageWithRetry(ctx, u.LineUserID, messageText)
			status := models.NotifStatusSuccess
			var errStr *string
			if err != nil {
				status = models.NotifStatusFailed
				msg := err.Error()
				errStr = &msg
				log.Printf("[notification_service] Failed to send push message to %s: %v", u.LineUserID, err)
			}

			// Log this notification in database
			if logErr := s.LogNotification(u.ID, u.LineUserID, models.NotifTypeDrawResult, &drawID, status, errStr); logErr != nil {
				log.Printf("[notification_service] Failed to write notification log: %v", logErr)
			}
		}
	}

	return nil
}

// LogNotification writes a record to the notification_logs table.
func (s *NotificationService) LogNotification(userID uuid.UUID, lineUserID string, notifType models.NotificationType, drawID *uuid.UUID, status models.NotificationStatus, errStr *string) error {
	logRecord := &models.NotificationLog{
		UserID:           userID,
		LineUserID:       lineUserID,
		NotificationType: notifType,
		DrawID:           drawID,
		Status:           status,
		ErrorMessage:     errStr,
	}
	return s.db.Create(logRecord).Error
}

func (s *NotificationService) pushMessageWithRetry(ctx context.Context, to, text string) error {
	req := &messaging_api.PushMessageRequest{
		To: to,
		Messages: []messaging_api.MessageInterface{
			&messaging_api.TextMessage{
				Text: text,
			},
		},
	}

	var lastErr error
	backoff := 100 * time.Millisecond

	for attempt := 1; attempt <= 3; attempt++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		_, err := s.bot.PushMessage(req, "")
		if err == nil {
			return nil
		}

		lastErr = err
		log.Printf("[notification_service] PushMessage attempt %d failed for %s: %v", attempt, to, err)

		if attempt < 3 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
				backoff *= 2
			}
		}
	}

	return lastErr
}


func formatMoney(amount int) string {
	str := strconv.Itoa(amount)
	var result []string
	n := len(str)
	for i, r := range str {
		if (n-i)%3 == 0 && i != 0 {
			result = append(result, ",")
		}
		result = append(result, string(r))
	}
	return strings.Join(result, "")
}
