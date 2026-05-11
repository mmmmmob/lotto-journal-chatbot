package repository

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"lotto-journal/api/internal/models"
)

type DrawRepository struct {
	db *gorm.DB
}

func NewDrawRepository(db *gorm.DB) *DrawRepository {
	return &DrawRepository{db: db}
}

// FindByDate returns the draw for the given date, or nil + gorm.ErrRecordNotFound.
func (r *DrawRepository) FindByDate(date time.Time) (*models.Draw, error) {
	var draw models.Draw
	result := r.db.Where("draw_date = ?", date.Format("2006-01-02")).First(&draw)
	if result.Error != nil {
		return nil, result.Error
	}
	return &draw, nil
}

// FindOrCreate returns the existing draw for the given date or creates a new one.
//
// This implementation is atomic at the SQL level via INSERT ... ON CONFLICT,
// avoiding the SELECT+INSERT race in FirstOrCreate.
func (r *DrawRepository) FindOrCreate(date time.Time) (*models.Draw, error) {
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	draw := models.Draw{DrawDate: dateOnly}

	result := r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "draw_date"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			// no-op update to force PostgreSQL RETURNING on conflict
			"draw_date": gorm.Expr("draws.draw_date"),
		}),
	}).Create(&draw)
	if result.Error != nil {
		return nil, result.Error
	}

	return &draw, nil
}
