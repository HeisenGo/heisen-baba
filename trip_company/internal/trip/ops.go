package trip

import (
	"context"
	"time"
)

type Ops struct {
	repo Repo
	//penaltyRepo
}

func NewOps(repo Repo) *Ops {
	return &Ops{repo}
}

func (o *Ops) CompanyTrips(ctx context.Context, companyID uint, page, pageSize uint) ([]Trip, uint, error) {
	limit := pageSize
	offset := (page - 1) * pageSize

	return o.repo.GetCompanyTrips(ctx, companyID, limit, offset)
}

func (o *Ops) Create(ctx context.Context, trip *Trip) error {
	if trip.AgencyPrice <= 0{
		return ErrNegativePrice
	}
	if trip.UserPrice < trip.AgencyPrice {
		return ErrPrice
	}
	if trip.StartDate.Before(time.Now()) {
		return ErrStartTime
	}

	// delta := trip.StartDate.Sub(time.Now())
	// deltaDays := int(delta.Hours() / 24)

	// if deltaDays < int(trip.TripCancellingPenalty.FirstDays){
	// 	return
	// }

	if trip.UserReleaseDate.Before(trip.TourReleaseDate) {
		return ErrReleaseDate
	}

	if trip.TripCancellingPenalty.FirstDays <= trip.TripCancellingPenalty.SecondDays {
		return ErrFirstPenalty
	}

	if trip.TripCancellingPenalty.SecondDays <= trip.TripCancellingPenalty.ThirdDays {
		return ErrSecondPenalty
	}

	trip.TripCancellingPenalty.FirstDate = trip.StartDate.AddDate(0, 0, int(-trip.TripCancellingPenalty.FirstDays))
	trip.TripCancellingPenalty.SecondDate = trip.StartDate.AddDate(0, 0, int(-trip.TripCancellingPenalty.SecondDays))
	trip.TripCancellingPenalty.ThirdDate = trip.StartDate.AddDate(0, 0, int(-trip.TripCancellingPenalty.ThirdDays))

	return o.repo.Insert(ctx, trip)
}
