package service

import (
	"context"
	"errors"
	"hotel/internal/hotel"
	"hotel/internal/room"
)

var (
	ErrPermissionDenied = errors.New("permission denied")
	ErrOwnerExists      = errors.New("owner already exists")
	ErrAMember          = errors.New("user already is a member")
)

type HotelService struct {
	hotelOps *hotel.Ops
	roomOps  *room.Ops
}

func NewHotelService(hotelOps *hotel.Ops, roomOps *room.Ops) *HotelService {
	return &HotelService{
		hotelOps: hotelOps,
		roomOps:         roomOps,}
}

func (s *HotelService) CreateHotel(ctx context.Context, h *hotel.Hotel) error {

	err := s.hotelOps.Create(ctx, h)
	if err != nil {
		return err
	}
	return nil
}
