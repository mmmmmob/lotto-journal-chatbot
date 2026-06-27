package repository

import (
	"fmt"

	"lotto-journal/api/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	TableUserWinning           = "user_winnings"
	TableTickets               = "tickets"
	UserWinningColID           = "id"
	UserWinningColTicketID     = "ticket_id"
	UserWinningColDrawResultID = "draw_result_id"
	UserWinningColPrizeMoney   = "prize_money"
	UserWinningColCreatedAt    = "created_at"
)

type UserWinningRepository struct {
	db *gorm.DB
}

func NewUserWinningRepository(db *gorm.DB) *UserWinningRepository {
	return &UserWinningRepository{db: db}
}

func (r *UserWinningRepository) CreateInBatches(winnings []*models.UserWinning) error {
	if len(winnings) == 0 {
		return nil
	}
	return r.db.CreateInBatches(winnings, 100).Error
}

func (r *UserWinningRepository) CreateInBatchesInTransaction(tx *gorm.DB, winnings []*models.UserWinning) error {
	if len(winnings) == 0 {
		return nil
	}
	return tx.CreateInBatches(winnings, 100).Error
}

func (r *UserWinningRepository) DeleteByDrawIDInTransaction(tx *gorm.DB, drawID uuid.UUID) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE %s IN (SELECT %s FROM %s WHERE %s = ?)", TableUserWinning, UserWinningColTicketID, TicketColID, TableTickets, TicketColDrawID)
	return tx.Exec(query, drawID).Error
}
