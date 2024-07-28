package entities

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	//ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"uniqueIndex;not null"`
	Email        string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Roles        []Role    `gorm:"many2many:user_roles;"`
	Companies    []Company `gorm:"many2many:user_companies;"`
}

type Role struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"uniqueIndex;not null"`
	Description string
	Permissions []Permission `gorm:"many2many:role_permissions;"`
}

type Permission struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"uniqueIndex;not null"`
	Description string
}

type Company struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Users     []User `gorm:"many2many:user_companies;"`
}

type UserCompany struct {
	UserID    uint   `gorm:"primaryKey"`
	CompanyID uint   `gorm:"primaryKey"`
	Role      string `gorm:"not null"`
}
