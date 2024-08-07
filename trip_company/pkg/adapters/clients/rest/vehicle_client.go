package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"terminalpathservice/pkg/ports"
	"tripcompanyservice/internal/vehicle"
	vehiclerequest "tripcompanyservice/internal/vehicle_request"
)

type RestVehicleClient struct {
	ServiceRegistry        ports.IServiceRegistry
	VehicleServiceName string
}

func NewVehicleClient(serviceRegistry ports.IServiceRegistry, vServiceName string) *RestVehicleClient {
	return &RestVehicleClient{ServiceRegistry: serviceRegistry, VehicleServiceName: vServiceName}
}

func (r *RestVehicleClient) SelectVehicles(vr vehiclerequest.VehicleRequest) (vehicle.FullVehicleResponse, error) {

	url := fmt.Sprintf("http://%s/api/v1/companies/vehicles/select%v", )

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