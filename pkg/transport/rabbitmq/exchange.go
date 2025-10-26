package rabbitmqhndl

import "github.com/rabbitmq/amqp091-go"

func (h *RabbitMQHandler) createRabbitMqExchange() error {
	createFleetExchange(h.rabbitMqClient.GetChannel(), h.config.RabbitMQ.Exchange.Fleet.Name, h.config.RabbitMQ.Exchange.Fleet.Kind)

	return nil
}

func createFleetExchange(ch *amqp091.Channel, name, kind string) error {
	ch.ExchangeDeclare(
		name,
		kind,
		true,
		false,
		false,
		false,
		nil,
	)

	return nil
}
