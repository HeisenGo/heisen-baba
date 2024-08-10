package mappers

import (
	tripcancellingpenalty "tripcompanyservice/internal/trip_cancelling_penalty"
	"tripcompanyservice/pkg/adapters/storage/entities"
	"tripcompanyservice/pkg/fp"
)

func PenaltyEntityToDomain(penaltyEntity entities.TripCancellingPenalty) tripcancellingpenalty.TripCancelingPenalty {
	return tripcancellingpenalty.TripCancelingPenalty{
		ID:                           penaltyEntity.ID,
		FirstDate:                    penaltyEntity.FirstDate,
		FirstDays:                    penaltyEntity.FirstDays,
		FirstCancellationPercentage:  penaltyEntity.FirstCancellationPercentage,
		SecondDate:                   penaltyEntity.SecondDate,
		SecondDays:                   penaltyEntity.SecondDays,
		SecondCancellationPercentage: penaltyEntity.SecondCancellationPercentage,
		ThirdDate:                    penaltyEntity.ThirdDate,
		ThirdDays:                    penaltyEntity.ThirdDays,
		ThirdCancellationPercentage:  penaltyEntity.ThirdCancellationPercentage,
		TripID:                       penaltyEntity.TripID,
	}
}

func PenaltyEntitiesToDomain(penaltyEntities []entities.TripCancellingPenalty) []tripcancellingpenalty.TripCancelingPenalty {
	return fp.Map(penaltyEntities, PenaltyEntityToDomain)
}

func PenaltyDomainToEntity(p *tripcancellingpenalty.TripCancelingPenalty) *entities.TripCancellingPenalty {
	return &entities.TripCancellingPenalty{
		FirstDate:                    p.FirstDate,
		FirstDays:                    p.FirstDays,
		FirstCancellationPercentage:  p.FirstCancellationPercentage,
		SecondDate:                   p.SecondDate,
		SecondDays:                   p.SecondDays,
		SecondCancellationPercentage: p.SecondCancellationPercentage,
		ThirdDate:                    p.ThirdDate,
		ThirdDays:                    p.ThirdDays,
		ThirdCancellationPercentage:  p.ThirdCancellationPercentage,
	}
}
