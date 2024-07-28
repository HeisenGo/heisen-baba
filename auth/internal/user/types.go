package user

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"regexp"
	"strings"
)

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrInvalidEmail          = errors.New("invalid email format")
	ErrInvalidPassword       = errors.New("invalid password format")
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrInvalidAuthentication = errors.New("email and password doesn't match")
)

type Repo interface {
	Create(ctx context.Context, user *User) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
}

// type Role uint8

// func (ur Role) String() string {
// 	switch ur {
// 	case RoleUser:
// 		return "user"
// 	case RoleAdmin:
// 		return "admin"
// 	default:
// 		return "unknown"
// 	}
// }

// const (
// 	RoleUser Role = iota + 1
// 	RoleAdmin
// )

type User struct {
	ID        uuid.UUID
	Username     string
	Email        string 
	PasswordHash string 
	IsSuperAdmin bool   
	IsAdmin      bool   
	Roles        []Role 
	IsBlocked    bool
}

func (u *User) SetPassword(password string) {
	u.PasswordHash = password
}

func (u *User) PasswordIsValid(pass string) bool {
	h := sha256.New()
	h.Write([]byte(pass))
	passSha256 := h.Sum(nil)
	return fmt.Sprintf("%x", passSha256) == u.PasswordHash
}

func ValidateEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`)
	isMatched := emailRegex.MatchString(email)
	if !isMatched {
		return ErrInvalidEmail
	}
	return nil
}

func ValidatePasswordWithFeedback(password string) error {
	tests := []struct {
		pattern string
		message string
	}{
		{".{7,}", "Password must be at least 7 characters long"},
		{"[a-z]", "Password must contain at least one lowercase letter"},
		{"[A-Z]", "Password must contain at least one uppercase letter"},
		{"[0-9]", "Password must contain at least one digit"},
		{"[^\\d\\w]", "Password must contain at least one special character"},
	}

	var errMessages []string

	for _, test := range tests {
		match, _ := regexp.MatchString(test.pattern, password)
		if !match {
			errMessages = append(errMessages, test.message)
		}
	}

	if len(errMessages) > 0 {
		feedback := strings.Join(errMessages, "\n")
		return errors.Join(ErrInvalidPassword, fmt.Errorf(feedback))
	}

	return nil
}

func LowerCaseEmail(email string) string {
	return strings.ToLower(email)
}

type Role struct {
	ID        uint
	Name        string `gorm:"uniqueIndex;not null"`
	Description string
	Permissions []Permission `gorm:"many2many:role_permissions;"`
}

type Permission struct {
	ID uint
	Name        string `gorm:"uniqueIndex;not null"`
	Description string
}

type CompanyUserRole struct {
	ID uint 
	UserID     uuid.UUID  
	CompanyID uint   
	Role      string 
}

type HotelUserRole struct {
	ID uint
	UserID   uuid.UUID   
	HotelID uint   
	Role    string 
}

type AgencyUserRole struct {
	ID uint
	UserID    uuid.UUID
	AgencyID uint   
	Role     string 
}

type UserRole struct {
	ID uint
	UserID  uuid.UUID
	RoleID uint 
}

type RolePermission struct {
	ID uint
	RoleID       uint
	PermissionID uint 
}

