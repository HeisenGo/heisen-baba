package service

import (
	"hotel/config"
	"hotel/internal/hotel"
	"hotel/internal/room"
	"hotel/pkg/adapters/storage"
	"log"

	"gorm.io/gorm"
)

type AppContainer struct {
	cfg          config.Config
	dbConn       *gorm.DB
	hotelService *HotelService
	roomService  *RoomService
}

func NewAppContainer(cfg config.Config) (*AppContainer, error) {
	app := &AppContainer{
		cfg: cfg,
	}

	app.mustInitDB()
	app.setHotelService()
	app.setRoomService()

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

func (a *AppContainer) HotelService() *HotelService {
	return a.hotelService
}

func (a *AppContainer) setHotelService() {
	if a.hotelService != nil {
		return
	}

	a.hotelService = NewHotelService(hotel.NewOps(storage.NewHotelRepo(a.dbConn)),room.NewOps(storage.NewRoomRepo(a.dbConn)) )
}

func (a *AppContainer) RoomService() *RoomService {
	return a.roomService
}

func (a *AppContainer) setRoomService() {
	if a.roomService != nil {
		return
	}
	a.roomService = NewRoomService(room.NewOps(storage.NewRoomRepo(a.dbConn)))
}