package entities

import "gorm.io/gorm"

type Employee struct {
	gorm.Model
	UserID             uint `gorm:"not null"`
	TransportCompanyID uint
	TransportCompany   TransportCompany `gorm:"foreignKey:TransportCompanyID"`
	Role               string `gorm:"type:varchar(50)"`
}
