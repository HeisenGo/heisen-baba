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
