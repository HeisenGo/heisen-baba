package vehiclerequest


type Repo interface {
}

type VehicleRequest struct {
	ID                    uint
	TripID                uint
	VehicleType           string
	MinCapacity           int
	ProductionYearMin     int
	Status                string
	MatchedVehicleID      uint
	VehicleReservationFee float64
	VehicleProductionYear int
	VehicleName           string
}
