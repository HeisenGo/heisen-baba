package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	ID           uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Username     string         `gorm:"uniqueIndex;not null"`
	Email        string         `gorm:"uniqueIndex;not null"`
	PasswordHash string         `gorm:"not null"`
	IsSuperAdmin bool           `gorm:"not null; default:false"`
	IsAdmin      bool           `gorm:"not null; default:false"`
	Roles        []Role         `gorm:"many2many:user_roles;"`
	IsBlocked    bool           `gorm:"not null; default:false"`
}

type Role struct {
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	ID          uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name        string         `gorm:"uniqueIndex;not null"`
	Description string
	Permissions []Permission `gorm:"many2many:role_permissions;"`
}

type Permission struct {
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	ID          uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name        string         `gorm:"uniqueIndex;not null"`
	Description string
}

type CompanyUserRole struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	uuid.UUID `gorm:"type:uuid;not null"`
	CompanyID uint   `gorm:"not null"`
	Role      string `gorm:"not null"`
}

type HotelUserRole struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID uuid.UUID `gorm:"type:uuid;not null"`
	HotelID   uint           `gorm:"not null"`
	Role      string         `gorm:"not null"`
}

type AgencyUserRole struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID     uuid.UUID `gorm:"type:uuid;not null"`
	AgencyID  uint           `gorm:"not null"`
	Role      string         `gorm:"not null"`
}

type UserRole struct {
	UserID  uuid.UUID `gorm:"type:uuid;not null; primaryKey"`
	RoleID uint `gorm:"primaryKey"`
}

type RolePermission struct {
	RoleID        uuid.UUID `gorm:"type:uuid;not null; primaryKey"`
	PermissionID uint `gorm:"primaryKey"`
}
