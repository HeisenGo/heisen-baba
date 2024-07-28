package entities

import "gorm.io/gorm"

type Room struct {
	gorm.Model
	HotelID     uint    
	Hotel       *Hotel   `gorm:"foreignKey:HotelID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name        string
	AgencyPrice uint64	
	UserPrice   uint64
	Facilities  string
	Capacity    uint8
	IsAvailable bool
}