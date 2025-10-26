package rabbitmq

import (
	"fmt"

	"github.com/elzestia/fleet/configs"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type RabbitMQClient struct {
	conn   *amqp.Connection
	config *configs.RabbitMQConfig
	logger *zap.Logger
}

// CreateRabbitMQConnection creates and returns a new RabbitMQ connection
func CreateRabbitMQConnection(config *configs.RabbitMQConfig, logger *zap.Logger) *RabbitMQClient {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/%s", config.User, config.Password, config.Host, config.Port, config.Vhost)

	conn, err := amqp.Dial(url)
	if err != nil {
		logger.Fatal("Failed to connect to RabbitMQ", zap.Error(err))
	}

	logger.Info("RabbitMQ client connected", zap.String("host", config.Host), zap.String("port", config.Port))

	return &RabbitMQClient{
		conn:   conn,
		config: config,
		logger: logger,
	}
}

// GetConnection returns the AMQP connection
func (r *RabbitMQClient) GetConnection() *amqp.Connection {
	return r.conn
}

// Close closes the RabbitMQ connection
func (r *RabbitMQClient) Close() error {
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}
