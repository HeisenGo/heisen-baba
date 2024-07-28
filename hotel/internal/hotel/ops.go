package hotel

import (
	"context"
	
)

type Ops struct {
	repo Repo
}

func NewOps(repo Repo) *Ops {
	return &Ops{repo}
}

func (o *Ops) Create(ctx context.Context, hotel *Hotel) error {
	if err := ValidateHotelName(hotel.Name); err != nil {
		return ErrInvalidHotelName
	}

	if err := ValidateName(hotel.City); err != nil {
		return ErrInvalidName
	}
	
	if err := ValidateName(hotel.Country); err != nil {
		return ErrInvalidName
	}
	return o.repo.CreateHotel(ctx, hotel)
}