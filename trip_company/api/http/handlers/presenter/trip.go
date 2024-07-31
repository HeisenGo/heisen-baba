package presenter

import (
	"time"
	"tripcompanyservice/internal/trip"
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
	StartDate          Timestamp `json:"start_date" validate:"required"` // should be given by trip generator
	// example: "2024-07-30 23:27:09"
	//end date should be calculated according to the vehicle speed and path distance
	TripCancellingPenalty *CreateTripCancelingPenaltyReq `json:"penalty" validate:"required"`
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
	}
}

func TripToCreateTripReq(t *trip.Trip) *CreateTripReq {
	p := TripCancelingPenaltyToTripCancellingPenaltyReq(t.TripCancellingPenalty)

	// Handle nil pointers for StartDate and EndDate
	var startDate, endDate Timestamp
	if t.StartDate != nil {
		startDate = Timestamp(*t.StartDate)
	}
	if t.EndDate != nil {
		endDate = Timestamp(*t.EndDate)
	}

	return &CreateTripReq{
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
	}
}
