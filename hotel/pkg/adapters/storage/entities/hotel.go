package entities

import "gorm.io/gorm"

type Hotel struct {
	gorm.Model
	OwnerID      uint   `gorm:"uniqueIndex:idx_owner_name"`
	Name         string `gorm:"uniqueIndex:idx_owner_name"`
	City         string
	Country      string
	Details      string
	IsBlocked    bool
	Rooms        []Room
}