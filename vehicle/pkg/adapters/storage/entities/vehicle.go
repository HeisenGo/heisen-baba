package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Vehicle struct {
	gorm.Model
	ID                    uint
	Name                  string
	OwnerID               uuid.UUID
	PricePerHour          float64
	MotorNumber           string
	MinRequiredTechPerson uint
	IsActive              bool
	Capacity              uint
	IsBlocked             bool
	Type                  string // rail, road, air, sailing
	Speed                 float64
	ProductionYear        uint
	IsConfirmedByAdmin    bool
	EntryDate             time.Time
}
