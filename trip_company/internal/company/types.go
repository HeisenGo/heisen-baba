package company

type TransportCompany struct {
	ID          uint
	Name        string
	Description string
	OwnerID     uint
	Address     string
	PhoneNumber string
	Email       string
	// relationships
	//Employees   []Employee
	//Trips       []Trip
	//TechTeams   []TechTeam
}
