package handlers

import (
	"errors"
	"tripcompanyservice/api/http/handlers/presenter"
	"tripcompanyservice/internal"
	"tripcompanyservice/internal/company"
	"tripcompanyservice/service"

	"github.com/gofiber/fiber/v2"
)

//
func CreateTransportCompany(companyService *service.TransportCompanyService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var req presenter.CompanyReq

		if err := c.BodyParser(&req); err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}

		//userClaims, ok := c.Locals(UserClaimKey).(*jwt.UserClaims)
		// if !ok {
		// 	return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		// }
		// get owner and owner_ID existance check!!!!!! TODO:
		err := BodyValidator(req)
		if err != nil {
			return presenter.BadRequest(c, err)
		}
		tCompany := presenter.CompanyReqToCompanyDomain(&req)

		if err := companyService.CreateTransportCompany(c.UserContext(), tCompany); err != nil {
			if errors.Is(err, company.ErrDuplication) || errors.Is(err, internal.ErrConsecutiveSpaces) || errors.Is(err, internal.ErrExceedsMaxLength) || errors.Is(err, internal.ErrInvalidCharacters) {
				return presenter.BadRequest(c, err)
			}

			return presenter.InternalServerError(c, err)
		}
		res := presenter.CompanyToCompanyRes(*tCompany)
		return presenter.Created(c, "Company created successfully", res)
	}
}


func GetUserCompanies(companyService *service.TransportCompanyService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// userClaims, ok := c.Locals(UserClaimKey).(*jwt.UserClaims)
		// if !ok {
		// 	return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		// }
		//query parameter
		ownerID, err :=  //uuid.Parse(c.Params("ownerID"))
		// if err != nil {
		// 	return presenter.BadRequest(c, errors.New("given owner_id format in path is not correct"))
		// }
		page, pageSize := PageAndPageSize(c)

		boards, total, err := companyService.GetUserTransportCompanies(c.UserContext(), ownerID, uint(page), uint(pageSize))
		if err != nil {
			status := fiber.StatusInternalServerError
			if errors.Is(err, user.ErrUserNotFound) {
				status = fiber.StatusBadRequest
			}
			return SendError(c, err, status)
		}
		data := presenter.NewPagination(
			presenter.BatchBoardsToUserBoard(boards),
			uint(page),
			uint(pageSize),
			total,
		)
		return presenter.OK(c, "boards successfully fetched.", data)
	}
}