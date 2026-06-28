package service

import (
	"context"

	"github.com/google/uuid"
	"lotto-journal/api/internal/models"
)

type UserServiceInterface interface {
	FindOrCreate(lineUserID string) (*models.User, bool, error)
	Deactivate(lineUserID string) error
	Reactivate(lineUserID string) error
}

type TicketServiceInterface interface {
	SubmitTickets(ownerID uuid.UUID, text string) ([]ParsedTicket, []string, uuid.UUID, error)
	ListTickets(ownerID uuid.UUID) ([]*models.Ticket, error)
}

type NotificationServiceInterface interface {
	SendDrawNotifications(ctx context.Context, drawID uuid.UUID, drawDateStr string) error
	LogNotification(userID uuid.UUID, lineUserID string, notifType string, drawID *uuid.UUID, status string, errStr *string) error
}
