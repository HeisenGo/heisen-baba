package mappers

import (
	"tripcompanyservice/internal/techteam"
	"tripcompanyservice/pkg/adapters/storage/entities"
	"tripcompanyservice/pkg/fp"
)

func MemberDomainToEntity(m *techteam.TechTeamMember) *entities.TechTeamMember {
	return &entities.TechTeamMember{
		Email:      m.Email,
		Role:       m.Role,
		UserID:     m.UserID,
		TechTeamID: m.TechTeamID,
	}
}

func TechTeamDomainToEntity(t *techteam.TechTeam) *entities.TechTeam {
	return &entities.TechTeam{
		Name:               t.Name,
		Description:        t.Description,
		TripType:           t.TripType,
		TransportCompanyID: t.TransportCompanyID,
	}
}

func MemberEntityToMemberDomain(m entities.TechTeamMember) techteam.TechTeamMember {
	return techteam.TechTeamMember{
		ID:         m.ID,
		Email:      m.Email,
		UserID:     m.UserID,
		TechTeamID: m.TechTeamID,
		Role:       m.Role,
	}
}

func BatchMembersEntitiesToMembersDomain(m []entities.TechTeamMember) []techteam.TechTeamMember {
	return fp.Map(m, MemberEntityToMemberDomain)
}

func TechTeamEntityToDomain(teamEntity entities.TechTeam) techteam.TechTeam {
	m := make([]techteam.TechTeamMember, len(teamEntity.Members))
	if len(teamEntity.Members) > 0 {
		m = BatchMembersEntitiesToMembersDomain(teamEntity.Members)
	}
	return techteam.TechTeam{
		ID:                 teamEntity.ID,
		TransportCompanyID: teamEntity.TransportCompanyID,
		Name:               teamEntity.Name,
		Description:        teamEntity.Description,
		TripType:           teamEntity.TripType,
		Members:            m,
	}
}

func BatchTecTeamEnToDo(t []entities.TechTeam) []techteam.TechTeam {
	return fp.Map(t, TechTeamEntityToDomain)
}
