package handlers

import (
	"errors"
	"tripcompanyservice/api/http/handlers/presenter"
	"tripcompanyservice/internal/company"
	"tripcompanyservice/internal/techteam"
	"tripcompanyservice/service"

	"github.com/gofiber/fiber/v2"
)

func CreateTechTeam(techService *service.TechTeamService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var req presenter.TechTeamRe

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
		team := presenter.TechTeamReqToTechTeam(&req)
		// TODO: from context
		// requester
		creatorID := uint(1)
		if err := techService.CreateTechTeam(c.UserContext(), team, creatorID); err != nil {
			if errors.Is(err, service.ErrForbidden) {
				return presenter.Forbidden(c, err)
			}
			if errors.Is(err, company.ErrCompanyNotFound) || errors.Is(err, techteam.ErrDuplication) {
				return presenter.BadRequest(c, err)
			}
			return presenter.InternalServerError(c, err)
		}
		res := presenter.TechTeamToTechTeamRe(*team)
		return presenter.Created(c, "team created successfully", res)
	}
}

func CreateTechMember(techService *service.TechTeamService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var req presenter.TechTeamMemberRe

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
		m := presenter.TechTeamMemberReToTechTeamMember(&req)
		// TODO: from context
		// requester
		creatorID := uint(1)
		if err := techService.CreateTechTeamMember(c.UserContext(), m, creatorID); err != nil {
			if errors.Is(err, service.ErrForbidden) {
				return presenter.Forbidden(c, err)
			}
			if errors.Is(err, company.ErrCompanyNotFound) || errors.Is(err, techteam.ErrTeamNotFound) {
				return presenter.BadRequest(c, err)
			}
			return presenter.InternalServerError(c, err)
		}
		res := presenter.TechMemberToTechTeamMemberRe(*m)
		return presenter.Created(c, "Member created successfully", res)
	}
}

func GetTechTeamsOfCompany(techService *service.TechTeamService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// userClaims, ok := c.Locals(UserClaimKey).(*jwt.UserClaims)
		// if !ok {
		// 	return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		// }
		//query parameter
		page, pageSize := PageAndPageSize(c)
		companyID, err := c.ParamsInt("companyID")
		if err != nil {
			return presenter.BadRequest(c, errWrongIDType)
		}

		if companyID < 0 {
			return presenter.BadRequest(c, errWrongIDType)
		}
		// TODO: requester // admin/owner/operator can do this
		requesterID := uint(1) // get from contex
		teams, total, err := techService.GetTechTeamsOfCompany(c.UserContext(), uint(companyID), requesterID, uint(page), uint(pageSize))
		if err != nil {
			if errors.Is(err, service.ErrForbidden) {
				return presenter.Forbidden(c, err)
			}
			if errors.Is(err, company.ErrCompanyNotFound) {
				presenter.BadRequest(c, err)
			}
			return presenter.InternalServerError(c, err)
		}
		data := presenter.NewPagination(
			presenter.BatchTeamToTechTeamRe(teams),
			uint(page),
			uint(pageSize),
			total,
		)
		return presenter.OK(c, "Teams fetched successfully", data)
	}
}

func DeleteTeam(teamService *service.TechTeamService) fiber.Handler {
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
		teamID, err := c.ParamsInt("teamID")
		if err != nil {
			return presenter.BadRequest(c, errWrongIDType)
		}

		//TO DO: check whether it has unfinihed trips if so do not delete that
		// tO DO add requesterID only owner can delete company
		requesterID := uint(1)
		err = teamService.DeleteTeam(c.UserContext(), uint(teamID), requesterID)
		if err != nil {
			if errors.Is(err, service.ErrForbidden){
				return presenter.Forbidden(c, err)
			}
			if errors.Is(err, techteam.ErrDeleteTeam) || errors.Is(err, techteam.ErrTeamNotFound) {
				return presenter.BadRequest(c, err)
			}
			return presenter.InternalServerError(c, err)
		}
		return presenter.NoContent(c)
	}
}
