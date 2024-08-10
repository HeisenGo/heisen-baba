package service

import (
	"context"
	"log"
	"terminalpathservice/config"
	"terminalpathservice/internal/path"
	"terminalpathservice/internal/terminal"
	"terminalpathservice/pkg/adapters/clients/grpc"
	"terminalpathservice/pkg/adapters/clients/rest"
	"terminalpathservice/pkg/adapters/consul"
	"terminalpathservice/pkg/adapters/storage"
	"terminalpathservice/pkg/ports"
	"terminalpathservice/pkg/ports/clients/clients"
	"terminalpathservice/pkg/valuecontext"

	"gorm.io/gorm"
)

type AppContainer struct {
	cfg               config.Config
	dbConn            *gorm.DB
	pathService       *PathService
	terminalService   *TerminalService
	serviceRegistry   ports.IServiceRegistry
	authClient        clients.IAuthClient
	tripCompanyClient clients.ITripCompanyClient
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

	// Initialize SQLite connection for city/country validation
	//err = db_helper.InitDB(cfg.SQLite.Path)
	//if err != nil {
	//	log.Fatal("Failed to initialize SQLite database: ", err)
	//}

	// service registry
	app.mustRegisterService(cfg.Server)
	app.setAuthClient(cfg.Server.ServiceRegistry.AuthServiceName)
	app.setTripCompanyClient(cfg.Server.ServiceRegistry.TripCompanyServiceName)

	app.setTerminalService()
	app.setPathService()
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
	err = storage.Migrate(a.dbConn)
	if err != nil {
		log.Fatal("Migration failed: ", err)
	}
}

func (a *AppContainer) TerminalService() *TerminalService {
	return a.terminalService
}

func (a *AppContainer) AuthClient() clients.IAuthClient {
	return a.authClient
}

func (a *AppContainer) TripCompanyClient() clients.ITripCompanyClient {
	return a.tripCompanyClient
}

func (a *AppContainer) TerminalServiceFromCtx(ctx context.Context) *TerminalService {
	tx, ok := valuecontext.TryGetTxFromContext(ctx)
	if !ok {
		return a.terminalService
	}

	gc, ok := tx.Tx().(*gorm.DB)
	if !ok {
		return a.terminalService
	}

	return NewTerminalService(
		terminal.NewOps(storage.NewTerminalRepo(gc)),
		path.NewOps(storage.NewPathRepo(gc)),
	)
}

func (a *AppContainer) setTerminalService() {
	if a.terminalService != nil {
		return
	}
	a.terminalService = NewTerminalService(terminal.NewOps(storage.NewTerminalRepo(a.dbConn)), path.NewOps(storage.NewPathRepo(a.dbConn)))
}

func (a *AppContainer) setAuthClient(authServiceName string) {
	if a.authClient != nil {
		return
	}
	a.authClient = grpc.NewGRPCAuthClient(a.serviceRegistry, authServiceName)
}

func (a *AppContainer) setTripCompanyClient(tripCompanyServiceName string) {
	if a.tripCompanyClient != nil {
		return
	}
	a.tripCompanyClient = rest.NewRestTripCompanyClient(a.serviceRegistry, tripCompanyServiceName)
}

func (a *AppContainer) PathService() *PathService {
	return a.pathService
}

func (a *AppContainer) PathServiceFromCtx(ctx context.Context) *PathService {
	tx, ok := valuecontext.TryGetTxFromContext(ctx)
	if !ok {
		return a.pathService
	}

	gc, ok := tx.Tx().(*gorm.DB)
	if !ok {
		return a.pathService
	}

	return NewPathService(
		path.NewOps(storage.NewPathRepo(gc)),
		terminal.NewOps(storage.NewTerminalRepo(gc)),
		a.TripCompanyClient(),
	)
}

func (a *AppContainer) setPathService() {
	if a.pathService != nil {
		return
	}
	a.pathService = NewPathService(path.NewOps(storage.NewPathRepo(a.dbConn)), terminal.NewOps(storage.NewTerminalRepo(a.dbConn)), a.TripCompanyClient())
}

func (a *AppContainer) mustRegisterService(srvCfg config.Server) {
	registry := consul.NewConsul(srvCfg.ServiceRegistry.Address)
	err := registry.RegisterService(srvCfg.ServiceRegistry.ServiceName, srvCfg.ServiceHostAddress, srvCfg.ServiceHTTPPrefixPath, srvCfg.ServiceHTTPHealthPath, srvCfg.GRPCPort, srvCfg.HTTPPort)
	if err != nil {
		log.Fatalf("Failed to register service with Consul: %v", err)
	}
	a.serviceRegistry = registry
}
