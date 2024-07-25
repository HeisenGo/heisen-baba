package service

import (
	"context"
	"log"
	"terminalpathservice/config"
	"terminalpathservice/internal/path"
	"terminalpathservice/internal/terminal"
	"terminalpathservice/pkg/adapters/consul"
	"terminalpathservice/pkg/adapters/storage"
	"terminalpathservice/pkg/ports"
	"terminalpathservice/pkg/valuecontext"

	"gorm.io/gorm"
)

type AppContainer struct {
	cfg             config.Config
	dbConn          *gorm.DB
	pathService     *PathService
	terminalService *TerminalService
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
	app.mustRegisterService(cfg.Server)

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
	)
}

func (a *AppContainer) setPathService() {
	if a.pathService != nil {
		return
	}
	a.pathService = NewPathService(path.NewOps(storage.NewPathRepo(a.dbConn)), terminal.NewOps(storage.NewTerminalRepo(a.dbConn)))
}

func (a *AppContainer) mustRegisterService(srvCfg config.Server) {
	registry := consul.NewConsul(srvCfg.ServiceRegistry.Address)
	err := registry.RegisterService(srvCfg.ServiceHostName, srvCfg.ServiceHTTPPrefixPath, srvCfg.ServiceHTTPHealthPath, srvCfg.HttpPort)
	if err != nil {
		log.Fatalf("Failed to register service with Consul: %v", err)
	}
}
