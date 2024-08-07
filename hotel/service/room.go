package service

import (
	"context"
	"errors"
	"hotel/internal/hotel"
	"hotel/internal/invoice"
	"hotel/internal/reservation"
	"hotel/internal/room"
	"hotel/pkg/ports/clients/clients"
	"time"
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
	invoiceOps     *invoice.Ops
	bankClient     clients.IBankClient
}

func NewRoomService(roomOps *room.Ops, reservationOps *reservation.Ops, hotelOps *hotel.Ops, invoiceOps *invoice.Ops, bankClient clients.IBankClient) *RoomService {
	return &RoomService{
		roomOps:        roomOps,
		reservationOps: reservationOps,
		hotelOps:       hotelOps,
		invoiceOps:     invoiceOps,
		bankClient:     bankClient,
	}
}

func (s *RoomService) CreateRoomReservation(ctx context.Context, res *reservation.Reservation) error {
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
	err = s.reservationOps.Create(ctx, res)
	if err != nil {
		return err
	}
	inv := &invoice.Invoice{
		ReservationID: res.ID,
		IssueDate:     time.Now(),
		Amount:        res.TotalPrice,
		UserID:        res.UserID,
		OwnerID:       res.OwnerID,
		Paid:          false,
	}

	err = s.invoiceOps.Create(ctx, inv)
	if err != nil {
		return err
	}
	isSuccess, err := s.bankClient.Transfer(inv.UserID.String(), inv.OwnerID.String(), false, inv.Amount)
	if err != nil {
		return err
	}
	if isSuccess {
		inv.Paid = true
	}
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
