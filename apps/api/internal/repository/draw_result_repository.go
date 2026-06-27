package repository

import (
	"gorm.io/gorm"
	"lotto-journal/api/internal/models"
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
