package mappers

import (
	"vehicle/internal/vehicle"
	"vehicle/pkg/adapters/storage/entities"
	"vehicle/pkg/fp"
)

func VehicleEntityToDomain(vehicleEntity entities.Vehicle) vehicle.Vehicle {
	return vehicle.Vehicle{
		ID:                  vehicleEntity.ID,
		Name:                vehicleEntity.Name,
		OwnerID:             vehicleEntity.OwnerID,
		PricePerHour:        vehicleEntity.PricePerHour,
		MotorNumber:         vehicleEntity.MotorNumber,
		MinRequiredTechPerson: vehicleEntity.MinRequiredTechPerson,
		IsActive:            vehicleEntity.IsActive,
		Capacity:            vehicleEntity.Capacity,
		IsBlocked:           vehicleEntity.IsBlocked,
		Type:                vehicleEntity.Type,
		Speed:               vehicleEntity.Speed,
		ProductionYear:      vehicleEntity.ProductionYear,
		IsConfirmedByAdmin:  vehicleEntity.IsConfirmedByAdmin,
	}
}

func BatchVehicleEntitiesToDomain(vehicleEntities []entities.Vehicle) []vehicle.Vehicle {
	return fp.Map(vehicleEntities, VehicleEntityToDomain)
}

func VehicleDomainToEntity(v *vehicle.Vehicle) *entities.Vehicle {
	return &entities.Vehicle{
		Name:                v.Name,
		OwnerID:             v.OwnerID,
		PricePerHour:        v.PricePerHour,
		MotorNumber:         v.MotorNumber,
		MinRequiredTechPerson: v.MinRequiredTechPerson,
		IsActive:            v.IsActive,
		Capacity:            v.Capacity,
		IsBlocked:           v.IsBlocked,
		Type:                v.Type,
		Speed:               v.Speed,
		ProductionYear:      v.ProductionYear,
		IsConfirmedByAdmin:  v.IsConfirmedByAdmin,
	}
}