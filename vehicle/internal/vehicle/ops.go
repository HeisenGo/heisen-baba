package vehicle

import (
	"context"

	"github.com/google/uuid"
)

type Ops struct {
	repo Repo
}

func NewOps(repo Repo) *Ops {
	return &Ops{repo}
}

func (o *Ops) Create(ctx context.Context, vehicle *Vehicle) error {
	if err := ValidateVehicleName(vehicle.Name); err != nil {
		return ErrInvalidVehicleName
	}

	if err := ValidateType(vehicle.Type); err != nil {
		return ErrInvalidType
	}

	if err := ValidateProductionYear(vehicle.ProductionYear); err != nil {
		return ErrInvalidProductionYear
	}

	if err := ValidateMotorNumber(vehicle.MotorNumber); err != nil {
		return ErrInvalidMotorNumber
	}

	if err := ValidateCapacity(vehicle.Capacity); err != nil {
		return ErrInvalidCapacity
	}
	return o.repo.CreateVehicle(ctx, vehicle)
}

func (o *Ops) GetVehicleByID(ctx context.Context, id uint) (*Vehicle, error) {
	return o.repo.GetVehicleByID(ctx, id)
}

func (o *Ops) GetVehicles(ctx context.Context, filters VehicleFilters, page, pageSize int) ([]Vehicle, uint, error) {
	return o.repo.GetVehicles(ctx, filters, page, pageSize)
}

func (o *Ops) GetVehiclesByOwnerID(ctx context.Context, ownerID uuid.UUID, page, pageSize int) ([]Vehicle, int, error) {
	filters := VehicleFilters{OwnerID: ownerID}
	vehicles, total, err := o.repo.GetVehicles(ctx, filters, page, pageSize)
	return vehicles, int(total), err
}

func (o *Ops) Update(ctx context.Context, vehicle *Vehicle) error {
	// Ensure vehicle exists before updating
	existingVehicle, err := o.repo.GetVehicleByID(ctx, vehicle.ID)
	if err != nil {
		return err
	}
	if existingVehicle == nil {
		return ErrRecordNotFound
	}

	if err := ValidateVehicleName(vehicle.Name); err != nil {
		return ErrInvalidVehicleName
	}

	if err := ValidateType(vehicle.Type); err != nil {
		return ErrInvalidType
	}

	if err := ValidateProductionYear(vehicle.ProductionYear); err != nil {
		return ErrInvalidProductionYear
	}

	if err := ValidateMotorNumber(vehicle.MotorNumber); err != nil {
		return ErrInvalidMotorNumber
	}

	if err := ValidateCapacity(vehicle.Capacity); err != nil {
		return ErrInvalidCapacity
	}
	return o.repo.UpdateVehicle(ctx, vehicle)
}

func (o *Ops) Delete(ctx context.Context, id uint) error {
	// Ensure vehicle exists before deleting
	existingVehicle, err := o.repo.GetVehicleByID(ctx, id)
	if err != nil {
		return err
	}
	if existingVehicle == nil {
		return ErrRecordNotFound
	}
	return o.repo.DeleteVehicle(ctx, id)
}

func (o *Ops) ApproveVehicle(ctx context.Context, id uint) error {
	// Ensure vehicle exists before approving
	existingVehicle, err := o.repo.GetVehicleByID(ctx, id)
	if err != nil {
		return err
	}
	if existingVehicle == nil {
		return ErrRecordNotFound
	}
	return o.repo.ApproveVehicle(ctx, id)
}

func (o *Ops) SetVehicleStatus(ctx context.Context, id uint, isActive bool) error {
	// Ensure vehicle exists before changing status
	existingVehicle, err := o.repo.GetVehicleByID(ctx, id)
	if err != nil {
		return err
	}
	if existingVehicle == nil {
		return ErrRecordNotFound
	}
	return o.repo.SetVehicleStatus(ctx, id, isActive)
}

func (o *Ops) SelectVehicles(ctx context.Context, numPassengers uint, cost float64) ([]Vehicle, error) {
	return o.repo.SelectVehicles(ctx, numPassengers, cost)
}