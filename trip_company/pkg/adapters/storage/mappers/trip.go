package mappers

import (
	"tripcompanyservice/internal/trip"
	"tripcompanyservice/pkg/adapters/storage/entities"
	"tripcompanyservice/pkg/fp"
)

func TripEntityToDomain(tripEntity entities.Trip) trip.Trip {
	path := &trip.Path{
		ID:   tripEntity.PathID,
		Name: tripEntity.PathName,
		Type: tripEntity.TripType,
		FromTerminal: &trip.Terminal{
			City: tripEntity.Origin,
			Name: tripEntity.FromTerminalName,
		},
		ToTerminal: &trip.Terminal{
			City: tripEntity.Destination,
			Name: tripEntity.ToTerminalName,
		},
	}
	penalty := PenaltyEntityToDomain(*tripEntity.TripCancelingPenalty)
	company := CompanyEntityToDomain(*tripEntity.TransportCompany)
	return trip.Trip{
		ID:                     tripEntity.ID,
		TransportCompanyID:     tripEntity.TransportCompanyID,
		TransportCompany:       company,
		TripCancellingPenalty:  &penalty,
		TripType:               trip.TripType(tripEntity.TripType),
		UserReleaseDate:        tripEntity.UserReleaseDate,
		TourReleaseDate:        tripEntity.TourReleaseDate,
		UserPrice:              tripEntity.UserPrice,
		AgencyPrice:            tripEntity.AgencyPrice,
		PathID:                 tripEntity.PathID,
		Origin:                 tripEntity.Origin,
		Destination:            tripEntity.Destination,
		Path:                   path,
		Status:                 tripEntity.Status,
		MinPassengers:          tripEntity.MinPassengers,
		TechTeamID:             tripEntity.TechTeamID,
		VehicleRequestID:       tripEntity.VehicleRequestID,
		TripCancelingPenaltyID: tripEntity.TripCancellingPenaltyID,
		MaxTickets:             tripEntity.MaxTickets,
		VehicleID:              tripEntity.VehicleID,
		IsCanceled:             tripEntity.IsCanceled,
		IsFinished:             tripEntity.IsFinished,
		StartDate:              tripEntity.StartDate,
		EndDate:                tripEntity.EndDate,
		SoldTickets:            tripEntity.SoldTickets,
	}
}

func TripFullEntityToDomain(tripEntity entities.Trip) trip.Trip {
	path := &trip.Path{
		ID:   tripEntity.PathID,
		Name: tripEntity.PathName,
		FromTerminal: &trip.Terminal{
			City: tripEntity.Origin,
			Name: tripEntity.FromTerminalName,
		},
		ToTerminal: &trip.Terminal{
			City: tripEntity.Destination,
			Name: tripEntity.ToTerminalName,
		},
	}
	penalty := PenaltyEntityToDomain(*tripEntity.TripCancelingPenalty)
	company := CompanyEntityToDomain(*tripEntity.TransportCompany)
	return trip.Trip{
		ID:                     tripEntity.ID,
		TransportCompanyID:     tripEntity.TransportCompanyID,
		TransportCompany:       company,
		TripType:               trip.TripType(tripEntity.TripType),
		UserReleaseDate:        tripEntity.UserReleaseDate,
		TourReleaseDate:        tripEntity.TourReleaseDate,
		UserPrice:              tripEntity.UserPrice,
		AgencyPrice:            tripEntity.AgencyPrice,
		PathID:                 tripEntity.PathID,
		Origin:                 tripEntity.Origin,
		Destination:            tripEntity.Destination,
		Path:                   path,
		Status:                 tripEntity.Status,
		MinPassengers:          tripEntity.MinPassengers,
		TechTeamID:             tripEntity.TechTeamID,
		VehicleRequestID:       tripEntity.VehicleRequestID,
		TripCancelingPenaltyID: tripEntity.TripCancellingPenaltyID,
		MaxTickets:             tripEntity.MaxTickets,
		VehicleID:              tripEntity.VehicleID,
		IsCanceled:             tripEntity.IsCanceled,
		IsFinished:             tripEntity.IsFinished,
		StartDate:              tripEntity.StartDate,
		EndDate:                tripEntity.EndDate,
		SoldTickets:            tripEntity.SoldTickets,
		TripCancellingPenalty:  &penalty,
	}
}

func TripEntitiesToDomain(tripEntities []entities.Trip) []trip.Trip {
	return fp.Map(tripEntities, TripEntityToDomain)
}

func TripDomainToEntity(t *trip.Trip) *entities.Trip {
	tCp := PenaltyDomainToEntity(t.TripCancellingPenalty)
	return &entities.Trip{
		TransportCompanyID:   t.TransportCompanyID,
		TripType:             string(t.TripType),
		UserReleaseDate:      t.UserReleaseDate,
		TourReleaseDate:      t.TourReleaseDate,
		UserPrice:            t.UserPrice,
		AgencyPrice:          t.AgencyPrice,
		PathID:               t.PathID,
		Origin:               t.Origin,
		Destination:          t.Destination,
		FromTerminalName:     t.Path.FromTerminal.Name,
		ToTerminalName:       t.Path.ToTerminal.Name,
		PathName:             t.Path.Name,
		Status:               t.Status,
		MinPassengers:        t.MinPassengers,
		TechTeamID:           t.TechTeamID,
		VehicleRequestID:     t.VehicleRequestID,
		TripCancelingPenalty: tCp,
		MaxTickets:           t.MaxTickets,
		VehicleID:            t.VehicleID,
		IsCanceled:           t.IsCanceled,
		IsFinished:           t.IsFinished,
		StartDate:            t.StartDate,
		EndDate:              t.EndDate,
		SoldTickets:          t.SoldTickets,
	}
}
