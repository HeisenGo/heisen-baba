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
	//PhoneNumber string `gorm:"type:varchar(20);uniqueIndex:idx_phone_email"`
	//Email       string `gorm:"type:varchar(100);uniqueIndex:idx_phone_email"`
	IsBlocked   bool   `gorm:"not null;default:false"`
	// relationships
	Employees []Employee `gorm:"foreignKey:TransportCompanyID; constraint:OnDelete:CASCADE;"`
	Trips     []Trip     `gorm:"foreignKey:TransportCompanyID; constraint:OnDelete:CASCADE;"`
	TechTeams []TechTeam `gorm:"foreignKey:TransportCompanyID; constraint:OnDelete:CASCADE;"`
}
