package handlers

import (
	"errors"
	"tripcompanyservice/api/http/handlers/presenter"
	"tripcompanyservice/internal/company"
	"tripcompanyservice/internal/techteam"
	"tripcompanyservice/internal/user"
	"tripcompanyservice/pkg/valuecontext"
	"tripcompanyservice/service"

	"github.com/gofiber/fiber/v2"
)

func CreateTechTeam(techService *service.TechTeamService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var req presenter.TechTeamRe

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
		team := presenter.TechTeamReqToTechTeam(&req)
		// TODO: from context
		// requester
		if err := techService.CreateTechTeam(c.UserContext(), team, userReq.ID); err != nil {
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

		err := BodyValidator(req)
		if err != nil {
			return presenter.BadRequest(c, err)
		}
		m := presenter.TechTeamMemberReToTechTeamMember(&req)
		userReq, ok := c.Locals(valuecontext.UserClaimKey).(*user.User)
		if !ok {
			return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		}

		if err := techService.CreateTechTeamMember(c.UserContext(), m, userReq.ID); err != nil {
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

		page, pageSize := PageAndPageSize(c)
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

		teams, total, err := techService.GetTechTeamsOfCompany(c.UserContext(), uint(companyID), userReq.ID, uint(page), uint(pageSize))
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

		teamID, err := c.ParamsInt("teamID")
		if err != nil {
			return presenter.BadRequest(c, errWrongIDType)
		}

		// tO DO add requesterID only owner can delete company
		userReq, ok := c.Locals(valuecontext.UserClaimKey).(*user.User)
		if !ok {
			return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		}

		err = teamService.DeleteTeam(c.UserContext(), uint(teamID), userReq.ID)
		if err != nil {
			if errors.Is(err, service.ErrForbidden) {
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
