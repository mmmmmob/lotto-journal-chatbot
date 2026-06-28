package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"lotto-journal/api/internal/models"
)

const (
	TableDrawResults           = "draw_results"
	DrawResultColID            = "id"
	DrawResultColDrawID        = "draw_id"
	DrawResultColPrizeCategory = "prize_category"
	DrawResultColWinningNumber = "winning_number"
	DrawResultColPrizeAmount   = "prize_amount"
)

type DrawResultRepository struct {
	db *gorm.DB
}

func NewDrawResultRepository(db *gorm.DB) *DrawResultRepository {
	return &DrawResultRepository{db: db}
}

func (r *DrawResultRepository) CreateInBatches(results []*models.DrawResult) error {
	if len(results) == 0 {
		return nil
	}
	return r.db.CreateInBatches(results, 100).Error
}

func (r *DrawResultRepository) CreateInBatchesInTransaction(tx *gorm.DB, results []*models.DrawResult) error {
	if len(results) == 0 {
		return nil
	}
	return tx.CreateInBatches(results, 100).Error
}

func (r *DrawResultRepository) DeleteByDrawIDInTransaction(tx *gorm.DB, drawID uuid.UUID) error {
	return tx.Where(DrawResultColDrawID+" = ?", drawID).Delete(&models.DrawResult{}).Error
}

func (r *DrawResultRepository) FindSpecialResultByDrawID(drawID uuid.UUID) (*models.DrawResult, error) {
	var result models.DrawResult
	err := r.db.Where(DrawResultColDrawID+" = ? AND "+DrawResultColPrizeCategory+" = ?", drawID, "n3_special").First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}
