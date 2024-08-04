package entities

import "gorm.io/gorm"

type Ticket struct {
    gorm.Model
    TripID     uint            `gorm:"not null"`
    Trip       *Trip            `gorm:"foreignKey:TripID; constraint:OnDelete:CASCADE;"`
    UserID     *uint           `gorm:"default:NULL"` // Use `default:NULL` for nullable field
    AgencyID   *uint           `gorm:"default:NULL"`
    Quantity   int
    TotalPrice float64
    Status     string          `gorm:"type:varchar(20);default:'pending'"` //confirmed
    Penalty    float64   `gorm:"default:0"`
    Invoice         *Invoice    `gorm:"foreignKey:TicketID; constraint:OnDelete:CASCADE;"`
}
