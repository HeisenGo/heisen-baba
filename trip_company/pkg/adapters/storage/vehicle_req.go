package storage

import (
	"context"
	"fmt"
	"strings"
	"tripcompanyservice/internal/company"
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

func (r *vehicleReqRepo) Delete(ctx context.Context, vRID uint) error {
	if err := r.db.WithContext(ctx).Delete(&entities.VehicleRequest{}, vRID).Error; err != nil {
		return fmt.Errorf("%w %w", vehiclerequest.ErrDeleteVehicleReq, err)
	} else {
		return nil
	}

}

func (r *vehicleReqRepo) GetByID(ctx context.Context, id uint) (*vehiclerequest.VehicleRequest, error) {
	var t entities.VehicleRequest

	err := r.db.WithContext(ctx).Model(&entities.VehicleRequest{}).Where("id = ?", id).First(&t).Error
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, nil
		}
		return nil, err
	}
	dC := mappers.VehicleReqEntityToVehicleReqDomain(t)
	return &dC, nil
}

func (r *vehicleReqRepo) GetTransportCompanyByVehicleRequestID(ctx context.Context, vehicleRequestID uint) (*company.TransportCompany, error) {
	var vehicleRequest entities.VehicleRequest

	if err := r.db.WithContext(ctx).First(&vehicleRequest, vehicleRequestID).Error; err != nil {
		return nil, fmt.Errorf("could not find vehicle request: %w", err)
	}

	var trip entities.Trip

	if err := r.db.WithContext(ctx).First(&trip, vehicleRequest.TripID).Error; err != nil {
		return nil, fmt.Errorf("could not find trip: %w", err)
	}

	var transportCompany entities.TransportCompany

	if err := r.db.WithContext(ctx).First(&transportCompany, trip.TransportCompanyID).Error; err != nil {
		return nil, fmt.Errorf("could not find transport company: %w", err)
	}
	dt := mappers.CompanyEntityToDomain(transportCompany)
	return &dt, nil
}
