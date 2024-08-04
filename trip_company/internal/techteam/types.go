package techteam

import "context"

type Repo interface {
	GetTechTeamMemberByUserIDAndTechTeamID(ctx context.Context, userID uint, techTeamID uint) (*TechTeamMember, error) 
	GetTechTeamByID(ctx context.Context, id uint) (*TechTeam, error)
	Insert(ctx context.Context, t *TechTeam) error
	InsertMember(ctx context.Context, t *TechTeamMember) error
	GetTechTeamsOfCompany(ctx context.Context, companyId uint, limit, offset uint) ([]TechTeam, uint, error)
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
