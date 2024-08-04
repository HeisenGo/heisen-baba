package invoice

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Repo interface {
	CreateInvoice(ctx context.Context, invoice *Invoice) error
	GetInvoicesByHotelID(ctx context.Context, hotelID uint, page, pageSize int) ([]Invoice, int, error)
	GetInvoicesByUserID(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]Invoice, int, error)
	GetInvoiceByID(ctx context.Context, id uint) (*Invoice, error)
	UpdateInvoice(ctx context.Context, invoice *Invoice) error
	DeleteInvoice(ctx context.Context, id uint) error
}

type Invoice struct {
	ID            uint
	UserID        uuid.UUID
	OwnerID       uuid.UUID
	ReservationID uint
	IssueDate     time.Time
	Amount        uint64
	Paid          bool
}

var (
	ErrInvalidAmount  = errors.New("invalid amount: must be a positive number")
	ErrRecordNotFound = errors.New("record not found")
)

func ValidateAmount(amount uint64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	return nil
}
