package storage

import (
	"context"
	"errors"
	"vehicle/internal/vehicle"
	"vehicle/pkg/adapters/storage/entities"
	"vehicle/pkg/adapters/storage/mappers"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type vehicleRepo struct {
	db *gorm.DB
}

func NewVehicleRepo(db *gorm.DB) vehicle.Repo {
	return &vehicleRepo{
		db: db,
	}
}

func (r *vehicleRepo) CreateVehicle(ctx context.Context, v *vehicle.Vehicle) error {
	vehicleEntity := mappers.VehicleDomainToEntity(v)
	if err := r.db.WithContext(ctx).Create(&vehicleEntity).Error; err != nil {
		return err
	}
	v.ID = vehicleEntity.ID
	return nil
}

func (r *vehicleRepo) GetVehicles(ctx context.Context, vehicleType string, capacity uint, page, pageSize int) ([]vehicle.Vehicle, uint, error) {
	var vehicles []entities.Vehicle
	var int64Total int64

	query := r.db.Model(&entities.Vehicle{})

	// Filters
	if vehicleType != "" {
		query = query.Where("type = ?", vehicleType)
	}
	if capacity > 0 {
		query = query.Where("capacity >= ?", capacity)
	}

	// Count total records for pagination
	query.Count(&int64Total)

	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	if err := query.Find(&vehicles).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, nil
		}
		return nil, 0, err
	}

	total := uint(int64Total)
	domainVehicles := mappers.BatchVehicleEntitiesToDomain(vehicles)
	return domainVehicles, total, nil
}

func (r *vehicleRepo) GetVehiclesByOwnerID(ctx context.Context, ownerID uuid.UUID, page, pageSize int) ([]vehicle.Vehicle, int, error) {
	var vehicleEntities []entities.Vehicle
	var total int64

	query := r.db.WithContext(ctx).Model(&entities.Vehicle{}).Where("owner_id = ?", ownerID)

	query.Count(&total)

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&vehicleEntities).Error; err != nil {
		return nil, 0, err
	}

	domainVehicles := mappers.BatchVehicleEntitiesToDomain(vehicleEntities)
	return domainVehicles, int(total), nil
}

func (r *vehicleRepo) GetVehicleByID(ctx context.Context, id uint) (*vehicle.Vehicle, error) {
	var vehicleEntity entities.Vehicle
	if err := r.db.WithContext(ctx).First(&vehicleEntity, id).Error; err != nil {
		return nil, err
	}
	domainVehicle := mappers.VehicleEntityToDomain(vehicleEntity)
	return &domainVehicle, nil
}

func (r *vehicleRepo) UpdateVehicle(ctx context.Context, v *vehicle.Vehicle) error {
	vehicleEntity := mappers.VehicleDomainToEntity(v)
	if err := r.db.WithContext(ctx).Model(&entities.Vehicle{}).Where("id = ?", v.ID).Updates(vehicleEntity).Error; err != nil {
		return err
	}
	return nil
}

func (r *vehicleRepo) DeleteVehicle(ctx context.Context, id uint) error {
	var vehicleEntity entities.Vehicle
	if err := r.db.WithContext(ctx).First(&vehicleEntity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return vehicle.ErrRecordNotFound
		}
		return err
	}

	if err := r.db.WithContext(ctx).Delete(&vehicleEntity).Error; err != nil {
		return err
	}
	return nil
}