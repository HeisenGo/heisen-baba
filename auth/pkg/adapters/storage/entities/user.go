package entities

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	//ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"uniqueIndex;not null"`
	Email        string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	IsSuperAdmin bool   `gorm:"not null; default:false"`
	IsAdmin      bool   `gorm:"not null; default:false"`
	Roles        []Role `gorm:"many2many:user_roles;"`
	IsBlocked    bool   `gorm:"not null; default:false"`
}

type Role struct {
	gorm.Model
	Name        string `gorm:"uniqueIndex;not null"`
	Description string
	Permissions []Permission `gorm:"many2many:role_permissions;"`
}

type Permission struct {
	gorm.Model
	Name        string `gorm:"uniqueIndex;not null"`
	Description string
}

type CompanyUserRole struct {
	gorm.Model
	UserID    uint   `gorm:"not null"`
	CompanyID uint   `gorm:"not null"`
	Role      string `gorm:"not null"`
}

type HotelUserRole struct {
	gorm.Model
	UserID  uint   `gorm:"not null"`
	HotelID uint   `gorm:"not null"`
	Role    string `gorm:"not null"`
}

type AgencyUserRole struct {
	gorm.Model
	UserID   uint   `gorm:"not null"`
	AgencyID uint   `gorm:"not null"`
	Role     string `gorm:"not null"`
}

type UserRole struct {
	UserID uint `gorm:"primaryKey"`
	RoleID uint `gorm:"primaryKey"`
}

type RolePermission struct {
	RoleID       uint `gorm:"primaryKey"`
	PermissionID uint `gorm:"primaryKey"`
}
