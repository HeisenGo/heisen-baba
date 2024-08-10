package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Reservation struct {
	gorm.Model
	TourID     uint
	Room       *Tour `gorm:"foreignKey:TourID;constraint:OnUpdate:CASCADE;"`
	UserID     uuid.UUID
	CheckIn    time.Time
	CheckOut   time.Time
	TotalPrice uint64
	Status     string // e.g., "booked", "checked_in", "checked_out", "canceled"
}
