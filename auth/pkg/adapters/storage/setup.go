package storage

import (
	"authservice/config"
	"authservice/internal/user"
	"authservice/pkg/adapters/storage"
	"authservice/pkg/adapters/storage/entities"
	"authservice/service"
	"fmt"

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

	err := migrator.AutoMigrate(&entities.User{})
	if err != nil {
		return err
	}
	initAdminDB(db)
	return nil
}

func initAdminDB( db *gorm.DB) error {
	newUser := &entities.User{
		Email: "admin@gmail.com",
		Password: "$2a$10$D/gJTXTNteoOggH/wRQBUeB7.fSsYx8xMm4yC1Q97HoIcI8jAK4U.",
		IsAdmin: true,
		IsSuperAdmin: true,
	}
	err := db.FirstOrCreate(&newUser).Error
	if err != nil {
		return err
	}
	return nil
}

