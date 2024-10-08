package service

import (
	"context"
	"errors"
	"fmt"
	"time"
	"tripcompanyservice/internal/invoice"
	"tripcompanyservice/internal/ticket"
	"tripcompanyservice/internal/trip"
	"tripcompanyservice/pkg/ports/clients/clients"

	"github.com/google/uuid"
)

var (
	ErrImpossibleToBuy = errors.New("not possible to buy")
	ErrUnableToCancel  = errors.New("unable to cancel")
	ErrUnAuthorized    = errors.New("not logged in")
)

type TicketService struct {
	ticketOps  *ticket.Ops
	tripOps    *trip.Ops
	invoiceOps *invoice.Ops
	bankClient clients.IBankClient
}

func NewTicketService(ticketOps *ticket.Ops, tripOps *trip.Ops, invoiceOps *invoice.Ops, bclient clients.IBankClient) *TicketService {
	return &TicketService{
		ticketOps:  ticketOps,
		tripOps:    tripOps,
		invoiceOps: invoiceOps,
		bankClient: bclient,
	}
}

func (s *TicketService) ProcessAgencyTicket(ctx context.Context, t *ticket.Ticket) error {
	trp, err := s.tripOps.GetFullTripByID(ctx, t.TripID)
	if err != nil {
		return err
	}
	t.Trip = trp
	// check agency exist !!! TODO
	// get agency owner wallet!
	if trp.TourReleaseDate.After(time.Now()) {
		return ErrImpossibleToBuy
	}
	if trp.SoldTickets+uint(t.Quantity) > trp.MaxTickets {
		return ErrImpossibleToBuy
	}
	t.TotalPrice = trp.AgencyPrice * float64(t.Quantity)
	// t.UserID agency owner ID
	err = s.ticketOps.Create(ctx, t)
	if err != nil {
		return err
	}
	info := fmt.Sprintf("Agency %v bought %d trips of from %s %s to %s %s in %s for date %v", t.AgencyID, t.Quantity, trp.Origin, trp.Path.FromTerminal.Name, trp.Destination, trp.Path.ToTerminal.Name, trp.TripType, trp.StartDate)

	newInvoice := invoice.NewInvoice(t.ID, time.Now(), info, trp.AgencyPrice, t.TotalPrice, 0)

	err = s.invoiceOps.Create(ctx, newInvoice)
	t.Invoice = *newInvoice
	if err != nil {
		return err
	}
	//*********************************
	y, err := s.bankClient.Transfer(trp.TransportCompany.OwnerID.String(), "", true, uint64(t.TotalPrice))
	if err != nil {
		return err
	}
	if !y {
		return errors.New("unsuccessful pay")
	}
	//if invoice was successfull // TODO
	//****************************************
	trp.SoldTickets = trp.SoldTickets + uint(t.Quantity)
	newTrip := trip.NewTripTOUpdateSoldTickets(trp.SoldTickets)
	s.tripOps.UpdateTrip(ctx, trp.ID, newTrip, trp)
	t.Status = "Paid"
	newInvoice.Status = "Paid"                                   //sold tickets
	s.ticketOps.UpdateTicketStatus(ctx, t.ID, "Paid")            //state
	s.invoiceOps.UpdateInvoiceStatus(ctx, newInvoice.ID, "Paid") //state
	// notif with bank?
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
	s.ticketOps.UpdateTicketStatus(ctx, t.ID, "Paid")            //state
	s.invoiceOps.UpdateInvoiceStatus(ctx, newInvoice.ID, "Paid") //state
	// notif with bank?
	return nil
}

func (s *TicketService) GetTicketsByUserOrAgency(ctx context.Context, userID uuid.UUID, agencyID uint, page, pageSize uint) ([]ticket.Ticket, uint, error) {
	if userID == uuid.Nil && agencyID == 0 {
		return nil, 0, ErrUnAuthorized
	}
	if userID == uuid.Nil {
		// TODO is this person from the agency?
		return s.ticketOps.GetTicketsByUserOrAgency(ctx, nil, &agencyID, page, pageSize)
	} else {
		return s.ticketOps.GetTicketsByUserOrAgency(ctx, &userID, nil, page, pageSize)
	}
}

func (s *TicketService) CancelTicket(ctx context.Context, ticketID uint, userID *uuid.UUID, agencyID *uint) (*invoice.Invoice, error) {
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
		if fullTicket.Status == "Canceled" {
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
		if fullTicket.Status == "Canceled" {
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
	err = s.tripOps.UpdateTripTechTimID(ctx, fullTicket.TripID, trip_updates)
	if err != nil {
		return inv, err
	}
	// send to bank TODO: from alibaba to UserID/ AgnecyID owner: id by company owner !
	return inv, nil
}
