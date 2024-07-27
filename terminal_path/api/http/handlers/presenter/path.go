package presenter

import (
	"terminalpathservice/internal/path"
	"terminalpathservice/pkg/fp"
)

type PathRequest struct {
	FromTerminalID uint    `json:"from_terminal_id" validate:"required"`
	ToTerminalID   uint    `json:"to_terminal_id" validate:"required"`
	DistanceKM     float64 `json:"distance_km" validate:"required"` // in kilometers
	Code           string  `json:"code" validate:"required"`
	Name           string  `json:"name" validate:"required"`
}

type PathResponse struct {
	ID               uint    `json:"id"`
	DistanceKM       float64 `json:"distance_km"`
	Code             string  `json:"code"`
	Name             string  `json:"name"`
	Type             string  `json:"type"`
	FromCountry      string  `json:"from_country"`
	FromCity         string  `json:"from"`
	FromTerminalName string  `json:"from_terminal"`
	ToCountry        string  `json:"to_country"`
	ToCity           string  `json:"to"`
	ToTerminalName   string  `json:"to_terminal"`
}

func PathToPathResponse(p path.Path) PathResponse {
	return PathResponse{
		ID:               p.ID,
		DistanceKM:       p.DistanceKM, // in kilometers
		Code:             p.Code,
		Name:             p.Name,
		Type:             string(p.Type),
		FromCountry:      p.FromTerminal.Country,
		ToCountry:        p.ToTerminal.Country,
		FromCity:         p.FromTerminal.City,
		ToCity:           p.ToTerminal.City,
		FromTerminalName: p.FromTerminal.Name,
		ToTerminalName:   p.ToTerminal.Name,
	}
}

type UpdatePathRequest struct {
	ID             uint    `json:"terminal_id"`
	Name           string  `json:"name"`
	Code           string  `json:"code"`
	DistanceKM     float64 `json:"distance_km"`
	FromTerminalID uint    `json:"from_terminal_id"`
	ToTerminalID   uint    `json:"to_terminal_id"`
}

func UpdatePathReqToPath(p *UpdatePathRequest, id uint) *path.Path {
	return &path.Path{
		ID:             id,
		Name:           p.Name,
		FromTerminalID: p.FromTerminalID,
		ToTerminalID:   p.ToTerminalID,
		DistanceKM:     p.DistanceKM,
		Code:           p.Code,
	}
}

func BatchPathsToPathResponse(paths []path.Path) []PathResponse {
	return fp.Map(paths, PathToPathResponse)
}

func PathRequestToPath(pathReq *PathRequest) *path.Path {
	return &path.Path{
		FromTerminalID: pathReq.FromTerminalID,
		ToTerminalID:   pathReq.ToTerminalID,
		DistanceKM:     pathReq.DistanceKM,
		Code:           pathReq.Code,
		Name:           pathReq.Name,
	}
}
