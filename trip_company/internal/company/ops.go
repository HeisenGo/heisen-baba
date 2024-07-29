package company

import (
	"context"
	"tripcompanyservice/internal"
)

type Ops struct {
	repo Repo
}

func NewOps(repo Repo) *Ops {
	return &Ops{
		repo: repo,
	}
}

func (o *Ops) Create(ctx context.Context, c *TransportCompany) error {
	// if c.Email != "" {
	// 	c.Email = LowerCaseEmail(c.Email)
	// 	if err := ValidateEmail(c.Email); err != nil {
	// 		return err
	// 	}
	// }
	if err := internal.ValidateName(c.Name, MaxNameLength); err != nil {
		return err
	}

	if c.Description != "" {
		if err := internal.ValidateName(c.Description, MaxDescriptionLength); err != nil {
			return err
		}

	}
	if c.Address != "" {
		if err := internal.ValidateName(c.Address, MaxAddressLength); err != nil {
			return err
		}
	}

	return o.repo.Insert(ctx, c)
}

func (o *Ops) GetUserTransportCompanies(ctx context.Context, ownerID uint, page, pageSize uint) ([]TransportCompany, uint, error) {
	limit := pageSize
	offset := (page - 1) * pageSize
	return o.repo.GetUserTransportCompanies(ctx, ownerID, limit, offset)
}

func (o *Ops) GetTransportCompanies(ctx context.Context, page, pageSize uint) ([]TransportCompany, uint, error) {
	limit := pageSize
	offset := (page - 1) * pageSize
	return o.repo.GetTransportCompanies(ctx, limit, offset)
}
