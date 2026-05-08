package models

import (
	"time"

	"github.com/google/uuid"
)

type Draw struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	DrawDate   time.Time `gorm:"type:date;uniqueIndex;not null"                  json:"draw_date"`
	IsVerified bool      `gorm:"type:boolean;default:false"                     json:"is_verified"`
	CreatedAt  time.Time `gorm:"type:timestamp;autoCreateTime"                  json:"created_at"`
	UpdatedAt  time.Time `gorm:"type:timestamp;autoUpdateTime"                  json:"updated_at"`
}
