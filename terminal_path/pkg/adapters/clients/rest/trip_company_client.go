package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"terminalpathservice/pkg/ports"
)

type RestTripCompanyClient struct {
	ServiceRegistry        ports.IServiceRegistry
	TripCompanyServiceName string
}

func NewRestTripCompanyClient(serviceRegistry ports.IServiceRegistry, tripCompanyServiceName string) *RestTripCompanyClient {
	return &RestTripCompanyClient{ServiceRegistry: serviceRegistry, TripCompanyServiceName: tripCompanyServiceName}
}

func (r *RestTripCompanyClient) GetCountPathUnfinishedTrips(pathID uint) (uint, error) {

	url := fmt.Sprintf("http://%s/api/v1/companies/path-trips/%v", r.TripCompanyServiceName, pathID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("Error creating request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("Error reading response: %v", err)
	}

	var responseData struct {
		Count int `json:"count"`
	}

	if err := json.Unmarshal(body, &responseData); err != nil {
		return 0, fmt.Errorf("Error unmarshalling JSON: %v", err)
	}

	return uint(responseData.Count), nil
}
