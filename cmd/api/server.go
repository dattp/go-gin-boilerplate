package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-gin-boilerplate/internal/app/api/controller"
	"go-gin-boilerplate/internal/app/api/router"
	"go-gin-boilerplate/internal/config"
	"go-gin-boilerplate/internal/database"
	"go-gin-boilerplate/internal/eventbus"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Server handles the HTTP server lifecycle
type Server struct {
	config           *config.Config
	logger           *logrus.Logger
	router           *router.Router
	httpServer       *http.Server
	healthController *controller.HealthController
	redisClient      *database.RedisClient
	mongoDBClient    *database.MongoDBClient
	eventBus         *eventbus.EventBus
}

// New creates a new server instance
func New(
	cfg *config.Config,
	logger *logrus.Logger,
	router *router.Router,
	healthController *controller.HealthController,
	redisClient *database.RedisClient,
	mongoDBClient *database.MongoDBClient,
	eventBus *eventbus.EventBus,
) *Server {
	return &Server{
		config:           cfg,
		logger:           logger,
		router:           router,
		healthController: healthController,
		redisClient:      redisClient,
		mongoDBClient:    mongoDBClient,
		eventBus:         eventBus,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	// Set Gin mode
	if s.config.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Setup router
	engine := s.router.Setup()

	// Create HTTP server
	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%s", s.config.Port),
		Handler: engine,
	}

	// Start server in a goroutine
	go func() {
		s.logger.Infof("Starting server on %s", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	return nil
}

// Stop gracefully shuts down the server
func (s *Server) Stop() error {
	s.logger.Info("Shutting down server...")

	// Give outstanding requests 5 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %v", err)
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

	s.logger.Info("Server exiting")
	return nil
}

// Run starts the server and waits for shutdown signals
func (s *Server) Run() error {
	if err := s.Start(); err != nil {
		return err
	}

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	return s.Stop()
}
