package ticket

import (
	"context"
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
