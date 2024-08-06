package service

import (
	"context"
	"hotel/config"
	"hotel/internal/hotel"
	"hotel/internal/invoice"
	"hotel/internal/reservation"
	"hotel/internal/room"
	"hotel/pkg/adapters/clients/grpc"
	"hotel/pkg/adapters/consul"
	"hotel/pkg/adapters/storage"
	"hotel/pkg/ports"
	"hotel/pkg/ports/clients/clients"
	"hotel/pkg/valuecontext"
	"log"

	"gorm.io/gorm"
)

type AppContainer struct {
	cfg                config.Config
	serviceRegistry    ports.IServiceRegistry
	authClient         clients.IAuthClient
	dbConn             *gorm.DB
	hotelService       *HotelService
	roomService        *RoomService
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
	app.setHotelService()
	app.setRoomService()
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

func (a *AppContainer) setAuthClient(authServiceName string) {
	if a.authClient != nil {
		return
	}
	a.authClient = grpc.NewGRPCAuthClient(a.serviceRegistry, authServiceName)
}

func (a *AppContainer) HotelService() *HotelService {
	return a.hotelService
}

func (a *AppContainer) HotelServiceFromCtx(ctx context.Context) *HotelService {
	tx, ok := valuecontext.TryGetTxFromContext(ctx)
	if !ok {
		return a.hotelService
	}

	gc, ok := tx.Tx().(*gorm.DB)
	if !ok {
		return a.hotelService
	}

	return NewHotelService(
		hotel.NewOps(storage.NewHotelRepo(gc)),
		room.NewOps(storage.NewRoomRepo(gc)),
	)
}

func (a *AppContainer) setHotelService() {
	if a.hotelService != nil {
		return
	}

	a.hotelService = NewHotelService(
		hotel.NewOps(storage.NewHotelRepo(a.dbConn)),
		room.NewOps(storage.NewRoomRepo(a.dbConn)),
	)
}

func (a *AppContainer) RoomService() *RoomService {
	return a.roomService
}

func (a *AppContainer) RoomServiceFromCtx(ctx context.Context) *RoomService {
	tx, ok := valuecontext.TryGetTxFromContext(ctx)
	if !ok {
		return a.roomService
	}

	gc, ok := tx.Tx().(*gorm.DB)
	if !ok {
		return a.roomService
	}

	return NewRoomService(
		room.NewOps(storage.NewRoomRepo(gc)),
		reservation.NewOps(storage.NewReservationRepo(gc)), hotel.NewOps(storage.NewHotelRepo(gc)))
}

func (a *AppContainer) setRoomService() {
	if a.roomService != nil {
		return
	}
	a.roomService = NewRoomService(
		room.NewOps(storage.NewRoomRepo(a.dbConn)),
		reservation.NewOps(storage.NewReservationRepo(a.dbConn)), hotel.NewOps(storage.NewHotelRepo(a.dbConn)))
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
		reservation.NewOps(storage.NewReservationRepo(gc)),
		invoice.NewOps(storage.NewInvoiceRepo(gc)),
	)
}

func (a *AppContainer) setReservationService() {
	if a.reservationService != nil {
		return
	}
	a.reservationService = NewReservationService(
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

func (a *AppContainer) mustRegisterService(srvCfg config.Server) {
	registry := consul.NewConsul(srvCfg.ServiceRegistry.Address)
	err := registry.RegisterService(srvCfg.ServiceRegistry.ServiceName, srvCfg.ServiceHostAddress, srvCfg.ServiceHTTPPrefixPath, srvCfg.ServiceHTTPHealthPath, srvCfg.HTTPPort)
	if err != nil {
		log.Fatalf("Failed to register service with Consul: %v", err)
	}
	a.serviceRegistry = registry
}
