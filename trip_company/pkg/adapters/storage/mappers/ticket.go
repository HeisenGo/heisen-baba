package mappers

import (
	"tripcompanyservice/internal/ticket"
	"tripcompanyservice/pkg/adapters/storage/entities"
	"tripcompanyservice/pkg/fp"
)

func TicketEntityToDomainWithTrip(ticketEntity entities.Ticket) ticket.Ticket {
	trip := SimpleTripEntityToDomain(*ticketEntity.Trip)
	return ticket.Ticket{
		ID:         ticketEntity.ID,
		TripID:     ticketEntity.TripID,
		Trip:       &trip,
		UserID:     ticketEntity.UserID,
		AgencyID:   ticketEntity.AgencyID,
		Quantity:   ticketEntity.Quantity,
		TotalPrice: ticketEntity.TotalPrice,
		Status:     ticketEntity.Status,
	}
}

func TicketEntitiesToDomainWithTrips(ticketEntities []entities.Ticket) []ticket.Ticket {
	return fp.Map(ticketEntities, TicketEntityToDomainWithTrip)
}

func TicketEntityToDomainWithTripWithCompanyWithPenalty(ticketEntity entities.Ticket) ticket.Ticket {
	trip := TripFullEntityToDomain(*ticketEntity.Trip)
	return ticket.Ticket{
		ID:         ticketEntity.ID,
		TripID:     ticketEntity.TripID,
		Trip:       &trip,
		UserID:     ticketEntity.UserID,
		AgencyID:   ticketEntity.AgencyID,
		Quantity:   ticketEntity.Quantity,
		TotalPrice: ticketEntity.TotalPrice,
		Status:     ticketEntity.Status,
	}
}

func TicketDomainToEntity(t *ticket.Ticket) *entities.Ticket {
	return &entities.Ticket{
		TripID:     t.TripID,
		UserID:     t.UserID,
		AgencyID:   t.AgencyID,
		Quantity:   t.Quantity,
		TotalPrice: t.TotalPrice,
		Status:     t.Status,
	}
}


func SimpleTicketEntityToDomain(ticketEntity entities.Ticket) ticket.Ticket {
	return ticket.Ticket{
		ID:         ticketEntity.ID,
		TripID:     ticketEntity.TripID,
		UserID:     ticketEntity.UserID,
		AgencyID:   ticketEntity.AgencyID,
		Quantity:   ticketEntity.Quantity,
		TotalPrice: ticketEntity.TotalPrice,
		Status:     ticketEntity.Status,
	}
}
