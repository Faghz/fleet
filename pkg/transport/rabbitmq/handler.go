package rabbitmqhndl

import (
	"github.com/elzestia/fleet/configs"
	"github.com/elzestia/fleet/pkg/external/rabbitmq"
)

type RabbitMQHandler struct {
	rabbitMqClient *rabbitmq.RabbitMQClient
	config         *configs.Config
}

func CreateRabbitMQHandler(rabbitMqClient *rabbitmq.RabbitMQClient, config *configs.Config) *RabbitMQHandler {
	rabbitMqHandler := &RabbitMQHandler{
		rabbitMqClient: rabbitMqClient,
		config:         config,
	}

	rabbitMqHandler.createRabbitMqExchange()

	return rabbitMqHandler
}

func (h *RabbitMQHandler) Shutdown() {
	h.rabbitMqClient.Close()
}
