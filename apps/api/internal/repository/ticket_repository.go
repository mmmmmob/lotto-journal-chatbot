package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"lotto-journal/api/internal/models"
)

type TicketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) *TicketRepository {
	return &TicketRepository{db: db}
}

// Create inserts a new ticket record.
func (r *TicketRepository) Create(ticket *models.Ticket) error {
	return r.db.Create(ticket).Error
}

// List all tickets user hold for each draw
func (r *TicketRepository) List(drawID uuid.UUID, userID uuid.UUID) ([]*models.Ticket, error) {
	var tickets []*models.Ticket
	result := r.db.Where("owner_id = ? AND draw_id = ?", userID, drawID).Find(&tickets)
	return tickets, result.Error
}
