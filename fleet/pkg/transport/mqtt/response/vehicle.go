package response

// VehicleLocationResponse represents the response for vehicle location operations
type VehicleLocationResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}
