package repository

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"lotto-journal/api/internal/models"
)

const (
	DrawColID         = "id"
	DrawColDrawDate   = "draw_date"
	DrawColIsVerified = "is_verified"
)

type DrawRepository struct {
	db *gorm.DB
}

func NewDrawRepository(db *gorm.DB) *DrawRepository {
	return &DrawRepository{db: db}
}

// FindNextDraw returns the first unverified draw scheduled on or after the given date.
func (r *DrawRepository) FindNextDraw(fromDate time.Time) (*models.Draw, error) {
	var draw models.Draw
	dateStr := fromDate.Format("2006-01-02")
	result := r.db.Where(DrawColDrawDate+" >= ? AND "+DrawColIsVerified+" = false", dateStr).Order(DrawColDrawDate + " ASC").First(&draw)
	if result.Error != nil {
		return nil, result.Error
	}
	return &draw, nil
}

// FindByDate returns the draw for the given date, or nil + gorm.ErrRecordNotFound.
func (r *DrawRepository) FindByDate(date time.Time) (*models.Draw, error) {
	var draws []models.Draw
	result := r.db.Where(DrawColDrawDate+" = ?", date.Format("2006-01-02")).Limit(1).Find(&draws)
	if result.Error != nil {
		return nil, result.Error
	}
	if len(draws) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &draws[0], nil
}

// FindOrCreate returns the existing draw for the given date or creates a new one.
//
// This implementation is atomic at the SQL level via INSERT ... ON CONFLICT,
// avoiding the SELECT+INSERT race in FirstOrCreate.
func (r *DrawRepository) FindOrCreate(date time.Time) (*models.Draw, error) {
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	draw := models.Draw{DrawDate: dateOnly}

	result := r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: DrawColDrawDate}},
		// no-op update to force PostgreSQL RETURNING on conflict
		DoUpdates: clause.AssignmentColumns([]string{DrawColDrawDate}),
	}).Create(&draw)
	if result.Error != nil {
		return nil, result.Error
	}

	// Re-fetch the complete record. GORM's Create with OnConflict does not populate
	// unmodified fields (like IsVerified) in the Go struct when a conflict occurs.
	var fullDraw models.Draw
	if err := r.db.Where("id = ?", draw.ID).First(&fullDraw).Error; err != nil {
		return nil, err
	}

	return &fullDraw, nil
}

// FindLatestUnverified returns the most recent draw on or before the given date that is not yet verified.
func (r *DrawRepository) FindLatestUnverified(date time.Time) (*models.Draw, error) {
	var draw models.Draw
	dateStr := date.Format("2006-01-02")
	result := r.db.Where(DrawColDrawDate+" <= ? AND "+DrawColIsVerified+" = false", dateStr).Order(DrawColDrawDate + " DESC").First(&draw)
	if result.Error != nil {
		return nil, result.Error
	}
	return &draw, nil
}

// MarkVerifiedInTransaction marks a draw as verified inside a transaction.
func (r *DrawRepository) MarkVerifiedInTransaction(tx *gorm.DB, drawID uuid.UUID) error {
	return tx.Model(&models.Draw{}).Where("id = ?", drawID).Update(DrawColIsVerified, true).Error
}
