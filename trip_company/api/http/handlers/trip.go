package handlers

import (
	"tripcompanyservice/api/http/handlers/presenter"
	"tripcompanyservice/service"

	"github.com/gofiber/fiber/v2"
)

func CreateTrip(tripService *service.TripService)fiber.Handler {//serviceFactory ServiceFactory[*service.TripService])fiber.Handler {
	return func(c *fiber.Ctx) error {
		//tripService := serviceFactory(c.UserContext())

		var req presenter.CreateTripReq

		if err := c.BodyParser(&req); err != nil {
			return presenter.BadRequest(c, err)
		}

		err := BodyValidator(req)
		if err != nil {
			return presenter.BadRequest(c, err)
		}

		
		//userClaims, ok := c.Locals(UserClaimKey).(*jwt.UserClaims)
		// if !ok {
		// 	return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		// }

		t := presenter.CreateTripReqToTrip(&req)
		creatorID := uint(1) // TO DO

		if err := tripService.CreateTrip(c.UserContext(), t, creatorID); err != nil {
			// if errors.Is(err, trip.CompanyNotExist) || errors.Is(err, trip.ErrPathNotExist) || errors.Is(err, trip.ErrDuplication) || errors.Is(err, trip.ErrWrongPrice) || errors.Is(err, trip.ErrWrongReleaseDate) {
			// 	return presenter.BadRequest(c, err)
			// }
			//err := errors.New("Error")
			// apply trace ID here .... TODO
			return presenter.InternalServerError(c, err)
		}

		res := presenter.TripToCreateTripReq(t)
		return presenter.Created(c, "Trip created successfully", res)
	}
}
