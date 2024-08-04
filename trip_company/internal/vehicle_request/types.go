package vehiclerequest

import "context"

type Repo interface {
	Insert(ctx context.Context, vR *VehicleRequest) error
	UpdateVehicleReq(ctx context.Context, id uint, updates map[string]interface{}) error
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
}
