package services

import (
	"bankservice/config"
	creditCard "bankservice/internal/wallet/credit_card"
	"bankservice/internal/wallet/wallet"
	"bankservice/pkg/adapters/clients/grpc"
	"bankservice/pkg/adapters/consul"
	"bankservice/pkg/adapters/storage"
	"bankservice/pkg/ports"
	"bankservice/pkg/ports/clients/clients"
	"bankservice/pkg/valuecontext"
	"context"
	"log"

	"gorm.io/gorm"
)

type AppContainer struct {
	cfg             config.Config
	dbConn          *gorm.DB
	walletService   *WalletService
	serviceRegistry ports.IServiceRegistry
	authClient      clients.IAuthClient
}

func NewAppContainer(cfg config.Config) (*AppContainer, error) {
	app := &AppContainer{
		cfg: cfg,
	}

	app.mustInitDB()
	app.mustRegisterService(cfg.Server)
	app.setAuthClient(cfg.Server.ServiceRegistry.AuthServiceName)

	app.setWalletService()

	return app, nil
}

func (a *AppContainer) RawDBConnection() *gorm.DB {
	return a.dbConn
}

func (a *AppContainer) mustInitDB() {
	if a.dbConn != nil {
		return
	}

	db, err := storage.NewPostgresGormConnection(a.cfg.DB)
	if err != nil {
		log.Fatal(err)
	}

	a.dbConn = db

	err = storage.AddExtension(a.dbConn)
	if err != nil {
		log.Fatal("Create extension failed: ", err)
	}

	err = storage.Migrate(a.dbConn)
	if err != nil {
		log.Fatal("Migration failed: ", err)
	}

	err = storage.InitDBRecords(a.cfg.DB, a.dbConn)
	if err != nil {
		log.Fatal("cannot init db records", err)
	}
}

func (a *AppContainer) setWalletService() {
	if a.walletService != nil {
		return
	}
	a.walletService = NewWalletService(wallet.NewWalletOps(storage.NewWalletRepo(a.RawDBConnection())), creditCard.NewCreditCardOps(storage.NewCreditCardRepo(a.RawDBConnection())))
}

func (a *AppContainer) WalletService() *WalletService {
	return a.walletService
}

func (a *AppContainer) WalletServiceFromCtx(ctx context.Context) *WalletService {
	tx, ok := valuecontext.TryGetTxFromContext(ctx)
	if !ok {
		return a.walletService
	}

	gc, ok := tx.Tx().(*gorm.DB)
	if !ok {
		return a.walletService
	}

	return NewWalletService(
		wallet.NewWalletOps(storage.NewWalletRepo(gc)),
		creditCard.NewCreditCardOps(storage.NewCreditCardRepo(gc)),
	)
}

func (a *AppContainer) mustRegisterService(srvCfg config.Server) {
	registry := consul.NewConsul(srvCfg.ServiceRegistry.Address)
	err := registry.RegisterService(srvCfg.ServiceRegistry.ServiceName, srvCfg.ServiceHostAddress, srvCfg.ServiceHTTPPrefixPath, srvCfg.ServiceHTTPHealthPath, srvCfg.HTTPPort)
	if err != nil {
		log.Fatalf("Failed to register service with Consul: %v", err)
	}
	a.serviceRegistry = registry
}

func (a *AppContainer) AuthClient() clients.IAuthClient {
	return a.authClient
}

func (a *AppContainer) setAuthClient(authServiceName string) {
	if a.authClient != nil {
		return
	}
	a.authClient = grpc.NewGRPCAuthClient(a.serviceRegistry, authServiceName)
}
