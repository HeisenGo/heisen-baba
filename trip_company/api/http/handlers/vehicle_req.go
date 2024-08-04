package handlers

import (
	"tripcompanyservice/api/http/handlers/presenter"
	"tripcompanyservice/service"

	"github.com/gofiber/fiber/v2"
)

// TransactionaL ? TODO
func CreateVehicleRequest(vehicleREqService *service.VehicleReService) fiber.Handler { //serviceFactory ServiceFactory[*service.TripService])fiber.Handler {
	return func(c *fiber.Ctx) error {
		//tripService := serviceFactory(c.UserContext())

		// TODO: auth owner can do this and operator
		var req presenter.CreateVehicleRequest

		if err := c.BodyParser(&req); err != nil {
			return presenter.BadRequest(c, err)
		}

		err := BodyValidator(req)
		if err != nil {
			return presenter.BadRequest(c, err)
		}

		v := presenter.CreateVehicleRequestToVehicleRequest(&req)
		creatorID := uint(1) // TO DO

		if err := vehicleREqService.CreateVehicleReq(c.UserContext(), v, creatorID); err != nil {
			// if errors.Is(err, trip.CompanyNotExist) || errors.Is(err, trip.ErrPathNotExist) || errors.Is(err, trip.ErrDuplication) || errors.Is(err, trip.ErrWrongPrice) || errors.Is(err, trip.ErrWrongReleaseDate) {
			// 	return presenter.BadRequest(c, err)
			// }
			//err := errors.New("Error")
			// apply trace ID here .... TODO
			return presenter.InternalServerError(c, err)
		}
		res := presenter.VehicleToCreateVehicleRes(*v)
		return presenter.Created(c, "Vehicle Request created successfully", res)
	}
}
