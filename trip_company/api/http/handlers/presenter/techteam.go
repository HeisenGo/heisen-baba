package presenter

import (
	"tripcompanyservice/internal/techteam"
	"tripcompanyservice/pkg/fp"

	"github.com/google/uuid"
)

type TechTeamRe struct {
	ID                 uint               `json:"id"`
	Name               string             `json:"name" validate:"required"`
	Description        string             `json:"desc"`
	TripType           string             `json:"type" validate:"required"`
	TransportCompanyID uint               `json:"company_id" validate:"required"`
	Members            []TechTeamMemberRe `json:"members"`
}

type TechTeamMemberRe struct {
	ID         uint   `json:"id"`
	TechTeamID uint   `json:"team_id"`
	UserID     uuid.UUID   `json:"user_id"`
	Email      string `json:"email"`
	Role       string `json:"role"`
}

func TechMemberToTechTeamMemberRe(m techteam.TechTeamMember) TechTeamMemberRe {
	return TechTeamMemberRe{
		ID:         m.ID,
		TechTeamID: m.TechTeamID,
		UserID:     m.UserID,
		Email:      m.Email,
		Role:       m.Role,
	}
}

// type TechTeamMemberRe struct {
// 	ID uint  `json:"id"`
// 	TechTeamID uint `json:"team_id"`
// 	UserID     uint `json:"user_id"`
// 	Role       string  `json:"role"`
// }

func BatchTeamToTechTeamRe(m []techteam.TechTeam) []TechTeamRe {
	return fp.Map(m, TechTeamToTechTeamRe)
}

func BatchTeamMemberToTechTeamMemberRe(m []techteam.TechTeamMember) []TechTeamMemberRe {
	return fp.Map(m, TechMemberToTechTeamMemberRe)
}

func TechTeamToTechTeamRe(t techteam.TechTeam) TechTeamRe {
	m := make([]TechTeamMemberRe, len(t.Members))
	if len(t.Members) > 0 {
		m = BatchTeamMemberToTechTeamMemberRe(t.Members)
	}

	return TechTeamRe{
		ID:                 t.ID,
		Members:            m,
		Name:               t.Name,
		Description:        t.Description,
		TripType:           t.TripType,
		TransportCompanyID: t.TransportCompanyID,
	}
}

func TechTeamReqToTechTeam(t *TechTeamRe) *techteam.TechTeam {
	return &techteam.TechTeam{
		Name:               t.Name,
		Description:        t.Description,
		TransportCompanyID: t.TransportCompanyID,
		TripType:           t.TripType,
	}
}

func TechTeamMemberReToTechTeamMember(m *TechTeamMemberRe) *techteam.TechTeamMember {
	return &techteam.TechTeamMember{
		Email:      m.Email,
		UserID:     m.UserID,
		TechTeamID: m.TechTeamID,
		Role:       m.Role,
	}
}
