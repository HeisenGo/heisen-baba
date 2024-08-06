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
	"tripcompanyservice/pkg/ports/clients/clients"

	"github.com/google/uuid"
)

var (
	ErrAlreadyHasATeam   = errors.New("already has a team")
	ErrInvalidAssignment = errors.New("types are different")

	ErrSameState    = errors.New("same state")
	ErrFinishedTrip = errors.New("this trip finished")

	ErrCanceled     = errors.New("trip is canceled")
	ErrUnConfirmed  = errors.New("trip is not confirmed")
	ErrFutureTrip   = errors.New("you can not finish future trips")
	ErrPathNotFound = errors.New("path not found")
)

type TripService struct {
	tripOps     *trip.Ops
	companyOps  *company.Ops
	techTeamOps *techteam.Ops
	ticketOps   *ticket.Ops
	invoiceOps  *invoice.Ops
	pathClient  clients.IPathClient
}

func NewTripService(tripOps *trip.Ops, companyOps *company.Ops, 
	techTeamOps *techteam.Ops, ticketOps *ticket.Ops, 
	invoiceOps *invoice.Ops, pathClient clients.IPathClient) *TripService {
	return &TripService{
		tripOps:     tripOps,
		companyOps:  companyOps,
		techTeamOps: techTeamOps,
		ticketOps:   ticketOps,
		invoiceOps:  invoiceOps,
		pathClient: pathClient,
	}
}

func (s *TripService) GetCountPathUnfinishedTrips(ctx context.Context, pathID uint) (uint, error) {
	return s.tripOps.GetCountPathUnfinishedTrips(ctx, pathID)
}
func (s *TripService) GetUpcomingUnconfirmedTripIDsToCancel(ctx context.Context) error {
	// TODO : get them one by one and cancel them move money from libaba to the buyers wallet
	tripIDs, _ := s.tripOps.GetUpcomingUnconfirmedTripIDsToCancel(ctx)
	if len(tripIDs) > 0 {
		updates := make(map[string]interface{})
		for i := range len(tripIDs) {
			updates["is_canceled"] = true
			updates["status"] = "Canceled"

			err := s.tripOps.UpdateTripTechTimID(ctx, tripIDs[i], updates)
			if err != nil {
				return err
			}
			ticks, err := s.ticketOps.GetTicketsWithInvoicesByTripID(ctx, tripIDs[i])
			if err != nil {
				return err
			}

			var state string
			var refund, penalty float64
			//var receiverID uint
			var userRecieverID uuid.UUID
			state = "Canceled"
			penalty = 0
			for i := range len(ticks) {
				if ticks[i].AgencyID == nil {
					userRecieverID = *ticks[i].UserID
				} else {
					//receiverID = *ticks[i].AgencyID
					// get owner of agency id !!! TO Do
					userRecieverID = uuid.Nil
				}
				if ticks[i].Penalty != 0 { // cause already canceled by user
					refund = ticks[i].Penalty
				} else {
					refund = ticks[i].TotalPrice
				}
				s.UpdateInvoiceTicket(ctx, ticks[i], state, penalty, refund, userRecieverID)

			}
		}

	}
	return nil
}

func (s *TripService) GetCompanyTrips(ctx context.Context, originCity, destinationCity, pathType string, startDate *time.Time, requesterType string, companyID uint, userID uuid.UUID, page, pageSize uint) ([]trip.Trip, uint, string, error) {
	tCompany, err := s.companyOps.GetByID(ctx, companyID)
	if err != nil {
		return nil, 0, string(company.UserRole), err
	}

	if tCompany == nil {
		return nil, 0, string(company.UserRole), company.ErrCompanyNotFound
	}

	var role company.CompanyRoleType
	if requesterType == "agency" {
		role = company.CompanyRoleType(requesterType)
	} else if requesterType == "admin" {
		role = company.CompanyRoleType(requesterType)
	} else {
		role, _ = s.GetLoggedInUserRole(ctx, companyID, userID)
	}

	trips, total, err := s.tripOps.CompanyTrips(ctx, originCity, destinationCity, pathType, startDate, string(role), companyID, page, pageSize)
	if err != nil {
		return nil, 0, string(company.UserRole), err
	}
	return trips, total, string(role), nil
}

