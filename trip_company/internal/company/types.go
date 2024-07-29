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
	ErrDeleteCompany   = errors.New("error deleting company")
	ErrCanNotDelete    = errors.New("company can not be deleted due to unfinished trips")
	ErrFailedToBlock   = errors.New("failed to block the company")
)

type Repo interface {
	GetByID(ctx context.Context, id uint) (*TransportCompany, error)
	GetTransportCompanies(ctx context.Context, limit, offset uint) ([]TransportCompany, uint, error)
	Insert(ctx context.Context, company *TransportCompany) error
	GetUserTransportCompanies(ctx context.Context, ownerID uint, limit, offset uint) ([]TransportCompany, uint, error)
	Delete(ctx context.Context, companyID uint) error
	BlockCompany(ctx context.Context, companyID uint, isBlocked bool) error
}

type TransportCompany struct {
	ID          uint
	Name        string
	Description string
	OwnerID     uint
	Address     string
	IsBlocked   bool
	//PhoneNumber string
	//Email       string
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
