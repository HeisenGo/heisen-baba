package entities

import "gorm.io/gorm"

type Terminal struct {
	gorm.Model
	Name      string       `gorm:"type:varchar(100);not null;uniqueIndex"`
	Type      string `gorm:"type:varchar(20);not null"`
	City      string       `gorm:"type:varchar(100);not null"`
	Country   string       `gorm:"type:varchar(100);not null"`
	OutgoingPaths []Path   `gorm:"foreignKey:FromTerminalID"`
	IncomingPaths []Path   `gorm:"foreignKey:ToTerminalID"`
}