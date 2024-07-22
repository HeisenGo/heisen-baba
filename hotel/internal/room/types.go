package room

type Room struct {
	ID          uint
	HotelID     uint
	AgencyPrice uint64
	UserPrice   uint64
	Facilities  string
	Capacity    uint8
	IsAvailable bool
}
