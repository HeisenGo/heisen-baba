package handlers

import (
	"errors"
	"tripcompanyservice/api/http/handlers/presenter"
	"tripcompanyservice/internal/company"
	vehiclerequest "tripcompanyservice/internal/vehicle_request"
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
		creatorID := uint(1) // TO DO

		if err := vehicleREqService.CreateVehicleReq(c.UserContext(), v, creatorID); err != nil {
			if errors.Is(err,service.ErrForbidden){
				return presenter.Forbidden(c, err)
			}
			if errors.Is(err, company.ErrCompanyNotFound) || errors.Is(err, vehiclerequest.ErrNotFound){
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
		// userClaims, ok := c.Locals(UserClaimKey).(*jwt.UserClaims)
		// if !ok {
		// 	return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		// }
		//query parameter
		// TODo: check it is admin!
		//ownerID, err :=  //uuid.Parse(c.Params("ownerID"))
		// if err != nil {
		// 	return presenter.BadRequest(c, errors.New("given owner_id format in path is not correct"))
		// }
		vehicleREqService := serviceFactory(c.UserContext())

		vrID, err := c.ParamsInt("vRID")
		if err != nil {
			return presenter.BadRequest(c, errWrongIDType)
		}

		//TO DO: check whether it has unfinihed trips if so do not delete that
		// tO DO add requesterID only owner can delete company
		creatorID := uint(1)
		err = vehicleREqService.Delete(c.UserContext(), uint(vrID), creatorID)
		if err != nil {
			if errors.Is(err,service.ErrForbidden){
				return presenter.Forbidden(c, err)
			}
			if errors.Is(err, company.ErrCompanyNotFound) || errors.Is(err, vehiclerequest.ErrNotFound){
				return presenter.BadRequest(c, err)
			}
			return presenter.InternalServerError(c, err)
		}
		return presenter.NoContent(c)
	}
}
