package service

import (
	"tripcompanyservice/internal/trip"
	vehiclerequest "tripcompanyservice/internal/vehicle_request"
)

type VehicleReService struct {
	vehicleReqOps *vehiclerequest.Ops
	tripOps       *trip.Ops
}

func NewVehicleReService(vehicleReqOps *vehiclerequest.Ops, tripOps *trip.Ops) *VehicleReService {
	return &VehicleReService{
		vehicleReqOps: vehicleReqOps,
		tripOps:       tripOps,
	}
}
