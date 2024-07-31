package presenter

import tripcancellingpenalty "tripcompanyservice/internal/trip_cancelling_penalty"

type CreateTripCancelingPenaltyReq struct {
	ID                           uint `json:"id"`
	FirstDays                    uint `json:"first_days" validate:"required"`
	FirstCancellationPercentage  uint `json:"first_percentage" validate:"required"`
	SecondDays                   uint `json:"second_days" validate:"required"`
	SecondCancellationPercentage uint `json:"second_percentage" validate:"required"`
	ThirdCancellationPercentage  uint `json:"third_percentage" validate:"required"`
	ThirdDays                    uint `json:"third_days" validate:"required"`
}

func CreateTripCancelingPenaltyReqToTripCancellingPenalty(req *CreateTripCancelingPenaltyReq) *tripcancellingpenalty.TripCancelingPenalty {

	return &tripcancellingpenalty.TripCancelingPenalty{
		FirstDays:                    req.FirstDays,
		FirstCancellationPercentage:  req.FirstCancellationPercentage,
		SecondDays:                   req.SecondDays,
		SecondCancellationPercentage: req.SecondCancellationPercentage,
		ThirdDays:                    req.ThirdDays,
		ThirdCancellationPercentage:  req.ThirdCancellationPercentage,
	}
}

func TripCancelingPenaltyToTripCancellingPenaltyReq(tCp *tripcancellingpenalty.TripCancelingPenalty) *CreateTripCancelingPenaltyReq {
	return &CreateTripCancelingPenaltyReq{
		ID:                           tCp.ID,
		FirstDays:                    tCp.FirstDays,
		FirstCancellationPercentage:  tCp.FirstCancellationPercentage,
		SecondDays:                   tCp.SecondDays,
		SecondCancellationPercentage: tCp.SecondCancellationPercentage,
		ThirdDays:                    tCp.ThirdDays,
		ThirdCancellationPercentage:  tCp.ThirdCancellationPercentage,
	}
}
