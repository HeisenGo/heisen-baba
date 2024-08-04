package agency

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

func (o *Ops) CreateAgency(ctx context.Context, agency *Agency) error {
	if err := ValidateAgencyName(agency.Name); err != nil {
		return ErrInvalidAgencyName
	}

	return o.repo.CreateAgency(ctx, agency)
}

func (o *Ops) GetAgencyByID(ctx context.Context, id uint) (*Agency, error) {
	return o.repo.GetAgencyByID(ctx, id)
}

func (o *Ops) GetAgencies(ctx context.Context, name string, page, pageSize int) ([]Agency, uint, error) {
	return o.repo.GetAgencies(ctx, name, page, pageSize)
}

func (o *Ops) GetAgenciesByOwnerID(ctx context.Context, ownerID uuid.UUID, page, pageSize int) ([]Agency, int, error) {
	return o.repo.GetAgenciesByOwnerID(ctx, ownerID, page, pageSize)
}

func (o *Ops) UpdateAgency(ctx context.Context, agency *Agency) error {
	// Ensure agency exists before updating
	existingAgency, err := o.repo.GetAgencyByID(ctx, agency.ID)
	if err != nil {
		return err
	}
	if existingAgency == nil {
		return ErrRecordNotFound
	}

	if err := ValidateAgencyName(agency.Name); err != nil {
		return ErrInvalidAgencyName
	}

	return o.repo.UpdateAgency(ctx, agency)
}

func (o *Ops) DeleteAgency(ctx context.Context, id uint) error {
	// Ensure agency exists before deleting
	existingAgency, err := o.repo.GetAgencyByID(ctx, id)
	if err != nil {
		return err
	}
	if existingAgency == nil {
		return ErrRecordNotFound
	}

	return o.repo.DeleteAgency(ctx, id)
}

func (o *Ops) BlockAgency(ctx context.Context, agencyID uint) error {
	return o.repo.BlockAgency(ctx, agencyID)
}
