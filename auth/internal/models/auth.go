package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email        string `gorm:"uniqueIndex;not null"`
	Password     string `gorm:"not null"`
	IsSuperAdmin bool   `gorm:"default:false"`
	WalletID     uint
}

func (u User) validate() {

}
