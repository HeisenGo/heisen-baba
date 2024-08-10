package presenter

import (
	"tripcompanyservice/internal/company"
	"tripcompanyservice/pkg/fp"

	"github.com/google/uuid"
)

type BlockCompany struct{
	IsBlocked bool `json:"is_blocked"`
}

type UpdateCompanyReq struct {
	Name        string `json:"name"`
	Description string `json:"desc"`
	NewOwnerEmail    string   `json:"new_owner_email"` // in order to withdraw
	Address     string `json:"address"`
	//PhoneNumber string `json:"phone"`
	//Email       string `json:"email"`
	// relationships
	//Employees   []Employee
	//Trips       []Trip
	//TechTeams   []TechTeam
}

func UpdateCompanyToCompany(req *UpdateCompanyReq, id uint) *company.TransportCompany{
	return &company.TransportCompany{
		ID: id,
		Name: req.Name,
		Description: req.Description,
		Address: req.Address,
	}
}

type CompanyReq struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"desc"`
	//OwnerID     uuid.UUID   `json:"owner_id" validate:"required"`
	Address     string `json:"address"`
	//PhoneNumber string `json:"phone"`
	//Email       string `json:"email"`
	// relationships
	//Employees   []Employee
	//Trips       []Trip
	//TechTeams   []TechTeam
}

type CompanyRes struct {
	ID          uint   `json:"id"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"desc"`
	OwnerID     uuid.UUID   `json:"owner_id" validate:"required"`
	Address     string `json:"address"`
	IsBlocked   bool   `json:"is_blocked"`
	//PhoneNumber string `json:"phone"`
	//Email       string `json:"email"`
	// relationships
	//Employees   []Employee
	//Trips       []Trip
	//TechTeams   []TechTeam
}


func CompanyReqToCompanyDomain(req *CompanyReq) *company.TransportCompany {
	return &company.TransportCompany{
		Name:        req.Name,
		Description: req.Description,
		//OwnerID:     req.OwnerID,
		Address:     req.Address,

		//PhoneNumber: req.PhoneNumber,
		//Email:       req.Email,
	}
}

func CompanyToCompanyRes(c company.TransportCompany) CompanyRes {
	return CompanyRes{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		OwnerID:     c.OwnerID,
		IsBlocked: c.IsBlocked,
		// owner: c.Owner !!!
		Address:     c.Address,
		//PhoneNumber: c.PhoneNumber,
		//Email:       c.Email,
	}
}

func BatchCompaniesToCompanies(companies []company.TransportCompany) []CompanyRes{
	return fp.Map(companies, CompanyToCompanyRes)
}
