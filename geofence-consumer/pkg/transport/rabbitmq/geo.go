package rabbitmqhndl

import (
	"context"
	"encoding/json"

	"github.com/elzestia/fleet/pkg/transport/rabbitmq/request"
	"go.uber.org/zap"
)

func (h *RabbitMQHandler) createGeoFenceConsumer() error {
	if !h.config.RabbitMQ.Consumer.GeoFence.Enabled {
		h.logger.Debug("GeoFence consumer is disabled, skipping creation")
		return nil
	}

	h.logger.Debug("Creating GeoFence consumer", zap.String("exchange", h.config.RabbitMQ.Consumer.GeoFence.Exchange), zap.String("queue", h.config.RabbitMQ.Consumer.GeoFence.Queue))
	msgs, err := h.rabbitMqClient.CreateConsumer(
		h.config.RabbitMQ.Consumer.GeoFence.Exchange,
		h.config.RabbitMQ.Consumer.GeoFence.Queue,
		h.config.RabbitMQ.Consumer.GeoFence.RoutingKey,
		h.config.RabbitMQ.Consumer.GeoFence.ConsumerTag,
	)
	if err != nil {
		h.logger.Error("Failed to create GeoFence consumer", zap.Error(err))
		return err
	}

	for msg := range msgs {
		h.logger.Debug("Received message from GeoFence consumer", zap.String("exchange", msg.Exchange), zap.String("routingKey", msg.RoutingKey), zap.String("consumerTag", msg.ConsumerTag))

		var geoFenceEvent request.ReachedNearestPointOfInterestEvent
		if err := json.Unmarshal(msg.Body, &geoFenceEvent); err != nil {
			h.logger.Error("Failed to unmarshal GeoFence event", zap.Error(err))
			continue
		}

		err = h.services.GeoFenceService.ConsumeGeoFenceEvents(context.Background(), geoFenceEvent)
		if err != nil {
			h.logger.Error("Failed to process GeoFence event", zap.Error(err), zap.Any("event", geoFenceEvent))
			continue
		}
		h.logger.Info("Processed GeoFence event", zap.Any("event", geoFenceEvent))
	}

	return nil
}
