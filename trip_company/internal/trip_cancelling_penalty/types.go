package tripcancellingpenalty

import (
	"context"
	"time"
)

type Repo interface {
	GetByID(ctx context.Context, id uint) (*TripCancelingPenalty, error)
}

type TripCancelingPenalty struct {
	ID                           uint
	FirstDate                    time.Time
	FirstDays                    uint
	FirstCancellationPercentage  uint
	SecondDate                   time.Time
	SecondDays                   uint
	SecondCancellationPercentage uint
	ThirdDate                    time.Time
	ThirdCancellationPercentage  uint
	ThirdDays                    uint
	TripID                       uint
}
