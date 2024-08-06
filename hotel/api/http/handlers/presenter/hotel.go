package presenter

import (
	"hotel/internal/hotel"
	"hotel/internal/room"
	"hotel/pkg/fp"

	"github.com/google/uuid"
)

type CreateHotelReq struct {
	Name    string `json:"name" validate:"required" example:"myhotel"`
	City    string `json:"city" validate:"required" example:"Los Angles"`
	Country string `json:"country" validate:"required" example:"United States America"`
	Details string `json:"details" validate:"required" example:"5 Star Beach Palm Hotel"`
}

type HotelResp struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type FullHotelResponse struct {
	ID        uint        `json:"hotel_id" example:"12"`
	OwnerID   uuid.UUID   `json:"owner_id" example:"1"`
	Name      string      `json:"name" example:"myhotel"`
	City      string      `json:"city" example:"Los Angles"`
	Country   string      `json:"country" example:"United States America"`
	Details   string      `json:"details" example:"5 Star Beach Palm Hotel"`
	IsBlocked bool        `json:"is_blocked" example:"false"`
	Rooms     []room.Room `json:"rooms"`
}

type HotelRoomResp struct {
	ID          uint   `json:"hotel_id"`
	Name        string `json:"name"`
	AgencyPrice uint64 `json:"agencyprice"`
	UserPrice   uint64 `json:"userprice"`
	Facilities  string `json:"facilities"`
	Capacity    uint8  `json:"capacity"`
	IsAvailable bool   `json:"is_available"`
}
type UpdateHotelReq struct {
	Name      *string `json:"name" example:"myhotel"`
	City      *string `json:"city" example:"Los Angles"`
	Country   *string `json:"country" example:"United States America"`
	Details   *string `json:"details" example:"5 Star Beach Palm Hotel"`
	IsBlocked *bool   `json:"is_blocked" example:"false"`
}

type CreateHotelResponse struct {
	ID      uint            `json:"hotel_id"`
	OwnerID uuid.UUID       `json:"owner_id"`
	Name    string          `json:"name"`
	City    string          `json:"city"`
	Country string          `json:"country"`
	Details string          `json:"details"`
	Rooms   []HotelRoomResp `json:"rooms"`
}

func CreateHotelRequest(sampleHotel *CreateHotelReq) *hotel.Hotel {
	h := &hotel.Hotel{
		Name:    sampleHotel.Name,
		City:    sampleHotel.City,
		Country: sampleHotel.Country,
		Details: sampleHotel.Details,
	}
	return h
}

func BatchRoomToHotelResp(r []room.Room) []HotelRoomResp {
	return fp.Map(r, roomToHotelRoomResp)
}
func BatchHotelsToHotelResponse(hotels []hotel.Hotel) []FullHotelResponse {
	return fp.Map(hotels, HotelToFullHotelResponse)
}

func roomToHotelRoomResp(r room.Room) HotelRoomResp {
	return HotelRoomResp{
		ID:          r.ID,
		Name:        r.Name,
		AgencyPrice: r.AgencyPrice,
		UserPrice:   r.UserPrice,
		Facilities:  r.Facilities,
		Capacity:    r.Capacity,
		IsAvailable: r.IsAvailable,
	}
}

func HotelToCreateHotelResponse(h *hotel.Hotel) *CreateHotelResponse {
	rooms := BatchRoomToHotelResp(h.Rooms)
	return &CreateHotelResponse{
		ID:      h.ID,
		OwnerID: h.OwnerID,
		Name:    h.Name,
		City:    h.City,
		Country: h.Country,
		Details: h.Details,
		Rooms:   rooms,
	}
}

func HotelToFullHotelResponse(h hotel.Hotel) FullHotelResponse {
	return FullHotelResponse{
		ID:        h.ID,
		OwnerID:   h.OwnerID,
		Name:      h.Name,
		City:      h.City,
		Country:   h.Country,
		Details:   h.Details,
		IsBlocked: h.IsBlocked,
		Rooms:     h.Rooms,
	}
}

func UpdateHotelRequestToDomain(req *UpdateHotelReq) *hotel.Hotel {
	h := &hotel.Hotel{}
	if req.Name != nil {
		h.Name = *req.Name
	}
	if req.City != nil {
		h.City = *req.City
	}
	if req.Country != nil {
		h.Country = *req.Country
	}
	if req.Details != nil {
		h.Details = *req.Details
	}
	if req.IsBlocked != nil {
		h.IsBlocked = *req.IsBlocked
	}
	return h
}
