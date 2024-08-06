package entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TechTeam struct {
	gorm.Model
	Name               string           `gorm:"type:varchar(100);not null;uniqueIndex:idx_name_co_unique"`
	Description        string           `gorm:"type:text"`
	TripType           string           `gorm:"type:varchar(20);not null"`
	TransportCompanyID uint             `gorm:"not null; uniqueIndex:idx_name_co_unique"`
	TransportCompany   TransportCompany `gorm:"foreignKey:TransportCompanyID;constraint:OnDelete:CASCADE;"`
	Members            []TechTeamMember `gorm:"foreignKey:TechTeamID;constraint:OnDelete:CASCADE;"`
}

type TechTeamMember struct {
	gorm.Model
	TechTeamID uint     `gorm:"not null; uniqueIndex:idx_tech_team_user"`
	TechTeam   TechTeam `gorm:"foreignKey:TechTeamID; constraint:OnDelete:CASCADE;"`
	UserID     uuid.UUID     `gorm:"type:uuid;not null; uniqueIndex:idx_tech_team_user"`
	Role       string   `gorm:"type:varchar(50); default:'technician'"`
	Email      string
}
