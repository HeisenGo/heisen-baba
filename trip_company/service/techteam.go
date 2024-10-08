package service

import (
	"context"
	"tripcompanyservice/internal/company"
	"tripcompanyservice/internal/techteam"
	"tripcompanyservice/internal/trip"

	"github.com/google/uuid"
)

type TechTeamService struct {
	techTeamOps *techteam.Ops
	tripOps     *trip.Ops
	companyOps  *company.Ops
}

func NewTechTeamService(techTeamOps *techteam.Ops, tripOps *trip.Ops, companyOps *company.Ops) *TechTeamService {
	return &TechTeamService{
		techTeamOps: techTeamOps,
		tripOps:     tripOps,
		companyOps:  companyOps,
	}
}

func (s *TechTeamService) CreateTechTeam(ctx context.Context, t *techteam.TechTeam, creatorID uuid.UUID) error {
	c, err := s.companyOps.GetByID(ctx, t.TransportCompanyID)
	if err != nil {
		return err
	}
	if c.OwnerID != creatorID {
		return ErrForbidden
	}
	return s.techTeamOps.CreateTechTeam(ctx, t)
}

func (s *TechTeamService) CreateTechTeamMember(ctx context.Context, m *techteam.TechTeamMember, creatorID uuid.UUID) error {

	m.Role = "technician"
	team, err := s.techTeamOps.GetTechTeamByID(ctx, m.TechTeamID)

	if err != nil {
		return err
	}
	c, err := s.companyOps.GetByID(ctx, team.TransportCompanyID)
	if err != nil {
		return err
	}
	if c.OwnerID != creatorID {
		return ErrForbidden
	}

	return s.techTeamOps.CreateTechTeamMember(ctx, m)
}

func (s *TechTeamService) GetTechTeamByID(ctx context.Context, id uint) (*techteam.TechTeam, error) {
	return s.techTeamOps.GetTechTeamByID(ctx, id)
}

func (s *TechTeamService) GetTechTeamsOfCompany(ctx context.Context, companyId uint, requesterID uuid.UUID, page, pageSize uint) ([]techteam.TechTeam, uint, error) {
	c, err := s.companyOps.GetByID(ctx, companyId)
	if err != nil {
		return nil, 0, err
	}
	if c.OwnerID != requesterID {
		return nil, 0, ErrForbidden
	}
	return s.techTeamOps.GetTechTeamsOfCompany(ctx, companyId, page, pageSize)
}

func (s *TechTeamService) DeleteTeam(ctx context.Context, teamID uint, creatorID uuid.UUID) error {
	team, err := s.techTeamOps.GetFullTechTeamByID(ctx, teamID)
	if err != nil {
		return err
	}

	if team.TransportCompany.OwnerID != creatorID {
		return ErrForbidden
	}

	err = s.techTeamOps.Delete(ctx, teamID)
	if err != nil {
		return err
	}
	return nil
}
