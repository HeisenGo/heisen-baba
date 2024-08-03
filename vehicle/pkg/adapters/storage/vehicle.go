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

func (r *vehicleRepo) GetVehicleByID(ctx context.Context, id uint) (*vehicle.Vehicle, error) {
	var vehicleEntity entities.Vehicle
	if err := r.db.WithContext(ctx).First(&vehicleEntity, id).Error; err != nil {
		return nil, err
	}
	domainVehicle := mappers.VehicleEntityToDomain(vehicleEntity)
	return &domainVehicle, nil
}

func (r *vehicleRepo) GetVehicles(ctx context.Context, filters vehicle.VehicleFilters, page, pageSize int) ([]vehicle.Vehicle, uint, error) {
	var vehicles []entities.Vehicle
	var int64Total int64

	query := r.db.Model(&entities.Vehicle{})

	// Filters
	if filters.Type != "" {
		query = query.Where("type = ?", filters.Type)
	}
	if filters.OwnerID != uuid.Nil {
		query = query.Where("owner_id = ?", filters.OwnerID)
	}

	// Count total records for pagination
	if err := query.Count(&int64Total).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, nil
		}
		return nil, 0, err
	}

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

func (r *vehicleRepo) ApproveVehicle(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Model(&entities.Vehicle{}).Where("id = ?", id).Update("is_confirmed_by_admin", true).Error; err != nil {
		return err
	}
	return nil
}

func (r *vehicleRepo) SetVehicleStatus(ctx context.Context, id uint, isActive bool) error {
	if err := r.db.WithContext(ctx).Model(&entities.Vehicle{}).Where("id = ?", id).Update("is_active", isActive).Error; err != nil {
		return err
	}
	return nil
}

func (r *vehicleRepo) SelectVehicles(ctx context.Context, numPassengers uint, cost float64) ([]vehicle.Vehicle, error) {
	var vehicles []entities.Vehicle

	query := r.db.Model(&entities.Vehicle{}).Where("is_confirmed_by_admin = ?", true).Where("capacity >= ?", numPassengers).Where("price_per_hour <= ?", cost)
	
	if err := query.Order("capacity desc, price_per_hour asc, production_year desc, created_at asc").Find(&vehicles).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	domainVehicles := mappers.BatchVehicleEntitiesToDomain(vehicles)
	return domainVehicles, nil
}
