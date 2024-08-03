package presenter

import (
	"vehicle/internal/vehicle"
	"vehicle/pkg/fp"

	"github.com/google/uuid"
)

type CreateVehicleReq struct {
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

type VehicleResp struct {
	ID   uint   `json:"id"`
	Type string `json:"type"`
}

type FullVehicleResponse struct {
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

type UpdateVehicleReq struct {
	Name                  string    `json:"name,omitempty" example:"Volvo RG"`
	OwnerID               uuid.UUID `json:"owner_id,omitempty" example:"aba3b3ed-e3d8-4403-9751-1f04287c9d65"`
	PricePerHour          float64   `json:"priceperhour,omitempty" example:"250"`
	MotorNumber           string    `json:"motornumber,omitempty" example:"aba3b-e3d8-4403-9751-1f04287c9d65"`
	MinRequiredTechPerson uint      `json:"min_required_tech_person,omitempty" example:"4"`
	IsActive              bool      `json:"is_active,omitempty" example:"true"`
	Capacity              uint      `json:"capacity,omitempty" example:"25"`
	IsBlocked             bool      `json:"is_blocked,omitempty" example:"false"`
	Type                  string    `json:"type,omitempty" example:"road"` // rail, road, air, sailing
	Speed                 float64   `json:"speed,omitempty" example:"60"`
	ProductionYear        uint      `json:"productionyear,omitempty" example:"1998"`
	IsConfirmedByAdmin    bool      `json:"is_confirmed_by_admin,omitempty" example:"true"`
}

type CreateVehicleResponse struct {
	Name                  string    `json:"name"`
	OwnerID               uuid.UUID `json:"owner_id"`
	PricePerHour          float64   `json:"priceperhour"`
	MotorNumber           string    `json:"motornumber"`
	MinRequiredTechPerson uint      `json:"min_required_tech_person"`
	IsActive              bool      `json:"is_active"`
	Capacity              uint      `json:"capacity"`
	IsBlocked             bool      `json:"is_blocked"`
	Type                  string    `json:"type"` // rail, road, air, sailing
	Speed                 float64   `json:"speed"`
	ProductionYear        uint      `json:"productionyear"`
	IsConfirmedByAdmin    bool      `json:"is_confirmed_by_admin"`
}

func CreateVehicleRequest(req *CreateVehicleReq) *vehicle.Vehicle {
	return &vehicle.Vehicle{
		Name:                  req.Name,
		OwnerID:               req.OwnerID,
		PricePerHour:          req.PricePerHour,
		MotorNumber:           req.MotorNumber,
		MinRequiredTechPerson: req.MinRequiredTechPerson,
		IsActive:              req.IsActive,
		Capacity:              req.Capacity,
		IsBlocked:             req.IsBlocked,
		Type:                  req.Type,
		Speed:                 req.Speed,
		ProductionYear:        req.ProductionYear,
		IsConfirmedByAdmin:    req.IsConfirmedByAdmin,
	}
}

func VehicleToCreateVehicleResponse(v *vehicle.Vehicle) *CreateVehicleResponse {
	return &CreateVehicleResponse{
		Name:                  v.Name,
		OwnerID:               v.OwnerID,
		PricePerHour:          v.PricePerHour,
		MotorNumber:           v.MotorNumber,
		MinRequiredTechPerson: v.MinRequiredTechPerson,
		IsActive:              v.IsActive,
		Capacity:              v.Capacity,
		IsBlocked:             v.IsBlocked,
		Type:                  v.Type,
		Speed:                 v.Speed,
		ProductionYear:        v.ProductionYear,
		IsConfirmedByAdmin:    v.IsConfirmedByAdmin,
	}
}

func VehicleToFullVehicleResponse(v vehicle.Vehicle) FullVehicleResponse {
	return FullVehicleResponse{
		Name:                  v.Name,
		OwnerID:               v.OwnerID,
		PricePerHour:          v.PricePerHour,
		MotorNumber:           v.MotorNumber,
		MinRequiredTechPerson: v.MinRequiredTechPerson,
		IsActive:              v.IsActive,
		Capacity:              v.Capacity,
		IsBlocked:             v.IsBlocked,
		Type:                  v.Type,
		Speed:                 v.Speed,
		ProductionYear:        v.ProductionYear,
		IsConfirmedByAdmin:    v.IsConfirmedByAdmin,
	}
}

func BatchVehiclesToVehicleResponse(vehicles []vehicle.Vehicle) []FullVehicleResponse {
	return fp.Map(vehicles, VehicleToFullVehicleResponse)
}

func UpdateVehicleRequestToDomain(req *UpdateVehicleReq) *vehicle.Vehicle {
	v := &vehicle.Vehicle{
		Name:                  req.Name,
		OwnerID:               req.OwnerID,
		PricePerHour:          req.PricePerHour,
		MotorNumber:           req.MotorNumber,
		MinRequiredTechPerson: req.MinRequiredTechPerson,
		IsActive:              req.IsActive,
		Capacity:              req.Capacity,
		IsBlocked:             req.IsBlocked,
		Type:                  req.Type,
		Speed:                 req.Speed,
		ProductionYear:        req.ProductionYear,
		IsConfirmedByAdmin:    req.IsConfirmedByAdmin,
	}
	return v
}
