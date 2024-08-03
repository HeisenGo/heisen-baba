package handlers

import (
	"encoding/json"
	"fmt"
	"tripcompanyservice/api/http/handlers/presenter"
	"tripcompanyservice/service"

	"github.com/gofiber/fiber/v2"
)

// TODO: transactional
func BuyTicket(ticketService *service.TicketService) fiber.Handler { //serviceFactory ServiceFactory[*service.TripService])fiber.Handler {
	return func(c *fiber.Ctx) error {
		//tripService := serviceFactory(c.UserContext())

		var body map[string]interface{}
		if err := c.BodyParser(&body); err != nil {
			return presenter.BadRequest(c, fmt.Errorf("Invalid req body"))
		}
		var res interface{}

		if _, ok := body["agency_id"]; ok {

			var req *presenter.AgencyTicketReq
			if err := json.Unmarshal(c.Body(), &req); err != nil {
				return presenter.BadRequest(c, fmt.Errorf("Invalid req body"))
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
				return presenter.BadRequest(c, fmt.Errorf("Invalid req body"))
			}

			//userID := getUserID(c) // TODO: user ID from authentication
			ticket := presenter.UserTicketReqToTicket(req)
			if err := ticketService.ProcessUserTicket(c.UserContext(), ticket); err != nil {
				// err = "failed to process agency ticket"
				return presenter.BadRequest(c, err)
			}
			res = presenter.TicketToUserTicket(*ticket)
		}
		return presenter.Created(c, "Ticket created successfully", res)
	}
}

func GetTickets(ticketService *service.TicketService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// userClaims, ok := c.Locals(UserClaimKey).(*jwt.UserClaims)
		// if !ok {
		// 	return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		// }
		//query parameter
		page, pageSize := PageAndPageSize(c)
		// get from auth!!!!! TODO:
		// Parse dates
		//var err error
		//companyID, err := c.ParamsInt("companyID")
		UserID := uint(1)
		var data interface{}
		if UserID != 0 { // TODO :

			tickets, total, err := ticketService.GetTicketsByUserOrAgency(c.UserContext(), &UserID, nil, uint(page), uint(pageSize))
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
		} else {
			agencyID := uint(1)
			tickets, total, err := ticketService.GetTicketsByUserOrAgency(c.UserContext(), nil, &agencyID, uint(page), uint(pageSize))
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
		}

		return presenter.OK(c, "Tickets fetched successfully", data)
	}
}

func CancelTicketByID(ticketService *service.TicketService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// userClaims, ok := c.Locals(UserClaimKey).(*jwt.UserClaims)
		// if !ok {
		// 	return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		// }
		//query parameter
		ticketID, err := c.ParamsInt("ticketID")
		if err != nil {
			return presenter.BadRequest(c, errWrongIDType)
		}

		if ticketID < 0 {
			return presenter.BadRequest(c, errWrongIDType)
		}
		// User ID is needed!!! TODO!!

		var body map[string]interface{}
		if err := c.BodyParser(&body); err != nil {
			return presenter.BadRequest(c, fmt.Errorf("invalid req body"))
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
			// USER ID is needed ! TODO :
			userID := uint(1)
			invoice, err := ticketService.CancelTicket(c.UserContext(), uint(ticketID), &userID, nil)
			if err != nil { //err = "failed to process agency ticket"
				return presenter.BadRequest(c, err)
			}
			res = presenter.InvoiceToAgencyInvoice(*invoice)
		}

		return presenter.OK(c, "Ticket Canceled successfully", res)
	}
}
