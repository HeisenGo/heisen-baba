package service

import (
	"context"
	"errors"
	"time"
	"tripcompanyservice/internal/company"
	"tripcompanyservice/internal/invoice"
	"tripcompanyservice/internal/techteam"
	"tripcompanyservice/internal/ticket"
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
	ticketOps   *ticket.Ops
	invoiceOps  *invoice.Ops
}

func NewTripService(tripOps *trip.Ops, companyOps *company.Ops, techTeamOps *techteam.Ops, ticketOps *ticket.Ops, invoiceOps *invoice.Ops) *TripService {
	return &TripService{
		tripOps:     tripOps,
		companyOps:  companyOps,
		techTeamOps: techTeamOps,
		ticketOps:   ticketOps,
		invoiceOps:  invoiceOps,
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

func (s *TripService) UpdateInvoiceTicket(ctx context.Context, t ticket.Ticket, state string, penalty float64, refund float64, receiverID uint) error {

	updates := make(map[string]interface{})
	updates["status"] = state
	updates["penalty"] = penalty
	err := s.ticketOps.UpdateTicket(ctx, t.ID, updates)

	if err != nil {
		return err
	}
	err = s.invoiceOps.UpdateInvoice(ctx, t.Invoice.ID, updates)

	if err != nil {
		return err
	}
	// send to band : TODO from alibaba to user/aganecyId
	return nil
}

func (s *TripService) CancelTrip(ctx context.Context, tripID uint, requesterID uint, isCanceled bool) (*trip.Trip, error) {
	t, err := s.tripOps.GetFullTripByID(ctx, tripID)

	if err != nil {
		return nil, err
	}

	if t.TransportCompany.OwnerID != requesterID {
		return nil, ErrForbidden
	}

	if t.IsCanceled == isCanceled {
		return nil, errors.New("same state")
	}

	if t.IsFinished {
		return nil, errors.New("this trip finished")
	}

	updates := make(map[string]interface{})

	updates["is_canceled"] = isCanceled
	if isCanceled {
		updates["status"] = "Canceled"
	}

	err = s.tripOps.UpdateTripTechTimID(ctx, tripID, updates)
	if err != nil {
		return nil, err
	}

	ticks, err := s.ticketOps.GetTicketsWithInvoicesByTripID(ctx, tripID)
	if err != nil {
		return nil, err
	}
	var state string
	var refund, penalty float64
	var receiverID uint
	state = "Canceled"
	penalty = 0
	for i := range len(ticks) {
		if ticks[i].AgencyID == nil {
			receiverID = *ticks[i].UserID
		} else {
			receiverID = *ticks[i].AgencyID
		}
		if ticks[i].Penalty != 0 { // already canceled by user
			refund = ticks[i].Penalty
		} else {
			refund = ticks[i].TotalPrice
		}
		s.UpdateInvoiceTicket(ctx, ticks[i], state, penalty, refund, receiverID)
	}
	return t, nil

}

func (s *TripService) ConfirmTrip(ctx context.Context, tripID uint, requesterID uint, isConfirmed bool) (*trip.Trip, error) {
	t, err := s.tripOps.GetFullTripByID(ctx, tripID)

	if err != nil {
		return nil, err
	}

	// if t.TransportCompany.OwnerID != requesterID{
	// 	return nil, ErrForbidden
	// }
	_, err = s.techTeamOps.GetTechTeamMemberByUserIDAndTechTeamID(ctx, requesterID, *t.TechTeamID)
	if err != nil {
		return nil, ErrForbidden
	}
	if t.IsCanceled {
		return nil, errors.New("trip is canceled")
	}

	if t.IsConfirmed == isConfirmed {
		return nil, errors.New("same state")
	}

	if t.IsFinished {
		return nil, errors.New("this trip finished")
	}

	updates := make(map[string]interface{})

	updates["is_confirmed"] = isConfirmed
	if isConfirmed {
		updates["status"] = "Confirmed"
	}
	err = s.tripOps.UpdateTripTechTimID(ctx, tripID, updates)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (s *TripService) FinishTrip(ctx context.Context, tripID uint, requesterID uint, isFinished bool) (*trip.Trip, error) {
	t, err := s.tripOps.GetFullTripByID(ctx, tripID)

	if err != nil {
		return nil, err
	}

	if t.TransportCompany.OwnerID != requesterID {
		return nil, ErrForbidden
	}

	if !t.IsConfirmed {
		return nil, errors.New("trip is not confirmed")
	}

	if t.IsCanceled {
		return nil, errors.New("trip is canceled")
	}
	if t.IsFinished == isFinished {
		return nil, errors.New("same state")
	}

	updates := make(map[string]interface{})

	updates["is_finished"] = isFinished
	if isFinished {
		updates["status"] = "Finished"
	}

	profit, err := s.invoiceOps.CalculateCompanyProfitForTrip(ctx, tripID) // profit
	if err != nil {
		return nil, err
	}
	updates["profit"] = profit
	err = s.tripOps.UpdateEndDateTrip(ctx, tripID, updates)
	if err != nil {
		return nil, err
	}
	t.Profit = profit
	//TODO : bank
	// calculate profit tell alibaba to move money from alibaba to owner id wallet: profit = totalprice (status canceled nis) + penalty (status cancel and )

	return t, nil
}
