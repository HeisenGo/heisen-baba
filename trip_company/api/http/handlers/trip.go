package handlers

import (
	"errors"
	"time"
	"tripcompanyservice/api/http/handlers/presenter"
	"tripcompanyservice/internal/company"
	"tripcompanyservice/internal/trip"
	"tripcompanyservice/service"

	"github.com/gofiber/fiber/v2"
)

func CreateTrip(serviceFactory ServiceFactory[*service.TripService]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tripService := serviceFactory(c.UserContext())
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
			if errors.Is(err, company.ErrCompanyNotFound) || errors.Is(err, trip.ErrInvalidPercentage) || errors.Is(err, trip.ErrStartTime) || errors.Is(err, trip.ErrSecondPenalty) || errors.Is(err, trip.ErrInvalidPercentage) || errors.Is(err, trip.ErrDuplication) || errors.Is(err, trip.ErrFirstPenalty) || errors.Is(err, trip.ErrNegativePrice) {
				return presenter.BadRequest(c, err)
			}
			//err := errors.New("Error")
			// apply trace ID here .... TODO
			return presenter.InternalServerError(c, err)
		}
		res := presenter.TripToOwnerAdminTechTeamOperatorTripResponse(*t)
		return presenter.Created(c, "Trip created successfully", res)
	}
}

func GetCompanyTrips(tripService *service.TripService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// userClaims, ok := c.Locals(UserClaimKey).(*jwt.UserClaims)
		// if !ok {
		// 	return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		// }
		//query parameter
		page, pageSize := PageAndPageSize(c)
		startDateStr := c.Query("start_date")
		originCity := c.Query("from")
		destinationCity := c.Query("to")
		pathType := c.Query("type")

		var startDate time.Time
		var err error

		if startDateStr != "" {
			startDate, err = time.Parse("2006-01-02", startDateStr)
			if err != nil {
				return presenter.BadRequest(c, errors.New("invalid start date format"))
			}
		}
		// TO DO check requester!!!
		companyID, err := c.ParamsInt("companyID")
		if err != nil {
			return presenter.BadRequest(c, err)
		}
		//requester // how to show the result
		trips, total, err := tripService.GetCompanyTrips(c.UserContext(), originCity, destinationCity, pathType, &startDate, uint(companyID), uint(page), uint(pageSize))
		if err != nil {
			if errors.Is(err, trip.ErrRecordsNotFound) {
				return presenter.BadRequest(c, err)
			}
			err := errors.New("Error")
			return presenter.InternalServerError(c, err)
		}
		requester := "owner" //user // admin // agency // unknown
		var data interface{}
		if requester == "owner" || requester == "operator" || requester == "technician" || requester == "admin" {
			data = presenter.NewPagination(
				presenter.BatchTripToOwnerAdminTechTeamOperatorTripResponse(trips),
				uint(page),
				uint(pageSize),
				total,
			)
		} else if requester == "agency" {
			data = presenter.NewPagination(
				presenter.BatchTripToAgencyTripResponse(trips),
				uint(page),
				uint(pageSize),
				total,
			)
		} else {
			data = presenter.NewPagination(
				presenter.BatchTripToUserTripResponse(trips),
				uint(page),
				uint(pageSize),
				total,
			)
		}

		return presenter.OK(c, "Trips fetched successfully", data)
	}
}

func GetFullTripByID(tripService *service.TripService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// userClaims, ok := c.Locals(UserClaimKey).(*jwt.UserClaims)
		// if !ok {
		// 	return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		// }
		//query parameter
		tripID, err := c.ParamsInt("tripID")
		if err != nil {
			return presenter.BadRequest(c, errWrongIDType)
		}

		if tripID < 0 {
			return presenter.BadRequest(c, errWrongIDType)
		}
		// requester!!!!!!!!!!! check what to show
		requester := "owner" // "user" // "admin" // "operator" // "employee" // "techteam"
		t, err := tripService.GetFullTripByID(c.UserContext(), uint(tripID))
		if err != nil {
			if errors.Is(err, trip.ErrRecordNotFound) {
				return presenter.BadRequest(c, err)
			}
			err := errors.New("Error")
			return presenter.InternalServerError(c, err)
		}
		var data interface{}
		if requester == "owner" || requester == "operator" || requester == "technician" || requester == "admin" {
			data = presenter.TripToOwnerAdminTechTeamOperatorTripResponse(*t)
		} else if requester == "agency" {
			data = presenter.TripToAgencyTripResponse(*t)
		} else {
			data = presenter.TripToUserTripResponse(*t)
		}
		// TO DO implement else
		return presenter.OK(c, "Trip fetched successfully", data)
	}
}

