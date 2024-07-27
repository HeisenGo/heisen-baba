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
