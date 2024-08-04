package invoice

import (
	"context"
	"errors"
	"time"
)

var(
	ErrFailedToGetInvoice = errors.New("failed to get invoice")
	ErrRecordNotFound = errors.New("invoice not found")
	ErrFailedToUpdate = errors.New("failed to update invoice")

)

type Repo interface {
	GetInvoicesByUserOrAgency(ctx context.Context, userID *uint, agencyID *uint, limit, offset uint) ([]Invoice,uint, error)
	GetInvoiceByTicketID(ctx context.Context, ticketID uint) (*Invoice, error)
	Insert(ctx context.Context, i *Invoice) error
	UpdateInvoiceStatus(ctx context.Context, invoiceID uint, status string) error
	UpdateInvoice(ctx context.Context, id uint, updates map[string]interface{}) error
	CalculateCompanyProfitForTrip(ctx context.Context, tripID uint) (float64, error)

}

type Invoice struct {
	ID             uint
	TicketID       uint
	IssuedDate     time.Time
	Info           string // Detailed information for the invoice
	PerAmountPrice float64
	TotalPrice     float64
	Status         string
	Penalty     float64

}

func NewInvoice(ticket_id uint, issuedDate time.Time, info string, perAmountPrice, totalPrice, penalty float64) *Invoice {
	return &Invoice{
		TicketID:       ticket_id,
		IssuedDate:     issuedDate,
		Info:           info,
		PerAmountPrice: perAmountPrice,
		TotalPrice:     totalPrice,
		Penalty: penalty,
	}
}