func GetTrips(tripService *service.TripService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// userClaims, ok := c.Locals(UserClaimKey).(*jwt.UserClaims)
		// if !ok {
		// 	return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		// }
		//query parameter
		page, pageSize := PageAndPageSize(c)
		startDateStr := c.Query("start_date")
		originCity := c.Query("from")
		destinationCity := c.Query("to")
		pathType := c.Query("type")
		requesterType := c.Query("requester_type") // get from auth!!!!! TODO:

		// Parse dates
		var startDate time.Time
		var err error

		if startDateStr != "" {
			startDate, err = time.Parse("2006-01-02", startDateStr)
			if err != nil {
				return presenter.BadRequest(c, errors.New("invalid start date format"))
			}
		}

		trips, total, err := tripService.GetTrips(c.UserContext(), originCity, destinationCity, pathType, &startDate, requesterType, uint(page), uint(pageSize))
		if err != nil {
			if errors.Is(err, trip.ErrRecordsNotFound) {
				return presenter.BadRequest(c, err)
			}
			err := errors.New("Error")
			return presenter.InternalServerError(c, err)
		}
		requester := "owner" //user // admin // agency // unknown
		var data interface{}
		if requester == "owner" || requester == "operator" || requester == "technician" || requester == "admin" {
			data = presenter.NewPagination(
				presenter.BatchTripToOwnerAdminTechTeamOperatorTripResponse(trips),
				uint(page),
				uint(pageSize),
				total,
			)
		} else if requester == "agency" {
			data = presenter.NewPagination(
				presenter.BatchTripToAgencyTripResponse(trips),
				uint(page),
				uint(pageSize),
				total,
			)
		} else {
			data = presenter.NewPagination(
				presenter.BatchTripToUserTripResponse(trips),
				uint(page),
				uint(pageSize),
				total,
			)
		}

		return presenter.OK(c, "Trips fetched successfully", data)
	}
}

func PatchTrip(tripService *service.TripService) fiber.Handler { // tansactional!!!! TO DO:
	return func(c *fiber.Ctx) error {

		var req presenter.UpdateTripRequest

		if err := c.BodyParser(&req); err != nil {
			return presenter.BadRequest(c, err)
		}

		// userClaims, ok := c.Locals(UserClaimKey).(*jwt.UserClaims)
		// if !ok {
		// 	return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		// }
		tripID, err := c.ParamsInt("tripID")
		if err != nil {
			return presenter.BadRequest(c, errWrongIDType)
		}

		if tripID < 0 {
			return presenter.BadRequest(c, errWrongIDType)
		}

		newTrip := presenter.UpdateTripReqToTrip(&req)
		// only operator and owner can do this
		changedTrip, err := tripService.UpdateTrip(c.UserContext(), uint(tripID), newTrip)

		if err != nil {
			if errors.Is(err, trip.ErrCanNotUpdate) || errors.Is(err, trip.ErrNotUpdated) || errors.Is(err, trip.ErrRecordNotFound) {
				return presenter.BadRequest(c, err)
			}
			// trace ID : TODO
			return presenter.InternalServerError(c, err)
		}
		res := presenter.TripToOwnerAdminTechTeamOperatorTripResponse(*changedTrip)
		return presenter.OK(c, "Trip updated successfully", res)
	}
}

func SetTechTeamToTrip(serviceFactory ServiceFactory[*service.TripService]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tripService := serviceFactory(c.UserContext())
		var req presenter.SetTechTeamRequest

		if err := c.BodyParser(&req); err != nil {
			return presenter.BadRequest(c, err)
		}

		// userClaims, ok := c.Locals(UserClaimKey).(*jwt.UserClaims)
		// if !ok {
		// 	return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		// }
		tripID, err := c.ParamsInt("tripID")
		if err != nil {
			return presenter.BadRequest(c, errWrongIDType)
		}

		if tripID < 0 {
			return presenter.BadRequest(c, errWrongIDType)
		}

		// only owner and operator
		changedTrip, err := tripService.SetTechTeamToTrip(c.UserContext(), uint(tripID), req.TechTeamID)

		if err != nil {
			if errors.Is(err, trip.ErrCanNotUpdate) || errors.Is(err, trip.ErrNotUpdated) || errors.Is(err, trip.ErrRecordNotFound) {
				return presenter.BadRequest(c, err)
			}
			return presenter.InternalServerError(c, err)
		}
		res := presenter.TripToOwnerAdminTechTeamOperatorTripResponse(*changedTrip)
		return presenter.OK(c, "Team set successfully", res)
	}
}

