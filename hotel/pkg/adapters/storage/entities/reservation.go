package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Reservation struct {
	gorm.Model
	RoomID     uint
	Room       *Room `gorm:"foreignKey:RoomID;constraint:OnUpdate:CASCADE;"`
	UserID     uuid.UUID
	CheckIn    time.Time
	CheckOut   time.Time
	TotalPrice uint64
	Status     string // e.g., "booked", "checked_in", "checked_out", "canceled"
}
