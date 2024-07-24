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
	ID         uint    `json:"id"`
	DistanceKM float64 `json:"DistanceKM"` // in kilometers
	Code       string  `json:"code"`
	Name       string  `json:"name"`
	Type       string  `json:"type"`
}

func PathToPathResponse(p path.Path) PathResponse {
	return PathResponse{
		ID:         p.ID,
		DistanceKM: p.DistanceKM, // in kilometers
		Code:       p.Code,
		Name:       p.Name,
		Type:       string(p.Type),
	}
}

func PathsToPathRequests(paths []path.Path) []PathResponse {
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
