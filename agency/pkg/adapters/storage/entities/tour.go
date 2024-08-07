package entities

import (
	"time"

	"gorm.io/gorm"
)

type Tour struct {
	gorm.Model
	AgencyID     uint
	Agency       *Agency `gorm:"foreignKey:AgencyID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	GoTicketID   uint
	BackTicketID uint
	HotelID      uint
	Capacity     uint
	UserPrice    uint64
	ReleaseDate  time.Time
	IsApproved   bool
	IsActive     bool
}
