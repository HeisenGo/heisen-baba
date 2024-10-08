package service

import (
	"context"
	"errors"
	"tripcompanyservice/internal/company"
	"tripcompanyservice/internal/trip"

	"github.com/google/uuid"
)

var (
	ErrForbidden = errors.New("you are not allowed to do this function")
)

type TransportCompanyService struct {
	companyOps *company.Ops
	tripOps    *trip.Ops
}

func NewTransportCompanyService(companyOps *company.Ops, tripOps *trip.Ops) *TransportCompanyService {
	return &TransportCompanyService{
		companyOps: companyOps,
		tripOps:    tripOps,
	}
}

func (s *TransportCompanyService) CreateTransportCompany(ctx context.Context, company *company.TransportCompany) error {
	return s.companyOps.Create(ctx, company)
}

func (s *TransportCompanyService) GetUserTransportCompanies(ctx context.Context, ownerID uuid.UUID, page, pageSize uint) ([]company.TransportCompany, uint, error) {
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

func (s *TransportCompanyService) DeleteCompany(ctx context.Context, companyID uint, requesterID uuid.UUID) error {
	co, err := s.companyOps.GetByID(ctx, companyID)
	if err != nil {
		return err
	}

	nActiveTrips, err := s.tripOps.GetCountCompanyUnfinishedUncanceledTrips(ctx, companyID)

	if err!=nil{
		return err
	}
	if nActiveTrips > 0 {
		return company.ErrCanNotDelete
	}
	if co.OwnerID != requesterID{
		return ErrForbidden
	}
	err = s.companyOps.Delete(ctx, companyID)
	if err != nil {
		return err
	}
	return nil
}

func (s *TransportCompanyService) PatchCompanyByOwner(ctx context.Context, updatedCompany *company.TransportCompany, userID uuid.UUID, newOwnerEmail string) (*company.TransportCompany, error) {
	originalCompany, err := s.companyOps.GetByID(ctx, updatedCompany.ID)
	if err != nil {
		return nil, err
	}
	if originalCompany.OwnerID != userID {
		return nil, ErrForbidden
	}

	if newOwnerEmail != "" {
		//******
		// TO DO: check this user exists and get it!!!!
		//updatedCompany.OwnerID = 1
		return nil, errors.New("we need to find user's uuid to withdraw the company")
	}

	err = s.companyOps.PatchCompanyByOwner(ctx, updatedCompany, originalCompany)

	return originalCompany, err
}
