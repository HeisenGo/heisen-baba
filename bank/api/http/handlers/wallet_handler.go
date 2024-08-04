package handlers

import (
	presenter "bankservice/api/http/handlers/presentor"
	"bankservice/internal/user"
	creaditCard "bankservice/internal/wallet/credit_card"
	"bankservice/internal/wallet/wallet"
	"bankservice/pkg/valuecontext"
	"bankservice/service"
	"errors"
	"github.com/gofiber/fiber/v2"
)

func AddCardToWallet(serviceFactory ServiceFactory[*services.WalletService]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		walletService := serviceFactory(c.UserContext())

		var req presenter.AddCardToWalletReq

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

		newCard := presenter.AddCardToWalletReqToCard(&req)

		createdCard, err := walletService.AddCardToWalletByUserID(c.UserContext(), newCard, userReq.ID)

		if err != nil {
			if errors.Is(err, creaditCard.ErrInvalidCardNumber) || errors.Is(err, creaditCard.ErrCardAlreadyExists) {
				return presenter.BadRequest(c, err)
			}
			err := errors.New("Error")
			// apply trace ID here .... TODO
			return presenter.InternalServerError(c, err)
		}

		res := presenter.CardToAddCardToWalletResp(*createdCard)
		return presenter.Created(c, "Card successfully added to wallet.", res)
	}
}

func WalletCards(walletService *services.WalletService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userReq, ok := c.Locals(valuecontext.UserClaimKey).(*user.User)
		if !ok {
			return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		}
		userWalletCards, err := walletService.GetUserWalletCards(c.UserContext(), userReq.ID)

		if err != nil {
			//if errors.Is(err, wallet.ErrInvalidCardNumber) || errors.Is(err, wallet.ErrCardAlreadyExists) {
			//	return presenter.BadRequest(c, err)
			//}
			err := errors.New("Error")
			// apply trace ID here .... TODO
			return presenter.InternalServerError(c, err)
		}

		res := presenter.CardsToWalletCardsResp(userWalletCards)
		return presenter.OK(c, "Cards successfully fetched", res)
	}
}

func Deposit(serviceFactory ServiceFactory[*services.WalletService]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		walletService := serviceFactory(c.UserContext())

		var req presenter.DepositReq

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

		card := presenter.DepositReqNumToCard(req.CardNumber)
		userWallet, err := walletService.Deposit(c.UserContext(), card, req.Amount, userReq.ID)

		if err != nil {
			//if errors.Is(err, wallet.ErrInvalidCardNumber) || errors.Is(err, wallet.ErrCardAlreadyExists) {
			//	return presenter.BadRequest(c, err)
			//}
			err := errors.New("Error")
			// apply trace ID here .... TODO
			return presenter.InternalServerError(c, err)
		}

		res := presenter.WalletToDepositResp(*userWallet)
		return presenter.OK(c, "deposit successfully done", res)
	}
}

func Withdraw(serviceFactory ServiceFactory[*services.WalletService]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		walletService := serviceFactory(c.UserContext())

		var req presenter.WithdrawReq

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

		card := presenter.WithdrawReqNumToCard(req.CardNumber)
		userWallet, err := walletService.Withdraw(c.UserContext(), card, req.Amount, userReq.ID)

		if err != nil {
			if errors.Is(err, wallet.ErrNotEnoughBalance) {
				return presenter.BadRequest(c, err)
			}
			err := errors.New("Error")
			// apply trace ID here .... TODO
			return presenter.InternalServerError(c, err)
		}

		res := presenter.WalletToWithdrawResp(*userWallet)
		return presenter.OK(c, "withdraw successfully done", res)
	}
}

func GetWallet(walletService *services.WalletService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userReq, ok := c.Locals(valuecontext.UserClaimKey).(*user.User)
		if !ok {
			return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		}
		userWallet, err := walletService.GetWallet(c.UserContext(), userReq.ID)

		if err != nil {
			//if errors.Is(err, wallet.ErrInvalidCardNumber) || errors.Is(err, wallet.ErrCardAlreadyExists) {
			//	return presenter.BadRequest(c, err)
			//}
			err := errors.New("Error")
			// apply trace ID here .... TODO
			return presenter.InternalServerError(c, err)
		}

		res := presenter.WalletToGetWalletResp(*userWallet)
		return presenter.OK(c, "wallet successfully fetched", res)
	}
}
