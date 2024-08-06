package ticket

import (
	"context"

	"github.com/google/uuid"
)

type Ops struct {
	repo Repo
}

func NewOps(repo Repo) *Ops {
	return &Ops{repo}
}

func (o *Ops) Create(ctx context.Context, t *Ticket) error {
	return o.repo.Insert(ctx, t)
}

func (o *Ops) UpdateTicketStatus(ctx context.Context, ticketID uint, status string) error {
	return o.repo.UpdateTicketStatus(ctx, ticketID, status)
}

func (o *Ops) GetTicketsByUserOrAgency(ctx context.Context, userID *uuid.UUID, agencyID *uint, page, pageSize uint) ([]Ticket, uint, error) {
	limit := pageSize
	offset := (page - 1) * pageSize
	return o.repo.GetTicketsByUserOrAgency(ctx, userID, agencyID, limit, offset)
}

func (o *Ops) GetFullTicketByID(ctx context.Context, id uint) (*Ticket, error) {

	t, err := o.repo.GetFullTicketByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if t == nil {
		return nil, ErrRecordNotFound
	}
	return t, nil
}

func (o *Ops) UpdateTicket(ctx context.Context, id uint, updates map[string]interface{}) error {
	return o.repo.UpdateTicket(ctx, id, updates)
}


func (o *Ops)GetTicketsWithInvoicesByTripID(ctx context.Context, tripID uint) ([]Ticket, error) {
	return o.repo.GetTicketsWithInvoicesByTripID(ctx, tripID)
}