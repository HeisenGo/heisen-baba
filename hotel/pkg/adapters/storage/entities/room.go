package entities

import "gorm.io/gorm"

type Room struct {
	gorm.Model
	HotelID     uint   `gorm:"uniqueIndex:idx_hotel_room_name"`
	Hotel       *Hotel `gorm:"foreignKey:HotelID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name        string `gorm:"uniqueIndex:idx_hotel_room_name"`
	AgencyPrice uint64
	UserPrice   uint64
	Facilities  string
	Capacity    uint8
	IsAvailable bool
}
