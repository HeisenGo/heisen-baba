package invoice

import (
	"time"
	"tripcompanyservice/internal/ticket"
)

type Repo interface {
}

type Invoice struct {
	ID         uint
	TicketID   uint
	Ticket     ticket.Ticket
	IssuedDate time.Time
	Info       string // Detailed information for the invoice
	PerAmountPrice float64
	TotalPrice float64
}
