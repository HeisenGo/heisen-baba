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

// TO Do implement other layers!!!
func (o *Ops) CompanyTrips(ctx context.Context, companyID uint, page, pageSize uint) ([]Trip, uint, error) {
	limit := pageSize
	offset := (page - 1) * pageSize

	return o.repo.GetCompanyTrips(ctx, companyID, limit, offset)
}

func (o *Ops) Create(ctx context.Context, trip *Trip) error {
	if trip.AgencyPrice <= 0 {
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

func (o *Ops) GetFullTripByID(ctx context.Context, id uint) (*Trip, error) {
	p, err := o.repo.GetFullTripByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if p == nil {
		return nil, ErrTripNotFound
	}

	return p, nil
}

func (o *Ops) GetTrips(ctx context.Context, originCity, destinationCity, pathType string, startDate *time.Time, requesterType string, pageSize, page uint) ([]Trip, uint, error) {
	limit := pageSize
	offset := (page - 1) * pageSize

	trips, total, err := o.repo.GetTrips(ctx, originCity, destinationCity, pathType, startDate, requesterType, limit, offset)

	return trips, total, err

}

func (o *Ops) UpdateTrip(ctx context.Context, id uint, newTrip, oldTrip *Trip) error {
	updates := make(map[string]interface{})

	if newTrip.Path != nil {
		if newTrip.TripType != oldTrip.TripType {
			if oldTrip.SoldTickets != 0 {
				return ErrCanNotUpdate
			}
		} else {
			updates["trip_type"] = newTrip.TripType
			updates["path_id"] = newTrip.PathID
			updates["path_name"] = newTrip.Path.Name
			updates["origin"] = newTrip.Path.FromTerminal.City
			updates["from_terminal_name"] = newTrip.Path.FromTerminal.Name
			updates["to_terminal_name"] = newTrip.Path.ToTerminal.Name
			updates["destination"] = newTrip.Path.ToTerminal.City
			oldTrip.Path = newTrip.Path
		}
	}

	if newTrip.UserPrice != 0 && newTrip.AgencyPrice == 0 {
		if newTrip.UserPrice > oldTrip.AgencyPrice {
			updates["user_price"] = newTrip.UserPrice
			oldTrip.UserPrice = newTrip.UserPrice

		} else {
			return ErrCanNotUpdate
		}
	}

	if newTrip.UserPrice == 0 && newTrip.AgencyPrice != 0 {
		if oldTrip.UserPrice > newTrip.AgencyPrice {
			updates["agency_price"] = newTrip.UserPrice
			oldTrip.AgencyPrice = newTrip.AgencyPrice
		} else {
			return ErrCanNotUpdate
		}
	}

	if newTrip.UserPrice != 0 && newTrip.AgencyPrice != 0 {
		if newTrip.UserPrice > newTrip.AgencyPrice {
			updates["agency_price"] = newTrip.UserPrice
			updates["user_price"] = newTrip.UserPrice

			oldTrip.AgencyPrice = newTrip.AgencyPrice
			oldTrip.UserPrice = newTrip.UserPrice

		} else {
			return ErrCanNotUpdate
		}
	}

	if newTrip.UserReleaseDate.IsZero() && !newTrip.TourReleaseDate.IsZero() {
		if oldTrip.UserReleaseDate.After(newTrip.TourReleaseDate) {
			updates["tour_release_date"] = newTrip.TourReleaseDate
			oldTrip.TourReleaseDate = newTrip.TourReleaseDate
		} else {
			return ErrCanNotUpdate
		}
	}

	if !newTrip.UserReleaseDate.IsZero() && newTrip.TourReleaseDate.IsZero() {
		if newTrip.UserReleaseDate.After(oldTrip.TourReleaseDate) {
			updates["user_release_date"] = newTrip.UserReleaseDate
			oldTrip.UserReleaseDate = newTrip.UserReleaseDate
		} else {
			return ErrCanNotUpdate
		}
	}

	if !newTrip.UserReleaseDate.IsZero() && !newTrip.TourReleaseDate.IsZero() {
		if newTrip.UserReleaseDate.After(newTrip.TourReleaseDate) {
			updates["user_release_date"] = newTrip.UserReleaseDate
			updates["user_release_date"] = newTrip.UserReleaseDate
			oldTrip.UserReleaseDate = newTrip.UserReleaseDate
			oldTrip.UserReleaseDate = newTrip.UserReleaseDate
		} else {
			return ErrCanNotUpdate
		}
	}

	if newTrip.MinPassengers != 0 {
		if oldTrip.VehicleID != nil {
			return ErrCanNotUpdate
		} else {
			updates["min_passengers"] = newTrip.MinPassengers
			oldTrip.MinPassengers = newTrip.MinPassengers
		}
	}

	if newTrip.TechTeamID != nil {
		//TO DO: check exsistance of that tech team and their availablity!!!!
		updates["tech_team_id"] = newTrip.TechTeamID
		oldTrip.TechTeamID = newTrip.TechTeamID
	}

	if newTrip.MaxTickets != 0 {
		if oldTrip.VehicleID != nil {
			return ErrCanNotUpdate
		} else {
			updates["max_tickets"] = newTrip.MaxTickets
			oldTrip.MaxTickets = newTrip.MaxTickets
		}
	}

	if newTrip.IsCanceled {
		updates["is_canceled"] = newTrip.IsCanceled
		oldTrip.IsCanceled = newTrip.IsCanceled
		// To DO notify sold tickets and move money
	}

	if newTrip.IsFinished {
		updates["is_finished"] = newTrip.IsFinished
		oldTrip.IsFinished = newTrip.IsFinished
		// TO DO move money from alibaba to company's owner wallet
	}

	if !newTrip.StartDate.IsZero() {
		if oldTrip.SoldTickets != 0 {
			return ErrCanNotUpdate
		} else {
			updates["start_date"] = newTrip.StartDate
			oldTrip.StartDate = newTrip.StartDate
		}
	}

	return o.repo.UpdateTrip(ctx, id, updates)

}
