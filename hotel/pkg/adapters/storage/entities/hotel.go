package entities

import "gorm.io/gorm"

type Hotel struct {
	gorm.Model
	OwnerID		  uint
	Name          string
	City          string
	Country       string
	Details       string
	IsBlocked     bool
	Rooms         []Room
}