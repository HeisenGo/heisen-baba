package service

import (
	"context"
	"agency/internal/tour"
)

type TourService struct {
	tourOps *tour.Ops
}

func NewTourService(tourOps *tour.Ops) *TourService {
	return &TourService{
		tourOps: tourOps,
	}
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
