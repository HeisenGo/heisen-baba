package tour

import (
	"context"
)

type Ops struct {
	repo Repo
}

func NewOps(repo Repo) *Ops {
	return &Ops{repo}
}

func (o *Ops) CreateTour(ctx context.Context, tour *Tour) error {
	if err := ValidateTourCapacity(tour.Capacity); err != nil {
		return ErrInvalidTourCapacity
	}
	return o.repo.CreateTour(ctx, tour)
}

func (o *Ops) GetTourByID(ctx context.Context, id uint) (*Tour, error) {
	return o.repo.GetTourByID(ctx, id)
}

func (o *Ops) GetTours(ctx context.Context, agencyID uint, page, pageSize int) ([]Tour, uint, error) {
	return o.repo.GetTours(ctx, agencyID, page, pageSize)
}

func (o *Ops) GetToursByAgencyID(ctx context.Context, agencyID uint, page, pageSize int) ([]Tour, int, error) {
	return o.repo.GetToursByAgencyID(ctx, agencyID, page, pageSize)
}

func (o *Ops) UpdateTour(ctx context.Context, tour *Tour) error {
	// Ensure tour exists before updating
	existingTour, err := o.repo.GetTourByID(ctx, tour.ID)
	if err != nil {
		return err
	}
	if existingTour == nil {
		return ErrRecordNotFound
	}

	if err := ValidateTourCapacity(tour.Capacity); err != nil {
		return ErrInvalidTourCapacity
	}

	return o.repo.UpdateTour(ctx, tour)
}

func (o *Ops) DeleteTour(ctx context.Context, id uint) error {
	// Ensure tour exists before deleting
	existingTour, err := o.repo.GetTourByID(ctx, id)
	if err != nil {
		return err
	}
	if existingTour == nil {
		return ErrRecordNotFound
	}

	return o.repo.DeleteTour(ctx, id)
}

func (o *Ops) ApproveTour(ctx context.Context, tourID uint) error {
	return o.repo.ApproveTour(ctx, tourID)
}

func (o *Ops) SetTourStatus(ctx context.Context, tourID uint, isActive bool) error {
	return o.repo.SetTourStatus(ctx, tourID, isActive)
}
