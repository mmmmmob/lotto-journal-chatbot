package service

import (
	"github.com/google/uuid"
	"lotto-journal/api/internal/models"
)

type UserServiceInterface interface {
	FindOrCreate(lineUserID string) (*models.User, bool, error)
	Deactivate(lineUserID string) error
	Reactivate(lineUserID string) error
}

type TicketServiceInterface interface {
	SubmitTickets(ownerID uuid.UUID, text string) ([]ParsedTicket, []string, error)
	ListTickets(ownerID uuid.UUID) ([]*models.Ticket, error)
}
