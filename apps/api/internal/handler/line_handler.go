package handler

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"unicode"

	"github.com/gofiber/fiber/v3"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"

	"lotto-journal/api/internal/models"
	"lotto-journal/api/internal/repository"
	"lotto-journal/api/internal/service"
)

// LineHandler handles all LINE Messaging API webhook events.
type LineHandler struct {
	channelSecret string
	bot           *messaging_api.MessagingApiAPI
	userSvc       *service.UserService
	ticketSvc     *service.TicketService
	webhookRepo   *repository.WebhookEventRepository
}

func NewLineHandler(
	channelSecret string,
	bot *messaging_api.MessagingApiAPI,
	userSvc *service.UserService,
	ticketSvc *service.TicketService,
	webhookRepo *repository.WebhookEventRepository,
) *LineHandler {
	return &LineHandler{
		channelSecret: channelSecret,
		bot:           bot,
		userSvc:       userSvc,
		ticketSvc:     ticketSvc,
		webhookRepo:   webhookRepo,
	}
}

// Handle is the Fiber route handler for POST /webhook.
// It builds a synthetic *http.Request so the LINE SDK can verify the signature
// and parse the event payload (the SDK expects net/http; Fiber uses fasthttp).
func (h *LineHandler) Handle(c fiber.Ctx) error {
	req := &http.Request{
		Method: "POST",
		Header: http.Header{
			"X-Line-Signature": []string{c.Get("X-Line-Signature")},
			"Content-Type":     []string{"application/json"},
		},
		Body: io.NopCloser(bytes.NewReader(c.Body())),
	}

	cb, err := webhook.ParseRequest(h.channelSecret, req)
	if err != nil {
		if errors.Is(err, webhook.ErrInvalidSignature) {
			log.Println("[webhook] invalid signature")
			return c.Status(fiber.StatusBadRequest).SendString("invalid signature")
		}
		log.Printf("[webhook] parse error: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("failed to parse request")
	}

	for _, event := range cb.Events {
		h.dispatch(event)
	}

	return c.SendStatus(fiber.StatusOK)
}

// dispatch deduplicates and routes a single LINE event to the appropriate handler.
func (h *LineHandler) dispatch(event webhook.EventInterface) {
	eventID := webhookEventID(event)
	if eventID == "" {
		// webhookEventID only extracts IDs for event types we support.
		// Anything else (postback, join, leave, …) is unsupported and skipped.
		log.Printf("[webhook] unsupported event type %T — skipping", event)
		return
	}

	isNew, err := h.webhookRepo.MarkProcessed(eventID)
	if err != nil {
		log.Printf("[webhook] MarkProcessed %s: %v", eventID, err)
		return
	}
	if !isNew {
		log.Printf("[webhook] duplicate event %s — skipping", eventID)
		return
	}

	switch e := event.(type) {
	case webhook.FollowEvent:
		h.handleFollow(e)
	case webhook.UnfollowEvent:
		h.handleUnfollow(e)
	case webhook.MessageEvent:
		h.handleMessage(e)
		// No default: unsupported types never reach here — webhookEventID returns "" for them.
	}
}

// --- event handlers ---

func (h *LineHandler) handleFollow(e webhook.FollowEvent) {
	lineUserID := sourceUserID(e.Source)
	if lineUserID == "" {
		log.Println("[follow] no userId in source")
		return
	}

	_, isNew, err := h.userSvc.FindOrCreate(lineUserID)
	if err != nil {
		log.Printf("[follow] FindOrCreate %s: %v", lineUserID, err)
		return
	}

	if isNew {
		log.Printf("[follow] new user created: %s", lineUserID)
	} else {
		// User was previously inactive (unfollowed) — restore active status.
		if err := h.userSvc.Reactivate(lineUserID); err != nil {
			log.Printf("[follow] Reactivate %s: %v", lineUserID, err)
			// Non-fatal: still send the welcome reply even if reactivation fails.
		}
		log.Printf("[follow] existing user re-followed: %s", lineUserID)
	}

	displayName := h.getDisplayName(lineUserID)
	h.replyText(e.ReplyToken, buildWelcomeMessage(displayName))
}

func (h *LineHandler) handleUnfollow(e webhook.UnfollowEvent) {
	lineUserID := sourceUserID(e.Source)
	if lineUserID == "" {
		log.Println("[unfollow] no userId in source")
		return
	}

	if err := h.userSvc.Deactivate(lineUserID); err != nil {
		log.Printf("[unfollow] Deactivate %s: %v", lineUserID, err)
		return
	}
	log.Printf("[unfollow] user marked inactive: %s", lineUserID)
	// No reply — LINE does not allow replying to unfollow events.
}

func (h *LineHandler) handleMessage(e webhook.MessageEvent) {
	// Only handle plain text messages; ignore stickers, images, etc.
	textMsg, ok := e.Message.(webhook.TextMessageContent)
	if !ok {
		log.Printf("[message] non-text content %T — ignoring", e.Message)
		return
	}

	lineUserID := sourceUserID(e.Source)
	if lineUserID == "" {
		log.Println("[message] no userId in source")
		return
	}

	// Best-effort loading indicator to show user we're processing the request.
	// Run asynchronously so a slow LINE API call does not add webhook latency.
	go h.showLoading(lineUserID, 5)

	// Ensure the user record exists (edge case: message before follow event).
	user, _, err := h.userSvc.FindOrCreate(lineUserID)
	if err != nil {
		log.Printf("[message] FindOrCreate %s: %v", lineUserID, err)
		return
	}

	if isTicketListCmd(textMsg.Text) {
		userTickets, err := h.ticketSvc.ListTickets(user.ID)
		if err != nil {
			log.Printf("[message] error retrieving %s tickets: %v", user.ID, err)
			h.replyText(e.ReplyToken, "เกิดข้อผิดพลาด กรุณาลองใหม่อีกครั้ง 🙏")
			return
		}
		h.replyText(e.ReplyToken, buildTicketListReply(userTickets))
	} else {
		saved, invalid, err := h.ticketSvc.SubmitTickets(user.ID, textMsg.Text)
		if err != nil {
			log.Printf("[message] SubmitTickets for %s: %v", lineUserID, err)
			h.replyText(e.ReplyToken, "เกิดข้อผิดพลาด กรุณาลองใหม่อีกครั้ง 🙏")
			return
		}
		h.replyText(e.ReplyToken, buildReply(saved, invalid))
	}
}

// --- helpers ---

func (h *LineHandler) replyText(replyToken, text string) {
	if _, err := h.bot.ReplyMessage(&messaging_api.ReplyMessageRequest{
		ReplyToken: replyToken,
		Messages: []messaging_api.MessageInterface{
			messaging_api.TextMessage{Text: text},
		},
	}); err != nil {
		log.Printf("[reply] error: %v", err)
	}
}

func (h *LineHandler) getDisplayName(lineUserID string) string {
	profile, err := h.bot.GetProfile(lineUserID)
	if err != nil {
		log.Printf("[profile] GetProfile %s: %v", lineUserID, err)
		return ""
	}
	return strings.TrimSpace(profile.DisplayName)
}

func buildWelcomeMessage(displayName string) string {
	greeting := "👋 สวัสดี!"
	if displayName != "" {
		greeting = fmt.Sprintf("👋 สวัสดีคุณ %s!", displayName)
	}

	return greeting + "\n\n" +
		"🎟️ ยินดีต้อนรับสู่ Lotto Journal!\n" +
		"ตัวอย่าง: 123456 หรือ 456\n" +
		"ส่งหลายเลขได้ เช่น 123456 789012\n" +
		"ระบุจำนวนตั๋วด้วย x เช่น 123456x2\n\n" +
		"📝 หากต้องการดูสลากที่บันทึกไว้ พิมพ์ 'โพย'"
}

func (h *LineHandler) showLoading(chatID string, loadingSeconds int32) {
	if chatID == "" {
		return
	}
	if loadingSeconds < 5 {
		loadingSeconds = 5
	}
	// LINE requires loadingSeconds to be in 5-second increments.
	if rem := loadingSeconds % 5; rem != 0 {
		loadingSeconds += 5 - rem
	}
	if loadingSeconds > 60 {
		loadingSeconds = 60
	}

	if _, err := h.bot.ShowLoadingAnimation(&messaging_api.ShowLoadingAnimationRequest{
		ChatId:         chatID,
		LoadingSeconds: loadingSeconds,
	}); err != nil {
		log.Printf("[loading] error: %v", err)
	}
}

// buildReply constructs the confirmation (or error) text for a ticket submission.
func buildReply(saved []service.ParsedTicket, invalid []string) string {
	if len(saved) == 0 && len(invalid) == 0 {
		// No digit tokens found at all — unrecognised message.
		return "🎟️ ส่งเลขสลากของคุณมาได้เลย\n\n" +
			"ตัวอย่าง: 123456 หรือ 456\n" +
			"ส่งหลายเลขได้ เช่น 123456 789012\n" +
			"ระบุจำนวนตั๋วด้วย x เช่น 123456x2\n\n" +
			"📝 หากต้องการดูสลากที่บันทึกไว้ พิมพ์ 'โพย'"
	}

	if len(saved) == 0 {
		return fmt.Sprintf(
			"ไม่พบเลขสลากที่ถูกต้อง ❌\nกรุณาส่งเลข 3 หรือ 6 หลักเท่านั้น\nเลขที่ไม่ถูกต้อง: %s",
			strings.Join(invalid, ", "),
		)
	}

	lines := []string{"บันทึกสลากเรียบร้อย ✅"}
	for _, t := range saved {
		if t.Quantity > 1 {
			lines = append(lines, fmt.Sprintf("  • %s x%d (%s)", t.Number, t.Quantity, t.Type))
		} else {
			lines = append(lines, fmt.Sprintf("  • %s (%s)", t.Number, t.Type))
		}
	}
	if len(invalid) > 0 {
		lines = append(lines, fmt.Sprintf("\nเลขที่ไม่ถูกต้อง (ข้ามไป): %s", strings.Join(invalid, ", ")))
	}

	return strings.Join(lines, "\n")
}

// sourceUserID extracts the LINE userId from a SourceInterface.
// Returns "" if the source is nil or is not a user source (e.g. group/room).
func sourceUserID(src webhook.SourceInterface) string {
	if src == nil {
		return ""
	}
	if us, ok := src.(webhook.UserSource); ok {
		return us.UserId
	}
	return ""
}

// Return if message sent is command for "List all tickets" or not.
// Accepts extra spaces (including internal/Unicode spaces), e.g. "โ พย".
func isTicketListCmd(text string) bool {
	normalized := strings.Map(func(r rune) rune {
		switch {
		case unicode.IsSpace(r):
			return -1
		case r == '\u200B' || r == '\u200C' || r == '\u200D' || r == '\uFEFF':
			return -1
		default:
			return r
		}
	}, strings.TrimSpace(text))

	return normalized == "โพย"
}

func buildTicketListReply(tickets []*models.Ticket) string {
	if len(tickets) == 0 {
		return "คุณยังไม่ได้บันทึกสลากในงวดนี้"
	}

	lines := []string{"สลากที่คุณบันทึกไว้ในงวดนี้ 📝"}

	for _, t := range tickets {
		if t.Quantity > 1 {
			lines = append(lines, fmt.Sprintf("  • %s x%d (%s)", t.Number, t.Quantity, t.Type))
		} else {
			lines = append(lines, fmt.Sprintf("  • %s (%s)", t.Number, t.Type))
		}
	}
	return strings.Join(lines, "\n")
}

// webhookEventID extracts the webhookEventId from the event types we support.
// Returns "" for all other types — dispatch will log and skip those.
func webhookEventID(event webhook.EventInterface) string {
	switch e := event.(type) {
	case webhook.FollowEvent:
		return e.WebhookEventId
	case webhook.UnfollowEvent:
		return e.WebhookEventId
	case webhook.MessageEvent:
		return e.WebhookEventId
	default:
		return ""
	}
}
