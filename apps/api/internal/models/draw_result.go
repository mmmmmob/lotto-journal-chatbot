package models

import (
	"github.com/google/uuid"
)

type DrawResult struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	DrawID        uuid.UUID `gorm:"type:uuid;not null;index"                       json:"draw_id"`
	PrizeCategory string    `gorm:"type:prize_type;not null"                       json:"prize_category"`
	WinningNumber string    `gorm:"type:varchar(12);not null"                      json:"winning_number"`
	PrizeAmount   int       `gorm:"type:int4;not null"                             json:"prize_amount"`
}
