package service

import (
	"time"

	"github.com/sirupsen/logrus"
	"go-gin-boilerplate/internal/eventbus"
)

// HealthService handles health check related operations
type HealthService interface {
	Check() map[string]interface{}
}

type healthService struct {
	startTime time.Time
	logger    *logrus.Logger
	eventBus  *eventbus.EventBus
}

func NewHealthService(logger *logrus.Logger, eventBus *eventbus.EventBus) HealthService {
	service := &healthService{
		startTime: time.Now(),
		logger:    logger,
		eventBus:  eventBus,
	}

	// Subscribe to health check events
	service.eventBus.Subscribe("health:check", service.handleHealthCheckEvent)

	return service
}

func (s *healthService) Check() map[string]interface{} {
	// Publish health check event
	s.eventBus.Publish("health:check", time.Now())

	return map[string]interface{}{
		"status":    "ok",
		"uptime":    time.Since(s.startTime).String(),
		"timestamp": time.Now().Unix(),
	}
}

func (s *healthService) handleHealthCheckEvent(timestamp time.Time) {
	s.logger.WithField("timestamp", timestamp).Info("Health check event received")
} 