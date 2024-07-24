package entities

import (
	"gorm.io/gorm"
)

type Path struct {
	gorm.Model
	FromTerminalID uint     `gorm:"not null"`
	ToTerminalID   uint     `gorm:"not null"`
	FromTerminal   Terminal `gorm:"foreignKey:FromTerminalID"`
	ToTerminal     Terminal `gorm:"foreignKey:ToTerminalID"`
	DistanceKM     float64  `gorm:"type:decimal(10,2);not null"` // in kilometers
	Code           string   `gorm:"type:varchar(50);not null;uniqueIndex"`
	Name           string   `gorm:"type:varchar(100)"`
	Type           string   `gorm:"type:varchar(20);not null"`
}
