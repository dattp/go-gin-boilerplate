//go:build wireinject
// +build wireinject

package main

import (
	"go-gin-boilerplate/internal/app/worker/handler"
	"go-gin-boilerplate/internal/config"
	"go-gin-boilerplate/internal/database"
	"go-gin-boilerplate/internal/eventbus"
	"go-gin-boilerplate/internal/logger"
	"go-gin-boilerplate/internal/service"

	"github.com/google/wire"
)

// InitializeWorker creates all dependencies for the worker
func InitializeWorker() (*Server, error) {
	wire.Build(
		config.LoadConfig,
		logger.GetLogger,
		eventbus.NewEventBus,
		database.NewRedisClient,
		database.NewMongoDBClient,
		handler.NewNoteHandler,
		service.NewWorkerService,
		New,
	)
	return nil, nil
}
