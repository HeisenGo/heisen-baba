package handlers

import (
	"tripcompanyservice/api/http/handlers/presenter"
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
		creatorID := uint(3)
		if err := techService.CreateTechTeam(c.UserContext(), team, creatorID); err != nil {
			return presenter.InternalServerError(c, err)
			//TODO error handeling
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
		creatorID := uint(3)
		if err := techService.CreateTechTeamMember(c.UserContext(), m, creatorID); err != nil {
			return presenter.InternalServerError(c, err)
			//TODO error handeling
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
		requesterID := uint(1) // get from contex
		teams, total, err := techService.GetTechTeamsOfCompany(c.UserContext(), uint(companyID), requesterID, uint(page), uint(pageSize))
		if err != nil {
			
			return presenter.InternalServerError(c, err)
		}
		data := presenter.NewPagination(
			presenter.BatchTeamToTechTeamRe(teams),
			uint(page),
			uint(pageSize),
			total,
		)
		return presenter.OK(c, "Trips fetched successfully", data)
	}
}

