package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	LineUserID string    `gorm:"type:varchar;uniqueIndex;not null"               json:"line_user_id"`
	Status     string    `gorm:"type:account_status;default:'active'"            json:"status"`
	Language   string    `gorm:"type:varchar(10);default:'en';not null"          json:"language"`
	CreatedAt  time.Time `gorm:"type:timestamp;autoCreateTime"                   json:"created_at"`
	UpdatedAt  time.Time `gorm:"type:timestamp;autoUpdateTime"                   json:"updated_at"`
}
