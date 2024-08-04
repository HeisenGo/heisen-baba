package storage

import (
	"bankservice/config"
	"bankservice/pkg/adapters/storage/entities"
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

	err := migrator.AutoMigrate(&entities.Wallet{}, &entities.CreditCard{}, entities.WalletCreditCard{}, entities.TransferTransaction{}, entities.Commission{})
	if err != nil {
		return err
	}
	return nil
}

func InitDBRecords(dbConfig config.DB, db *gorm.DB) error {
	err := initCommissionDB(dbConfig, db)
	if err != nil {
		return err
	}
	err = initWalletDB(dbConfig, db)
	if err != nil {
		return err
	}
	return nil
}

func initCommissionDB(dbConfig config.DB, db *gorm.DB) error {
	commission := &entities.Commission{
		AppCommissionPercentage: dbConfig.AppCommission,
	}
	err := db.FirstOrCreate(&commission).Error
	if err != nil {
		return err
	}
	return nil
}

func initWalletDB(dbConfig config.DB, db *gorm.DB) error {
	var systemWallet *entities.Wallet
	err := db.Where("is_system_wallet = ?", true).First(&systemWallet).Error
	if err != nil {
		systemWallet = &entities.Wallet{
			IsSystemWallet: true,
			Balance:        0,
		}
		err := db.Create(&systemWallet).Error
		if err != nil {
			return err
		}
	}
	return nil
}
