package service

import (
	"context"

	"go-gin-boilerplate/internal/app/worker/handler"
	"go-gin-boilerplate/internal/config"

	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
)

// WorkerService handles background tasks using asynq
type WorkerService struct {
	logger      *logrus.Logger
	server      *asynq.Server
	inspector   *asynq.Inspector
	noteHandler *handler.NoteHandler
}

// NewWorkerService creates a new worker service
func NewWorkerService(
	logger *logrus.Logger,
	cfg *config.Config,
	noteHandler *handler.NoteHandler,
) *WorkerService {
	// Create Redis connection
	redisOpt := asynq.RedisClientOpt{
		Addr: cfg.Redis.Addr,
	}

	// Create asynq server
	srv := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				config.QueueDefault: 6,
				config.QueueHigh:    3,
				config.QueueLow:     1,
			},
		},
	)

	// Create inspector for monitoring
	inspector := asynq.NewInspector(redisOpt)

	return &WorkerService{
		logger:      logger,
		server:      srv,
		inspector:   inspector,
		noteHandler: noteHandler,
	}
}

// Start starts the worker service
func (s *WorkerService) Start(ctx context.Context) error {
	s.logger.Info("Starting worker service...")

	// Register task handlers
	mux := asynq.NewServeMux()
	mux.HandleFunc(handler.TypeNoteCreated, s.noteHandler.HandleNoteCreatedTask)

	// Start the server
	return s.server.Run(mux)
}

// Stop stops the worker service
func (s *WorkerService) Stop() error {
	s.logger.Info("Stopping worker service...")
	s.server.Shutdown()
	return nil
}

// NewAsynqClient creates a new Asynq client
func NewAsynqClient(cfg *config.Config) *asynq.Client {
	return asynq.NewClient(asynq.RedisClientOpt{
		Addr: cfg.Redis.Addr,
	})
}
