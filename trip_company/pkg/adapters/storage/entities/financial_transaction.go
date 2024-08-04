package entities

import "gorm.io/gorm"

type FinancialTransaction struct {
	gorm.Model
	TicketID         uint
	Ticket           Ticket `gorm:"foreignKey:TicketID"`
	CompanyShare     float64
	TourShare        float64
	AppProviderShare float64
	Status           string `gorm:"type:varchar(20);default:'pending'"`
}
