package repository

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"lotto-journal/api/internal/models"
)

const (
	TableTickets       = "tickets"
	TicketColID        = "id"
	TicketColOwnerID   = "owner_id"
	TicketColDrawID    = "draw_id"
	TicketColType      = "type"
	TicketColNumber    = "number"
	TicketColQuantity  = "quantity"
	TicketColIsChecked = "is_checked"
)

type DrawTicketWithOwner struct {
	ID         uuid.UUID
	OwnerID    uuid.UUID
	Type       string
	Number     string
	Quantity   int
	LineUserID string
	Language   string
}

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

func (r *TicketRepository) FindDrawTicketsWithOwners(drawID uuid.UUID) ([]DrawTicketWithOwner, error) {
	var tickets []DrawTicketWithOwner

	selectFields := fmt.Sprintf("%s.%s, %s.%s, %s.%s, %s.%s, %s.%s, %s.%s, %s.%s",
		TableTickets, TicketColID,
		TableTickets, TicketColOwnerID,
		TableTickets, TicketColType,
		TableTickets, TicketColNumber,
		TableTickets, TicketColQuantity,
		TableUsers, UserColLineUserID,
		TableUsers, UserColLanguage,
	)

	joinClause := fmt.Sprintf("JOIN %s ON %s.%s = %s.%s",
		TableUsers,
		TableTickets, TicketColOwnerID,
		TableUsers, UserColID,
	)

	whereClause := fmt.Sprintf("%s.%s = ? AND %s.%s = ?",
		TableTickets, TicketColDrawID,
		TableUsers, UserColStatus,
	)

	err := r.db.Table(TableTickets).
		Select(selectFields).
		Joins(joinClause).
		Where(whereClause, drawID, "active").
		Scan(&tickets).Error

	return tickets, err
}


