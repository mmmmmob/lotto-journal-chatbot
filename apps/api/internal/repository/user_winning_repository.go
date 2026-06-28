package repository

import (
	"fmt"

	"lotto-journal/api/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	TableUserWinning           = "user_winnings"
	UserWinningColID           = "id"
	UserWinningColTicketID     = "ticket_id"
	UserWinningColDrawResultID = "draw_result_id"
	UserWinningColPrizeMoney   = "prize_money"
	UserWinningColCreatedAt    = "created_at"
)

type DrawWinningDetail struct {
	TicketID      uuid.UUID
	PrizeMoney    int
	PrizeCategory string
}

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

func (r *UserWinningRepository) FindDrawWinnings(drawID uuid.UUID) ([]DrawWinningDetail, error) {
	var winnings []DrawWinningDetail

	selectFields := fmt.Sprintf("%s.%s, %s.%s, %s.%s",
		TableUserWinning, UserWinningColTicketID,
		TableUserWinning, UserWinningColPrizeMoney,
		TableDrawResults, DrawResultColPrizeCategory,
	)

	joinClause := fmt.Sprintf("JOIN %s ON %s.%s = %s.%s",
		TableDrawResults,
		TableUserWinning, UserWinningColDrawResultID,
		TableDrawResults, DrawResultColID,
	)

	whereClause := fmt.Sprintf("%s.%s = ?",
		TableDrawResults, DrawResultColDrawID,
	)

	err := r.db.Table(TableUserWinning).
		Select(selectFields).
		Joins(joinClause).
		Where(whereClause, drawID).
		Scan(&winnings).Error

	return winnings, err
}
