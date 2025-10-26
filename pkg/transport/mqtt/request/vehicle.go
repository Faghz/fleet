package request

// VehicleLocationRequest represents the payload for vehicle location updates
type VehicleLocationRequest struct {
	VehicleID string  `json:"vehicle_id" validate:"required"`
	Latitude  float64 `json:"latitude" validate:"required,min=-90,max=90"`
	Longitude float64 `json:"longitude" validate:"required,min=-180,max=180"`
	Timestamp int64   `json:"timestamp" validate:"required"`
}
