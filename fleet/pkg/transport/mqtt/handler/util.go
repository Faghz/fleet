package handler

import "strings"

// extractVehicleID extracts vehicle_id from the MQTT topic
// Topic format: fleet/vehicle/{vehicle_id}/location
func extractVehicleID(topic string) string {
	parts := strings.Split(topic, "/")
	if len(parts) >= 3 {
		return parts[2]
	}
	return ""
}
