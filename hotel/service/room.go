package service

import (
	"context"
	"errors"
	"hotel/internal/hotel"
	"hotel/internal/reservation"
	"hotel/internal/room"
)

var (
	ErrRoomNotFound          = errors.New("room not found")
	ErrReservationNotCreated = errors.New("reservation not created")
	ErrInvoiceNotCreated     = errors.New("invoice not created")
)

type RoomService struct {
	roomOps        *room.Ops
	reservationOps *reservation.Ops
	hotelOps       *hotel.Ops
}

func NewRoomService(roomOps *room.Ops, reservationOps *reservation.Ops, hotelOps *hotel.Ops) *RoomService {
	return &RoomService{
		roomOps:        roomOps,
		reservationOps: reservationOps,
		hotelOps:       hotelOps,
	}
}

func (s *RoomService) CreateReservation(ctx context.Context, res *reservation.Reservation) error {
	// Ensure room exists before creating a reservation
	existingRoom, err := s.roomOps.GetRoomByID(ctx, res.RoomID)
	if err != nil {
		return ErrRoomNotFound
	}
	if existingRoom == nil {
		return ErrRoomNotFound
	}
	res.TotalPrice = existingRoom.UserPrice
	hotel, err := s.hotelOps.GetHotelsByID(ctx, existingRoom.HotelID)
	if err != nil {
		return err
	}
	
	res.OwnerID = hotel.OwnerID
	// Validate and create the reservation
	if err := s.reservationOps.Create(ctx, res); err != nil {
		return ErrReservationNotCreated

	}

	return nil
}

func (s *RoomService) CreateRoom(ctx context.Context, r *room.Room) (*room.Room, error) {
	return s.roomOps.CreateRoom(ctx, r)
}

func (s *RoomService) GetRooms(ctx context.Context, page, pageSize int) ([]room.Room, int, error) {
	return s.roomOps.GetRooms(ctx, page, pageSize)
}

func (s *RoomService) GetRoomByID(ctx context.Context, id uint) (*room.Room, error) {
	return s.roomOps.GetRoomByID(ctx, id)
}

func (s *RoomService) UpdateRoom(ctx context.Context, r *room.Room) error {
	return s.roomOps.UpdateRoom(ctx, r)
}

func (s *RoomService) DeleteRoom(ctx context.Context, id uint) error {
	return s.roomOps.DeleteRoom(ctx, id)
}
