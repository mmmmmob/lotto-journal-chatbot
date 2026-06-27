package models

import (
	"time"

	"github.com/google/uuid"
)

type UserWinning struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	TicketID     uuid.UUID `gorm:"type:uuid;not null;index"                       json:"ticket_id"`
	DrawResultID uuid.UUID `gorm:"type:uuid;not null;index"                       json:"draw_result_id"`
	PrizeMoney   int       `gorm:"type:int4;not null"                             json:"prize_money"`
	UserID       uuid.UUID `gorm:"type:uuid;not null;index"                       json:"user_id"`
	CreatedAt    time.Time `gorm:"type:timestamp;autoCreateTime"                  json:"created_at"`
}
