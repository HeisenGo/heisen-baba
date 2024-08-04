package presenter

import vehiclerequest "tripcompanyservice/internal/vehicle_request"



type CreateVehicleRequest struct{
	TripID            uint `json:"trip_id"` 
	MinCapacity       int `json:"min_capacity"`
	ProductionYearMin  int `json:"product_year"`
}


type CreateVehicleRes struct{
	ID uint  `json:"id"`
	TripID            uint `json:"trip_id"` 
	VehicleType       string  `json:"type"`
	MinCapacity       int `json:"min_capacity"`
	ProductionYearMin  int `json:"product_year"`
	Status            string  `json:"status"`
	MatchedVehicleID  uint `json:"matched_vehicle_id"`
	VehicleReservationFee float64 `json:"vehicle_fee"`
	VehicleProductionYear int `json:"vehicle_product_year"`
	VehicleName        string `json:"name"`
}

func CreateVehicleRequestToVehicleRequest(r *CreateVehicleRequest) *vehiclerequest.VehicleRequest{
	return &vehiclerequest.VehicleRequest{
		TripID: r.TripID,
		MinCapacity: r.MinCapacity,
		ProductionYearMin: r.ProductionYearMin,
	}
}

func VehicleToCreateVehicleRes(v vehiclerequest.VehicleRequest) CreateVehicleRes{
	return CreateVehicleRes{
		ID: v.ID,
		TripID: v.TripID,
		VehicleType: v.VehicleType,
		VehicleName: v.VehicleName,
		VehicleReservationFee: v.VehicleReservationFee,
		VehicleProductionYear: v.VehicleProductionYear,
		MatchedVehicleID: v.MatchedVehicleID,
		MinCapacity: v.MinCapacity,
		ProductionYearMin: v.ProductionYearMin,
	}
}