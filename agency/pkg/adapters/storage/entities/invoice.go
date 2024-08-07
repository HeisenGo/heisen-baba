package entities

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Invoice struct {
	gorm.Model
	ReservationID uint
	Reservation   *Reservation `gorm:"foreignKey:ReservationID"`
	UserID        uuid.UUID
	OwnerID       uuid.UUID
	IssueDate     time.Time
	Amount        uint64
	Paid          bool
}
