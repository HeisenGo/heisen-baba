package service

import (
	"log"
	"terminalpathservice/config"
	"terminalpathservice/pkg/adapters/storage"

	"gorm.io/gorm"
)

type AppContainer struct {
	cfg          config.Config
	dbConn       *gorm.DB
	//userService  *UserService
	//authService  *AuthService
	//orderService *OrderService
}

func NewAppContainer(cfg config.Config) (*AppContainer, error) {
	app := &AppContainer{
		cfg: cfg,
	}

	app.mustInitDB()
	storage.Migrate(app.dbConn)

	//app.setUserService()
	//app.setAuthService()
	//app.setOrderService()

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