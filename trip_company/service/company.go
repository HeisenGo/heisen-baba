package service

import (
	"context"
	"tripcompanyservice/internal/company"
)

type TransportCompanyService struct {
	companyOps *company.Ops
}

func NewTransportCompanyService(companyOps *company.Ops) *TransportCompanyService {
	return &TransportCompanyService{
		companyOps: companyOps,
	}
}

func (s *TransportCompanyService) CreateTransportCompany(ctx context.Context, company *company.TransportCompany) error {
	return s.companyOps.Create(ctx, company)
}

func (s *TransportCompanyService) GetUserTransportCompanies(ctx context.Context, ownerID uint, page, pageSize uint) ([]company.TransportCompany, uint, error) {
	// user, err := s.userOps.GetUserByID(ctx, userID) check by the other service
	// if user == nil {
	// 	return nil, 0, u.ErrUserNotFound
	// }

	return s.companyOps.GetUserTransportCompanies(ctx, ownerID, page, pageSize)
}

func (s *TransportCompanyService) GetTransportCompanies(ctx context.Context, page, pageSize uint) ([]company.TransportCompany, uint, error) {
	// user, err := s.userOps.GetUserByID(ctx, userID) check by the other service
	// if user == nil {
	// 	return nil, 0, u.ErrUserNotFound
	// }

	return s.companyOps.GetTransportCompanies(ctx, page, pageSize)
}

func (s *TransportCompanyService) BlockCompany(ctx context.Context, companyID uint, isBlocked bool) (*company.TransportCompany, error) {
	transportCompany, err := s.companyOps.GetByID(ctx, companyID)
	if err != nil {
		return nil, err
	}

	err = s.companyOps.BlockUnBlockCompany(ctx, companyID, isBlocked)
	if err != nil {
		return nil, err
	}
	transportCompany.IsBlocked = isBlocked
	return transportCompany, nil
}

func (s *TransportCompanyService) DeleteCompany(ctx context.Context, companyID uint) error {
	_, err := s.companyOps.GetByID(ctx, companyID)
	if err != nil {
		return err
	}
	err = s.companyOps.Delete(ctx, companyID)
	if err != nil {
		return err
	}
	return nil
}
