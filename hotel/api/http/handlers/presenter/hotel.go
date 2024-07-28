package presenter

import (
	"hotel/internal/hotel"
	"hotel/internal/room"
	"hotel/pkg/fp"
)

type CreateHotelReq struct {
	OwnerID uint   `json:"owner_id" example:"12"`
	Name    string `json:"name" example:"myhotel"`
	City    string `json:"city" example:"Los Angles"`
	Country string `json:"country" example:"United States America"`
	Details string `json:"details" example:"5 Star Beach Palm Hotel"`
}

type HotelResp struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type FullHotelResponse struct {
	ID        uint        `json:"hotel_id" example:"12"`
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

func CreateHotelRequest(samplehotel *CreateHotelReq) *hotel.Hotel {
	h := &hotel.Hotel{
		OwnerID: samplehotel.OwnerID,
		Name:    samplehotel.Name,
		City:    samplehotel.City,
		Country: samplehotel.Country,
		Details: samplehotel.Details,
	}
	return h
}

func BatchRoomToHotelResp(r []room.Room) []HotelRoomResp {
	return fp.Map(r, roomToHotelRoomResp)
}

type CreateHotelResponse struct {
	ID      uint            `json:"hotel_id"`
	Name    string          `json:"name"`
	City    string          `json:"city"`
	Country string          `json:"country"`
	Details string          `json:"details"`
	Rooms   []HotelRoomResp `json:"rooms"`
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
		Name:    h.Name,
		City:    h.City,
		Country: h.Country,
		Details: h.Details,
		Rooms:   rooms,
	}
}

func HotelToFullHotelResponse(h *hotel.Hotel) *FullHotelResponse {
	return &FullHotelResponse{
		ID:        h.ID,
		Name:      h.Name,
		City:      h.City,
		Country:   h.Country,
		Details:   h.Details,
		IsBlocked: h.IsBlocked,
		Rooms:     h.Rooms,
	}
}
