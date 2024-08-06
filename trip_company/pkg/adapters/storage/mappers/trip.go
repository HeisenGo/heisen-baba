package mappers

import (
	"tripcompanyservice/internal/company"
	"tripcompanyservice/internal/techteam"
	"tripcompanyservice/internal/trip"
	tripcancellingpenalty "tripcompanyservice/internal/trip_cancelling_penalty"
	vehiclerequest "tripcompanyservice/internal/vehicle_request"
	"tripcompanyservice/pkg/adapters/storage/entities"
	"tripcompanyservice/pkg/fp"
)

func TripEntityToDomain(tripEntity entities.Trip) trip.Trip {
	path := &trip.Path{
		ID:         tripEntity.PathID,
		Name:       tripEntity.PathName,
		Type:       tripEntity.TripType,
		DistanceKM: tripEntity.PathDistanceKM,
		FromTerminal: &trip.Terminal{
			City: tripEntity.Origin,
			Name: tripEntity.FromTerminalName,
		},
		ToTerminal: &trip.Terminal{
			City: tripEntity.Destination,
			Name: tripEntity.ToTerminalName,
		},
	}
	var penalty tripcancellingpenalty.TripCancelingPenalty
	if tripEntity.TripCancelingPenalty != nil {
		penalty = PenaltyEntityToDomain(*tripEntity.TripCancelingPenalty)
	}
	var company company.TransportCompany
	if tripEntity.TransportCompany != nil {
		company = CompanyEntityToDomain(*tripEntity.TransportCompany)
	}
	var vR vehiclerequest.VehicleRequest
	if tripEntity.VehicleRequest != nil {
		vR = VehicleReqEntityToVehicleReqDomain(*tripEntity.VehicleRequest)
	}
	var team techteam.TechTeam
	if tripEntity.TechTeam != nil {
		team = TechTeamEntityToDomain(*tripEntity.TechTeam)
	}
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
		IsConfirmed:            tripEntity.IsConfirmed,
		TechTeam:               &team,
		VehicleRequest:         &vR,
		Profit:                 tripEntity.Profit,
	}
}

func TripFullEntityToDomain(tripEntity entities.Trip) trip.Trip {
	path := &trip.Path{
		ID:         tripEntity.PathID,
		Name:       tripEntity.PathName,
		DistanceKM: tripEntity.PathDistanceKM,
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
	var vR vehiclerequest.VehicleRequest
	if tripEntity.VehicleRequest != nil {
		vR = VehicleReqEntityToVehicleReqDomain(*tripEntity.VehicleRequest)
	}
	var team techteam.TechTeam
	if tripEntity.TechTeam != nil {
		team = TechTeamEntityToDomain(*tripEntity.TechTeam)
	}
	return trip.Trip{
		ID:                     tripEntity.ID,
		TransportCompanyID:     tripEntity.TransportCompanyID,
		TransportCompany:       company,
		VehicleRequest:         &vR,
		VehicleRequestID:       tripEntity.VehicleRequestID,
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
		TechTeam:               &team,
		TripCancelingPenaltyID: tripEntity.TripCancellingPenaltyID,
		MaxTickets:             tripEntity.MaxTickets,
		VehicleID:              tripEntity.VehicleID,
		IsCanceled:             tripEntity.IsCanceled,
		IsFinished:             tripEntity.IsFinished,
		StartDate:              tripEntity.StartDate,
		EndDate:                tripEntity.EndDate,
		SoldTickets:            tripEntity.SoldTickets,
		TripCancellingPenalty:  &penalty,
		IsConfirmed:            tripEntity.IsConfirmed,
		Profit:                 tripEntity.Profit,
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
		PathDistanceKM:       t.Path.DistanceKM,
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
		IsConfirmed:          t.IsConfirmed,
		Profit:               t.Profit,
	}
}

func SimpleTripEntityToDomain(tripEntity entities.Trip) trip.Trip {
	path := &trip.Path{
		ID:         tripEntity.PathID,
		Name:       tripEntity.PathName,
		Type:       tripEntity.TripType,
		DistanceKM: tripEntity.PathDistanceKM,
		FromTerminal: &trip.Terminal{
			City: tripEntity.Origin,
			Name: tripEntity.FromTerminalName,
			Type: tripEntity.TripType,
		},
		ToTerminal: &trip.Terminal{
			City: tripEntity.Destination,
			Name: tripEntity.ToTerminalName,
			Type: tripEntity.TripType,
		},
	}
	return trip.Trip{
		ID:                     tripEntity.ID,
		TransportCompanyID:     tripEntity.TransportCompanyID,
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
		IsConfirmed:            tripEntity.IsConfirmed,
		Profit:                 tripEntity.Profit,
	}
}

func SimpleTripEntityToDomainWithPenalty(tripEntity entities.Trip) trip.Trip {
	path := &trip.Path{
		ID:         tripEntity.PathID,
		Name:       tripEntity.PathName,
		Type:       tripEntity.TripType,
		DistanceKM: tripEntity.PathDistanceKM,
		FromTerminal: &trip.Terminal{
			City: tripEntity.Origin,
			Name: tripEntity.FromTerminalName,
			Type: tripEntity.TripType,
		},
		ToTerminal: &trip.Terminal{
			City: tripEntity.Destination,
			Name: tripEntity.ToTerminalName,
			Type: tripEntity.TripType,
		},
	}
	penalty := PenaltyEntityToDomain(*tripEntity.TripCancelingPenalty)
	return trip.Trip{
		TripCancellingPenalty:  &penalty,
		ID:                     tripEntity.ID,
		TransportCompanyID:     tripEntity.TransportCompanyID,
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
		IsConfirmed:            tripEntity.IsConfirmed,
		Profit:                 tripEntity.Profit,
	}
}
