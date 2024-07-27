package service

import (
	"context"
	"log"
	"tripcompanyservice/config"
	"tripcompanyservice/internal/company"
	"tripcompanyservice/pkg/adapters/consul"
	"tripcompanyservice/pkg/adapters/storage"
	"tripcompanyservice/pkg/ports"
	"tripcompanyservice/pkg/valuecontext"

	"gorm.io/gorm"
)

type AppContainer struct {
	cfg            config.Config
	dbConn         *gorm.DB
	companyService *TransportCompanyService
	//terminalService *TerminalService
	serviceRegistry *ports.IServiceRegistry
}

func NewAppContainer(cfg config.Config) (*AppContainer, error) {
	app := &AppContainer{
		cfg: cfg,
	}

	app.mustInitDB()
	err := storage.Migrate(app.dbConn)
	if err != nil {
		log.Fatal("Migration failed: ", err)
	}
	// service registry
	//app.mustRegisterService(cfg.Server)

	app.setCompanyService()
	//app.setPathService()
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

}

func (a *AppContainer) CompanyService() *TransportCompanyService {
	return a.companyService
}

func (a *AppContainer) CompanyServiceFromCtx(ctx context.Context) *TransportCompanyService {
	tx, ok := valuecontext.TryGetTxFromContext(ctx)
	if !ok {
		return a.companyService
	}

	gc, ok := tx.Tx().(*gorm.DB)
	if !ok {
		return a.companyService
	}

	return NewTransportCompanyService(
		company.NewOps(storage.NewTransportCompanyRepo(gc)),
	)
}

func (a *AppContainer) setCompanyService() {
	if a.companyService != nil {
		return
	}
	a.companyService = NewTransportCompanyService(company.NewOps(storage.NewTransportCompanyRepo(a.dbConn)))
}

func (a *AppContainer) mustRegisterService(srvCfg config.Server) {
	registry := consul.NewConsul(srvCfg.ServiceRegistry.Address)
	err := registry.RegisterService(srvCfg.ServiceHostName, srvCfg.ServiceHTTPPrefixPath, srvCfg.ServiceHTTPHealthPath, srvCfg.HttpPort)
	if err != nil {
		log.Fatalf("Failed to register service with Consul: %v", err)
	}
}