func CancelTrip(serviceFactory ServiceFactory[*service.TripService]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tripService := serviceFactory(c.UserContext())
		var req presenter.CancelTripReq

		if err := c.BodyParser(&req); err != nil {
			return presenter.BadRequest(c, err)
		}

		// userClaims, ok := c.Locals(UserClaimKey).(*jwt.UserClaims)
		// if !ok {
		// 	return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		// }
		tripID, err := c.ParamsInt("tripID")
		if err != nil {
			return presenter.BadRequest(c, errWrongIDType)
		}

		if tripID < 0 {
			return presenter.BadRequest(c, errWrongIDType)
		}

		// USERID from context TODO:
		// only owner and operator
		requesterID := uint(2)
		changedTrip, err := tripService.CancelTrip(c.UserContext(), uint(tripID), requesterID, req.IsCanceled)

		if err != nil {
			if errors.Is(err, service.ErrForbidden) {
				return presenter.Unauthorized(c, err)
			}
			if errors.Is(err, service.ErrFinishedTrip) || errors.Is(err, service.ErrSameState) || errors.Is(err, trip.ErrCanNotUpdate) || errors.Is(err, trip.ErrNotUpdated) || errors.Is(err, trip.ErrRecordNotFound) {
				return presenter.BadRequest(c, err)
			}
			return presenter.InternalServerError(c, err)
		}
		res := presenter.TripToOwnerAdminTechTeamOperatorTripResponse(*changedTrip)
		return presenter.OK(c, "Trip canceled successfully", res)
	}
}

func ConfirmTrip(serviceFactory ServiceFactory[*service.TripService]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tripService := serviceFactory(c.UserContext())
		var req presenter.ConfirmTripReq

		if err := c.BodyParser(&req); err != nil {
			return presenter.BadRequest(c, err)
		}

		// userClaims, ok := c.Locals(UserClaimKey).(*jwt.UserClaims)
		// if !ok {
		// 	return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		// }
		tripID, err := c.ParamsInt("tripID")
		if err != nil {
			return presenter.BadRequest(c, errWrongIDType)
		}

		if tripID < 0 {
			return presenter.BadRequest(c, errWrongIDType)
		}
		// check role of requester
		requesterID := uint(5)
		changedTrip, err := tripService.ConfirmTrip(c.UserContext(), uint(tripID), requesterID, req.IsConfirmed)

		if err != nil {
			if errors.Is(err, service.ErrForbidden) {
				return presenter.Unauthorized(c, err)
			}
			if errors.Is(err, service.ErrFinishedTrip) || errors.Is(err, service.ErrSameState) || errors.Is(err, trip.ErrCanNotUpdate) || errors.Is(err, trip.ErrNotUpdated) || errors.Is(err, trip.ErrRecordNotFound) {
				return presenter.BadRequest(c, err)
			}

			return presenter.InternalServerError(c, err)
		}
		res := presenter.TripToOwnerAdminTechTeamOperatorTripResponse(*changedTrip)
		return presenter.OK(c, "Trip updated successfully", res)
	}
}

func FinishTrip(serviceFactory ServiceFactory[*service.TripService]) fiber.Handler { // tansactional!!!! TO DO:
	return func(c *fiber.Ctx) error {
		tripService := serviceFactory(c.UserContext())
		var req presenter.FinishTripReq

		if err := c.BodyParser(&req); err != nil {
			return presenter.BadRequest(c, err)
		}

		// userClaims, ok := c.Locals(UserClaimKey).(*jwt.UserClaims)
		// if !ok {
		// 	return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		// }
		tripID, err := c.ParamsInt("tripID")
		if err != nil {
			return presenter.BadRequest(c, errWrongIDType)
		}

		if tripID < 0 {
			return presenter.BadRequest(c, errWrongIDType)
		}
		// TO DO
		requesterID := uint(1)
		changedTrip, err := tripService.FinishTrip(c.UserContext(), uint(tripID), requesterID, req.IsFinished)

		if err != nil {

			if errors.Is(err, service.ErrForbidden) {
				return presenter.Unauthorized(c, err)
			}
			if errors.Is(err, service.ErrUnConfirmed) || errors.Is(err, service.ErrFutureTrip) || errors.Is(err, service.ErrFinishedTrip) || errors.Is(err, service.ErrSameState) || errors.Is(err, trip.ErrCanNotUpdate) || errors.Is(err, trip.ErrNotUpdated) || errors.Is(err, trip.ErrRecordNotFound) {
				return presenter.BadRequest(c, err)
			}

			return presenter.InternalServerError(c, err)
		}
		res := presenter.TripToOwnerAdminTechTeamOperatorTripResponse(*changedTrip)
		return presenter.OK(c, "Trip updated successfully", res)
	}
}

// REST FOR TERMINAL PATH SERVICE
func GetCountPathUnfinishedTrips(tripService *service.TripService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pathID, err := c.ParamsInt("pathID")
		if err != nil {
			return presenter.BadRequest(c, errWrongIDType)
		}

		if pathID < 0 {
			return presenter.BadRequest(c, errWrongIDType)
		}

		total, err := tripService.GetCountPathUnfinishedTrips(c.UserContext(), uint(pathID))
		if err != nil {
			if errors.Is(err, trip.ErrRecordsNotFound) {
				return presenter.BadRequest(c, err)
			}
			err := errors.New("Error")
			return presenter.InternalServerError(c, err)
		}

		return presenter.OK(c, "Number of Trips fetched successfully", fiber.Map{"count": total})
	}
}