func (s *TripService) CreateTrip(ctx context.Context, t *trip.Trip, creatorID uuid.UUID) error {
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
	if tCompany.OwnerID != creatorID {
		return ErrForbidden
	}
	// GET PATH TODO: //****************************
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
	pp, err := s.pathClient.GetFullPathByID(uint32(t.PathID))
	if err != nil {
		return err
	}
	t.Path = pp
	//v := uint(1)
	//t.VehicleID = &v
	//********************************************************
	if err := s.tripOps.Create(ctx, t); err != nil {
		return err
	}

	return nil
}

func (s *TripService) GetTrips(ctx context.Context, originCity, destinationCity, pathType string, startDate *time.Time, requesterType string, page, pageSize uint) ([]trip.Trip, uint, error) {
	//check claim and role!!!
	if requesterType == "" {
		requesterType = "user"
	}
	return s.tripOps.GetTrips(ctx, originCity, destinationCity, pathType, startDate, requesterType, pageSize, page)
}

func (s *TripService) GetFullTripByID(ctx context.Context, id uint, requesterID uuid.UUID, requester string) (*trip.Trip, string, error) {
	t, err := s.tripOps.GetFullTripByID(ctx, id)
	if err != nil {
		return nil, string(company.UserRole), err
	}
	if requester == "admin" {
		return t, requester, nil
	}
	if requester == "agency" {
		if t.TourReleaseDate.After(time.Now()) {
			return nil, requester, trip.ErrTripUnAvailable
		}
		if t.TransportCompany.IsBlocked {
			return nil, requester, trip.ErrTripUnAvailable
		}
		return t, requester, nil

	}
	role, _ := s.GetLoggedInUserRole(ctx, t.TransportCompanyID, requesterID)

	if role == company.UserRole {
		if t.UserReleaseDate.After(time.Now()) {
			return nil, string(role), trip.ErrTripUnAvailable
		}

		if t.TransportCompany.IsBlocked {
			return nil, string(role), trip.ErrTripUnAvailable
		}
		return t, string(role), nil
	}
	return t, string(role), nil
}

func (s *TripService) UpdateTrip(ctx context.Context, id uint, newTrip *trip.Trip, requesterID uuid.UUID) (*trip.Trip, error) {
	oldTrip, err := s.tripOps.GetFullTripByID(ctx, id)

	if err != nil {
		return nil, err
	}

	// TO DO check permissions and roles ex: tech team can only update is_confirmed
	// admin can update is_finished
	// owner/operator can update everything according to the conditions of trip
	if oldTrip.TransportCompany.OwnerID != requesterID {
		return nil, ErrForbidden
	}
	err = s.tripOps.UpdateTrip(ctx, id, newTrip, oldTrip)
	if err != nil {
		return nil, err
	}
	return oldTrip, nil
}

func (s *TripService) SetTechTeamToTrip(ctx context.Context, tripID, techteamID uint, requesterID uuid.UUID) (*trip.Trip, error) {
	t, err := s.tripOps.GetFullTripByID(ctx, tripID)
	if err != nil {
		return nil, err
	}

	if requesterID != t.TransportCompany.OwnerID {
		return nil, ErrForbidden
	}

	if t.TechTeamID != nil {
		return nil, ErrAlreadyHasATeam
	}

	techteam, err := s.techTeamOps.GetTechTeamByID(ctx, techteamID)

	if techteam.TransportCompanyID != t.TransportCompanyID {
		return nil, ErrForbidden
	}

	if err != nil {
		return nil, err
	}

	if t.TripType != trip.TripType(techteam.TripType) {
		return nil, ErrInvalidAssignment
	}

	endDate := t.StartDate.Add(5 * time.Hour)
	// check team availability
	isAvailable, err := s.tripOps.CheckAvailabilityTechTeam(ctx, tripID, techteamID, *t.StartDate, endDate)
	if err != nil {
		return nil, err
	}
	if !isAvailable {
		return nil, errors.New("team is unavailable at this date")
	}
	updates := make(map[string]interface{})

	updates["tech_team_id"] = techteamID

	err = s.tripOps.UpdateTripTechTimID(ctx, tripID, updates)
	if err != nil {
		return nil, err
	}
	t.TechTeamID = &techteamID
	t.TechTeam = techteam
	return t, nil
}

