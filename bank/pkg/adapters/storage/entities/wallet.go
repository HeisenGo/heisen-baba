package entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Wallet struct {
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	IsSystemWallet bool      `gorm:"default:false"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
	UserID         *uuid.UUID     `gorm:"uniqueIndex"`
	Balance        uint
	CreditCards    []*CreditCard `gorm:"many2many:wallet_credit_cards;"`
}

type CreditCard struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Number    string         `gorm:"uniqueIndex;not null"`
	Wallets   []*Wallet      `gorm:"many2many:wallet_credit_cards;"`
}
type WalletTransaction struct {
	gorm.Model
	Amount             uint
	Type               string
	Status             string
	WalletCreditCardID uint
	WalletCreditCard   *WalletCreditCard `gorm:"foreignKey:WalletCreditCardID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
type WalletCreditCard struct {
	ID           uuid.UUID   `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	WalletID     uint        `gorm:"index:idx_together_wallet_card,unique"`
	Wallet       *Wallet     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreditCardID uint        `gorm:"index:idx_together_wallet_card,unique"`
	CreditCard   *CreditCard `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type TransferTransactionStatus string

// Define constants for the Status type
const (
	TransactionSuccess TransferTransactionStatus = "success"
	Failed             TransferTransactionStatus = "failed"
)

type TransferTransaction struct {
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
	Amount         uint
	Status         TransferTransactionStatus
	FromWalletID   *uuid.UUID
	FromWallet     *Wallet `gorm:"foreignKey:FromWalletID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ToWalletID     *uuid.UUID
	ToWallet       *Wallet `gorm:"foreignKey:ToWalletID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	IsPaidToSystem bool    `gorm:"default:false"`
}

type Commission struct {
	gorm.Model
	AppCommissionPercentage uint
}
