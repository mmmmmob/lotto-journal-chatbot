package models

import "time"

// WebhookEvent records every processed LINE webhookEventId.
// Used to detect and skip duplicate event deliveries from LINE.
type WebhookEvent struct {
	EventID     string    `gorm:"type:varchar;primaryKey;column:event_id"                   json:"event_id"`
	ProcessedAt time.Time `gorm:"type:timestamp;not null;default:now();column:processed_at" json:"processed_at"`
}
