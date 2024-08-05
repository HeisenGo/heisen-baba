package storage

import (
	"context"
	"fmt"
	vehiclerequest "tripcompanyservice/internal/vehicle_request"
	"tripcompanyservice/pkg/adapters/storage/entities"
	"tripcompanyservice/pkg/adapters/storage/mappers"

	"gorm.io/gorm"
)

type vehicleReqRepo struct {
	db *gorm.DB
}

func NewVehicleReqRepo(db *gorm.DB) vehiclerequest.Repo {
	return &vehicleReqRepo{db}
}

func (r *vehicleReqRepo) Insert(ctx context.Context, vR *vehiclerequest.VehicleRequest) error {
	vREntity := mappers.VehicleDomainToVehicleEntity(vR)
	if err := r.db.WithContext(ctx).Save(&vREntity).Error; err != nil {
		return err
	}

	vR.ID = vREntity.ID

	return nil
}

func (r *vehicleReqRepo) UpdateVehicleReq(ctx context.Context, id uint, updates map[string]interface{}) error {
	var t entities.VehicleRequest

	if err := r.db.WithContext(ctx).Model(&t).Where("id = ?", id).Updates(updates).Error; err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
