package trip

import (
	"context"
	"errors"
	"time"
	"tripcompanyservice/internal/company"
	"tripcompanyservice/internal/techteam"
	tripcancellingpenalty "tripcompanyservice/internal/trip_cancelling_penalty"
	vehiclerequest "tripcompanyservice/internal/vehicle_request"
)

type TripType string

const (
	RailTrip    TripType = "rail"
	AirTrip     TripType = "air"
	RoadTrip    TripType = "road"
	SailingTrip TripType = "sailing"
)

var (
	ErrTripUnAvailable = errors.New("trip is unavailable")
	ErrInvalidPercentage = errors.New("invalid percent")
	ErrCanNotUpdate    = errors.New("can not update")
	ErrNotUpdated      = errors.New("could not update trip")
	ErrRecordsNotFound = errors.New("no records found")
	ErrTripNotFound    = errors.New("trip not found")
	ErrFailedToGetTrip = errors.New("failed to get trip")
	ErrRecordNotFound  = errors.New("records not found")
	ErrNegativePrice   = errors.New("price should be more than 0")
	ErrPrice           = errors.New("agency price should not be greater than uer price")
	ErrReleaseDate     = errors.New("agency release should not be greater than uer release date")

	ErrStartTime   = errors.New("wrong start time")
	ErrDuplication = errors.New("duplicated trip")

	ErrFirstPenalty    = errors.New("first penalty days should be more than the second penalty days")
	ErrSecondPenalty   = errors.New("second penalty days should be more than the third penalty days")
	ErrFailedToRestore = errors.New("failed to restore duplicated deleted record")
)

type Repo interface {
	GetCompanyTrips(ctx context.Context, originCity, destinationCity, pathType string, startDate *time.Time,requesterType string ,companyID uint, limit, offset uint) ([]Trip, uint, error)
	Insert(ctx context.Context, t *Trip) error

	GetFullTripByID(ctx context.Context, id uint) (*Trip, error)
	GetTrips(ctx context.Context, originCity, destinationCity, pathType string, startDate *time.Time, requesterType string,limit, offset uint) ([]Trip, uint, error)
	UpdateTrip(ctx context.Context, id uint, updates map[string]interface{}) error

	GetCountPathUnfinishedTrips(ctx context.Context, pathID uint) (uint, error)
	GetUpcomingUnconfirmedTripIDsToCancel(ctx context.Context) ([]uint, error)

	GetCountCompanyUnfinishedUncanceledTrips(ctx context.Context, companyID uint) (uint, error)
	CheckAvailabilityTechTeam(ctx context.Context, tripID uint, techTeamID uint, startDate time.Time, endDate time.Time) (bool, error)

	// TODO: remove if no tickets are sold
}

type Trip struct {
	ID                    uint
	TransportCompanyID    uint
	TransportCompany      company.TransportCompany
	TripType              TripType
	UserReleaseDate       time.Time
	TourReleaseDate       time.Time
	UserPrice             float64
	AgencyPrice           float64
	PathID                uint
	Path                  *Path
	Origin                string
	Destination           string
	Status                string `gorm:"type:varchar(20);default:'pending'"`
	MinPassengers         uint
	SoldTickets           uint
	TechTeamID            *uint
	TechTeam              *techteam.TechTeam
	TripCancellingPenalty *tripcancellingpenalty.TripCancelingPenalty
	//TechTeam         TechTeam `gorm:"foreignKey:TechTeamID"`
	//VehicleRequestID uint
	//VehicleRequest   VehicleRequest `gorm:"foreignKey:TripID"`
	//Tickets          []Ticket       `gorm:"foreignKey:TripID"`
	TripCancelingPenaltyID *uint
	MaxTickets             uint
	VehicleID              *uint
	VehicleRequestID       *uint
	VehicleRequest *vehiclerequest.VehicleRequest
	IsCanceled             bool
	IsFinished             bool
	IsConfirmed            bool
	StartDate              *time.Time // should be given by trip generator
	EndDate                *time.Time // should be calculated according to the vehicle speed and path distance
    Profit      float64
}

type Path struct {
	ID             uint
	FromTerminalID uint
	ToTerminalID   uint
	FromTerminal   *Terminal
	ToTerminal     *Terminal
	DistanceKM     float64 // in kilometers
	Code           string
	Name           string
	Type           string
}

type Terminal struct {
	ID      uint
	Name    string
	Type    string
	City    string
	Country string
}

func NewTripTOUpdateSoldTickets(soldTickets uint) *Trip {
	return &Trip{
		SoldTickets: soldTickets,
	}
}
