package presenter

import "hotel/internal/room"

type CreateRoomReq struct {
	HotelID     uint   `json:"hotel_id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	AgencyPrice uint64 `json:"agency_price" validate:"required"`
	UserPrice   uint64 `json:"user_price" validate:"required"`
	Facilities  string `json:"facilities" validate:"required"`
	Capacity    uint8  `json:"capacity" validate:"required"`
	IsAvailable bool   `json:"is_available" validate:"required"`
}

type RoomResp struct {
	ID          uint   `json:"id"`
	HotelID     uint   `json:"hotel_id"`
	Name        string `json:"name"`
	AgencyPrice uint64 `json:"agency_price"`
	UserPrice   uint64 `json:"user_price"`
	Facilities  string `json:"facilities"`
	Capacity    uint8  `json:"capacity"`
	IsAvailable bool   `json:"is_available"`
}

type UpdateRoomReq struct {
	HotelID     *uint   `json:"hotel_id" validate:"required"`
	Name        *string `json:"name" `
	AgencyPrice *uint64 `json:"agency_price" `
	UserPrice   *uint64 `json:"user_price" `
	Facilities  *string `json:"facilities" `
	Capacity    *uint8  `json:"capacity" `
	IsAvailable *bool   `json:"is_available" `
}
func CreateRoomRequest(req *CreateRoomReq) *room.Room {
	return &room.Room{
		HotelID:     req.HotelID,
		Name:        req.Name,
		AgencyPrice: req.AgencyPrice,
		UserPrice:   req.UserPrice,
		Facilities:  req.Facilities,
		Capacity:    req.Capacity,
		IsAvailable: req.IsAvailable,
	}
}

func RoomToCreateRoomResponse(r *room.Room) *RoomResp {
	return &RoomResp{
		ID:          r.ID,
		HotelID:     r.HotelID,
		Name:        r.Name,
		AgencyPrice: r.AgencyPrice,
		UserPrice:   r.UserPrice,
		Facilities:  r.Facilities,
		Capacity:    r.Capacity,
		IsAvailable: r.IsAvailable,
	}
}

func RoomToFullRoomResponse(r *room.Room) *RoomResp {
	return RoomToCreateRoomResponse(r)
}

func UpdateRoomRequestToDomain(req *UpdateRoomReq) (*room.Room,error) {
	r := &room.Room{}
	if req.Name != nil {
		r.Name = *req.Name
	}
	if req.AgencyPrice != nil {
		r.AgencyPrice = *req.AgencyPrice
	}
	if req.UserPrice != nil {
		r.UserPrice = *req.UserPrice
	}
	if req.Capacity != nil {
		r.Capacity = *req.Capacity
	}
	if req.Facilities != nil {
		r.Facilities= *req.Facilities
	}
	if req.IsAvailable != nil {
		r.IsAvailable= *req.IsAvailable
	}
	if req.HotelID == nil {
		return nil, room.ErrEmptyHotelID
	}
	return r, nil
}