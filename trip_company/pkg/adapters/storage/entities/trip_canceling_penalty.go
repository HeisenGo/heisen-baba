package entities

import (
	"time"

	"gorm.io/gorm"
)

type TripCancellingPenalty struct {
	gorm.Model
	FirstDate                    time.Time
	FirstDays					uint
	FirstCancellationPercentage  uint
	SecondDate                   time.Time
	SecondDays                   uint
	SecondCancellationPercentage uint
	ThirdDate                    time.Time
	ThirdDays                    uint
	ThirdCancellationPercentage  uint
	TripID                       uint  `gorm:"not null"`
}
