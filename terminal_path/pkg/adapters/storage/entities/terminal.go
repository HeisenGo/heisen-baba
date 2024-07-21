package entities

import "gorm.io/gorm"

type TerminalType string

const (
	Air TerminalType = "ait"
	Rail TerminalType = "rail"
	Road TerminalType = "road"
	Sailing TerminalType = "sailing" // port
)

type Terminal struct {
	gorm.Model
	Name      string       `gorm:"type:varchar(100);not null;uniqueIndex"`
	Type      TerminalType `gorm:"type:varchar(20);not null"`
	City      string       `gorm:"type:varchar(100);not null"`
	Country   string       `gorm:"type:varchar(100);not null"`
	OutgoingPaths []Path   `gorm:"foreignKey:FromTerminalID"`
	IncomingPaths []Path   `gorm:"foreignKey:ToTerminalID"`
}