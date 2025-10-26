package mqtt

import (
	"encoding/json"
	"fmt"
	"time"

	pahomqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/elzestia/fleet/configs"
	"github.com/elzestia/fleet/pkg/transport/mqtt/response"
	"go.uber.org/zap"
)

type MQTTClient struct {
	client pahomqtt.Client
	config *configs.MQTTConfig
	logger *zap.Logger
}

// CreateMQTTConnection creates and returns a new MQTT client connection
func CreateMQTTConnection(config *configs.MQTTConfig, logger *zap.Logger) *MQTTClient {
	opts := pahomqtt.NewClientOptions()

	// Set broker address
	broker := fmt.Sprintf("tcp://%s:%s", config.Broker, config.Port)
	opts.AddBroker(broker)

	// Set client ID
	opts.SetClientID(config.ClientID)

	// Set credentials if provided
	if config.Username != "" {
		opts.SetUsername(config.Username)
	}
	if config.Password != "" {
		opts.SetPassword(config.Password)
	}

	// Set auto reconnect
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(10 * time.Second)
	opts.SetMaxReconnectInterval(1 * time.Minute)

	// Set connection handlers
	opts.SetOnConnectHandler(func(c pahomqtt.Client) {
		logger.Info("MQTT client connected to broker", zap.String("broker", broker))
	})

	opts.SetConnectionLostHandler(func(c pahomqtt.Client, err error) {
		logger.Error("MQTT connection lost", zap.Error(err))
	})

	opts.SetReconnectingHandler(func(c pahomqtt.Client, opts *pahomqtt.ClientOptions) {
		logger.Warn("MQTT client attempting to reconnect")
	})

	// Create client
	client := pahomqtt.NewClient(opts)

	// Connect to broker
	logger.Info("Connecting to MQTT broker", zap.String("broker", broker))
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		logger.Fatal("Failed to connect to MQTT broker", zap.Error(token.Error()))
	}

	return &MQTTClient{
		client: client,
		config: config,
		logger: logger,
	}
}

// Subscribe subscribes to a topic with a message handler
func (m *MQTTClient) Subscribe(topic string, handler pahomqtt.MessageHandler) error {
	qos := m.config.QoS
	if qos > 2 {
		qos = 0 // Default to QoS 0 if invalid
	}

	token := m.client.Subscribe(topic, qos, handler)
	token.Wait()

	if token.Error() != nil {
		m.logger.Error("Failed to subscribe to topic",
			zap.String("topic", topic),
			zap.Error(token.Error()))
		return token.Error()
	}

	m.logger.Info("Successfully subscribed to topic",
		zap.String("topic", topic),
		zap.Uint8("qos", qos))
	return nil
}

// Publish publishes a message to a topic
func (m *MQTTClient) Publish(topic string, payload interface{}) error {
	qos := m.config.QoS
	if qos > 2 {
		qos = 0
	}

	token := m.client.Publish(topic, qos, false, payload)
	token.Wait()

	if token.Error() != nil {
		m.logger.Error("Failed to publish to topic",
			zap.String("topic", topic),
			zap.Error(token.Error()))
		return token.Error()
	}

	return nil
}

// Unsubscribe unsubscribes from a topic
func (m *MQTTClient) Unsubscribe(topic string) error {
	token := m.client.Unsubscribe(topic)
	token.Wait()

	if token.Error() != nil {
		m.logger.Error("Failed to unsubscribe from topic",
			zap.String("topic", topic),
			zap.Error(token.Error()))
		return token.Error()
	}

	m.logger.Info("Successfully unsubscribed from topic", zap.String("topic", topic))
	return nil
}

// publishResponse publishes a response to the response topic
func (h *MQTTClient) PublishResponse(status string, requestTopic, message string, data interface{}) {
	responseTopic := fmt.Sprintf("%s/response", requestTopic)

	resp := response.VehicleLocationResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}

	respJSON, err := json.Marshal(resp)
	if err != nil {
		h.logger.Error("Failed to marshal success response", zap.Error(err))
		return
	}

	if err := h.Publish(responseTopic, respJSON); err != nil {
		h.logger.Error("Failed to publish success response",
			zap.String("topic", responseTopic),
			zap.Error(err))
	}
}

// Disconnect disconnects from the MQTT broker
func (m *MQTTClient) Disconnect(quiesce uint) {
	m.logger.Info("Disconnecting from MQTT broker")
	m.client.Disconnect(quiesce)
}

// IsConnected returns true if the client is connected
func (m *MQTTClient) IsConnected() bool {
	return m.client.IsConnected()
}

// GetClient returns the underlying MQTT client
func (m *MQTTClient) GetClient() pahomqtt.Client {
	return m.client
}
