package presenter

import "tripcompanyservice/internal/ticket"

type AgencyTicketReq struct {
	TripID   uint  `json:"trip_id" validate:"required"`
	AgencyID *uint `json:"agency_id" validate:"required"`
	Quantity int   `json:"quantity" validate:"required"`
}

type UserTicketReq struct {
	TripID   uint  `json:"trip_id" validate:"required"`
	UserID   *uint `json:"user_id" validate:"required"`
	Quantity int   `json:"quantity" validate:"required"`
}

func AgencyTicketReqToTicket(t *AgencyTicketReq) *ticket.Ticket {
	return &ticket.Ticket{
		TripID:   t.TripID,
		AgencyID: t.AgencyID,
		Quantity: t.Quantity,
	}
}

func UserTicketReqToAgency(t *UserTicketReq) *ticket.Ticket {
	return &ticket.Ticket{
		TripID:   t.TripID,
		UserID:   t.UserID,
		Quantity: t.Quantity,
	}
}

type AgencyTicket struct {
	ID         uint               `json:"id"`
	TripID     uint               `json:"trip_id"`
	Trip       AgencyTripResponse `json:"trip"`
	AgencyID   *uint              `json:"agency_id"`
	Quantity   int                `json:"quantity"`
	TotalPrice float64            `json:"total_price"`
	Status     string             `json:"status"`
}

type UserTicket struct {
	ID         uint             `json:"id"`
	TripID     uint             `json:"trip_id"`
	Trip       UserTripResponse `json:"trip"`
	UserID     *uint            `json:"user_id"`
	Quantity   int              `json:"quantity"`
	TotalPrice float64          `json:"total_price"`
	Status     string           `json:"status"`
}

func TicketToAgencyTicket(t ticket.Ticket) AgencyTicket {
	trip := TripToAgencyTripResponse(t.Trip)
	return AgencyTicket{
		ID:         t.ID,
		TripID:     t.TripID,
		Trip:       trip,
		AgencyID:   t.AgencyID,
		Quantity:   t.Quantity,
		TotalPrice: t.TotalPrice,
		Status:     t.Status,
	}
}

func TicketToUserTicket(t ticket.Ticket) UserTicket {
	trip := TripToUserTripResponse(t.Trip)
	return UserTicket{
		ID:         t.ID,
		TripID:     t.TripID,
		Trip:       trip,
		UserID:     t.UserID,
		Quantity:   t.Quantity,
		TotalPrice: t.TotalPrice,
		Status:     t.Status,
	}
}
