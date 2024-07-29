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
func (o *Ops) GetHotelsByID(ctx context.Context, id uint) (*Hotel, error) {
	return o.repo.GetHotelsByID(ctx, id)
}
func (o *Ops) GetHotels(ctx context.Context, city, country string, capacity, page, pageSize int) ([]Hotel, uint, error) {
	return o.repo.GetHotels(ctx, city, country, capacity, page, pageSize)
}

func (o *Ops) GetHotelsByOwnerID(ctx context.Context, ownerID uint, page, pageSize int) ([]Hotel, int, error) {
	return o.repo.GetHotelsByOwnerID(ctx, ownerID, page, pageSize)
}

func (o *Ops) Update(ctx context.Context, hotel *Hotel) error {
	// Ensure hotel exists before updating
	existingHotel, err := o.repo.GetHotelsByID(ctx, hotel.ID)
	if err != nil {
		return err
	}
	if existingHotel == nil {
		return ErrRecordNotFound
	}

	if err := ValidateHotelName(hotel.Name); err != nil {
		return ErrInvalidHotelName
	}

	if err := ValidateName(hotel.City); err != nil {
		return ErrInvalidName
	}

	if err := ValidateName(hotel.Country); err != nil {
		return ErrInvalidName
	}
	return o.repo.UpdateHotel(ctx, hotel)
}

func (o *Ops) Delete(ctx context.Context, id uint) error {
	// Ensure hotel exists before deleting
	existingHotel, err := o.repo.GetHotelsByID(ctx, id)
	if err != nil {
		return err
	}
	if existingHotel == nil {
		return ErrRecordNotFound
	}
	return o.repo.DeleteHotel(ctx, id)
}
