package mappers

import (
	"tripcompanyservice/internal/trip"
	"tripcompanyservice/protobufs"
)

func GRPCPathToDomain(p *protobufs.GetFullPathByIDResponse) (*trip.Path, error) {
	return &trip.Path{
		FromTerminal: &trip.Terminal{
			City:    p.Data.FromCity,
			Country: p.Data.FromCountry,
			Name:    p.Data.FromTerminalName,
			Type:    p.Data.Type,
		},
		ToTerminal: &trip.Terminal{
			City:    p.Data.ToCity,
			Country: p.Data.ToCountry,
			Name:    p.Data.ToTerminalName,
			Type:    p.Data.Type,
		},
		ID:         uint(p.Data.Id),
		DistanceKM: p.Data.DistanceKm,
		Code:       p.Data.Code,
		Name:       p.Data.Name,
		Type:       p.Data.Type,
	}, nil
}
