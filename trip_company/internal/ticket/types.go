package ticket

import (
	"context"
	"errors"
	"tripcompanyservice/internal/invoice"
	"tripcompanyservice/internal/trip"
)

var(
    ErrRecordNotFound = errors.New("ticket not found")
    ErrFailedToGetTicket = errors.New("failed to ge ticket")
    ErrFailedToUpdate = errors.New("failed to update ticket")
)

type Repo interface {
	Insert(ctx context.Context, t *Ticket) error
    UpdateTicketStatus(ctx context.Context, ticketID uint, status string) error
    GetTicketsByUserOrAgency(ctx context.Context, userID *uint, agencyID *uint, limit, offset uint) ([]Ticket,uint, error) 
    GetFullTicketByID(ctx context.Context, id uint) (*Ticket, error)
    UpdateTicket(ctx context.Context, id uint, updates map[string]interface{}) error
	GetTicketsWithInvoicesByTripID(ctx context.Context, tripID uint) ([]Ticket, error)
}

type Ticket struct {
	ID         uint
	TripID     uint
	Trip       *trip.Trip
	UserID     *uint // Use `default:NULL` for nullable field
	AgencyID   *uint
	Quantity   int
	TotalPrice float64
	Status     string
    Penalty    float64
	Invoice   invoice.Invoice
}
