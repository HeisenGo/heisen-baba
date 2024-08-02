package presenter

import (
	"time"
	"tripcompanyservice/internal/trip"
	"tripcompanyservice/pkg/fp"
)

type CreateTripReq struct {
	ID                 uint      `json:"id"`
	TransportCompanyID uint      `json:"company_id"  validate:"required"`
	UserReleaseDate    Timestamp `json:"user_date"  validate:"required"`
	TourReleaseDate    Timestamp `json:"tour_date"  validate:"required"`
	UserPrice          float64   `json:"user_price"  validate:"required"`
	AgencyPrice        float64   `json:"agency_price"  validate:"required"`
	PathID             uint      `json:"path_id"  validate:"required"`
	MinPassengers      uint      `json:"min_pass"`
	TechTeamID         *uint     `json:"tech_team_id"`
	MaxTickets         uint      `json:"max_tickets" validate:"required"`
	SoldTickets        uint      `json:"sold_tickets"`
	StartDate          Timestamp `json:"start_date" validate:"required"` // should be given by trip generator
	// example: "2024-07-30 23:27:09"
	//end date should be calculated according to the vehicle speed and path distance
	TripCancellingPenalty *CreateTripCancelingPenaltyReq `json:"penalty" validate:"required"`
	EndDate               Timestamp                      `json:"end_date"`
}

type CreateTripRes struct {
	ID                 uint      `json:"id"`
	TransportCompanyID uint      `json:"company_id"`
	UserReleaseDate    Timestamp `json:"user_date" `
	TourReleaseDate    Timestamp `json:"tour_date"  `
	UserPrice          float64   `json:"user_price" `
	Origin             string    `json:"from"`
	PathName           string    `json:"path"`
	FromTerminalName   string    `json:"from_terminal"`
	Destination        string    `json:"to"`
	ToTerminalName     string    `json:"to_terminal"`
	Type               string    `json:"type"`
	AgencyPrice        float64   `json:"agency_price"  `
	PathID             uint      `json:"path_id" `
	MinPassengers      uint      `json:"min_pass"`
	TechTeamID         *uint     `json:"tech_team_id"`
	MaxTickets         uint      `json:"max_tickets" `
	SoldTickets        uint      `json:"sold_tickets"`
	StartDate          Timestamp `json:"start_date" ` // should be given by trip generator
	// example: "2024-07-30 23:27:09"
	//end date should be calculated according to the vehicle speed and path distance
	TripCancellingPenalty *CreateTripCancelingPenaltyReq `json:"penalty" `
	EndDate               Timestamp                      `json:"end_date"`
}

func CreateTripReqToTrip(req *CreateTripReq) *trip.Trip {
	penalty := CreateTripCancelingPenaltyReqToTripCancellingPenalty(req.TripCancellingPenalty)

	// Convert StartDate and EndDate
	var startDate, endDate *time.Time
	if !time.Time(req.StartDate).IsZero() {
		startDate = (*time.Time)(&req.StartDate)
	}
	if !time.Time(req.EndDate).IsZero() {
		endDate = (*time.Time)(&req.EndDate)
	}

	return &trip.Trip{
		TransportCompanyID:    req.TransportCompanyID,
		UserReleaseDate:       time.Time(req.UserReleaseDate),
		TourReleaseDate:       time.Time(req.TourReleaseDate),
		UserPrice:             req.UserPrice,
		AgencyPrice:           req.AgencyPrice,
		PathID:                req.PathID,
		MinPassengers:         req.MinPassengers,
		TechTeamID:            req.TechTeamID,
		MaxTickets:            req.MaxTickets,
		StartDate:             startDate,
		EndDate:               endDate,
		TripCancellingPenalty: penalty,
		SoldTickets:           req.SoldTickets,
	}
}

func TripToCreateTripRes(t *trip.Trip) *CreateTripRes {
	p := TripCancelingPenaltyToTripCancellingPenaltyReq(t.TripCancellingPenalty)
	var startDate, endDate Timestamp
	if t.StartDate != nil {
		startDate = Timestamp(*t.StartDate)
	}
	if t.EndDate != nil {
		endDate = Timestamp(*t.EndDate)
	}

	return &CreateTripRes{
		ID:                    t.ID,
		TransportCompanyID:    t.TransportCompanyID,
		UserReleaseDate:       Timestamp(t.UserReleaseDate),
		TourReleaseDate:       Timestamp(t.TourReleaseDate),
		UserPrice:             t.UserPrice,
		AgencyPrice:           t.AgencyPrice,
		PathID:                t.PathID,
		MinPassengers:         t.MinPassengers,
		TechTeamID:            t.TechTeamID,
		MaxTickets:            t.MaxTickets,
		StartDate:             startDate,
		EndDate:               endDate,
		TripCancellingPenalty: p,
		SoldTickets:           t.SoldTickets,
		Origin:                t.Path.FromTerminal.City,
		Destination:           t.Path.ToTerminal.City,
		FromTerminalName:      t.Path.FromTerminal.Name,
		ToTerminalName:        t.Path.ToTerminal.Name,
		PathName:              t.Path.Name,
		Type:                  t.Path.Type,
	}
}

