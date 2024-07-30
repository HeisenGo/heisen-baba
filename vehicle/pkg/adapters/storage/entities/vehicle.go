package entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Vehicle struct {
	gorm.Model
	Name                  string    `gorm:"not null"`
	OwnerID               uuid.UUID `gorm:"not null"`
	PricePerHour          float64   `gorm:"not null"`
	MotorNumber           string    `gorm:"unique;not null"`
	MinRequiredTechPerson uint      `gorm:"not null"`
	IsActive              bool      `gorm:"not null"`
	Capacity              uint      `gorm:"not null"`
	IsBlocked             bool      `gorm:"not null"`
	Type                  string    `gorm:"not null"` // rail, road, air, sailing
	Speed                 float64   `gorm:"not null"`
	ProductionYear        uint      `gorm:"not null"`
	IsConfirmedByAdmin    bool      `gorm:"not null"`
}
