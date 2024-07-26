package entities

import (
	"time"

	"gorm.io/gorm"
)

type TripType string

const (
	RailTrip    TripType = "rail"
	AirTrip     TripType = "air"
	RoadTrip    TripType = "road"
	SailingTrip TripType = "sailing"
)

type Trip struct {
	gorm.Model
	TransportCompanyID    uint
	TransportCompany      TransportCompany `gorm:"foreignKey:TransportCompanyID"`
	TripType              TripType `gorm:"type:varchar(20);not null"`
	UserReleaseDate       time.Time
	TourReleaseDate       time.Time
	Price                 float64
	PathID                uint   `gorm:"not null"`
	Origin                string `gorm:"type:varchar(100);not null"`
	Destination           string `gorm:"type:varchar(100);not null"`
	MinPassengers         int
	VehicleReservationFee float64
	VehicleProductionYear int
	Status                string `gorm:"type:varchar(20);default:'pending'"`
	TechTeamID            uint
	TechTeam              TechTeam `gorm:"foreignKey:TechTeamID"`
	VehicleRequestID      uint
	VehicleRequest        VehicleRequest `gorm:"foreignKey:TripID"`
	Tickets               []Ticket `gorm:"foreignKey:TripID"`
}


