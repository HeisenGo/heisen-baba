package presenter

import (
	"agency/internal/agency"
	"agency/pkg/fp"

	"github.com/google/uuid"
)

type CreateAgencyReq struct {
	OwnerID uuid.UUID `json:"owner_id" validate:"required" example:"aba3b3ed-e3d8-4403-9751-1f04287c9d65"`
	Name    string    `json:"name" validate:"required" example:"myagency"`
}

type AgencyResp struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type FullAgencyResponse struct {
	ID        uint      `json:"agency_id" example:"12"`
	OwnerID   uuid.UUID `json:"owner_id" example:"aba3b3ed-e3d8-4403-9751-1f04287c9d65"`
	Name      string    `json:"name" example:"myagency"`
	IsBlocked bool      `json:"is_blocked" example:"false"`
}

type UpdateAgencyReq struct {
	Name      *string `json:"name" example:"myagency"`
	IsBlocked *bool   `json:"is_blocked" example:"false"`
}

type CreateAgencyResponse struct {
	ID      uint      `json:"agency_id"`
	OwnerID uuid.UUID `json:"owner_id"`
	Name    string    `json:"name"`
}

func CreateAgencyRequest(req *CreateAgencyReq) *agency.Agency {
	a := &agency.Agency{
		OwnerID: req.OwnerID,
		Name:    req.Name,
	}
	return a
}

func AgencyToCreateAgencyResponse(a *agency.Agency) *CreateAgencyResponse {
	return &CreateAgencyResponse{
		ID:      a.ID,
		OwnerID: a.OwnerID,
		Name:    a.Name,
	}
}

func AgencyToFullAgencyResponse(a agency.Agency) FullAgencyResponse {
	return FullAgencyResponse{
		ID:        a.ID,
		OwnerID:   a.OwnerID,
		Name:      a.Name,
		IsBlocked: a.IsBlocked,
	}
}

func BatchAgenciesToAgencyResponse(agencies []agency.Agency) []FullAgencyResponse {
	return fp.Map(agencies, AgencyToFullAgencyResponse)
}

func UpdateAgencyRequestToDomain(req *UpdateAgencyReq) *agency.Agency {
	a := &agency.Agency{}
	if req.Name != nil {
		a.Name = *req.Name
	}
	if req.IsBlocked != nil {
		a.IsBlocked = *req.IsBlocked
	}
	return a
}
