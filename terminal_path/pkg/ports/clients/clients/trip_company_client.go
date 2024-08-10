package clients

type ITripCompanyClient interface {
	GetCountPathUnfinishedTrips(uint) (uint, error)
}
