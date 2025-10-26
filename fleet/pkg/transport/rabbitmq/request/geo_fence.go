package request

type ReachedNearestPointOfInterestEvent struct {
	VehicleID string                                `json:"vehicle_id"`
	Event     string                                `json:"event"`
	Location  ReachedNearestPointOfInterestLocation `json:"location"`
	Timestamp int64                                 `json:"timestamp"`
}

type ReachedNearestPointOfInterestLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
