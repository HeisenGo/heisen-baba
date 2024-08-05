package service

import (
	"authservice/config"
	"authservice/internal/user"
	"authservice/pkg/adapters/consul"
	"authservice/pkg/adapters/rabbitmq"
	"authservice/pkg/adapters/storage"
	"authservice/pkg/ports"
	"log"

	"gorm.io/gorm"
)

type AppContainer struct {
	cfg           config.Config
	messageBroker ports.IMessageBroker
	dbConn        *gorm.DB
	authService   *AuthService
}

func NewAppContainer(cfg config.Config) (*AppContainer, error) {
	app := &AppContainer{
		cfg: cfg,
	}

	app.setMessageBroker(cfg.MessageBroker)
	app.mustInitDB()

	// service registry
	app.mustRegisterService(cfg.Server)

	app.setAuthService()

	return app, nil
}

func (a *AppContainer) mustRegisterService(srvCfg config.Server) {
	registry := consul.NewConsul(srvCfg.ServiceRegistry.Address)

	err := registry.RegisterService(srvCfg.ServiceRegistry.ServiceName, srvCfg.ServiceHostAddress, srvCfg.ServiceHTTPPrefixPath, srvCfg.ServiceHTTPHealthPath, srvCfg.GRPCPort, srvCfg.HTTPPort)
	if err != nil {
		log.Fatalf("Failed to register service with Consul: %v", err)
	}
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
func (a *AppContainer) MessageBroker() ports.IMessageBroker {
	return a.messageBroker
}

func (a *AppContainer) setMessageBroker(messageBrokerCfg config.MessageBroker) {
	if a.messageBroker != nil {
		return
	}

	a.messageBroker = rabbitmq.NewRabbitMQ(messageBrokerCfg.Username, messageBrokerCfg.Password, messageBrokerCfg.Host, messageBrokerCfg.Port)
}

func (a *AppContainer) AuthService() *AuthService {
	return a.authService
}

func (a *AppContainer) setAuthService() {
	if a.authService != nil {
		return
	}

	a.authService = NewAuthService(user.NewOps(storage.NewUserRepo(a.dbConn)), a.MessageBroker(), []byte(a.cfg.Server.TokenSecret),
		a.cfg.Server.TokenExpMinutes,
		a.cfg.Server.RefreshTokenExpMinutes)
}
