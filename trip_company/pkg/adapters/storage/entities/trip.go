package entities

import (
	"time"

	"gorm.io/gorm"
)

type Trip struct {
	gorm.Model
	TransportCompanyID      uint `gorm:"not null"`//;uniqueIndex:idx_trip_unique"`
	TransportCompany        *TransportCompany `gorm:"foreignKey:TransportCompanyID; constraint:OnDelete:CASCADE;"`
	TripType                string            `gorm:"type:varchar(20);not null"`
	UserReleaseDate         time.Time
	TourReleaseDate         time.Time
	UserPrice               float64
	AgencyPrice             float64
	PathID                  uint   `gorm:"not null"`//;uniqueIndex:idx_trip_unique"`
	FromCountry      string
	ToCountry string
	Origin                  string `gorm:"type:varchar(100);not null"`
	FromTerminalName        string
	ToTerminalName          string
	Destination             string `gorm:"type:varchar(100);not null"`
	PathName                string
	PathDistanceKM          float64
	Status                  string `gorm:"type:varchar(20);default:'pending'"`
	MinPassengers           uint
	TechTeamID              *uint     `gorm:""`
	TechTeam                *TechTeam `gorm:"foreignKey:TechTeamID; constraint:OnDelete:CASCADE;"`
	VehicleRequestID        *uint
	VehicleRequest          *VehicleRequest `gorm:"foreignKey:TripID; constraint:OnDelete:CASCADE;"`
	Tickets                 []Ticket        `gorm:"foreignKey:TripID; constraint:OnDelete:CASCADE;"`
	SoldTickets 	uint `gorm:"default:0"`
	TripCancellingPenaltyID *uint
	TripCancelingPenalty    *TripCancellingPenalty `gorm:"foreignKey:TripCancellingPenaltyID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	MaxTickets              uint
	VehicleID               *uint
	VehicleName             string 
	IsCanceled              bool       `gorm:"default:false"`
	IsFinished              bool       `gorm:"default:false"`
	IsConfirmed             bool    `gorm:"default:false"`
	StartDate               *time.Time `gorm:"not null;"`//uniqueIndex:idx_trip_unique"` // should be given by trip generator
	EndDate                 *time.Time // should be calculated according to the vehicle speed and path distance
	Profit      float64  `gorm:"default:0"`
}
