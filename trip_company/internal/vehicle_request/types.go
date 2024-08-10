package vehiclerequest

import (
	"context"
	"errors"
	"tripcompanyservice/internal/company"
)

var(
	ErrDeleteVehicleReq = errors.New("error deleting vr")
	ErrNotFound = errors.New("not found vr")
)

type Repo interface {
	Insert(ctx context.Context, vR *VehicleRequest) error
	UpdateVehicleReq(ctx context.Context, id uint, updates map[string]interface{}) error
	Delete(ctx context.Context, vRID uint) error 
	GetByID(ctx context.Context, id uint) (*VehicleRequest, error) 
	GetTransportCompanyByVehicleRequestID(ctx context.Context, vehicleRequestID uint) (*company.TransportCompany, error) 
}

type VehicleRequest struct {
	ID                    uint
	TripID                uint
	VehicleType           string
	MinCapacity           int
	ProductionYearMin     int
	Status                string
	MatchedVehicleID      uint
	VehicleReservationFee float64
	VehicleProductionYear int
	VehicleName           string
	MatchVehicleSpeed float64
	MinCost            float64
}
