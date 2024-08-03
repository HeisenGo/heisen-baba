package service

import (
	"context"
	"errors"
	"vehicle/internal/vehicle"

	"github.com/google/uuid"
)

var (
	ErrPermissionDenied = errors.New("permission denied")
	ErrOwnerExists      = errors.New("owner already exists")
	ErrAMember          = errors.New("user already is a member")
)

type VehicleService struct {
	vehicleOps *vehicle.Ops
}

func NewVehicleService(vehicleOps *vehicle.Ops) *VehicleService {
	return &VehicleService{
		vehicleOps: vehicleOps,
	}
}

func (s *VehicleService) CreateVehicle(ctx context.Context, v *vehicle.Vehicle) error {
	return s.vehicleOps.Create(ctx, v)
}

func (s *VehicleService) GetVehicles(ctx context.Context, filters vehicle.VehicleFilters, page, pageSize int) ([]vehicle.Vehicle, uint, error) {
	return s.vehicleOps.GetVehicles(ctx, filters, page, pageSize)
}

func (s *VehicleService) GetVehiclesByOwnerID(ctx context.Context, ownerID uuid.UUID, page, pageSize int) ([]vehicle.Vehicle, int, error) {
	return s.vehicleOps.GetVehiclesByOwnerID(ctx, ownerID, page, pageSize)
}

func (s *VehicleService) UpdateVehicle(ctx context.Context, id uint, updates *vehicle.Vehicle) error {
	existingVehicle, err := s.vehicleOps.GetVehicleByID(ctx, id)
	if err != nil {
		return err
	}

	// Update only the fields that are provided
	if updates.Name != "" {
		existingVehicle.Name = updates.Name
	}
	if updates.PricePerHour != 0 {
		existingVehicle.PricePerHour = updates.PricePerHour
	}
	if updates.MotorNumber != "" {
		existingVehicle.MotorNumber = updates.MotorNumber
	}
	if updates.MinRequiredTechPerson != 0 {
		existingVehicle.MinRequiredTechPerson = updates.MinRequiredTechPerson
	}
	existingVehicle.IsActive = updates.IsActive
	if updates.Capacity != 0 {
		existingVehicle.Capacity = updates.Capacity
	}
	existingVehicle.IsBlocked = updates.IsBlocked
	if updates.Type != "" {
		existingVehicle.Type = updates.Type
	}
	if updates.Speed != 0 {
		existingVehicle.Speed = updates.Speed
	}
	if updates.ProductionYear != 0 {
		existingVehicle.ProductionYear = updates.ProductionYear
	}
	existingVehicle.IsConfirmedByAdmin = updates.IsConfirmedByAdmin

	return s.vehicleOps.Update(ctx, existingVehicle)
}

func (s *VehicleService) DeleteVehicle(ctx context.Context, id uint) error {
	// Check if the vehicle exists
	_, err := s.vehicleOps.GetVehicleByID(ctx, id)
	if err != nil {
		return err
	}

	return s.vehicleOps.Delete(ctx, id)
}

func (s *VehicleService) ApproveVehicle(ctx context.Context, id uint) error {
	return s.vehicleOps.ApproveVehicle(ctx, id)
}

func (s *VehicleService) SetVehicleStatus(ctx context.Context, id uint, isActive bool) error {
	return s.vehicleOps.SetVehicleStatus(ctx, id, isActive)
}

func (s *VehicleService) SelectVehicles(ctx context.Context, numPassengers uint, cost float64) ([]vehicle.Vehicle, error) {
	return s.vehicleOps.SelectVehicles(ctx, numPassengers, cost)
}
