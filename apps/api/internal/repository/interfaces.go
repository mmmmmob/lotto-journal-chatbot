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
	FindUncheckedInTransaction(tx *gorm.DB, drawID uuid.UUID) ([]*models.Ticket, error)
	MarkCheckedInTransaction(tx *gorm.DB, ticketIDs []uuid.UUID) error
	ResetCheckedStatusByDrawIDInTransaction(tx *gorm.DB, drawID uuid.UUID) error
}

type DrawRepositoryInterface interface {
	FindNextDraw(fromDate time.Time) (*models.Draw, error)
	FindByDate(date time.Time) (*models.Draw, error)
	FindOrCreate(date time.Time) (*models.Draw, error)
	FindLatestUnverified(date time.Time) (*models.Draw, error)
	MarkVerifiedInTransaction(tx *gorm.DB, drawID uuid.UUID) error
}

type DrawResultRepositoryInterface interface {
	CreateInBatches(results []*models.DrawResult) error
	CreateInBatchesInTransaction(tx *gorm.DB, results []*models.DrawResult) error
	DeleteByDrawIDInTransaction(tx *gorm.DB, drawID uuid.UUID) error
}

type UserWinningRepositoryInterface interface {
	CreateInBatches(winnings []*models.UserWinning) error
	CreateInBatchesInTransaction(tx *gorm.DB, winnings []*models.UserWinning) error
	DeleteByDrawIDInTransaction(tx *gorm.DB, drawID uuid.UUID) error
}
