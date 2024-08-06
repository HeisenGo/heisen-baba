package entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Employee struct {
	gorm.Model
	UserID    uuid.UUID `gorm:"type:uuid;not null; uniqueIndex"`
	TransportCompanyID uint
	TransportCompany   TransportCompany `gorm:"foreignKey:TransportCompanyID"`
	Role               string           `gorm:"type:varchar(50)"`
}
