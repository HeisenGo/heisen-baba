package entities

import (
	"time"

	"gorm.io/gorm"
)

type Invoice struct {
	gorm.Model
	TicketID       uint      `gorm:"not null"`
	Ticket         *Ticket    `gorm:"foreignKey:TicketID"`
	IssuedDate     time.Time `gorm:"not null"`
	Info           string    `gorm:"type:text"` // Detailed information for the invoice
	PerAmountPrice float64
	TotalPrice     float64
	Status         string `gorm:"type:varchar(20);default:'pending'"` // paid //canceled
	Penalty       float64 `gorm:"default:0"`
}
