package entities

import "gorm.io/gorm"

type Ticket struct {
    gorm.Model
    TripID     uint            `gorm:"not null"`
    Trip       *Trip            `gorm:"foreignKey:TripID"`
    UserID     *uint           `gorm:"default:NULL"` // Use `default:NULL` for nullable field
    AgencyID   *uint           `gorm:"default:NULL"`
    Quantity   int
    TotalPrice float64
    Status     string          `gorm:"type:varchar(20);default:'pending'"` //confirmed
}
