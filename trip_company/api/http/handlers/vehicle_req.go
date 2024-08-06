package handlers

import (
	"errors"
	"tripcompanyservice/api/http/handlers/presenter"
	"tripcompanyservice/internal/company"
	"tripcompanyservice/internal/user"
	vehiclerequest "tripcompanyservice/internal/vehicle_request"
	"tripcompanyservice/pkg/valuecontext"
	"tripcompanyservice/service"

	"github.com/gofiber/fiber/v2"
)

// TransactionaL ? TODO
func CreateVehicleRequest(serviceFactory ServiceFactory[*service.VehicleReService]) fiber.Handler { //serviceFactory ServiceFactory[*service.TripService])fiber.Handler {
	return func(c *fiber.Ctx) error {
		vehicleREqService := serviceFactory(c.UserContext())

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

		userReq, ok := c.Locals(valuecontext.UserClaimKey).(*user.User)
		if !ok {
			return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		}
		if err := vehicleREqService.CreateVehicleReq(c.UserContext(), v, userReq.ID); err != nil {
			if errors.Is(err, service.ErrForbidden) {
				return presenter.Forbidden(c, err)
			}
			if errors.Is(err, company.ErrCompanyNotFound) || errors.Is(err, vehiclerequest.ErrNotFound) {
				return presenter.BadRequest(c, err)
			}

			return presenter.InternalServerError(c, err)
		}
		res := presenter.VehicleToCreateVehicleRes(*v)
		return presenter.Created(c, "Vehicle Request created successfully", res)
	}
}

func DeleteVR(serviceFactory ServiceFactory[*service.VehicleReService]) fiber.Handler {
	return func(c *fiber.Ctx) error {

		userReq, ok := c.Locals(valuecontext.UserClaimKey).(*user.User)
		if !ok {
			return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		}
		vehicleREqService := serviceFactory(c.UserContext())

		vrID, err := c.ParamsInt("vRID")
		if err != nil {
			return presenter.BadRequest(c, errWrongIDType)
		}

		// only owner can delete company
		err = vehicleREqService.Delete(c.UserContext(), uint(vrID), userReq.ID)
		if err != nil {
			if errors.Is(err, service.ErrForbidden) {
				return presenter.Forbidden(c, err)
			}
			if errors.Is(err, company.ErrCompanyNotFound) || errors.Is(err, vehiclerequest.ErrNotFound) {
				return presenter.BadRequest(c, err)
			}
			return presenter.InternalServerError(c, err)
		}
		return presenter.NoContent(c)
	}
}
