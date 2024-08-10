package company

import (
	"context"
	"tripcompanyservice/internal"

	"github.com/google/uuid"
)

type Ops struct {
	repo Repo
}

func NewOps(repo Repo) *Ops {
	return &Ops{
		repo: repo,
	}
}

func (o *Ops) GetByID(ctx context.Context, id uint) (*TransportCompany, error) {
	t, err := o.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if t == nil {
		return nil, ErrCompanyNotFound
	}
	return t, nil
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

func (o *Ops) GetUserTransportCompanies(ctx context.Context, ownerID uuid.UUID, page, pageSize uint) ([]TransportCompany, uint, error) {
	limit := pageSize
	offset := (page - 1) * pageSize
	return o.repo.GetUserTransportCompanies(ctx, ownerID, limit, offset)
}

func (o *Ops) GetTransportCompanies(ctx context.Context, page, pageSize uint) ([]TransportCompany, uint, error) {
	limit := pageSize
	offset := (page - 1) * pageSize
	return o.repo.GetTransportCompanies(ctx, limit, offset)
}

func (o *Ops) Delete(ctx context.Context, companyID uint) error {
	return o.repo.Delete(ctx, companyID)
}

func (o *Ops) BlockUnBlockCompany(ctx context.Context, companyID uint, isBlocked bool) error {
	return o.repo.BlockCompany(ctx, companyID, isBlocked)
}

func(o *Ops)PatchCompanyByOwner(ctx context.Context, updatedCompany, originalCompany *TransportCompany) error{
	return o.repo.PatchCompany(ctx, updatedCompany, originalCompany)
}

func (o *Ops)IsUserOwnerOfCompany(ctx context.Context, companyID uint, userID uuid.UUID) (bool, error) {
	return o.repo.IsUserOwnerOfCompany(ctx, companyID, userID)
}
