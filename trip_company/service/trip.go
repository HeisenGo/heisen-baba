package service

import (
	"context"
	"tripcompanyservice/internal/company"
	"tripcompanyservice/internal/trip"
)

type TripService struct {
	tripOps    *trip.Ops
	companyOps *company.Ops
}

func NewTripService(tripOps *trip.Ops, companyOps *company.Ops) *TripService {
	return &TripService{
		tripOps:    tripOps,
		companyOps: companyOps,
	}
}

func (s *TripService) GetCompanyTrips(ctx context.Context, companyID uint, page, pageSize uint) ([]trip.Trip, uint, error) {
	tCompany, err := s.companyOps.GetByID(ctx, companyID)
	if err != nil {
		return nil, 0, err
	}

	if tCompany == nil {
		return nil, 0, company.ErrCompanyNotFound
	}

	return s.GetCompanyTrips(ctx, companyID, page, pageSize)
}

func (s *TripService) CreateTrip(ctx context.Context, t *trip.Trip, creatorID uint) error {
	// user, err := s.userOps.GetUserByID(ctx, o.UserID)
	// if err != nil {
	// 	return err
	// }

	// if user == nil {
	// 	return u.ErrUserNotFound
	// }

	tCompany, err := s.companyOps.GetByID(ctx, t.TransportCompanyID)
	if err != nil {
		return err
	}

	if tCompany == nil {
		return company.ErrCompanyNotFound
	}

	// GET PATH TODO:
	t.Path = &trip.Path{
		FromTerminal: &trip.Terminal{},
		ToTerminal:   &trip.Terminal{},
	}
	t.Path.Name = "jjdjdjlk"
	t.Path.FromTerminal.Name = "kjdkdkdkk"
	t.Origin = "Tehran"
	t.Destination = "Mashhad"
	t.Path.ToTerminalID = 2
	t.Path.FromTerminalID = 1
	t.Path.ToTerminal.Name = "central"
	t.Path.Type = "rail"
	t.TripType = "rail"
	if err := s.tripOps.Create(ctx, t); err != nil {
		return err
	}

	return nil
}
