package company

import (
	"context"
	"errors"
	"regexp"
	"strings"
)

const (
	MaxNameLength        = 100
	MaxDescriptionLength = 1500
	MaxAddressLength     = 800
)

var (
	ErrCompanyNotFound = errors.New("company not found")
	ErrFailedToRestore = errors.New("failed to restore company")
	ErrDuplication     = errors.New("company with same email/ ownerID-Name already exists")
	ErrInvalidEmail    = errors.New("invalid email")
)

type Repo interface {
	GetTransportCompanies(ctx context.Context, limit, offset uint) ([]TransportCompany, uint, error)
	Insert(ctx context.Context, company *TransportCompany) error
	GetUserTransportCompanies(ctx context.Context, ownerID uint, limit, offset uint) ([]TransportCompany, uint, error) 

}

type TransportCompany struct {
	ID          uint
	Name        string
	Description string
	OwnerID     uint
	Address     string
	PhoneNumber string
	Email       string
	// relationships
	//Employees   []Employee
	//Trips       []Trip
	//TechTeams   []TechTeam
}

func ValidateEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`)
	isMatched := emailRegex.MatchString(email)
	if !isMatched {
		return ErrInvalidEmail
	}
	return nil
}

func LowerCaseEmail(email string) string {
	return strings.ToLower(email)
}
