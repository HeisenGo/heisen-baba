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
	ErrUnableToCancel  = errors.New("unable to cancel")
	ErrUnAuthorized = errors.New("not logged in")
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
	info := fmt.Sprintf("User %v bought %d trips of from %s %s to %s %s in %s for date %v", t.UserID, t.Quantity, trp.Origin, trp.Path.FromTerminal.Name, trp.Destination, trp.Path.ToTerminal.Name, trp.TripType, trp.StartDate)

	newInvoice := invoice.NewInvoice(t.ID, time.Now(), info, trp.UserPrice, t.TotalPrice, 0)

	err = s.invoiceOps.Create(ctx, newInvoice)
	t.Invoice = *newInvoice
	if err != nil {
		return err
	}

	//if invoice was successfull // TODO
	//****************************************
	trp.SoldTickets = trp.SoldTickets + uint(t.Quantity)
	newTrip := trip.NewTripTOUpdateSoldTickets(trp.SoldTickets)
	s.tripOps.UpdateTrip(ctx, trp.ID, newTrip, trp)
	t.Status = "Paid"
	newInvoice.Status = "Paid"                                   //sold tickets
	s.ticketOps.UpdateTicketStatus(ctx, t.ID, "Paid")       //state
	s.invoiceOps.UpdateInvoiceStatus(ctx, newInvoice.ID, "Paid") //state
	// notif with bank?
	return nil
}

func (s *TicketService) GetTicketsByUserOrAgency(ctx context.Context, userID uint, agencyID uint, page, pageSize uint) ([]ticket.Ticket, uint, error) {
	if userID == 0 && agencyID == 0{
		return nil, 0, ErrUnAuthorized
	}
	return s.ticketOps.GetTicketsByUserOrAgency(ctx, &userID, &agencyID, page, pageSize)
}

func (s *TicketService) CancelTicket(ctx context.Context, ticketID uint, userID *uint, agencyID *uint) (*invoice.Invoice, error) {
	// check permisson of requester if AgencyID!=nil!!!
	fullTicket, err := s.ticketOps.GetFullTicketByID(ctx, ticketID)
	
	if err != nil {
		return nil, err
	}
	var perTripCost float64
	if agencyID != nil {
		// check validation of agency TODO:
		if *fullTicket.AgencyID != *agencyID {
			return nil, ErrForbidden
		}
		if fullTicket.Status == "Canceled"{
			return &fullTicket.Invoice, fmt.Errorf("ticket is already")
		}
		perTripCost = fullTicket.Trip.AgencyPrice
		// id = ownerID
	} else {
		// check user TODO:
		fmt.Println(userID, fullTicket.UserID, *fullTicket.UserID == *userID)
		if *fullTicket.UserID != *userID {
			return nil, ErrForbidden
		}
		if fullTicket.Status == "Canceled"{
			return &fullTicket.Invoice, fmt.Errorf("ticket is already")
		}
		perTripCost = fullTicket.Trip.UserPrice
		// id = ownerID
	}

	fullTicket.Trip.TripCancellingPenalty.FirstDate = fullTicket.Trip.StartDate.AddDate(0, 0, int(-fullTicket.Trip.TripCancellingPenalty.FirstDays))
	fullTicket.Trip.TripCancellingPenalty.SecondDate = fullTicket.Trip.StartDate.AddDate(0, 0, int(-fullTicket.Trip.TripCancellingPenalty.SecondDays))
	fullTicket.Trip.TripCancellingPenalty.ThirdDate = fullTicket.Trip.StartDate.AddDate(0, 0, int(-fullTicket.Trip.TripCancellingPenalty.ThirdDays))

	first := fullTicket.Trip.TripCancellingPenalty.FirstDate
	second := fullTicket.Trip.TripCancellingPenalty.SecondDate
	third := fullTicket.Trip.TripCancellingPenalty.ThirdDate

	date := time.Now()
	var p uint
	if date.Before(first) {
		p = fullTicket.Trip.TripCancellingPenalty.FirstCancellationPercentage
	} else if date.Before(second) {
		p = fullTicket.Trip.TripCancellingPenalty.SecondCancellationPercentage
	} else if date.Before(third) {
		p = fullTicket.Trip.TripCancellingPenalty.ThirdCancellationPercentage
	} else {
		return nil, ErrUnableToCancel
	}
	penalty := (perTripCost * float64(p) / 100) * float64(fullTicket.Quantity)

	fullTicket.Penalty = penalty
	fullTicket.Status = "Canceled"
	inv, err := s.invoiceOps.GetInvoiceByTicketID(ctx, ticketID)
	if err != nil {
		return nil, err
	}
	inv.Penalty = penalty
	inv.Status = "Canceled"

	updates := make(map[string]interface{})
	updates["status"] = "Canceled"
	updates["penalty"] = penalty
	inv.Info = fmt.Sprintf("Ticket %s %s Canceled By User on %v", fullTicket.Trip.Origin, fullTicket.Trip.Destination, time.Now())

	err = s.ticketOps.UpdateTicket(ctx, ticketID, updates)
	if err != nil {
		return nil, err
	}
	updates["info"] = inv.Info

	err = s.invoiceOps.UpdateInvoice(ctx, inv.ID, updates)
	if err != nil {
		return nil, err
	}
	trip_updates := make(map[string]interface{})
	trip_updates["sold_tickets"] = fullTicket.Trip.SoldTickets - uint(fullTicket.Quantity)
	err = s.tripOps.UpdateTripTechTimID(ctx,fullTicket.TripID, trip_updates)
	if err!=nil{
		return inv, err
	}
	// send to bank TODO: from alibaba to UserID/ AgnecyID owner: id by company owner !
	return inv, nil
}
