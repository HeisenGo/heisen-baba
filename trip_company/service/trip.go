package service

import (
	"context"
	"errors"
	"time"
	"tripcompanyservice/internal/company"
	"tripcompanyservice/internal/techteam"
	"tripcompanyservice/internal/trip"
)

var (
	ErrAlreadyHasATeam   = errors.New("already has a team")
	ErrInvalidAssignment = errors.New("types are different")
)

type TripService struct {
	tripOps     *trip.Ops
	companyOps  *company.Ops
	techTeamOps *techteam.Ops
}

func NewTripService(tripOps *trip.Ops, companyOps *company.Ops, techTeamOps *techteam.Ops) *TripService {
	return &TripService{
		tripOps:     tripOps,
		companyOps:  companyOps,
		techTeamOps: techTeamOps,
	}
}

func (s *TripService) GetCountPathUnfinishedTrips(ctx context.Context, pathID uint) (uint, error) {
	return s.tripOps.GetCountPathUnfinishedTrips(ctx, pathID)
}
func (s *TripService) GetUpcomingUnconfirmedTripIDsToCancel(ctx context.Context) ([]uint, error) {
	// TODO : get them one by one and cancel them move money from libaba to the buyers wallet
	return s.tripOps.GetUpcomingUnconfirmedTripIDsToCancel(ctx)
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
	//******************************************
	t.Path.Name = "jjdjdjlk"
	t.Path.FromTerminal.Name = "kjdkdkdkk"
	t.Origin = "Tehran"
	t.Destination = "Mashhad"
	t.Path.ToTerminalID = 2
	t.Path.FromTerminalID = 1
	t.Path.ToTerminal.Name = "central"
	t.Path.Type = "rail"
	t.TripType = trip.TripType(t.Path.Type)
	t.TripType = "rail"
	t.Path.DistanceKM = 220
	v := uint(1)
	t.VehicleID = &v
	//********************************************************
	if err := s.tripOps.Create(ctx, t); err != nil {
		return err
	}

	return nil
}

func (s *TripService) GetTrips(ctx context.Context, originCity, destinationCity, pathType string, startDate *time.Time, requesterType string, page, pageSize uint) ([]trip.Trip, uint, error) {
	//check claim and role!!!
	return s.tripOps.GetTrips(ctx, originCity, destinationCity, pathType, startDate, requesterType, pageSize, page)
}

func (s *TripService) GetFullTripByID(ctx context.Context, id uint) (*trip.Trip, error) {
	return s.tripOps.GetFullTripByID(ctx, id)
}

func (s *TripService) UpdateTrip(ctx context.Context, id uint, newTrip *trip.Trip) (*trip.Trip, error) {
	oldTrip, err := s.tripOps.GetFullTripByID(ctx, id)

	if err != nil {
		return nil, err
	}

	// TO DO check permissions and roles ex: tech team can only update is_confirmed
	// admin can update is_finished
	// owner/operator can update everything according to the conditions of trip

	err = s.tripOps.UpdateTrip(ctx, id, newTrip, oldTrip)
	if err != nil {
		return nil, err
	}
	return oldTrip, nil
}

func (s *TripService) SetTechTeamToTrip(ctx context.Context, tripID, techteamID uint) (*trip.Trip, error) {
	t, err := s.tripOps.GetFullTripByID(ctx, tripID)
	if err != nil {
		return nil, err
	}

	if t.TechTeamID != nil {
		return nil, ErrAlreadyHasATeam
	}

	techteam, err := s.techTeamOps.GetTechTeamByID(ctx, techteamID)

	if err != nil {
		return nil, err
	}

	if t.TripType != trip.TripType(techteam.TripType) {
		return nil, ErrInvalidAssignment
	}

	updates := make(map[string]interface{})

	updates["tech_team_id"] = techteamID

	err = s.tripOps.UpdateTripTechTimID(ctx, tripID, updates)
	if err != nil {
		return nil, err
	}
	t.TechTeamID = &techteamID
	return t, nil
}
