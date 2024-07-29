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

func (s *HotelService) GetHotels(ctx context.Context, city, country string, capacity, page, pageSize int) ([]hotel.Hotel, uint, error) {
	return s.hotelOps.GetHotels(ctx, city, country, capacity, page, pageSize)
}

func (s *HotelService) GetHotelsByOwnerID(ctx context.Context, ownerID uint, page, pageSize int) ([]hotel.Hotel, int, error) {
	return s.hotelOps.GetHotelsByOwnerID(ctx, ownerID, page, pageSize)
}
func (s *HotelService) UpdateHotel(ctx context.Context, id uint, updates *hotel.Hotel) error {
	existingHotel, err := s.hotelOps.GetHotelsByID(ctx, id)
	if err != nil {
		return err
	}

	// Update only the fields that are provided
	if updates.Name != "" {
		existingHotel.Name = updates.Name
	}
	if updates.City != "" {
		existingHotel.City = updates.City
	}
	if updates.Country != "" {
		existingHotel.Country = updates.Country
	}
	if updates.Details != "" {
		existingHotel.Details = updates.Details
	}
	existingHotel.IsBlocked = updates.IsBlocked

	return s.hotelOps.Update(ctx, existingHotel)
}
func (s *HotelService) DeleteHotel(ctx context.Context, id uint) error {
	// Check if the hotel exists
	existingHotel, err := s.hotelOps.GetHotelsByID(ctx, id)
	if err != nil {
		return err
	}

	// Delete all rooms associated with the hotel
	for _, room := range existingHotel.Rooms {
		if err := s.roomOps.DeleteRoom(ctx, room.ID); err != nil {
			return err
		}
	}

	return s.hotelOps.Delete(ctx, id)
}
