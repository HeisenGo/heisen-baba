package presenter

import (
	"hotel/internal/reservation"
	"time"

	"github.com/google/uuid"
)

type ReservationCreateReq struct {
	HotelID    uint      `json:"hotel_id" validate:"required" example:"3"`
	RoomID     uint      `json:"room_id" validate:"required" example:"1"`
	UserID     uuid.UUID `json:"user_id" validate:"required" example:"aba3b3ed-e3d8-4403-9751-1f04287c9d65"`
	CheckIn    time.Time `json:"check_in" validate:"required" example:"2024-08-01T00:00:00Z"`
	CheckOut   time.Time `json:"check_out" validate:"required" example:"2024-08-05T00:00:00Z"`
}

type ReservationResp struct {
	ID         uint      `json:"id"`
	HotelID    uint      `json:"hotel_id"`
	RoomID     uint      `json:"room_id"`
	UserID     uuid.UUID `json:"user_id"`
	CheckIn    time.Time `json:"check_in"`
	CheckOut   time.Time `json:"check_out"`
	TotalPrice uint64    `json:"total_price"`
	Status     string    `json:"status"`
}

type FullReservationResponse struct {
	ID         uint      `json:"reservation_id" example:"12"`
	HotelID    uint      `json:"hotel_id" validate:"required" example:"3"`
	RoomID     uint      `json:"room_id" example:"1"`
	UserID     uuid.UUID `json:"user_id" example:"aba3b3ed-e3d8-4403-9751-1f04287c9d65"`
	CheckIn    time.Time `json:"check_in" example:"2024-08-01T00:00:00Z"`
	CheckOut   time.Time `json:"check_out" example:"2024-08-05T00:00:00Z"`
	TotalPrice uint64    `json:"total_price" example:"50000"`
	Status     string    `json:"status" example:"booked"`
}

func ReservationReqToReservationDomain(req *ReservationCreateReq) *reservation.Reservation {
	return &reservation.Reservation{
		RoomID:     req.RoomID,
		UserID:     req.UserID,
		CheckIn:    req.CheckIn,
		CheckOut:   req.CheckOut,
	}
}

func ReservationToReservationResp(r *reservation.Reservation) *ReservationResp {
	return &ReservationResp{
		ID:         r.ID,
		RoomID:     r.RoomID,
		UserID:     r.UserID,
		CheckIn:    r.CheckIn,
		CheckOut:   r.CheckOut,
		TotalPrice: r.TotalPrice,
		Status:     r.Status,
	}
}

func BatchReservationsToReservationResponse(reservations []reservation.Reservation) []ReservationResp {
	var responses []ReservationResp
	for _, r := range reservations {
		responses = append(responses, *ReservationToReservationResp(&r))
	}
	return responses
}

func ReservationToFullReservationResponse(r *reservation.Reservation) *FullReservationResponse {
	return &FullReservationResponse{
		ID:         r.ID,
		RoomID:     r.RoomID,
		UserID:     r.UserID,
		CheckIn:    r.CheckIn,
		CheckOut:   r.CheckOut,
		TotalPrice: r.TotalPrice,
		Status:     r.Status,
	}
}
