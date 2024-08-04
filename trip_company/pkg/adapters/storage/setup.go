package storage

import (
	"fmt"
	"tripcompanyservice/config"
	"tripcompanyservice/pkg/adapters/storage/entities"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresGormConnection(dbConfig config.DB) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		dbConfig.Host, dbConfig.User, dbConfig.Pass, dbConfig.DBName, dbConfig.Port)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func AddExtension(db *gorm.DB) error {
	return db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error
}

func Migrate(db *gorm.DB) error {
	migrator := db.Migrator()

	err := migrator.AutoMigrate(&entities.TransportCompany{}, &entities.Trip{}, &entities.TripCancellingPenalty{}, 
	&entities.Ticket{}, &entities.Invoice{}, &entities.VehicleRequest{})
	if err != nil {
		return err
	}
	return nil
}
