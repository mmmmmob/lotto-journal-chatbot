package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"lotto-journal/api/internal/models"
)

const (
	TicketColID        = "id"
	TicketColOwnerID   = "owner_id"
	TicketColDrawID    = "draw_id"
	TicketColType      = "type"
	TicketColNumber    = "number"
	TicketColQuantity  = "quantity"
	TicketColIsChecked = "is_checked"
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
	result := r.db.Where(TicketColOwnerID+" = ? AND "+TicketColDrawID+" = ?", userID, drawID).Find(&tickets)
	return tickets, result.Error
}

// FindUnchecked retrieves all tickets for a draw that have not been checked yet.
func (r *TicketRepository) FindUnchecked(drawID uuid.UUID) ([]*models.Ticket, error) {
	var tickets []*models.Ticket
	result := r.db.Where(TicketColDrawID+" = ? AND "+TicketColIsChecked+" = false", drawID).Find(&tickets)
	return tickets, result.Error
}

// FindUncheckedInTransaction retrieves unchecked tickets using a transaction object.
func (r *TicketRepository) FindUncheckedInTransaction(tx *gorm.DB, drawID uuid.UUID) ([]*models.Ticket, error) {
	var tickets []*models.Ticket
	result := tx.Where(TicketColDrawID+" = ? AND "+TicketColIsChecked+" = false", drawID).Find(&tickets)
	return tickets, result.Error
}

// MarkCheckedInTransaction updates the status of the given tickets to checked.
func (r *TicketRepository) MarkCheckedInTransaction(tx *gorm.DB, ticketIDs []uuid.UUID) error {
	if len(ticketIDs) == 0 {
		return nil
	}
	return tx.Model(&models.Ticket{}).Where(TicketColID+" IN ?", ticketIDs).Update(TicketColIsChecked, true).Error
}

// ResetCheckedStatusByDrawIDInTransaction resets the checked status of all tickets for a draw.
func (r *TicketRepository) ResetCheckedStatusByDrawIDInTransaction(tx *gorm.DB, drawID uuid.UUID) error {
	return tx.Model(&models.Ticket{}).Where(TicketColDrawID+" = ?", drawID).Update(TicketColIsChecked, false).Error
}


