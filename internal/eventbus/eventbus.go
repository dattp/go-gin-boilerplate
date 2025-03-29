package eventbus

import (
	LibEventBus "github.com/asaskevich/EventBus"
	"github.com/sirupsen/logrus"
)

// EventBus wraps the EventBus library with additional functionality
type EventBus struct {
	bus    LibEventBus.Bus
	logger *logrus.Logger
}

// New creates a new EventBus instance
func NewEventBus(logger *logrus.Logger) *EventBus {
	return &EventBus{
		bus:    LibEventBus.New(),
		logger: logger,
	}
}

// Subscribe subscribes to a topic with a handler function
func (eb *EventBus) Subscribe(topic string, handler interface{}) error {
	err := eb.bus.Subscribe(topic, handler)
	if err != nil {
		eb.logger.WithFields(logrus.Fields{
			"topic": topic,
			"error": err,
		}).Error("Failed to subscribe to topic")
		return err
	}
	eb.logger.WithField("topic", topic).Debug("Subscribed to topic")
	return nil
}

// SubscribeAsync subscribes to a topic with an asynchronous handler
func (eb *EventBus) SubscribeAsync(topic string, handler interface{}, transactional bool) error {
	err := eb.bus.SubscribeAsync(topic, handler, transactional)
	if err != nil {
		eb.logger.WithFields(logrus.Fields{
			"topic": topic,
			"error": err,
		}).Error("Failed to subscribe to topic asynchronously")
		return err
	}
	eb.logger.WithField("topic", topic).Debug("Subscribed to topic asynchronously")
	return nil
}

// Publish publishes an event to a topic
func (eb *EventBus) Publish(topic string, args ...interface{}) error {
	eb.bus.Publish(topic, args...)
	eb.logger.WithField("topic", topic).Debug("Published to topic")
	return nil
}

// Unsubscribe unsubscribes from a topic
func (eb *EventBus) Unsubscribe(topic string, handler interface{}) error {
	err := eb.bus.Unsubscribe(topic, handler)
	if err != nil {
		eb.logger.WithFields(logrus.Fields{
			"topic": topic,
			"error": err,
		}).Error("Failed to unsubscribe from topic")
		return err
	}
	eb.logger.WithField("topic", topic).Debug("Unsubscribed from topic")
	return nil
}

// Close closes the event bus
func (eb *EventBus) Close() {
	eb.bus.WaitAsync()
	eb.logger.Info("Event bus closed")
}
