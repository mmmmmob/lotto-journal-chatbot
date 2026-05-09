package service

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"lotto-journal/api/internal/models"
	"lotto-journal/api/internal/repository"
)

var (
	// spaceXRe normalises "digits <space> xN" → "digitsxN" (e.g. "123456 x2" → "123456x2").
	spaceXRe = regexp.MustCompile(`(?i)(\d+)\s+x(\d+)`)
	// numXRe matches a normalised quantity token (e.g. "123456x2").
	numXRe = regexp.MustCompile(`(?i)^(\d+)x(\d+)$`)
	// numOnlyRe matches a plain digit-only token.
	numOnlyRe = regexp.MustCompile(`^\d+$`)
)

// ParsedTicket holds one parsed ticket entry from a LINE message.
type ParsedTicket struct {
	Number   string
	Quantity int
	Type     string // "L6" (6-digit) or "N3" (3-digit)
}

type TicketService struct {
	ticketRepo *repository.TicketRepository
	drawRepo   *repository.DrawRepository
}

func NewTicketService(
	ticketRepo *repository.TicketRepository,
	drawRepo *repository.DrawRepository,
) *TicketService {
	return &TicketService{
		ticketRepo: ticketRepo,
		drawRepo:   drawRepo,
	}
}

// ParseTicketInput parses a LINE message text and extracts lottery ticket entries.
//
// Accepted formats (spaces and commas are both valid separators):
//
//	123456          → L6 ticket, quantity 1
//	456             → N3 ticket, quantity 1
//	123456x2        → L6 ticket, quantity 2
//	123456 x2       → same (space before x is normalised)
//	123456, 654321  → two L6 tickets
//
// Returns (validTickets, invalidTokens). Non-digit text tokens are skipped
// silently so that Thai text in the same message does not trigger an error.
func ParseTicketInput(text string) ([]ParsedTicket, []string) {
	// Normalise: commas → spaces
	text = strings.ReplaceAll(text, ",", " ")

	// Merge "digits <space> xN" into "digitsxN" so a space before 'x' is allowed.
	text = spaceXRe.ReplaceAllString(text, "$1x$2")

	tokens := strings.Fields(text)

	var tickets []ParsedTicket
	var invalid []string

	for _, token := range tokens {
		var numStr string
		qty := 1

		if m := numXRe.FindStringSubmatch(token); m != nil {
			numStr = m[1]
			var err error
			qty, err = strconv.Atoi(m[2])
			if err != nil || qty <= 0 {
				invalid = append(invalid, token)
				continue
			}
		} else if numOnlyRe.MatchString(token) {
			numStr = token
		} else {
			// Non-numeric token (e.g. Thai text) — skip silently.
			continue
		}

		switch len(numStr) {
		case 6:
			tickets = append(tickets, ParsedTicket{Number: numStr, Quantity: qty, Type: "L6"})
		case 3:
			tickets = append(tickets, ParsedTicket{Number: numStr, Quantity: qty, Type: "N3"})
		default:
			invalid = append(invalid, numStr)
		}
	}

	return tickets, invalid
}

// SubmitTickets parses the message text, resolves (or creates) the upcoming draw,
// persists each valid ticket, and returns the saved tickets and any invalid tokens.
func (s *TicketService) SubmitTickets(ownerID uuid.UUID, text string) ([]ParsedTicket, []string, error) {
	parsed, invalid := ParseTicketInput(text)
	if len(parsed) == 0 {
		return nil, invalid, nil
	}

	draw, err := s.drawRepo.FindOrCreate(NextDrawDate(time.Now()))
	if err != nil {
		return nil, invalid, fmt.Errorf("resolve upcoming draw: %w", err)
	}

	for _, pt := range parsed {
		ticket := &models.Ticket{
			OwnerID:  ownerID,
			DrawID:   draw.ID,
			Type:     pt.Type,
			Number:   pt.Number,
			Quantity: pt.Quantity,
		}
		if err := s.ticketRepo.Create(ticket); err != nil {
			return nil, invalid, fmt.Errorf("save ticket %s: %w", pt.Number, err)
		}
	}

	return parsed, invalid, nil
}

// ListTickets find the nearest draw date from current time and call List method on ticketRepo to return list of tickets user holds
func (s *TicketService) ListTickets(ownerID uuid.UUID) ([]*models.Ticket, error) {
	draw, err := s.drawRepo.FindOrCreate(NextDrawDate(time.Now()))
	if err != nil {
		return nil, fmt.Errorf("resolve upcoming draw: %w", err)
	}
	tickets, err := s.ticketRepo.List(draw.ID, ownerID)
	return tickets, err
}
