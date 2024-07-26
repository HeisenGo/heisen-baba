package entities

import "gorm.io/gorm"

type Ticket struct { // Renamed from TripPurchase to Ticket
	gorm.Model
	TripID     uint
	Trip       Trip `gorm:"foreignKey:TripID"`
	UserID     uint `gorm:"not null"`
	TourID     uint
	Quantity   int
	TotalPrice float64
	Status     string `gorm:"type:varchar(20);default:'pending'"`
}
