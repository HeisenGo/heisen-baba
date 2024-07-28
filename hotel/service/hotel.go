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
		roomOps:  roomOps,
	}
}

func (s *HotelService) CreateHotel(ctx context.Context, h *hotel.Hotel) error {
	return s.hotelOps.Create(ctx, h)
}

func (s *HotelService) GetHotel(ctx context.Context, id uint) (*hotel.Hotel, error) {
	return s.hotelOps.GetByID(ctx, id)
}

func (s *HotelService) UpdateHotel(ctx context.Context, id uint, h *hotel.Hotel) error {
	existingHotel, err := s.hotelOps.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Update fields
	existingHotel.Name = h.Name
	existingHotel.City = h.City
	existingHotel.Country = h.Country
	existingHotel.Details = h.Details
	existingHotel.IsBlocked = h.IsBlocked

	return s.hotelOps.Update(ctx, existingHotel)
}

func (s *HotelService) DeleteHotel(ctx context.Context, id uint) error {
	return s.hotelOps.Delete(ctx, id)
}