package service

import (
	"agency/internal/agency"
	"agency/internal/invoice"
	"agency/internal/reservation"
	"agency/internal/tour"
	"agency/pkg/ports/clients/clients"
	"context"
	"errors"
	"time"
)


var (
	ErrTourNotFound          = errors.New("tour not found")
	ErrReservationNotCreated = errors.New("reservation not created")
	ErrInvoiceNotCreated     = errors.New("invoice not created")
)

type TourService struct {
	tourOps        *tour.Ops
	reservationOps *reservation.Ops
	agencyOps       *agency.Ops
	invoiceOps     *invoice.Ops
	bankClient     clients.IBankClient
}

func NewTourService(tourOps *tour.Ops, reservationOps *reservation.Ops, agencyOps *agency.Ops, invoiceOps *invoice.Ops, bankClient clients.IBankClient) *TourService {
	return &TourService{
		tourOps:        tourOps,
		reservationOps: reservationOps,
		agencyOps:       agencyOps,
		invoiceOps:     invoiceOps,
		bankClient:     bankClient,
	}
}
func (s *TourService) CreateTourReservation(ctx context.Context, res *reservation.Reservation) error {
	// Ensure Tour exists before creating a reservation
	existingTour, err := s.tourOps.GetTourByID(ctx, res.TourID)
	if err != nil {
		return ErrTourNotFound
	}
	if existingTour == nil {
		return ErrTourNotFound
	}
	res.TotalPrice = existingTour.UserPrice
	agency, err := s.agencyOps.GetAgencyByID(ctx, existingTour.AgencyID)
	if err != nil {
		return err
	}

	res.OwnerID = agency.OwnerID
	// Validate and create the reservation
	err = s.reservationOps.Create(ctx, res)
	if err != nil {
		return err
	}
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
	isSuccess, err := s.bankClient.Transfer(inv.UserID.String(), inv.OwnerID.String(), false, inv.Amount)
	if err != nil {
		return err
	}
	if isSuccess {
		inv.Paid = true
	}
	if err := s.reservationOps.Create(ctx, res); err != nil {
		return ErrReservationNotCreated

	}

	return nil
}
func (s *TourService) CreateTour(ctx context.Context, t *tour.Tour) error {
	return s.tourOps.CreateTour(ctx, t)
}

func (s *TourService) GetTours(ctx context.Context, agencyID uint, page, pageSize int) ([]tour.Tour, uint, error) {
	return s.tourOps.GetTours(ctx, agencyID, page, pageSize)
}

func (s *TourService) GetToursByAgencyID(ctx context.Context, agencyID uint, page, pageSize int) ([]tour.Tour, int, error) {
	return s.tourOps.GetToursByAgencyID(ctx, agencyID, page, pageSize)
}

func (s *TourService) UpdateTour(ctx context.Context, id uint, updates *tour.Tour) error {
	existingTour, err := s.tourOps.GetTourByID(ctx, id)
	if err != nil {
		return err
	}

	// Update only the fields that are provided
	if updates.GoTicketID != 0 {
		existingTour.GoTicketID = updates.GoTicketID
	}
	if updates.BackTicketID != 0 {
		existingTour.BackTicketID = updates.BackTicketID
	}
	if updates.HotelID != 0 {
		existingTour.HotelID = updates.HotelID
	}
	if updates.Capacity != 0 {
		existingTour.Capacity = updates.Capacity
	}
	existingTour.IsApproved = updates.IsApproved
	existingTour.IsActive = updates.IsActive

	return s.tourOps.UpdateTour(ctx, existingTour)
}

func (s *TourService) DeleteTour(ctx context.Context, id uint) error {
	return s.tourOps.DeleteTour(ctx, id)
}

func (s *TourService) ApproveTour(ctx context.Context, tourID uint) error {
	return s.tourOps.ApproveTour(ctx, tourID)
}

func (s *TourService) SetTourStatus(ctx context.Context, tourID uint, isActive bool) error {
	return s.tourOps.SetTourStatus(ctx, tourID, isActive)
}
