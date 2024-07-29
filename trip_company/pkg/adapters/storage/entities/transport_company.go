package entities

import (
	"gorm.io/gorm"
)

type TransportCompany struct {
	gorm.Model
	Name        string `gorm:"type:varchar(100);not null;uniqueIndex:idx_owner_name"`
	Description string `gorm:"type:text"`
	OwnerID     uint   `gorm:"not null;uniqueIndex:idx_owner_name"`
	Address     string `gorm:"type:varchar(255)"`
	PhoneNumber string `gorm:"type:varchar(20)"`
	Email       string `gorm:"type:varchar(100);uniqueIndex"`
	IsBlocked   bool   `gorm:"not null; default:false"`
	// relationships
	Employees   []Employee `gorm:"foreignKey:TransportCompanyID"`
	Trips       []Trip `gorm:"foreignKey:TransportCompanyID"`
	TechTeams   []TechTeam `gorm:"foreignKey:TransportCompanyID"`
}
