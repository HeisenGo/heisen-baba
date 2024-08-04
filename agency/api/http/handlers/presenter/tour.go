package presenter

import (
	"agency/internal/tour"
	"agency/pkg/fp"
	"time"
)

type CreateTourReq struct {
	AgencyID     uint      `json:"agency_id" validate:"required" example:"1"`
	GoTicketID   uint      `json:"go_ticket_id" validate:"required" example:"1"`
	BackTicketID uint      `json:"back_ticket_id" validate:"required" example:"2"`
	HotelID      uint      `json:"hotel_id" validate:"required" example:"1"`
	Capacity     uint      `json:"capacity" validate:"required" example:"20"`
	ReleaseDate  time.Time `json:"release_date" validate:"required" example:"2024-12-01T15:04:05Z"`
	IsApproved   bool      `json:"is_approved" example:"true"`
	IsActive     bool      `json:"is_active" example:"true"`
}

type TourResp struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type FullTourResponse struct {
	ID           uint      `json:"tour_id" example:"12"`
	AgencyID     uint      `json:"agency_id" example:"1"`
	GoTicketID   uint      `json:"go_ticket_id" example:"1"`
	BackTicketID uint      `json:"back_ticket_id" example:"2"`
	HotelID      uint      `json:"hotel_id" example:"1"`
	Capacity     uint      `json:"capacity" example:"20"`
	ReleaseDate  time.Time `json:"release_date" example:"2024-12-01T15:04:05Z"`
	IsApproved   bool      `json:"is_approved" example:"true"`
	IsActive     bool      `json:"is_active" example:"true"`
}

type UpdateTourReq struct {
	AgencyID     *uint      `json:"agency_id" example:"1"`
	GoTicketID   *uint      `json:"go_ticket_id" example:"1"`
	BackTicketID *uint      `json:"back_ticket_id" example:"2"`
	HotelID      *uint      `json:"hotel_id" example:"1"`
	Capacity     *uint      `json:"capacity" example:"20"`
	ReleaseDate  *time.Time `json:"release_date" example:"2024-12-01T15:04:05Z"`
	IsApproved   *bool      `json:"is_approved" example:"true"`
	IsActive     *bool      `json:"is_active" example:"true"`
}

func CreateTourRequest(req *CreateTourReq) *tour.Tour {
	return &tour.Tour{
		AgencyID:     req.AgencyID,
		GoTicketID:   req.GoTicketID,
		BackTicketID: req.BackTicketID,
		HotelID:      req.HotelID,
		Capacity:     req.Capacity,
		ReleaseDate:  req.ReleaseDate,
		IsApproved:   req.IsApproved,
		IsActive:     req.IsActive,
	}
}
func TourToCreateTourResponse(t *tour.Tour) *FullTourResponse {
	return &FullTourResponse{
		AgencyID:     t.AgencyID,
		GoTicketID:   t.GoTicketID,
		BackTicketID: t.BackTicketID,
		HotelID:      t.HotelID,
		Capacity:     t.Capacity,
		ReleaseDate:  t.ReleaseDate,
		IsApproved:   t.IsApproved,
		IsActive:     t.IsActive,
	}
}
func TourToFullTourResponse(t tour.Tour) FullTourResponse {
	return FullTourResponse{
		ID:           t.ID,
		AgencyID:     t.AgencyID,
		GoTicketID:   t.GoTicketID,
		BackTicketID: t.BackTicketID,
		HotelID:      t.HotelID,
		Capacity:     t.Capacity,
		ReleaseDate:  t.ReleaseDate,
		IsApproved:   t.IsApproved,
		IsActive:     t.IsActive,
	}
}

func BatchToursToTourResponse(tours []tour.Tour) []FullTourResponse {
	return fp.Map(tours, TourToFullTourResponse)
}

func UpdateTourRequestToDomain(req *UpdateTourReq) *tour.Tour {
	t := &tour.Tour{}
	if req.AgencyID != nil {
		t.AgencyID = *req.AgencyID
	}
	if req.GoTicketID != nil {
		t.GoTicketID = *req.GoTicketID
	}
	if req.BackTicketID != nil {
		t.BackTicketID = *req.BackTicketID
	}
	if req.HotelID != nil {
		t.HotelID = *req.HotelID
	}
	if req.Capacity != nil {
		t.Capacity = *req.Capacity
	}
	if req.ReleaseDate != nil {
		t.ReleaseDate = *req.ReleaseDate
	}
	if req.IsApproved != nil {
		t.IsApproved = *req.IsApproved
	}
	if req.IsActive != nil {
		t.IsActive = *req.IsActive
	}
	return t
}
