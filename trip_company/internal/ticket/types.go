package ticket

import "tripcompanyservice/internal/trip"


type Repo interface {
}

type Ticket struct {
    ID uint
    TripID     uint            
    Trip       trip.Trip         
    UserID     *uint            // Use `default:NULL` for nullable field
    AgencyID   *uint           
    Quantity   int
    TotalPrice float64
    Status     string         
}