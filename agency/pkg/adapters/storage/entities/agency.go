package entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Agency struct {
	gorm.Model
	OwnerID    `gorm:"uniqueIndex"`
	Name      string    `gorm:"uniqueIndex"`
	IsBlocked bool
}