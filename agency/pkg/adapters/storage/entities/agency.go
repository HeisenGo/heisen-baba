package entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Agency struct {
	gorm.Model
	OwnerID   uuid.UUID 
	Name      string    `gorm:"uniqueIndex"`
	IsBlocked bool
	Tours     []Tour
}
