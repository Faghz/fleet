package rabbitmqhndl

import (
	"github.com/elzestia/fleet/configs"
	"github.com/elzestia/fleet/pkg/external/rabbitmq"
	"github.com/elzestia/fleet/service"
	"go.uber.org/zap"
)

type RabbitMQHandler struct {
	rabbitMqClient *rabbitmq.RabbitMQClient
	config         *configs.Config
	logger         *zap.Logger
	services       *service.Services
}

func CreateRabbitMQHandler(rabbitMqClient *rabbitmq.RabbitMQClient, config *configs.Config, logger *zap.Logger, services *service.Services) *RabbitMQHandler {
	rabbitMqHandler := &RabbitMQHandler{
		rabbitMqClient: rabbitMqClient,
		config:         config,
		logger:         logger,
		services:       services,
	}

	rabbitMqHandler.createGeoFenceConsumer()

	return rabbitMqHandler
}

func (h *RabbitMQHandler) Shutdown() {
	h.rabbitMqClient.Close()
}
