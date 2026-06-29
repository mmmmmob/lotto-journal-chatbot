package models

import (
	"time"

	"github.com/google/uuid"
)

type NotificationType string

const (
	NotifTypeWelcome         NotificationType = "welcome"
	NotifTypeTicketSubmitted NotificationType = "ticket_submitted"
	NotifTypeTicketList      NotificationType = "ticket_list"
	NotifTypeDrawResult      NotificationType = "draw_result"
	NotifTypeLanguageChanged NotificationType = "language_changed"
	NotifTypeHelpAdd         NotificationType = "help_add"
	NotifTypeHelpNotify      NotificationType = "help_notify"
)

type NotificationStatus string

const (
	NotifStatusSuccess NotificationStatus = "success"
	NotifStatusFailed  NotificationStatus = "failed"
)

type NotificationLog struct {
	ID               uuid.UUID          `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID           uuid.UUID          `gorm:"type:uuid;not null;index"                       json:"user_id"`
	LineUserID       string             `gorm:"type:varchar;not null"                          json:"line_user_id"`
	NotificationType NotificationType   `gorm:"type:notification_type;not null"                json:"notification_type"`
	DrawID           *uuid.UUID         `gorm:"type:uuid;index"                                json:"draw_id,omitempty"`
	Status           NotificationStatus `gorm:"type:notification_status;not null"              json:"status"`
	ErrorMessage     *string            `gorm:"type:text"                                      json:"error_message,omitempty"`
	CreatedAt        time.Time          `gorm:"type:timestamp;autoCreateTime"                  json:"created_at"`
}
