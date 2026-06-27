package repository

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"lotto-journal/api/internal/models"
)

type UserRepositoryInterface interface {
	Create(user *models.User) error
	FindByLineUserID(lineUserID string) (*models.User, error)
	FindOrCreate(lineUserID string) (*models.User, bool, error)
	UpdateStatus(lineUserID string, status string) error
}

type TicketRepositoryInterface interface {
	Create(ticket *models.Ticket) error
	List(drawID uuid.UUID, userID uuid.UUID) ([]*models.Ticket, error)
	FindUnchecked(drawID uuid.UUID) ([]*models.Ticket, error)
	MarkCheckedInTransaction(tx *gorm.DB, ticketIDs []uuid.UUID) error
}

type DrawRepositoryInterface interface {
	FindNextDraw(fromDate time.Time) (*models.Draw, error)
	FindByDate(date time.Time) (*models.Draw, error)
	FindOrCreate(date time.Time) (*models.Draw, error)
}
