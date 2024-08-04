package mappers

import (
	vehiclerequest "tripcompanyservice/internal/vehicle_request"
	"tripcompanyservice/pkg/adapters/storage/entities"
)

func VehicleDomainToVehicleEntity(r *vehiclerequest.VehicleRequest) *entities.VehicleRequest {
	return &entities.VehicleRequest{

		TripID:            r.TripID,
		MinCapacity:       r.MinCapacity,
		ProductionYearMin: r.ProductionYearMin,
	}
}

func VehicleReqEntityToVehicleReqDomain(v entities.VehicleRequest) vehiclerequest.VehicleRequest {
	return vehiclerequest.VehicleRequest{
		ID:                    v.ID,
		TripID:                v.TripID,
		VehicleType:           v.VehicleType,
		VehicleName:           v.VehicleName,
		VehicleReservationFee: v.VehicleReservationFee,
		VehicleProductionYear: v.VehicleProductionYear,
		MatchedVehicleID:      v.MatchedVehicleID,
		MinCapacity:           v.MinCapacity,
		ProductionYearMin:     v.ProductionYearMin,
		MatchVehicleSpeed: v.MatchVehicleSpeed,
	}
}
