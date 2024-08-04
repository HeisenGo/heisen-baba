package entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Vehicle struct {
	gorm.Model
	ID                    uint
	Name                  string
	OwnerID               uuid.UUID		`gorm:"uniqueIndex:idx_owner_motor_number"`
	PricePerHour          float64
	MotorNumber           string		`gorm:"uniqueIndex:idx_owner_motor_number"`
	MinRequiredTechPerson uint
	IsActive              bool
	Capacity              uint
	IsBlocked             bool
	Type                  string // rail, road, air, sailing
	Speed                 float64
	ProductionYear        uint
	IsConfirmedByAdmin    bool
}
