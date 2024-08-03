package invoice

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

func (o *Ops) Create(ctx context.Context, invoice *Invoice) error {
	if err := ValidateAmount(invoice.Amount); err != nil {
		return err
	}
	return o.repo.CreateInvoice(ctx, invoice)
}

func (o *Ops) GetInvoiceByID(ctx context.Context, id uint) (*Invoice, error) {
	invoice, err := o.repo.GetInvoiceByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if invoice == nil {
		return nil, ErrRecordNotFound
	}
	return invoice, nil
}

func (o *Ops) GetInvoicesByHotelID(ctx context.Context, hotelID uint, page, pageSize int) ([]Invoice, int, error) {
	invoices, total, err := o.repo.GetInvoicesByHotelID(ctx, hotelID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	if len(invoices) == 0 {
		return nil, 0, ErrRecordNotFound
	}
	return invoices, total, nil
}

func (o *Ops) GetInvoicesByUserID(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]Invoice, int, error) {
	invoices, total, err := o.repo.GetInvoicesByUserID(ctx, userID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	if len(invoices) == 0 {
		return nil, 0, ErrRecordNotFound
	}
	return invoices, total, nil
}

func (o *Ops) Update(ctx context.Context, invoice *Invoice) error {
	// Ensure invoice exists before updating
	existingInvoice, err := o.repo.GetInvoiceByID(ctx, invoice.ID)
	if err != nil {
		return err
	}
	if existingInvoice == nil {
		return ErrRecordNotFound
	}

	if err := ValidateAmount(invoice.Amount); err != nil {
		return err
	}

	return o.repo.UpdateInvoice(ctx, invoice)
}

func (o *Ops) Delete(ctx context.Context, id uint) error {
	// Ensure invoice exists before deleting
	existingInvoice, err := o.repo.GetInvoiceByID(ctx, id)
	if err != nil {
		return err
	}
	if existingInvoice == nil {
		return ErrRecordNotFound
	}

	return o.repo.DeleteInvoice(ctx, id)
}
