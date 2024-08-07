package service

import (
	"context"
	"errors"
	"time"
	"tripcompanyservice/internal/trip"
	vehiclerequest "tripcompanyservice/internal/vehicle_request"
	"tripcompanyservice/pkg/ports/clients/clients"

	"github.com/google/uuid"
)

var (
	ErrUnableToVehicleReq = errors.New("unable to make vehicle req for this trip")
	ErrAlreadyHasVehicle  = errors.New("already has vehicle at first remove it")
)

type VehicleReService struct {
	vehicleReqOps *vehiclerequest.Ops
	tripOps       *trip.Ops
	vClient  clients.IVehicleClient

}

func NewVehicleReService(vehicleReqOps *vehiclerequest.Ops, tripOps *trip.Ops) *VehicleReService {
	return &VehicleReService{
		vehicleReqOps: vehicleReqOps,
		tripOps:       tripOps,
	}
}

func (s *VehicleReService) CreateVehicleReq(ctx context.Context, vR *vehiclerequest.VehicleRequest, creatorID uuid.UUID) error {
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

	if t.VehicleRequest != nil {
		return ErrAlreadyHasVehicle
	}
	vR.VehicleType = string(t.TripType)
	err = s.vehicleReqOps.Create(ctx, vR)
	if err != nil {
		return err
	}
	// send vr to rabbit MQ TODO: with distance and date
	vR.MatchedVehicleID = uint(2)
	vR.VehicleName = "new bmb"
	vR.VehicleProductionYear = 2021
	vR.MatchVehicleSpeed = 220
	vR.Status = "matched"
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
	updates2["match_vehicle_speed"] = vR.MatchVehicleSpeed
	updates2["vehicle_reservation_fee"] = vR.VehicleReservationFee
	updates2["status"] = "matched"

	err = s.vehicleReqOps.UpdateVehicleReq(ctx, vR.ID, updates2)
	if err != nil {
		return err
	}
	// sen to bank? TODO:
	return nil
}

func (s *VehicleReService) Delete(ctx context.Context, vRID uint, creatorID uuid.UUID) error {
	vr, err := s.vehicleReqOps.GetVehicleReqByID(ctx, vRID)
	if err != nil {
		return err
	}
	company, err := s.vehicleReqOps.GetTransportCompanyByVehicleRequestID(ctx, vRID)
	if err != nil {
		return err
	}

	if company.OwnerID != creatorID {
		return ErrForbidden
	}

	err = s.vehicleReqOps.Delete(ctx, vRID)
	if err != nil {
		return err
	}
	updates := make(map[string]interface{})
	updates["end_date"] = nil
	updates["vehicle_id"] = nil
	updates["vehicle_name"] = ""
	updates["vehicle_request_id"] = nil
	err = s.tripOps.UpdateEndDateTrip(ctx, vr.TripID, updates)
	if err != nil {
		return err
	}
	// TODO: send invoice to the bank and cancel vehicle with reservationfEE
	return nil
}
