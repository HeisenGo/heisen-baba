package entities

import "gorm.io/gorm"

type TechTeam struct {
	gorm.Model
	Name               string   `gorm:"type:varchar(100);not null"`
	Description        string   `gorm:"type:text"`
	TripType           string `gorm:"type:varchar(20);not null"`
	TransportCompanyID uint
	TransportCompany   TransportCompany `gorm:"foreignKey:TransportCompanyID"`
	Members            []TechTeamMember `gorm:"foreignKey:TechTeamID"`
}

type TechTeamMember struct {
	gorm.Model
	TechTeamID uint
	TechTeam   TechTeam `gorm:"foreignKey:TechTeamID"`
	UserID     uint   `gorm:"not null"`
	Role       string `gorm:"type:varchar(50)"`
}
