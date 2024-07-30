package entities

import (
	"time"

	"gorm.io/gorm"
)

type TripCancelingPenalty struct {
	gorm.Model
	FirstDate                    time.Time
	FirstCancellationPercentage  uint
	SecondDate                   time.Time
	SecondCancellationPercentage uint
	ThirdDate                    time.Time
	ThirdCancellationPercentage  uint
	TripID                       uint
}
