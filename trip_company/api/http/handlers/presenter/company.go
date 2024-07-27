package presenter

import "tripcompanyservice/internal/company"

type CompanyReq struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"desc"`
	OwnerID     uint   `json:"owner_id" validate:"required"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone"`
	Email       string `json:"email"`
	// relationships
	//Employees   []Employee
	//Trips       []Trip
	//TechTeams   []TechTeam
}

type CompanyRes struct {
	ID          uint   `json:"id"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"desc"`
	OwnerID     uint   `json:"owner_id" validate:"required"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone"`
	Email       string `json:"email"`
	// relationships
	//Employees   []Employee
	//Trips       []Trip
	//TechTeams   []TechTeam
}

func CompanyReqToCompanyDomain(req *CompanyReq) *company.TransportCompany {
	return &company.TransportCompany{
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     req.OwnerID,
		Address:     req.Address,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
	}
}

func CompanyToCompanyRes(c company.TransportCompany) CompanyRes {
	return CompanyRes{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		OwnerID:     c.OwnerID,
		// owner: c.Owner !!!
		Address:     c.Address,
		PhoneNumber: c.PhoneNumber,
		Email:       c.Email,
	}
}
