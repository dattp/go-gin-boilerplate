package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go-gin-boilerplate/internal/config"
	"go-gin-boilerplate/internal/database"
	"go-gin-boilerplate/internal/eventbus"
	"go-gin-boilerplate/internal/service"

	"github.com/sirupsen/logrus"
)

// Server handles the worker server lifecycle
type Server struct {
	config        *config.Config
	logger        *logrus.Logger
	workerService *service.WorkerService
	redisClient   *database.RedisClient
	mongoDBClient *database.MongoDBClient
	eventBus      *eventbus.EventBus
}

// New creates a new worker server instance
func New(
	cfg *config.Config,
	logger *logrus.Logger,
	workerService *service.WorkerService,
	redisClient *database.RedisClient,
	mongoDBClient *database.MongoDBClient,
	eventBus *eventbus.EventBus,
) *Server {
	return &Server{
		config:        cfg,
		logger:        logger,
		workerService: workerService,
		redisClient:   redisClient,
		mongoDBClient: mongoDBClient,
		eventBus:      eventBus,
	}
}

// Start starts the worker service
func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("Starting worker service...")
	return s.workerService.Start(ctx)
}

// Stop gracefully shuts down the worker service
func (s *Server) Stop() error {
	s.logger.Info("Stopping worker service...")

	// Stop the worker service
	if err := s.workerService.Stop(); err != nil {
		s.logger.Error("Failed to stop worker service:", err)
	}

	// Wait for any pending event bus operations
	s.eventBus.Close()

	// Close database connections
	if err := s.redisClient.Close(); err != nil {
		s.logger.Error("Failed to close Redis connection:", err)
	}
	if err := s.mongoDBClient.Close(); err != nil {
		s.logger.Error("Failed to close MongoDB connection:", err)
	}

	s.logger.Info("Worker service exiting")
	return nil
}

// Run starts the worker service and waits for shutdown signals
func (s *Server) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the worker service in a goroutine
	go func() {
		if err := s.Start(ctx); err != nil {
			s.logger.Fatalf("Failed to start worker service: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	return s.Stop()
}
