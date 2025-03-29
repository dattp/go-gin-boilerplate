package main

import (
	"log"

	"go-gin-boilerplate/internal/config"
	"go-gin-boilerplate/internal/logger"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Set up logger
	logger.SetLevel(cfg.LogLevel)
	logger.SetFormat(cfg.LogJSON)

	// Initialize dependencies using Wire
	server, err := InitializeAPI()
	if err != nil {
		log.Fatalf("Failed to initialize dependencies: %v", err)
	}

	// Run the server
	if err := server.Run(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
