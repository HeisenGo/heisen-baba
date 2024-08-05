package handlers

import (
	"errors"
	"tripcompanyservice/api/http/handlers/presenter"
	"tripcompanyservice/internal"
	"tripcompanyservice/internal/company"
	"tripcompanyservice/service"

	"github.com/gofiber/fiber/v2"
)

func CreateTransportCompany(companyService *service.TransportCompanyService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var req presenter.CompanyReq

		if err := c.BodyParser(&req); err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}
		// owner id and the user claimer id ==
		// owner_ID existance check!!!!!! TODO:
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
		// TO DO no need I get it from claims!!!!
		ownerID, err := c.ParamsInt("ownerID")
		if err != nil {
			return SendError(c, errWrongIDType, fiber.StatusBadRequest)
		}

		//ownerID, err :=  //uuid.Parse(c.Params("ownerID"))
		// if err != nil {
		// 	return presenter.BadRequest(c, errors.New("given owner_id format in path is not correct"))
		// }
		page, pageSize := PageAndPageSize(c)

		companies, total, err := companyService.GetUserTransportCompanies(c.UserContext(), uint(ownerID), uint(page), uint(pageSize))
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
		companyID, err := c.ParamsInt("companyID")
		if err != nil {
			return presenter.BadRequest(c, err)
		}

		//TO DO: check whether it has unfinihed trips if so do not delete that
		// tO DO add requesterID only owner can delete company
		err = companyService.DeleteCompany(c.UserContext(), uint(companyID))
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
		} //userClaims, ok := c.Locals(UserClaimKey).(*jwt.UserClaims)
		// if !ok {
		// 	return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		// }
		// get owner and owner_ID existance check!!!!!! TODO:

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

		// userClaims, ok := c.Locals(UserClaimKey).(*jwt.UserClaims)
		// if !ok {
		// 	return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		// }
		companyID, err := c.ParamsInt("companyID")
		if err != nil {
			return presenter.BadRequest(c, errWrongIDType)
		}

		if companyID < 0 {
			return presenter.BadRequest(c, errWrongIDType)
		}
		userID := uint(1) //ownerID
		//requesterID TODO:
		updatedCompany := presenter.UpdateCompanyToCompany(&req, uint(companyID))
		changedCompany, err := companyService.PatchCompanyByOwner(c.UserContext(), updatedCompany, userID, req.NewOwnerEmail)

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
