package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"tripcompanyservice/internal/vehicle"
	vehiclerequest "tripcompanyservice/internal/vehicle_request"
)

// RestVehicleClient defines the client used to interact with the vehicle service.
type RestVehicleClient struct {
	baseURL string
}

// NewRestVehicleClient initializes a new RestVehicleClient with the base URL of the service.
func NewRestVehicleClient(baseURL string) *RestVehicleClient {
	return &RestVehicleClient{baseURL: baseURL}
}

// SelectVehicles sends a GET request to fetch vehicles based on the VehicleRequest.
func (r *RestVehicleClient) SelectVehicles(vr vehiclerequest.VehicleRequest) (*vehicle.FullVehicleResponse, error) {
	url := fmt.Sprintf("%s/api/v1/companies/vehicles/select?id=%d", r.baseURL, vr.ID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK response status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	var response vehicle.FullVehicleResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	return &response, nil
}
