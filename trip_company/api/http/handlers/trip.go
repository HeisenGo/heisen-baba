package handlers

import (
	"errors"
	"time"
	"tripcompanyservice/api/http/handlers/presenter"
	"tripcompanyservice/internal/company"
	"tripcompanyservice/internal/trip"
	"tripcompanyservice/internal/user"
	"tripcompanyservice/pkg/valuecontext"
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

		userReq, ok := c.Locals(valuecontext.UserClaimKey).(*user.User)
		if !ok {
			return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		}

		t := presenter.CreateTripReqToTrip(&req)

		if err := tripService.CreateTrip(c.UserContext(), t, userReq.ID); err != nil {
			if errors.Is(err, service.ErrPathNotFound) || errors.Is(err, company.ErrCompanyNotFound) || errors.Is(err, trip.ErrInvalidPercentage) || errors.Is(err, trip.ErrStartTime) || errors.Is(err, trip.ErrSecondPenalty) || errors.Is(err, trip.ErrInvalidPercentage) || errors.Is(err, trip.ErrDuplication) || errors.Is(err, trip.ErrFirstPenalty) || errors.Is(err, trip.ErrNegativePrice) {
				return presenter.BadRequest(c, err)
			}
			if errors.Is(err, service.ErrForbidden) {
				return presenter.Forbidden(c, err)
			}
			return presenter.InternalServerError(c, err)
		}
		res := presenter.TripToOwnerAdminTechTeamOperatorTripResponse(*t)
		return presenter.Created(c, "Trip created successfully", res)
	}
}

