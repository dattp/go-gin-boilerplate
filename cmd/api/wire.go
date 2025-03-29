//go:build wireinject
// +build wireinject

package main

import (
	"go-gin-boilerplate/internal/app/api/controller"
	"go-gin-boilerplate/internal/app/api/router"
	"go-gin-boilerplate/internal/config"
	"go-gin-boilerplate/internal/database"
	"go-gin-boilerplate/internal/eventbus"
	"go-gin-boilerplate/internal/logger"
	"go-gin-boilerplate/internal/repository"
	"go-gin-boilerplate/internal/service"

	"github.com/google/wire"
)

// InitializeAPI creates all dependencies for the API server
func InitializeAPI() (*Server, error) {
	wire.Build(
		config.LoadConfig,
		logger.GetLogger,
		eventbus.NewEventBus,
		database.NewRedisClient,
		database.NewMongoDBClient,
		service.NewHealthService,
		controller.NewHealthController,
		repository.NewNoteRepository,
		service.NewAsynqClient,
		service.NewNoteService,
		controller.NewNoteController,
		router.New,
		New,
	)
	return nil, nil
}
