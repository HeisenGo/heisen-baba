package service

import (
	"context"
	"log"
	"agency/config"
	"agency/internal/agency"
	"agency/internal/tour"
	"agency/pkg/adapters/storage"
	"agency/pkg/valuecontext"

	"gorm.io/gorm"
)

type AppContainer struct {
	cfg             config.Config
	dbConn          *gorm.DB
	agencyService   *AgencyService
	tourService     *TourService
}

func NewAppContainer(cfg config.Config) (*AppContainer, error) {
	app := &AppContainer{
		cfg: cfg,
	}

	app.mustInitDB()
	app.setAgencyService()
	app.setTourService()

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

func (a *AppContainer) AgencyService() *AgencyService {
	return a.agencyService
}

func (a *AppContainer) AgencyServiceFromCtx(ctx context.Context) *AgencyService {
	tx, ok := valuecontext.TryGetTxFromContext(ctx)
	if !ok {
		return a.agencyService
	}

	gc, ok := tx.Tx().(*gorm.DB)
	if !ok {
		return a.agencyService
	}

	return NewAgencyService(
		agency.NewOps(storage.NewAgencyRepo(gc)),
	)
}

func (a *AppContainer) setAgencyService() {
	if a.agencyService != nil {
		return
	}

	a.agencyService = NewAgencyService(
		agency.NewOps(storage.NewAgencyRepo(a.dbConn)),
	)
}

func (a *AppContainer) TourService() *TourService {
	return a.tourService
}

func (a *AppContainer) TourServiceFromCtx(ctx context.Context) *TourService {
	tx, ok := valuecontext.TryGetTxFromContext(ctx)
	if !ok {
		return a.tourService
	}

	gc, ok := tx.Tx().(*gorm.DB)
	if !ok {
		return a.tourService
	}

	return NewTourService(
		tour.NewOps(storage.NewTourRepo(gc)),
	)
}

func (a *AppContainer) setTourService() {
	if a.tourService != nil {
		return
	}

	a.tourService = NewTourService(
		tour.NewOps(storage.NewTourRepo(a.dbConn)),
	)
}
