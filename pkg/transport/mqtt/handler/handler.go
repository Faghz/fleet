package handler

import (
	"github.com/elzestia/fleet/configs"
	"github.com/elzestia/fleet/pkg/external/mqtt"
	"github.com/elzestia/fleet/service/vehicle"
	"go.uber.org/zap"
)

type MQTTHandler struct {
	config         *configs.Config
	logger         *zap.Logger
	mqttClient     *mqtt.MQTTClient
	vehicleService *vehicle.VehicleService
}

// CreateMqttConsumer creates a new MQTT handler instance
func CreateMqttConsumer(config *configs.Config, logger *zap.Logger, mqttClient *mqtt.MQTTClient, vehicleService *vehicle.VehicleService) *MQTTHandler {
	mqttHandler := &MQTTHandler{
		config:         config,
		logger:         logger,
		mqttClient:     mqttClient,
		vehicleService: vehicleService,
	}

	// Initialize vehicle-related MQTT handlers
	if err := createVehicleHandler(mqttHandler); err != nil {
		logger.Fatal("[CreateMqttConsumer] Failed to create vehicle MQTT handler", zap.Error(err))
	}

	logger.Info("MQTT handlers initialized successfully")

	return mqttHandler
}

// Shutdown gracefully shuts down the MQTT handler
func (h *MQTTHandler) Shutdown() {
	h.logger.Info("Shutting down MQTT handler")
	h.mqttClient.Disconnect(250) // Wait 250ms for cleanup
}
