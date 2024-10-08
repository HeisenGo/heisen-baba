package handlers

import (
	"errors"
	"tripcompanyservice/api/http/handlers/presenter"
	"tripcompanyservice/internal"
	"tripcompanyservice/internal/company"
	"tripcompanyservice/internal/user"
	"tripcompanyservice/pkg/valuecontext"
	"tripcompanyservice/service"

	"github.com/gofiber/fiber/v2"
)

func CreateTransportCompany(companyService *service.TransportCompanyService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var req presenter.CompanyReq

		if err := c.BodyParser(&req); err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}
		userReq, ok := c.Locals(valuecontext.UserClaimKey).(*user.User)
		if !ok {
			return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		}
		err := BodyValidator(req)
		if err != nil {
			return presenter.BadRequest(c, err)
		}
		tCompany := presenter.CompanyReqToCompanyDomain(&req)
		tCompany.OwnerID = userReq.ID
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

		userReq, ok := c.Locals(valuecontext.UserClaimKey).(*user.User)
		if !ok {
			return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		}

		page, pageSize := PageAndPageSize(c)

		companies, total, err := companyService.GetUserTransportCompanies(c.UserContext(), userReq.ID, uint(page), uint(pageSize))
		if err != nil {
			status := fiber.StatusInternalServerError
			return SendError(c, err, status)
		}
		data := presenter.NewPagination(
			presenter.BatchCompaniesToCompanies(companies),
			uint(page),
			uint(pageSize),
			total,
		)
		return presenter.OK(c, "companies successfully fetched.", data)
	}
}

func GetCompanies(companyService *service.TransportCompanyService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		userReq, ok := c.Locals(valuecontext.UserClaimKey).(*user.User)
		if !ok {
			return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		}

		if !userReq.IsAdmin {
			return presenter.Unauthorized(c, errWrongClaimType)
		}
		page, pageSize := PageAndPageSize(c)

		companies, total, err := companyService.GetTransportCompanies(c.UserContext(), uint(page), uint(pageSize))
		if err != nil {
			status := fiber.StatusInternalServerError
			return SendError(c, err, status)
		}
		data := presenter.NewPagination(
			presenter.BatchCompaniesToCompanies(companies),
			uint(page),
			uint(pageSize),
			total,
		)
		return presenter.OK(c, "companies successfully fetched.", data)
	}
}

func DeleteCompany(companyService *service.TransportCompanyService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		companyID, err := c.ParamsInt("companyID")
		if err != nil {
			return presenter.BadRequest(c, err)
		}

		userReq, ok := c.Locals(valuecontext.UserClaimKey).(*user.User)
		if !ok {
			return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		}
		err = companyService.DeleteCompany(c.UserContext(), uint(companyID), userReq.ID)
		if err != nil {
			if errors.Is(err, company.ErrCompanyNotFound) {
				return presenter.BadRequest(c, err)
			}
			return presenter.InternalServerError(c, err)
		}
		return presenter.NoContent(c)
	}
}

func BlockCompany(companyService *service.TransportCompanyService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		//To DO: only admin
		var req presenter.BlockCompany

		if err := c.BodyParser(&req); err != nil {
			return presenter.BadRequest(c, err)
		}
		companyID, err := c.ParamsInt("companyID")
		if err != nil {
			return presenter.BadRequest(c, err)
		}

		userReq, ok := c.Locals(valuecontext.UserClaimKey).(*user.User)
		if !ok {
			return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		}

		if !userReq.IsAdmin {
			return presenter.Unauthorized(c, errWrongClaimType)
		}
		tCompany, err := companyService.BlockCompany(c.UserContext(), uint(companyID), req.IsBlocked)
		if err != nil {
			if errors.Is(err, company.ErrCompanyNotFound) {
				return presenter.BadRequest(c, err)
			}
			return presenter.InternalServerError(c, err)
		}
		res := presenter.CompanyToCompanyRes(*tCompany)
		return presenter.OK(c, "Company Updated successfully", res)
	}
}

func PatchCompany(companyService *service.TransportCompanyService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var req presenter.UpdateCompanyReq

		if err := c.BodyParser(&req); err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}

		companyID, err := c.ParamsInt("companyID")
		if err != nil {
			return presenter.BadRequest(c, errWrongIDType)
		}

		if companyID < 0 {
			return presenter.BadRequest(c, errWrongIDType)
		}
		userReq, ok := c.Locals(valuecontext.UserClaimKey).(*user.User)
		if !ok {
			return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		}

		updatedCompany := presenter.UpdateCompanyToCompany(&req, uint(companyID))
		changedCompany, err := companyService.PatchCompanyByOwner(c.UserContext(), updatedCompany, userReq.ID, req.NewOwnerEmail)

		if err != nil {
			if errors.Is(err, service.ErrForbidden) {
				return presenter.Forbidden(c, err)
			}
			if errors.Is(err, company.ErrCompanyNotFound) || errors.Is(err, company.ErrFailedToUpdate) {
				return presenter.BadRequest(c, err)
			}
			// trace ID : TODO
			return presenter.InternalServerError(c, err)
		}
		res := presenter.CompanyToCompanyRes(*changedCompany)
		return presenter.OK(c, "Company updated successfully", res)
	}
}
