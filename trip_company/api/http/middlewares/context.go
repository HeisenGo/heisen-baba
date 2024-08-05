package middlerwares

import (
	"tripcompanyservice/pkg/valuecontext"

	"github.com/gofiber/fiber/v2"
)

func SetTransaction(commiter valuecontext.Committer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		cm := commiter.Begin()
		valuecontext.SetTx(c.UserContext(), cm)

		//logger := valuecontext.GetLogger(c.UserContext())

		//logger.Info("starting transaction")
		if err := c.Next(); err != nil {
			//logger.Info("rollback on error", "error", err.Error())
			cm.Rollback()
			return err
		}

		err, ok := c.Locals(valuecontext.IsTxError).(error)
		if ok && err != nil {
			//logger.Info("rollback on not ok response", "error", err.Error())
			cm.Rollback()
			return nil
		}

		if err := cm.Commit(); err != nil {
			//logger.Info("commit error", "err", err.Error())
			cm.Rollback()
			return err
		}
		//logger.Info("ending transaction successfully")
		return nil
	}
}
