package service

import (
	"context"
	"errors"
	"time"
	"tripcompanyservice/internal/trip"
	vehiclerequest "tripcompanyservice/internal/vehicle_request"
)

var (
	ErrUnableToVehicleReq = errors.New("unable to make vehicle req for this trip")
)

type VehicleReService struct {
	vehicleReqOps *vehiclerequest.Ops
	tripOps       *trip.Ops
}

func NewVehicleReService(vehicleReqOps *vehiclerequest.Ops, tripOps *trip.Ops) *VehicleReService {
	return &VehicleReService{
		vehicleReqOps: vehicleReqOps,
		tripOps:       tripOps,
	}
}

func (s *VehicleReService) CreateVehicleReq(ctx context.Context, vR *vehiclerequest.VehicleRequest, creatorID uint) error {
	t, err := s.tripOps.GetFullTripByID(ctx, vR.TripID)
	if err != nil {
		return err
	}
	//if not operator
	//if creatorID!=t.TransportCompany.OwnerID
	if t.IsCanceled {
		return ErrUnableToVehicleReq
	}
	if t.IsFinished {
		return ErrUnableToVehicleReq
	}
	if t.VehicleID != nil { // TODO: should be able to remove vehicle !!!!!
		return ErrUnableToVehicleReq
	}
	vR.VehicleType = string(t.TripType)
	err = s.vehicleReqOps.Create(ctx, vR)
	if err != nil {
		return err
	}
	// send vr to rabbit MQ TODO:
	vR.MatchedVehicleID = uint(2)
	vR.VehicleName = "new bmb"
	vR.VehicleProductionYear = 2021
	vR.MatchVehicleSpeed = 220
	travelTimeHours := t.Path.DistanceKM / vR.MatchVehicleSpeed
	travelDuration := time.Duration(travelTimeHours * float64(time.Hour))
	endDate := t.StartDate.Add(travelDuration)
	t.EndDate = &endDate
	vR.VehicleReservationFee = 23000 * travelTimeHours
	updates := make(map[string]interface{})
	updates["end_date"] = endDate
	updates["vehicle_id"] = vR.MatchedVehicleID
	updates["vehicle_name"] = vR.VehicleName
	updates["vehicle_request_id"] = vR.ID
	err = s.tripOps.UpdateEndDateTrip(ctx, t.ID, updates)
	if err != nil {
		return err
	}
	updates2 := make(map[string]interface{})

	updates2["matched_vehicle_id"] = vR.MatchedVehicleID
	updates2["vehicle_name"] = vR.VehicleName
	updates2["vehicle_production_year"] = vR.VehicleProductionYear
	updates2["matched_vehicle_speed"] = vR.MatchVehicleSpeed
	updates2["vehicle_reservation_fee"] = vR.VehicleReservationFee

	err = s.vehicleReqOps.UpdateVehicleReq(ctx, vR.ID, updates2)
	if err != nil {
		return err
	}
	// update vR
	// trip update end_date
	// sen to bank? TODO:
	return nil
}
