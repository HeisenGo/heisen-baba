package service

import (
	"agency/config"
	"agency/internal/agency"
	"agency/internal/invoice"
	"agency/internal/reservation"
	"agency/internal/tour"
	"agency/pkg/adapters/clients/grpc"
	"agency/pkg/adapters/consul"
	"agency/pkg/adapters/storage"
	"agency/pkg/ports"
	"agency/pkg/ports/clients/clients"
	"agency/pkg/valuecontext"
	"context"
	"log"

	"gorm.io/gorm"
)

type AppContainer struct {
	cfg                config.Config
	serviceRegistry    ports.IServiceRegistry
	authClient         clients.IAuthClient
	bankClient         clients.IBankClient
	dbConn             *gorm.DB
	agencyService      *AgencyService
	tourService        *TourService
	reservationService *ReservationService
	invoiceService     *InvoiceService
}

func NewAppContainer(cfg config.Config) (*AppContainer, error) {
	app := &AppContainer{
		cfg: cfg,
	}

	app.mustInitDB()
	app.mustRegisterService(cfg.Server)
	app.setAuthClient(cfg.Server.ServiceRegistry.AuthServiceName)
	app.setBankClient(cfg.Server.ServiceRegistry.BankServiceName)
	app.setAgencyService()
	app.setTourService()
	app.setReservationService()
	app.setInvoiceService()

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

func (a *AppContainer) AuthClient() clients.IAuthClient {
	return a.authClient
}
func (a *AppContainer) BankClient() clients.IBankClient {
	return a.bankClient
}

func (a *AppContainer) setAuthClient(authServiceName string) {
	if a.authClient != nil {
		return
	}
	a.authClient = grpc.NewGRPCAuthClient(a.serviceRegistry, authServiceName)
}
func (a *AppContainer) setBankClient(bankServiceName string) {
	if a.bankClient != nil {
		return
	}
	a.bankClient = grpc.NewGRPCBankClient(a.serviceRegistry, bankServiceName)
}

func (a *AppContainer) mustRegisterService(srvCfg config.Server) {
	registry := consul.NewConsul(srvCfg.ServiceRegistry.Address)
	err := registry.RegisterService(srvCfg.ServiceRegistry.ServiceName, srvCfg.ServiceHostAddress, srvCfg.ServiceHTTPPrefixPath, srvCfg.ServiceHTTPHealthPath, srvCfg.HTTPPort)
	if err != nil {
		log.Fatalf("Failed to register service with Consul: %v", err)
	}
	a.serviceRegistry = registry
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
		reservation.NewOps(storage.NewReservationRepo(gc)), agency.NewOps(storage.NewAgencyRepo(gc)),invoice.NewOps(storage.NewInvoiceRepo(gc)),a.BankClient())
}

func (a *AppContainer) setTourService() {
	if a.tourService != nil {
		return
	}

	a.tourService = NewTourService(
		tour.NewOps(storage.NewTourRepo(a.dbConn)),
		reservation.NewOps(storage.NewReservationRepo(a.dbConn)), agency.NewOps(storage.NewAgencyRepo(a.dbConn)),invoice.NewOps(storage.NewInvoiceRepo(a.dbConn)),a.BankClient())
}

func (a *AppContainer) ReservationService() *ReservationService {
	return a.reservationService
}

func (a *AppContainer) ReservationServiceFromCtx(ctx context.Context) *ReservationService {
	tx, ok := valuecontext.TryGetTxFromContext(ctx)
	if !ok {
		return a.reservationService
	}

	gc, ok := tx.Tx().(*gorm.DB)
	if !ok {
		return a.reservationService
	}

	return NewReservationService(
		a.BankClient(),
		reservation.NewOps(storage.NewReservationRepo(gc)),
		invoice.NewOps(storage.NewInvoiceRepo(gc)),
	)
}

func (a *AppContainer) setReservationService() {
	if a.reservationService != nil {
		return
	}
	a.reservationService = NewReservationService(
		a.BankClient(),
		reservation.NewOps(storage.NewReservationRepo(a.dbConn)),
		invoice.NewOps(storage.NewInvoiceRepo(a.dbConn)),
	)
}

func (a *AppContainer) InvoiceService() *InvoiceService {
	return a.invoiceService
}

func (a *AppContainer) InvoiceServiceFromCtx(ctx context.Context) *InvoiceService {
	tx, ok := valuecontext.TryGetTxFromContext(ctx)
	if !ok {
		return a.invoiceService
	}

	gc, ok := tx.Tx().(*gorm.DB)
	if !ok {
		return a.invoiceService
	}

	return NewInvoiceService(
		invoice.NewOps(storage.NewInvoiceRepo(gc)),
	)
}

func (a *AppContainer) setInvoiceService() {
	if a.invoiceService != nil {
		return
	}
	a.invoiceService = NewInvoiceService(
		invoice.NewOps(storage.NewInvoiceRepo(a.dbConn)),
	)
}
