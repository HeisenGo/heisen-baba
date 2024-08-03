package invoice

import (
	"context"
	"time"
	"tripcompanyservice/internal/ticket"
)

type Repo interface {
	Insert(ctx context.Context, i *Invoice) error
	UpdateInvoiceStatus(ctx context.Context, invoiceID uint, status string) error
}

type Invoice struct {
	ID             uint
	TicketID       uint
	Ticket         *ticket.Ticket
	IssuedDate     time.Time
	Info           string // Detailed information for the invoice
	PerAmountPrice float64
	TotalPrice     float64
}

func NewInvoice(ticket_id uint, ticket *ticket.Ticket, issuedDate time.Time, info string, perAmountPrice, totalPrice float64) *Invoice {
	return &Invoice{
		TicketID:       ticket_id,
		Ticket:         ticket,
		IssuedDate:     issuedDate,
		Info:           info,
		PerAmountPrice: perAmountPrice,
		TotalPrice:     totalPrice,
	}
}
