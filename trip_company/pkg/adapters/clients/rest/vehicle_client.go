package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"tripcompanyservice/internal/vehicle"
	vehiclerequest "tripcompanyservice/internal/vehicle_request"
	"tripcompanyservice/pkg/ports"
)

type RestVehicleClient struct {
	ServiceRegistry        ports.IServiceRegistry
	VServiceName string
}

func NewRestVehicleClient(serviceRegistry ports.IServiceRegistry, vServiceName string) *RestVehicleClient {
	return &RestVehicleClient{ServiceRegistry: serviceRegistry, VServiceName: vServiceName}
}

func (r *RestVehicleClient) SelectVehicles(vr vehiclerequest.VehicleRequest) (*vehicle.FullVehicleResponse, error) {

	url := fmt.Sprintf("http://%s/api/v1/companies/vehicles/select?num_passengers=%v?cost=%v?production_year=%d", r.VServiceName, vr.MinCapacity, vr.VehicleReservationFee, vr.ProductionYearMin)

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
