package mappers

import (
	"tripcompanyservice/internal/ticket"
	"tripcompanyservice/pkg/adapters/storage/entities"
	"tripcompanyservice/pkg/fp"
)

func TicketEntityToDomainWithTrip(ticketEntity entities.Ticket) ticket.Ticket {
	trip := SimpleTripEntityToDomainWithPenalty(*ticketEntity.Trip)
	invoice := SimpleInvoiceEntityToDomain(*ticketEntity.Invoice)
	return ticket.Ticket{
		Invoice: invoice,
		ID:         ticketEntity.ID,
		TripID:     ticketEntity.TripID,
		Trip:       &trip,
		UserID:     ticketEntity.UserID,
		AgencyID:   ticketEntity.AgencyID,
		Quantity:   ticketEntity.Quantity,
		TotalPrice: ticketEntity.TotalPrice,
		Status:     ticketEntity.Status,
		Penalty: ticketEntity.Penalty,
	}
}


func TicketInvoiceEntityToDomain(ticketEntity entities.Ticket) ticket.Ticket {
	invoice := SimpleInvoiceEntityToDomain(*ticketEntity.Invoice)
	return ticket.Ticket{
		Invoice: invoice,
		ID:         ticketEntity.ID,
		TripID:     ticketEntity.TripID,
		UserID:     ticketEntity.UserID,
		AgencyID:   ticketEntity.AgencyID,
		Quantity:   ticketEntity.Quantity,
		TotalPrice: ticketEntity.TotalPrice,
		Status:     ticketEntity.Status,
		Penalty: ticketEntity.Penalty,
	}
}

func TicketEntitiesToDomainWithTrips(ticketEntities []entities.Ticket) []ticket.Ticket {
	return fp.Map(ticketEntities, TicketEntityToDomainWithTrip)
}

func TicketEntityToDomainWithTripWithCompanyWithPenaltyWithInvoice(ticketEntity entities.Ticket) ticket.Ticket {
	invoice := SimpleInvoiceEntityToDomain(*ticketEntity.Invoice)
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
		Penalty: ticketEntity.Penalty,
		Invoice: invoice,
	}
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
		Penalty: ticketEntity.Penalty,
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
		Penalty: t.Penalty,
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
		Penalty: ticketEntity.Penalty,
	}
}

func BatchTicketEntitiesToTickets(t []entities.Ticket)[]ticket.Ticket{
	return fp.Map(t, TicketInvoiceEntityToDomain)
}