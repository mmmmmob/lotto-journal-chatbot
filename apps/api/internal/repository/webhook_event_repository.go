package repository

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"lotto-journal/api/internal/models"
)

type WebhookEventRepository struct {
	db *gorm.DB
}

func NewWebhookEventRepository(db *gorm.DB) *WebhookEventRepository {
	return &WebhookEventRepository{db: db}
}

// MarkProcessed attempts to record the webhook event ID.
// Returns (isNew, error):
//   - isNew == true  → event was freshly recorded; caller should process it
//   - isNew == false → duplicate event ID; caller should skip processing
//
// Uses INSERT ... ON CONFLICT DO NOTHING so the check is atomic with the insert.
func (r *WebhookEventRepository) MarkProcessed(eventID string) (bool, error) {
	record := &models.WebhookEvent{
		EventID:     eventID,
		ProcessedAt: time.Now(),
	}
	result := r.db.
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(record)
	if result.Error != nil {
		return false, result.Error
	}
	// RowsAffected == 0 means the PK already existed (conflict suppressed).
	return result.RowsAffected > 0, nil
}
