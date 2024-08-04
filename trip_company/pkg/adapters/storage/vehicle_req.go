package storage

import (
	vehiclerequest "tripcompanyservice/internal/vehicle_request"

	"gorm.io/gorm"
)

type vehicleReqRepo struct {
	db *gorm.DB
}

func NewVehicleReqRepo(db *gorm.DB) vehiclerequest.Repo {
	return &vehicleReqRepo{db}
}
