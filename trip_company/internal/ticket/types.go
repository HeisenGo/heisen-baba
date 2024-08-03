package ticket

import (
	"context"
	"tripcompanyservice/internal/trip"
)

type Repo interface {
	Insert(ctx context.Context, t *Ticket) error
    UpdateTicketStatus(ctx context.Context, ticketID uint, status string) error
    GetTicketsByUserOrAgency(ctx context.Context, userID *uint, agencyID *uint, limit, offset uint) ([]Ticket,uint, error) 
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
}
