package clients

import (
	"tripcompanyservice/internal/vehicle"
	vehiclerequest "tripcompanyservice/internal/vehicle_request"
)

type IVehicleClient interface {
	SelectVehicles(vr *vehiclerequest.VehicleRequest) (*vehicle.FullVehicleResponse, error)
}