func (s *TripService) UpdateInvoiceTicket(ctx context.Context, t ticket.Ticket, state string, penalty float64, refund float64, receiverID uuid.UUID) error {

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
	// cancel vehicle
	//notification
	return nil
}

func (s *TripService) CancelTrip(ctx context.Context, tripID uint, requesterID uuid.UUID, isCanceled bool) (*trip.Trip, error) {
	t, err := s.tripOps.GetFullTripByID(ctx, tripID)

	if err != nil {
		return nil, err
	}

	if t.TransportCompany.OwnerID != requesterID {
		return nil, ErrForbidden
	}

	if t.IsCanceled == isCanceled {
		return nil, ErrSameState
	}

	if t.IsFinished {
		return nil, ErrFinishedTrip
	}

	updates := make(map[string]interface{})

	updates["is_canceled"] = isCanceled
	if isCanceled {
		updates["status"] = "Canceled"
	}

	err = s.tripOps.UpdateTripTechTimID(ctx, tripID, updates)
	t.IsCanceled = isCanceled
	if err != nil {
		return nil, err
	}

	ticks, err := s.ticketOps.GetTicketsWithInvoicesByTripID(ctx, tripID)
	if err != nil {
		return nil, err
	}
	var state string
	var refund, penalty float64
	//var receiverID uint
	var userRecieverID uuid.UUID
	state = "Canceled"
	penalty = 0
	for i := range len(ticks) {
		if ticks[i].AgencyID == nil {
			userRecieverID = *ticks[i].UserID
		} else {
			//receiverID = *ticks[i].AgencyID
			// get owner of agency id !!! TO Do
			userRecieverID = uuid.Nil
		}
		if ticks[i].Penalty != 0 { // cause already canceled by user
			refund = ticks[i].Penalty
		} else {
			refund = ticks[i].TotalPrice
		}
		s.UpdateInvoiceTicket(ctx, ticks[i], state, penalty, refund, userRecieverID)
	}
	return t, nil

}

func (s *TripService) ConfirmTrip(ctx context.Context, tripID uint, requesterID uuid.UUID, isConfirmed bool) (*trip.Trip, error) {
	t, err := s.tripOps.GetFullTripByID(ctx, tripID)

	if err != nil {
		return nil, err
	}

	_, err = s.techTeamOps.GetTechTeamMemberByUserIDAndTechTeamID(ctx, requesterID, *t.TechTeamID)
	if err != nil {
		return nil, ErrForbidden
	}
	if t.IsCanceled {
		return nil, ErrCanceled
	}

	if t.IsConfirmed == isConfirmed {
		return nil, ErrSameState
	}

	if t.IsFinished {
		return nil, ErrFinishedTrip
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
	t.IsConfirmed = true

	return t, nil
}

func (s *TripService) FinishTrip(ctx context.Context, tripID uint, requesterID uuid.UUID, isFinished bool) (*trip.Trip, error) {
	t, err := s.tripOps.GetFullTripByID(ctx, tripID)

	if err != nil {
		return nil, err
	}

	if t.TransportCompany.OwnerID != requesterID {
		return nil, ErrForbidden
	}

	if !t.IsConfirmed {
		return nil, ErrUnConfirmed
	}

	if t.IsCanceled {
		return nil, ErrCanceled
	}
	if t.IsFinished == isFinished {
		return nil, ErrSameState
	}
	if t.EndDate.After(time.Now()) {
		return nil, ErrFutureTrip
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
	//notification
	return t, nil
}

func (s *TripService) GetLoggedInUserRole(ctx context.Context, companyID uint, userID uuid.UUID) (company.CompanyRoleType, error) {
	isOwner, _ := s.companyOps.IsUserOwnerOfCompany(ctx, companyID, userID)
	if isOwner {
		return company.OwnerRole, nil
	}
	isTechnician, err := s.techTeamOps.IsUserTechnicianInCompany(ctx, companyID, userID)
	if err != nil {
		return company.UserRole, err
	}

	if isTechnician {
		return company.TechRole, nil
	}

	return company.UserRole, nil
}
