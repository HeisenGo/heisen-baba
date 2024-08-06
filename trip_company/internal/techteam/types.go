package techteam

import (
	"context"
	"errors"
)

var(
	ErrMemberNotFound = errors.New("member not found")
	ErrTeamNotFound =  errors.New("tech team not found")
	ErrDuplication = errors.New("team already exists")
	ErrFailedToFetchRecords = errors.New("records not found")
)

type Repo interface {
	GetTechTeamMemberByUserIDAndTechTeamID(ctx context.Context, userID uint, techTeamID uint) (*TechTeamMember, error) 
	GetTechTeamByID(ctx context.Context, id uint) (*TechTeam, error)
	Insert(ctx context.Context, t *TechTeam) error
	InsertMember(ctx context.Context, t *TechTeamMember) error
	GetTechTeamsOfCompany(ctx context.Context, companyId uint, limit, offset uint) ([]TechTeam, uint, error)
	IsUserTechnicianInCompany(ctx context.Context, companyID uint, userID uint) (bool, error)
}

type TechTeam struct {
	ID                 uint
	Name               string
	Description        string
	TripType           string
	TransportCompanyID uint
	Members            []TechTeamMember
}

type TechTeamMember struct {
	ID         uint
	TechTeamID uint
	TechTeam   TechTeam
	UserID     uint
	Role       string
	Email      string
}
