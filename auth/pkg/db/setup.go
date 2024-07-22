package db

import (
	"auth/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func newPostgresGormConnection(dsn string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func addExtension(db *gorm.DB) error {
	return nil
}

func migrate(db *gorm.DB) error {
	migrator := db.Migrator()

	err := migrator.AutoMigrate(models.User{})
	if err != nil {
		return err
	}
	return nil
}

func MustInitDB(dsn string) *gorm.DB {
	db, err := newPostgresGormConnection(dsn)
	if err != nil {
		log.Fatal(err)
	}

	err = addExtension(db)
	if err != nil {
		log.Fatal("Create extension failed: ", err)
	}

	err = migrate(db)
	if err != nil {
		log.Fatal("Migration failed: ", err)
	}
	return db
}
