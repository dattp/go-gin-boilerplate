package handler

import (
	"context"
	"time"

	"go-gin-boilerplate/internal/eventbus"

	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
)

// ExampleTaskHandler handles example tasks
type ExampleTaskHandler struct {
	logger   *logrus.Logger
	eventBus *eventbus.EventBus
}

// NewExampleTaskHandler creates a new example task handler
func NewExampleTaskHandler(logger *logrus.Logger, eventBus *eventbus.EventBus) *ExampleTaskHandler {
	return &ExampleTaskHandler{
		logger:   logger,
		eventBus: eventBus,
	}
}

// Handle processes the example task
func (h *ExampleTaskHandler) Handle(ctx context.Context, t *asynq.Task) error {
	// Publish task started event
	h.eventBus.Publish("task:started", t.Payload())

	h.logger.WithField("task_id", t.Payload()).Info("Processing example task")

	// Simulate work
	time.Sleep(2 * time.Second)

	h.logger.WithField("task_id", t.Payload()).Info("Completed example task")

	// Publish task completed event
	h.eventBus.Publish("task:completed", t.Payload())

	return nil
}
