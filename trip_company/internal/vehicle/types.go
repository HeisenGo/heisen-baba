package vehicle

import "github.com/google/uuid"

type FullVehicleResponse struct {
	ID                    uint      `json:"id" example:"12"`
	Name                  string    `json:"name" example:"Volvo RG"`
	OwnerID               uuid.UUID `json:"owner_id" example:"aba3b3ed-e3d8-4403-9751-1f04287c9d65"`
	PricePerHour          float64   `json:"priceperhour" example:"250"`
	MotorNumber           string    `json:"motornumber" example:"aba3b-e3d8-4403-9751-1f04287c9d65"`
	MinRequiredTechPerson uint      `json:"min_required_tech_person" example:"4"`
	IsActive              bool      `json:"is_active" example:"true"`
	Capacity              uint      `json:"capacity" example:"25"`
	IsBlocked             bool      `json:"is_blocked" example:"false"`
	Type                  string    `json:"type" example:"road"` // rail, road, air, sailing
	Speed                 float64   `json:"speed" example:"60"`
	ProductionYear        uint      `json:"productionyear" example:"1998"`
	IsConfirmedByAdmin    bool      `json:"is_confirmed_by_admin" example:"true"`
}
