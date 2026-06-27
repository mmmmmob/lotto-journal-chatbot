package repository

import (
	"gorm.io/gorm"
	"lotto-journal/api/internal/models"
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
