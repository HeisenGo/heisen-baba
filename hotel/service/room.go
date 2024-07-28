package service

import (
	"context"
	"hotel/internal/hotel"
	"hotel/internal/room"
	"hotel/pkg/adapters/storage/entities"
	"hotel/pkg/adapters/storage/mappers"
)


type RoomService struct {
	roomOps           *room.Ops
	hotelOps         *hotel.Ops
}

func NewRoomService(hotelOps *hotel.Ops, roomOps *room.Ops) *RoomService {
	return &RoomService{
		hotelOps: hotelOps,
		roomOps:         roomOps,
	}
}

func (s *RoomService) CreateRoom(ctx context.Context,hotelID uint,userPrice,agencyPrice uint64,capacity uint8, name,facilities string ,isAvailable bool) (*entities.Room, error) {
	room := &room.Room{
		Name:     name,
		HotelID:  hotelID,
		UserPrice: userPrice,
		AgencyPrice: agencyPrice,
		Facilities: facilities,
		Capacity: capacity,
		IsAvailable: isAvailable,
	}

	createdRoom,err := s.roomOps.CreateRoom(ctx, room)
	if err != nil {
		return nil , err
	}
	roomentity := mappers.RoomDomainToEntity(createdRoom)
	return roomentity ,nil
}