func GetCompanyTrips(tripService *service.TripService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		
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
		companyID, err := c.ParamsInt("companyID")
		if err != nil {
			return presenter.BadRequest(c, err)
		}
		userReq, ok := c.Locals(valuecontext.UserClaimKey).(*user.User)
		if !ok {
			return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		}
		var requester string
		userID := userReq.ID
		if userReq.IsAdmin{
			requester = "admin"
		}
		// check isOwner / isAdmin / isTechinician
		trips, total, role, err := tripService.GetCompanyTrips(c.UserContext(), originCity, destinationCity, pathType, &startDate, requester, uint(companyID),userID, uint(page), uint(pageSize))
		if err != nil {
			if errors.Is(err, trip.ErrRecordsNotFound) {
				return presenter.BadRequest(c, err)
			}
			err := errors.New("Error")
			return presenter.InternalServerError(c, err)
		}

		var data interface{}
		if role == string(company.OwnerRole) || role == string(company.OperatorRole) || role == string(company.TechRole) || role == "admin" {
			data = presenter.NewPagination(
				presenter.BatchTripToOwnerAdminTechTeamOperatorTripResponse(trips),
				uint(page),
				uint(pageSize),
				total,
			)
		} else if role == "agency" {
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

func GetCompanyAgencyTrips(tripService *service.TripService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		
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
		companyID, err := c.ParamsInt("companyID")
		if err != nil {
			return presenter.BadRequest(c, err)
		}

		userReq, ok := c.Locals(valuecontext.UserClaimKey).(*user.User)
		if !ok {
			return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		}
		var requester string
		requester = "agency" //no other type

		trips, total, _, err := tripService.GetCompanyTrips(c.UserContext(), originCity, destinationCity, pathType, &startDate, requester, uint(companyID),userReq.ID, uint(page), uint(pageSize))
		if err != nil {
			if errors.Is(err, trip.ErrRecordsNotFound) {
				return presenter.BadRequest(c, err)
			}
			return presenter.InternalServerError(c, err)
		}
		var data interface{}
		if requester == string(company.OwnerRole) || requester == string(company.OperatorRole) || requester == string(company.TechRole) || requester == "admin" {
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
		
		//query parameter
		tripID, err := c.ParamsInt("tripID")
		if err != nil {
			return presenter.BadRequest(c, errWrongIDType)
		}

		if tripID < 0 {
			return presenter.BadRequest(c, errWrongIDType)
		}
		
		userReq, ok := c.Locals(valuecontext.UserClaimKey).(*user.User)
		if !ok {
			return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		}
		var requester string
		if userReq.IsAdmin {
			requester = "admin"
		}
		t, requesterRole, err := tripService.GetFullTripByID(c.UserContext(), uint(tripID), userReq.ID, requester)
		if err != nil {
			if errors.Is(err, trip.ErrRecordNotFound) {
				return presenter.BadRequest(c, err)
			}
			err := errors.New("Error")
			return presenter.InternalServerError(c, err)
		}
		var data interface{}
		if requesterRole == string(company.OwnerRole) || requesterRole == string(company.OperatorRole) || requesterRole == string(company.TechRole) || requesterRole == "admin" {
			data = presenter.TripToOwnerAdminTechTeamOperatorTripResponse(*t)
		} else if requesterRole == "agency" {
			data = presenter.TripToAgencyTripResponse(*t)
		} else {
			data = presenter.TripToUserTripResponse(*t)
		}
		return presenter.OK(c, "Trip fetched successfully", data)
	}
}

func GetFullAgencyTripByID(tripService *service.TripService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		
		//query parameter
		tripID, err := c.ParamsInt("tripID")
		if err != nil {
			return presenter.BadRequest(c, errWrongIDType)
		}

		if tripID < 0 {
			return presenter.BadRequest(c, errWrongIDType)
		}
		//TO do from local context
		
		userReq, ok := c.Locals(valuecontext.UserClaimKey).(*user.User)
		if !ok {
			return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		}
		requester := "agency" // "user" // "admin" // "operator" // "employee" // "techteam"
		t, _, err := tripService.GetFullTripByID(c.UserContext(), uint(tripID), userReq.ID, requester)
		if err != nil {
			if errors.Is(err, trip.ErrRecordNotFound) {
				return presenter.BadRequest(c, err)
			}
			return presenter.InternalServerError(c, err)
		}
		var data interface{}
		if requester == string(company.OwnerRole) || requester == string(company.OperatorRole) || requester == string(company.TechRole) || requester == "admin" {
			data = presenter.TripToOwnerAdminTechTeamOperatorTripResponse(*t)
		} else if requester == "agency" {
			data = presenter.TripToAgencyTripResponse(*t)
		} else {
			data = presenter.TripToUserTripResponse(*t)
		}
		return presenter.OK(c, "Trip fetched successfully", data)
	}
}

func GetTrips(tripService *service.TripService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		
		userReq, ok := c.Locals(valuecontext.UserClaimKey).(*user.User)
		if !ok {
			return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		}
		var requester string
		//query parameter
		page, pageSize := PageAndPageSize(c)
		startDateStr := c.Query("start_date")
		originCity := c.Query("from")
		destinationCity := c.Query("to")
		pathType := c.Query("type")

		// Parse dates
		var startDate time.Time
		var err error

		if startDateStr != "" {
			startDate, err = time.Parse("2006-01-02", startDateStr)
			if err != nil {
				return presenter.BadRequest(c, errors.New("invalid start date format"))
			}
		}
		if userReq.IsAdmin{
			requester = "admin"
		}

		trips, total, err := tripService.GetTrips(c.UserContext(), originCity, destinationCity, pathType, &startDate, requester, uint(page), uint(pageSize))
		if err != nil {
			if errors.Is(err, trip.ErrRecordsNotFound) {
				return presenter.BadRequest(c, err)
			}
			err := errors.New("Error")
			return presenter.InternalServerError(c, err)
		}

		var data interface{}
		if requester == "admin" {
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

func GetAgencyTrips(tripService *service.TripService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		
		//query parameter
		page, pageSize := PageAndPageSize(c)
		startDateStr := c.Query("start_date")
		originCity := c.Query("from")
		destinationCity := c.Query("to")
		pathType := c.Query("type")


		_, ok := c.Locals(valuecontext.UserClaimKey).(*user.User)
		if !ok {
			return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		}
		// Parse dates
		var startDate time.Time
		var err error

		if startDateStr != "" {
			startDate, err = time.Parse("2006-01-02", startDateStr)
			if err != nil {
				return presenter.BadRequest(c, errors.New("invalid start date format"))
			}
		}
		requester := "agency"
		trips, total, err := tripService.GetTrips(c.UserContext(), originCity, destinationCity, pathType, &startDate, requester, uint(page), uint(pageSize))
		if err != nil {
			if errors.Is(err, trip.ErrRecordsNotFound) {
				return presenter.BadRequest(c, err)
			}
			return presenter.InternalServerError(c, err)
		}

		data := presenter.NewPagination(
			presenter.BatchTripToAgencyTripResponse(trips),
			uint(page),
			uint(pageSize),
			total,
		)

		return presenter.OK(c, "Trips fetched successfully", data)
	}
}

func PatchTrip(serviceFactory ServiceFactory[*service.TripService]) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var req presenter.UpdateTripRequest
		tripService := serviceFactory(c.UserContext())
		if err := c.BodyParser(&req); err != nil {
			return presenter.BadRequest(c, err)
		}

		
		userReq, ok := c.Locals(valuecontext.UserClaimKey).(*user.User)
		if !ok {
			return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		}
		tripID, err := c.ParamsInt("tripID")
		if err != nil {
			return presenter.BadRequest(c, errWrongIDType)
		}

		if tripID < 0 {
			return presenter.BadRequest(c, errWrongIDType)
		}

		newTrip := presenter.UpdateTripReqToTrip(&req)
		// only operator and owner can do this
		changedTrip, err := tripService.UpdateTrip(c.UserContext(), uint(tripID), newTrip, userReq.ID)

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

		
		userReq, ok := c.Locals(valuecontext.UserClaimKey).(*user.User)
		if !ok {
			return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		}
		tripID, err := c.ParamsInt("tripID")
		if err != nil {
			return presenter.BadRequest(c, errWrongIDType)
		}

		if tripID < 0 {
			return presenter.BadRequest(c, errWrongIDType)
		}

		// only owner and operator
		changedTrip, err := tripService.SetTechTeamToTrip(c.UserContext(), uint(tripID), req.TechTeamID, userReq.ID)

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

		
		userReq, ok := c.Locals(valuecontext.UserClaimKey).(*user.User)
		if !ok {
			return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		}
		tripID, err := c.ParamsInt("tripID")
		if err != nil {
			return presenter.BadRequest(c, errWrongIDType)
		}

		if tripID < 0 {
			return presenter.BadRequest(c, errWrongIDType)
		}

		// USERID from context TODO:
		// only owner and operator
		changedTrip, err := tripService.CancelTrip(c.UserContext(), uint(tripID), userReq.ID, req.IsCanceled)

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

		
		userReq, ok := c.Locals(valuecontext.UserClaimKey).(*user.User)
		if !ok {
			return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		}
		tripID, err := c.ParamsInt("tripID")
		if err != nil {
			return presenter.BadRequest(c, errWrongIDType)
		}

		if tripID < 0 {
			return presenter.BadRequest(c, errWrongIDType)
		}
		changedTrip, err := tripService.ConfirmTrip(c.UserContext(), uint(tripID), userReq.ID, req.IsConfirmed)

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

		
		userReq, ok := c.Locals(valuecontext.UserClaimKey).(*user.User)
		if !ok {
			return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		}
		tripID, err := c.ParamsInt("tripID")
		if err != nil {
			return presenter.BadRequest(c, errWrongIDType)
		}

		if tripID < 0 {
			return presenter.BadRequest(c, errWrongIDType)
		}
		changedTrip, err := tripService.FinishTrip(c.UserContext(), uint(tripID), userReq.ID, req.IsFinished)

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
			return presenter.InternalServerError(c, err)
		}

		return presenter.OK(c, "Number of Trips fetched successfully", fiber.Map{"count": total})
	}
}
