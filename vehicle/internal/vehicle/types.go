package vehicle

import (
	"context"
	"errors"
	"regexp"
	"time"

	"github.com/google/uuid"
)

type Repo interface {
	CreateVehicle(ctx context.Context, v *Vehicle) error
	GetVehicles(ctx context.Context, vehicleType string, capacity uint, page, pageSize int) ([]Vehicle, uint, error)
	GetVehiclesByOwnerID(ctx context.Context, ownerID uuid.UUID, page, pageSize int) ([]Vehicle, int, error)
	GetVehicleByID(ctx context.Context, id uint) (*Vehicle, error)
	UpdateVehicle(ctx context.Context, v *Vehicle) error
	DeleteVehicle(ctx context.Context, id uint) error
}

type Vehicle struct {
	ID                    uint
	Name                  string
	OwnerID               uuid.UUID
	PricePerHour          float64
	MotorNumber           string
	MinRequiredTechPerson uint
	IsActive              bool
	Capacity              uint
	IsBlocked             bool
	Type                  string
	Speed                 float64
	ProductionYear        uint
	IsConfirmedByAdmin    bool
}


var (
	ErrInvalidVehicleName     = errors.New("invalid vehicle name: must be 1-100 characters long and can only contain alphanumeric characters, spaces, hyphens, underscores, and periods")
	ErrInvalidType            = errors.New("invalid vehicle type: only 'rail', 'road', 'air', or 'sailing' are accepted")
	ErrInvalidProductionYear  = errors.New("invalid production year")
	ErrInvalidMotorNumber     = errors.New("invalid motor number: must be 1-50 characters long and can only contain alphanumeric characters, hyphens, and underscores")
	ErrInvalidCapacity        = errors.New("invalid capacity")
	ErrRecordNotFound         = errors.New("record not found")
)

func ValidateVehicleName(name string) error {
	var validVehicleName = regexp.MustCompile(`^[a-zA-Z0-9 ._-]{1,100}$`)
	if !validVehicleName.MatchString(name) {
		return ErrInvalidVehicleName
	}
	return nil
}

func ValidateType(vehicleType string) error {
	var validTypes = regexp.MustCompile(`^(train|bus|airplane|ship)$`)
	if !validTypes.MatchString(vehicleType) {
		return ErrInvalidType
	}
	return nil
}

func ValidateProductionYear(year uint) error {
	currentYear := time.Now().Year()
	if year < 1886 || year > uint(currentYear) { // 1886 is considered the birth year of modern vehicles
		return ErrInvalidProductionYear
	}
	return nil
}

func ValidateMotorNumber(motorNumber string) error {
	var validMotorNumber = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,50}$`)
	if !validMotorNumber.MatchString(motorNumber) {
		return ErrInvalidMotorNumber
	}
	return nil
}

func ValidateCapacity(capacity uint) error {
	if capacity <= 0 {
		return ErrInvalidCapacity
	}
	return nil
}