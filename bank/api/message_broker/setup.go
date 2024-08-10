package message_broker

import (
	"bankservice/api/message_broker/handler"
	services "bankservice/service"
)

func Run(app *services.AppContainer) {
	messageBroker := app.MessageBroker()
	walletHandler := handler.NewWalletHandler(app.WalletService())
	messageBroker.Consume("users", walletHandler.CreateWallet)
}
