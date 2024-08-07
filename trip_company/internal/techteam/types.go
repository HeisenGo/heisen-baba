package techteam

import (
	"context"
	"errors"
	"tripcompanyservice/internal/company"

	"github.com/google/uuid"
)

var(
	ErrMemberNotFound = errors.New("member not found")
	ErrTeamNotFound =  errors.New("tech team not found")
	ErrDuplication = errors.New("team already exists")
	ErrFailedToFetchRecords = errors.New("records not found")
	ErrDeleteTeam = errors.New("error in deleting team")
)

type Repo interface {
	GetTechTeamMemberByUserIDAndTechTeamID(ctx context.Context, userID uuid.UUID, techTeamID uint) (*TechTeamMember, error) 
	GetTechTeamByID(ctx context.Context, id uint) (*TechTeam, error)
	Insert(ctx context.Context, t *TechTeam) error
	InsertMember(ctx context.Context, t *TechTeamMember) error
	GetTechTeamsOfCompany(ctx context.Context, companyId uint, limit, offset uint) ([]TechTeam, uint, error)
	IsUserTechnicianInCompany(ctx context.Context, companyID uint, userID uuid.UUID) (bool, error)
	GetFullTechTeamByID(ctx context.Context, id uint) (*TechTeam, error) 
	Delete(ctx context.Context, tID uint) error
}

type TechTeam struct {
	ID                 uint
	Name               string
	Description        string
	TripType           string
	TransportCompanyID uint
	Members            []TechTeamMember
	TransportCompany   *company.TransportCompany
}

type TechTeamMember struct {
	ID         uint
	TechTeamID uint
	TechTeam   TechTeam
	UserID     uuid.UUID
	Role       string
	Email      string
}
