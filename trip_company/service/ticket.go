package service

import (
	"context"
	"errors"
	"fmt"
	"time"
	"tripcompanyservice/internal/invoice"
	"tripcompanyservice/internal/ticket"
	"tripcompanyservice/internal/trip"
)

var (
	ErrImpossibleToBuy = errors.New("not possible to buy")
)

type TicketService struct {
	ticketOps  *ticket.Ops
	tripOps    *trip.Ops
	invoiceOps *invoice.Ops
}

func NewTicketService(ticketOps *ticket.Ops, tripOps *trip.Ops, invoiceOps *invoice.Ops) *TicketService {
	return &TicketService{
		ticketOps:  ticketOps,
		tripOps:    tripOps,
		invoiceOps: invoiceOps,
	}
}

func (s *TicketService) ProcessAgencyTicket(ctx context.Context, t *ticket.Ticket) error {
	// 	trip, err := s.tripOps.GetFullTripByID(ctx, t.TripID)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	// check agency exist !!! TODO
	// 	// get agency owner wallet!
	// 	if trip.TourReleaseDate.After(time.Now()) {
	// 		return ErrImpossibleToBuy
	// 	}
	// 	t.TotalPrice = trip.AgencyPrice * float64(t.Quantity)
	// 	// t.UserID agency owner ID
	// 	err = s.ticketOps.CreateTicket(ctx, t)
	// 	agency_id := 1
	// 	info := fmt.Sprintf("Agency %s bought %d trips of from %s %s to %s %s in %s for date %v", agency_id, t.Quantity, trip.Origin, trip.Path.FromTerminal.Name, trip.Destination, trip.Path.ToTerminal.Name, trip.TripType, trip.StartDate)

	// 	newInvoice := invoice.NewInvoice(t.ID, t, time.Now(), info, trip.AgencyPrice, t.TotalPrice)

	// 	s.invoiceOps.CreateInvoice(ctx, newInvoice)

	// 	//if invoice was successfull // TODO
	// 	//****************************************
	// 	trip.SoldTickets = trip.SoldTickets + float64(t.Quantity)
	// 	newTrip := trip.NewTripTOUpdateSoldTickets(trip.SoldTickets)
	// 	s.tripOps.UpdateTrip(ctx, trip.ID, newTrip, trip)
	// 	s.ticketOps.UpdateTicket()
	// 	s.invoiceOps.UpdateInvoice()
	// 	// notif
	// 	return nil
	// }

	// func (s *TicketService) ProcessUserTicket(ctx context.Context, t *ticket.Ticket) error {
	return nil
}

func (s *TicketService) ProcessUserTicket(ctx context.Context, t *ticket.Ticket) error {
	trp, err := s.tripOps.GetFullTripByID(ctx, t.TripID)
	if err != nil {
		return err
	}

	t.Trip = trp

	// check user exist  and get user info!!! TODO
	if trp.UserReleaseDate.After(time.Now()) {
		return ErrImpossibleToBuy
	}

	if trp.SoldTickets+uint(t.Quantity) > trp.MaxTickets {
		return ErrImpossibleToBuy
	}
	t.TotalPrice = trp.UserPrice * float64(t.Quantity)
	// t.UserID agency owner ID
	err = s.ticketOps.Create(ctx, t)
	if err != nil {
		return err
	}
	username := "new_user"
	info := fmt.Sprintf("User %s bought %d trips of from %s %s to %s %s in %s for date %v", username, t.Quantity, trp.Origin, trp.Path.FromTerminal.Name, trp.Destination, trp.Path.ToTerminal.Name, trp.TripType, trp.StartDate)

	newInvoice := invoice.NewInvoice(t.ID, t, time.Now(), info, trp.UserPrice, t.TotalPrice)

	err = s.invoiceOps.Create(ctx, newInvoice)
	if err != nil {
		return err
	}

	//if invoice was successfull // TODO
	//****************************************
	trp.SoldTickets = trp.SoldTickets + uint(t.Quantity)
	newTrip := trip.NewTripTOUpdateSoldTickets(trp.SoldTickets)
	s.tripOps.UpdateTrip(ctx, trp.ID, newTrip, trp)
	t.Status = "Confirmed"
	newInvoice.Status = "Paid"                                   //sold tickets
	s.ticketOps.UpdateTicketStatus(ctx, t.ID, "Confirmed")       //state
	s.invoiceOps.UpdateInvoiceStatus(ctx, newInvoice.ID, "Paid") //state
	// notif with bank?
	return nil
}

func (s *TicketService) GetTicketsByUserOrAgency(ctx context.Context, userID *uint, agencyID *uint, page, pageSize uint) ([]ticket.Ticket, uint, error) {
	// check one of them should be nill !!!
	return s.ticketOps.GetTicketsByUserOrAgency(ctx, userID, agencyID, page, pageSize)
}
