package reservation

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

func (o *Ops) Create(ctx context.Context, reservation *Reservation) error {
	reservation.Status = "pending"
	// if err := ValidateReservationStatus(reservation.Status); err != nil {
	// 	return err
	// }

	if err := ValidateAmount(reservation.TotalPrice); err != nil {
		return err
	}

	if err := ValidateDates(reservation.CheckIn, reservation.CheckOut); err != nil {
		return err
	}




	return o.repo.CreateReservation(ctx, reservation)
}

func (o *Ops) GetReservationsByHotelOwner(ctx context.Context, ownerID uuid.UUID, page, pageSize int) ([]Reservation, int, error) {
	return o.repo.GetReservationsByHotelOwner(ctx, ownerID, page, pageSize)
}

func (o *Ops) GetReservationByUserID(ctx context.Context, userid uuid.UUID) ([]Reservation, error) {
	return o.repo.GetReservationByUserID(ctx, userid)
}

func (o *Ops) GetReservationByID(ctx context.Context, id uint) (*Reservation, error) {
	return o.repo.GetReservationByID(ctx, id)
}

func (o *Ops) Update(ctx context.Context, reservation *Reservation) error {
	// Ensure reservation exists before updating
	existingReservation, err := o.repo.GetReservationByID(ctx, reservation.ID)
	if err != nil {
		return err
	}
	if existingReservation == nil {
		return ErrRecordNotFound
	}

	if err := ValidateReservationStatus(reservation.Status); err != nil {
		return err
	}

	if err := ValidateAmount(reservation.TotalPrice); err != nil {
		return err
	}

	if err := ValidateDates(reservation.CheckIn, reservation.CheckOut); err != nil {
		return err
	}

	return o.repo.UpdateReservation(ctx, reservation)
}

func (o *Ops) Delete(ctx context.Context, id uint) error {
	// Ensure reservation exists before deleting
	existingReservation, err := o.repo.GetReservationByID(ctx, id)
	if err != nil {
		return err
	}
	if existingReservation == nil {
		return ErrRecordNotFound
	}

	return o.repo.DeleteReservation(ctx, id)
}
