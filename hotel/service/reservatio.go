package service

import (
	"context"
	"hotel/internal/invoice"
	"hotel/internal/reservation"
	"time"

	"github.com/google/uuid"
)

type ReservationService struct {
	reservationOps *reservation.Ops
	invoiceOps     *invoice.Ops
}

func NewReservationService(reservationOps *reservation.Ops, invoiceOps *invoice.Ops) *ReservationService {
	return &ReservationService{
		reservationOps: reservationOps,
		invoiceOps:     invoiceOps,
	}
}

func (s *ReservationService) CreateReservation(ctx context.Context, res *reservation.Reservation) error {
	// Create the reservation
	err := s.reservationOps.Create(ctx, res)
	if err != nil {
		return err
	}

	// Create the invoice for the reservation
	inv := &invoice.Invoice{
		ReservationID: res.ID,
		IssueDate:     time.Now(),
		Amount:        res.TotalPrice,
		UserID:        res.UserID,
		OwnerID:       res.OwnerID,
		Paid:          false,
	}

	err = s.invoiceOps.Create(ctx, inv)
	if err != nil {
		return err
	}
	// TO DO : Check the Paid status of invoice
	return nil
}

func (s *ReservationService) GetReservationsByHotelOwner(ctx context.Context, ownerID uuid.UUID, page, pageSize int) ([]reservation.Reservation, int, error) {
	return s.reservationOps.GetReservationsByHotelOwner(ctx, ownerID, page, pageSize)
}

func (s *ReservationService) GetReservationByUserID(ctx context.Context, userID uuid.UUID) ([]reservation.Reservation, error) {
	return s.reservationOps.GetReservationByUserID(ctx, userID)
}

func (s *ReservationService) GetReservationByID(ctx context.Context, id uint) (*reservation.Reservation, error) {
	return s.reservationOps.GetReservationByID(ctx, id)
}

func (s *ReservationService) UpdateReservation(ctx context.Context, id uint, updates *reservation.Reservation) error {
	existingReservation, err := s.reservationOps.GetReservationByID(ctx, id)
	if err != nil {
		return err
	}

	// Update only the fields that are provided
	if updates.RoomID != 0 {
		existingReservation.RoomID = updates.RoomID
	}
	if updates.UserID != (uuid.UUID{}) {
		existingReservation.UserID = updates.UserID
	}
	if !updates.CheckIn.IsZero() {
		existingReservation.CheckIn = updates.CheckIn
	}
	if !updates.CheckOut.IsZero() {
		existingReservation.CheckOut = updates.CheckOut
	}
	if updates.TotalPrice != 0 {
		existingReservation.TotalPrice = updates.TotalPrice
	}
	if updates.Status != "" {
		existingReservation.Status = updates.Status
	}

	return s.reservationOps.Update(ctx, existingReservation)
}

func (s *ReservationService) DeleteReservation(ctx context.Context, id uint) error {
	_, err := s.reservationOps.GetReservationByID(ctx, id)
	if err != nil {
		return err
	}
	return s.reservationOps.Delete(ctx, id)
}
