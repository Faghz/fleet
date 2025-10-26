package response

import "net/http"

var (
	ErrorVehicleNotFound = GenerateFailure(http.StatusNotFound, "Vehicle Not Found", "The requested vehicle does not exist")
)

type VehicleLocation struct {
	VehicleID string  `json:"vehicle_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp int64   `json:"timestamp"`
	UpdatedAt string  `json:"updated_at"`
}
