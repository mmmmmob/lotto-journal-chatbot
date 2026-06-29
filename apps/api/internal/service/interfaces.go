package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"lotto-journal/api/internal/models"
)

type UserServiceInterface interface {
	FindOrCreate(lineUserID string) (*models.User, bool, error)
	Deactivate(lineUserID string) error
	Reactivate(lineUserID string) error
	UpdateLanguage(lineUserID string, language string) error
}

type TicketServiceInterface interface {
	SubmitTickets(ownerID uuid.UUID, text string) ([]ParsedTicket, []string, uuid.UUID, error)
	ListTickets(ownerID uuid.UUID) ([]*models.Ticket, time.Time, error)
}

type NotificationServiceInterface interface {
	SendDrawNotifications(ctx context.Context, drawID uuid.UUID, drawDateStr string) error
	LogNotification(userID uuid.UUID, lineUserID string, notifType models.NotificationType, drawID *uuid.UUID, status models.NotificationStatus, errStr *string) error
}
