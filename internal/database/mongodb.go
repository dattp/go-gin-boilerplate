package database

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go-gin-boilerplate/internal/config"
)

// MongoDBClient wraps the MongoDB client
type MongoDBClient struct {
	client *mongo.Client
	logger *logrus.Logger
}

// NewMongoDBClient creates a new MongoDB client
func NewMongoDBClient(logger *logrus.Logger, cfg *config.Config) (*MongoDBClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoDB.URI))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Test connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	logger.Info("Successfully connected to MongoDB")
	return &MongoDBClient{
		client: client,
		logger: logger,
	}, nil
}

// Close closes the MongoDB connection
func (m *MongoDBClient) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := m.client.Disconnect(ctx); err != nil {
		m.logger.WithError(err).Error("Failed to close MongoDB connection")
		return err
	}
	m.logger.Info("MongoDB connection closed")
	return nil
}

// GetMongoDBClient returns the underlying MongoDB client
func (m *MongoDBClient) GetMongoDBClient() *mongo.Client {
	return m.client
} 