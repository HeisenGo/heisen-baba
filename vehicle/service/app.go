package service

import (
	"vehicle/config"
	"vehicle/internal/vehicle"
	"vehicle/pkg/adapters/storage"
	"log"
	"gorm.io/gorm"
)

type AppContainer struct {
	cfg          config.Config
	dbConn       *gorm.DB
	vehicleService *VehicleService
}

func NewAppContainer(cfg config.Config) (*AppContainer, error) {
	app := &AppContainer{
		cfg: cfg,
	}

	app.mustInitDB()
	app.setVehicleService()

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

func (a *AppContainer) VehicleService() *VehicleService {
	return a.vehicleService
}

func (a *AppContainer) setVehicleService() {
	if a.vehicleService != nil {
		return
	}

	a.vehicleService = NewVehicleService(vehicle.NewOps(storage.NewVehicleRepo(a.dbConn)))
}
