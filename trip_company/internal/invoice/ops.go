package invoice

import "context"

type Ops struct {
	repo Repo
}

func NewOps(repo Repo) *Ops {
	return &Ops{repo}
}

func (o *Ops) Create(ctx context.Context, i *Invoice) error {
	return o.repo.Insert(ctx, i)
}

func (o *Ops) UpdateInvoiceStatus(ctx context.Context, invoiceID uint, status string) error {
	return o.repo.UpdateInvoiceStatus(ctx, invoiceID, status)
}

func (o *Ops) GetInvoicesByUserOrAgency(ctx context.Context, userID *uint, agencyID *uint, page, pageSize uint) ([]Invoice, uint, error) {
	limit := pageSize
	offset := (page - 1) * pageSize

	return o.repo.GetInvoicesByUserOrAgency(ctx, userID, agencyID, limit, offset)
}

func (o *Ops) GetInvoiceByTicketID(ctx context.Context, ticketID uint) (*Invoice, error) {

	t, err := 	 o.repo.GetInvoiceByTicketID(ctx, ticketID)
	if err != nil {
		return nil, err
	}

	if t == nil {
		return nil, ErrRecordNotFound
	}
	return t, nil
}

func(o *Ops)UpdateInvoice(ctx context.Context, id uint, updates map[string]interface{}) error{
	return o.repo.UpdateInvoice(ctx,id, updates)
}
