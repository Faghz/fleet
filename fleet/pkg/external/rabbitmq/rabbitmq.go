package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/elzestia/fleet/configs"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type RabbitMQClient struct {
	conn    *amqp.Connection
	config  *configs.RabbitMQConfig
	channel *amqp.Channel
}

// CreateRabbitMQConnection creates and returns a new RabbitMQ connection
func CreateRabbitMQConnection(config *configs.RabbitMQConfig, logger *zap.Logger) *RabbitMQClient {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/%s", config.User, config.Password, config.Host, config.Port, config.Vhost)

	conn, err := amqp.Dial(url)
	if err != nil {
		logger.Fatal("Failed to connect to RabbitMQ", zap.Error(err))
		return nil
	}
	channel, err := conn.Channel()
	if err != nil {
		logger.Fatal("Failed to create RabbitMQ channel", zap.Error(err))
		return nil
	}

	logger.Info("RabbitMQ client connected", zap.String("host", config.Host), zap.String("port", config.Port))

	return &RabbitMQClient{
		conn:    conn,
		channel: channel,
		config:  config,
	}
}

// GetConnection returns the AMQP connection
func (r *RabbitMQClient) GetConnection() *amqp.Connection {
	return r.conn
}

func (r *RabbitMQClient) GetChannel() *amqp.Channel {
	return r.channel
}

func (r *RabbitMQClient) Publish(ctx context.Context, exchange string, routingKey string, body interface{}) error {
	// Serialize the body to JSON
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	err = r.channel.PublishWithContext(
		ctx,
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonBody,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// Close closes the RabbitMQ connection
func (r *RabbitMQClient) Close() error {
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}
