package handler

import (
	"context"
	"encoding/json"

	pahomqtt "github.com/eclipse/paho.mqtt.golang"
	inthttp "github.com/elzestia/fleet/pkg/transport/http"
	"github.com/elzestia/fleet/pkg/transport/mqtt/request"
	"go.uber.org/zap"
)

func createVehicleHandler(h *MQTTHandler) error {
	// Subscribe to vehicle location topic with wildcard for vehicle_id
	if err := h.mqttClient.Subscribe("fleet/vehicle/+/location", h.syncVehicleLocation); err != nil {
		h.logger.Error("Failed to subscribe to vehicle location topic",
			zap.String("topic", "fleet/vehicle/+/location"),
			zap.Error(err))
		return err
	}

	return nil
}

// syncVehicleLocation processes incoming vehicle location messages
func (h *MQTTHandler) syncVehicleLocation(client pahomqtt.Client, msg pahomqtt.Message) {
	topic := msg.Topic()
	payload := msg.Payload()

	h.logger.Info("[syncVehicleLocation] Received vehicle location message",
		zap.String("topic", topic),
		zap.Int("payload_size", len(payload)))

	var locationReq request.VehicleLocationRequest

	if err := json.Unmarshal(payload, &locationReq); err != nil {
		h.logger.Error("[syncVehicleLocation] Failed to unmarshal vehicle location message",
			zap.Error(err),
			zap.String("topic", topic))
		h.mqttClient.PublishResponse("error", topic, "Invalid JSON format", err.Error())
		return
	}

	// Extract vehicle ID from topic
	vehicleID := extractVehicleID(topic)
	if vehicleID == "" {
		h.logger.Error("[syncVehicleLocation] Failed to extract vehicle ID from topic",
			zap.String("topic", topic))
		h.mqttClient.PublishResponse("error", topic, "Invalid topic format", "Could not extract vehicle_id")
		return
	}

	// Ensure vehicle ID in payload matches topic
	if locationReq.VehicleID == "" {
		locationReq.VehicleID = vehicleID
	} else if locationReq.VehicleID != vehicleID {
		h.logger.Warn("[syncVehicleLocation] Vehicle ID mismatch between topic and payload",
			zap.String("topic_vehicle_id", vehicleID),
			zap.String("payload_vehicle_id", locationReq.VehicleID))
		h.mqttClient.PublishResponse("error", topic, "Validation error", "Vehicle ID mismatch")
		return
	}

	// Validate the request
	if err := inthttp.GetValidator().Validate(&locationReq); err != nil {
		h.logger.Error("[syncVehicleLocation] Validation failed for vehicle location",
			zap.Error(err),
			zap.String("topic", topic))
		h.mqttClient.PublishResponse("error", topic, "Validation failed", err.Error())
		return
	}

	// Process the location data (implement your business logic here)
	h.vehicleService.ProcessVehicleLocationSync(context.Background(), &locationReq)

	// Send success response
	h.mqttClient.PublishResponse("success", topic, "Location updated successfully", map[string]interface{}{
		"vehicle_id": locationReq.VehicleID,
		"timestamp":  locationReq.Timestamp,
	})
}