type OwnerAdminTechTeamOperatorTripResponse struct {
	ID                     uint                           `json:"id"`
	TransportCompanyID     uint                           `json:"company_id"`
	TransportCompanyName   string                         `json:"company_name"`
	TripType               string                         `json:"type"`
	UserReleaseDate        Timestamp                      `json:"user_release"`
	TourReleaseDate        Timestamp                      `json:"tour_release"`
	UserPrice              float64                        `json:"user_price"`
	AgencyPrice            float64                        `json:"agency_price"`
	PathID                 uint                           `json:"path_id"`
	PathName               string                         `json:"path_name"`
	FromTerminalName       string                         `json:"from_terminal"`
	Origin                 string                         `json:"from"`
	Destination            string                         `json:"to"`
	ToTerminalName         string                         `json:"to_terminal"`
	Status                 string                         `json:"status"`
	MinPassengers          uint                           `json:"min_pass"`
	SoldTickets            uint                           `json:"soled_tickets"`
	TechTeamID             *uint                          `json:"tech_id"`
	TripCancellingPenalty  *CreateTripCancelingPenaltyRes `json:"penalty"`
	TripCancelingPenaltyID *uint                          `json:"penalty_id"`
	MaxTickets             uint                           `json:"max_tickets"`
	VehicleID              *uint                          `json:"vehicle_id"`
	VehicleRequestID       *uint                          `json:"vehicle_req_id"`
	IsCanceled             bool                           `json:"is_canceled"`
	IsFinished             bool                           `json:"is_finished"`
	StartDate              Timestamp                      `json:"start_date"`
	EndDate                Timestamp                      `json:"end_date"`
}

func TripToOwnerAdminTechTeamOperatorTripResponse(t trip.Trip) OwnerAdminTechTeamOperatorTripResponse {
	// check ID is owner TODO:
	p := TripCancelingPenaltyToTripCancellingPenaltyRes(t.TripCancellingPenalty)
	var startDate, endDate Timestamp
	if t.StartDate != nil {
		startDate = Timestamp(*t.StartDate)
	}
	if t.EndDate != nil {
		endDate = Timestamp(*t.EndDate)
	}

	return OwnerAdminTechTeamOperatorTripResponse{
		ID:                     t.ID,
		TransportCompanyID:     t.TransportCompanyID,
		TransportCompanyName:   t.TransportCompany.Name,
		UserReleaseDate:        Timestamp(t.UserReleaseDate),
		TourReleaseDate:        Timestamp(t.TourReleaseDate),
		UserPrice:              t.UserPrice,
		AgencyPrice:            t.AgencyPrice,
		PathID:                 t.PathID,
		MinPassengers:          t.MinPassengers,
		TechTeamID:             t.TechTeamID,
		MaxTickets:             t.MaxTickets,
		StartDate:              startDate,
		EndDate:                endDate,
		TripCancellingPenalty:  p,
		SoldTickets:            t.SoldTickets,
		Origin:                 t.Path.FromTerminal.City,
		Destination:            t.Path.ToTerminal.City,
		FromTerminalName:       t.Path.FromTerminal.Name,
		ToTerminalName:         t.Path.ToTerminal.Name,
		PathName:               t.Path.Name,
		TripType:               string(t.TripType),
		Status:                 t.Status,
		TripCancelingPenaltyID: t.TripCancelingPenaltyID,
		VehicleID:              t.VehicleID,
		VehicleRequestID:       t.VehicleRequestID,
		IsCanceled:             t.IsCanceled,
		IsFinished:             t.IsFinished,
		//tech team name
	}
}

func BatchTripToOwnerAdminTechTeamOperatorTripResponse(trips []trip.Trip) []OwnerAdminTechTeamOperatorTripResponse {
	return fp.Map(trips, TripToOwnerAdminTechTeamOperatorTripResponse)

}

type UserTripResponse struct {
	ID                    uint                           `json:"id"`
	TripType              string                         `json:"trip_type"`
	UserPrice             float64                        `json:"user_price"`
	PathName              string                         `json:"path_name"`
	FromTerminalName      string                         `json:"from_terminal"`
	Origin                string                         `json:"from"`
	Destination           string                         `json:"to"`
	ToTerminalName        string                         `json:"to_terminal"`
	TripCancellingPenalty *CreateTripCancelingPenaltyRes `json:"penalty"`
	StartDate             Timestamp                      `json:"start_date"`
	EndDate               Timestamp                      `json:"end_date"`
}

func TripToUserTripResponse(t trip.Trip) UserTripResponse {
	// check ID is owner TODO:
	p := TripCancelingPenaltyToTripCancellingPenaltyRes(t.TripCancellingPenalty)
	var startDate, endDate Timestamp
	if t.StartDate != nil {
		startDate = Timestamp(*t.StartDate)
	}
	if t.EndDate != nil {
		endDate = Timestamp(*t.EndDate)
	}

	return UserTripResponse{
		ID:                    t.ID,
		UserPrice:             t.UserPrice,
		StartDate:             startDate,
		EndDate:               endDate,
		TripCancellingPenalty: p,
		Origin:                t.Path.FromTerminal.City,
		Destination:           t.Path.ToTerminal.City,
		FromTerminalName:      t.Path.FromTerminal.Name,
		ToTerminalName:        t.Path.ToTerminal.Name,
		PathName:              t.Path.Name,
		TripType:              t.Path.Type,
	}
}

func BatchTripToUserTripResponse(trips []trip.Trip) []UserTripResponse {
	return fp.Map(trips, TripToUserTripResponse)
}
