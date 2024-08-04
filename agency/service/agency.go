package service

import (
	"agency/internal/agency"
	"context"

	"github.com/google/uuid"
)

type AgencyService struct {
	agencyOps *agency.Ops
}

func NewAgencyService(agencyOps *agency.Ops) *AgencyService {
	return &AgencyService{
		agencyOps: agencyOps,
	}
}

func (s *AgencyService) CreateAgency(ctx context.Context, a *agency.Agency) error {
	return s.agencyOps.CreateAgency(ctx, a)
}

func (s *AgencyService) GetAgencies(ctx context.Context, name string, page, pageSize int) ([]agency.Agency, uint, error) {
	return s.agencyOps.GetAgencies(ctx, name, page, pageSize)
}

func (s *AgencyService) GetAgenciesByOwnerID(ctx context.Context, ownerID uuid.UUID, page, pageSize int) ([]agency.Agency, int, error) {
	return s.agencyOps.GetAgenciesByOwnerID(ctx, ownerID, page, pageSize)
}

func (s *AgencyService) UpdateAgency(ctx context.Context, id uint, updates *agency.Agency) error {
	existingAgency, err := s.agencyOps.GetAgencyByID(ctx, id)
	if err != nil {
		return err
	}

	// Update only the fields that are provided
	if updates.Name != "" {
		existingAgency.Name = updates.Name
	}
	existingAgency.IsBlocked = updates.IsBlocked

	return s.agencyOps.UpdateAgency(ctx, existingAgency)
}

func (s *AgencyService) DeleteAgency(ctx context.Context, id uint) error {
	return s.agencyOps.DeleteAgency(ctx, id)
}

func (s *AgencyService) BlockAgency(ctx context.Context, agencyID uint) error {
	return s.agencyOps.BlockAgency(ctx, agencyID)
}
