package clients

import "tripcompanyservice/internal/trip"



type IPathClient interface {
	GetFullPathByID(pathID uint32) (*trip.Path, error)
}
