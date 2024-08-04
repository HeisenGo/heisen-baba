package mappers


import (
	"hotel/internal/reservation"
	"hotel/pkg/adapters/storage/entities"
	"hotel/pkg/fp"
	"gorm.io/gorm"
)

func ReservationEntityToDomain(reservationEntity entities.Reservation) reservation.Reservation {
	return reservation.Reservation{
		ID:         reservationEntity.ID,
		RoomID:     reservationEntity.RoomID,
		UserID:     reservationEntity.UserID,
		CheckIn:    reservationEntity.CheckIn,
		CheckOut:   reservationEntity.CheckOut,
		TotalPrice: reservationEntity.TotalPrice,
		Status:     reservationEntity.Status,
	}
}

func BatchReservationEntitiesToDomain(reservationEntities []entities.Reservation) []reservation.Reservation {
	return fp.Map(reservationEntities, ReservationEntityToDomain)
}

func ReservationDomainToEntity(res *reservation.Reservation) entities.Reservation {
	return entities.Reservation{
		Model: gorm.Model{
			ID: res.ID,
		},
		RoomID:     res.RoomID,
		UserID:     res.UserID,
		CheckIn:    res.CheckIn,
		CheckOut:   res.CheckOut,
		TotalPrice: res.TotalPrice,
		Status:     res.Status,
	}
}