package entities

import "gorm.io/gorm"

type Terminal struct {
	gorm.Model
	Name          string `gorm:"type:varchar(100);not null;index:idx_name_type_city_country,priority:1"`
	Type          string `gorm:"type:varchar(20);not null;index:idx_name_type_city_country,priority:2"`
	City          string `gorm:"type:varchar(100);not null;index:idx_name_type_city_country,priority:3"`
	Country       string `gorm:"type:varchar(100);not null;index:idx_name_type_city_country,priority:4"`
	OutgoingPaths []Path `gorm:"foreignKey:FromTerminalID"`
	IncomingPaths []Path `gorm:"foreignKey:ToTerminalID"`

	// Composite unique constraint
	UniqueTerminal string `gorm:"uniqueIndex:idx_name_type_city_country"`
}
