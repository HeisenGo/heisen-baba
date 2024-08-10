package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"tripcompanyservice/api/http/handlers/presenter"
	"tripcompanyservice/internal/trip"
	"tripcompanyservice/internal/user"
	"tripcompanyservice/pkg/valuecontext"
	"tripcompanyservice/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func BuyTicket(serviceFactory ServiceFactory[*service.TicketService]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ticketService := serviceFactory(c.UserContext())

		userReq, ok := c.Locals(valuecontext.UserClaimKey).(*user.User)
		if !ok {
			return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		}

		var body map[string]interface{}
		if err := c.BodyParser(&body); err != nil {
			return presenter.BadRequest(c, fmt.Errorf("invalid req body/ body should be even {}"))
		}
		var res interface{}

		if _, ok := body["agency_id"]; ok {

			var req *presenter.AgencyTicketReq
			if err := json.Unmarshal(c.Body(), &req); err != nil {
				return presenter.BadRequest(c, fmt.Errorf("invalid req body"))
			}

			ticket := presenter.AgencyTicketReqToTicket(req)

			if err := ticketService.ProcessAgencyTicket(c.UserContext(), ticket); err != nil {
				//err = "failed to process agency ticket"
				return presenter.BadRequest(c, err)
			}
			res = presenter.TicketToAgencyTicket(*ticket)
		} else {
			var req *presenter.UserTicketReq
			if err := json.Unmarshal(c.Body(), &req); err != nil {
				return presenter.BadRequest(c, fmt.Errorf("invalid req body"))
			}
			//userID := getUserID(c) // TODO: user ID from authentication
			ticket := presenter.UserTicketReqToTicket(req)
			ticket.UserID = &userReq.ID
			if err := ticketService.ProcessUserTicket(c.UserContext(), ticket); err != nil {
				if errors.Is(err, trip.ErrTripNotFound) || errors.Is(err, service.ErrImpossibleToBuy) {
					return presenter.BadRequest(c, err)
				}
				return presenter.BadRequest(c, err)
			}
			res = presenter.TicketToUserTicket(*ticket)
		}
		return presenter.Created(c, "Ticket created successfully", res)
	}
}

func GetAgencyTickets(ticketService *service.TicketService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		agencyID, err := c.ParamsInt("agencyID")
		if err != nil {
			return presenter.BadRequest(c, errWrongIDType)
		}

		if agencyID < 0 {
			return presenter.BadRequest(c, errWrongIDType)
		}
		//query parameter
		page, pageSize := PageAndPageSize(c)

		// requester TODO IT Is from a user from agency!!!
		var data interface{}
		tickets, total, err := ticketService.GetTicketsByUserOrAgency(c.UserContext(), uuid.Nil, uint(agencyID), uint(page), uint(pageSize))
		if err != nil {
			//err := errors.New("Error")
			return presenter.InternalServerError(c, err)
		}
		res := presenter.BatchTicketsToUserTickets(tickets)

		data = presenter.NewPagination(
			res,
			uint(page),
			uint(pageSize),
			total,
		)

		return presenter.OK(c, "Tickets fetched successfully", data)
	}
}

func GetUserTickets(ticketService *service.TicketService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		page, pageSize := PageAndPageSize(c)

		userReq, ok := c.Locals(valuecontext.UserClaimKey).(*user.User)
		if !ok {
			return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		}
		var data interface{}

		tickets, total, err := ticketService.GetTicketsByUserOrAgency(c.UserContext(), userReq.ID, uint(0), uint(page), uint(pageSize))
		if err != nil {
			//err := errors.New("Error")
			return presenter.InternalServerError(c, err)
		}
		res := presenter.BatchTicketsToUserTickets(tickets)

		data = presenter.NewPagination(
			res,
			uint(page),
			uint(pageSize),
			total,
		)

		return presenter.OK(c, "Tickets fetched successfully", data)
	}
}

func CancelTicketByID(serviceFactory ServiceFactory[*service.TicketService]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ticketService := serviceFactory(c.UserContext())

		ticketID, err := c.ParamsInt("ticketID")
		if err != nil {
			return presenter.BadRequest(c, errWrongIDType)
		}

		if ticketID < 0 {
			return presenter.BadRequest(c, errWrongIDType)
		}

		userReq, ok := c.Locals(valuecontext.UserClaimKey).(*user.User)
		if !ok {
			return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		}
		var body map[string]interface{}
		if err := c.BodyParser(&body); err != nil {
			return presenter.BadRequest(c, fmt.Errorf("invalid req body - for user body should be {}"))
		}
		var res interface{}

		if _, ok := body["agency_id"]; ok {

			var req *presenter.CancelTicket
			if err := json.Unmarshal(c.Body(), &req); err != nil {
				return presenter.BadRequest(c, fmt.Errorf("invalid req body"))
			}

			invoice, err := ticketService.CancelTicket(c.UserContext(), uint(ticketID), nil, &req.AgencyID)
			if err != nil { //err = "failed to process agency ticket"
				return presenter.BadRequest(c, err)
			}
			res = presenter.InvoiceToAgencyInvoice(*invoice)
		} else {
			invoice, err := ticketService.CancelTicket(c.UserContext(), uint(ticketID), &userReq.ID, nil)
			if err != nil { //err = "failed to process agency ticket"
				return presenter.BadRequest(c, err)
			}
			res = presenter.InvoiceToAgencyInvoice(*invoice)
		}

		return presenter.OK(c, "Ticket Canceled successfully", res)
	}
}
