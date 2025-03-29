package database

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go-gin-boilerplate/internal/config"
)

// RedisClient wraps the Redis client
type RedisClient struct {
	client *redis.Client
	logger *logrus.Logger
}

// NewRedisClient creates a new Redis client
func NewRedisClient(logger *logrus.Logger, cfg *config.Config) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	logger.Info("Successfully connected to Redis")
	return &RedisClient{
		client: client,
		logger: logger,
	}, nil
}

// Close closes the Redis connection
func (r *RedisClient) Close() error {
	if err := r.client.Close(); err != nil {
		r.logger.WithError(err).Error("Failed to close Redis connection")
		return err
	}
	r.logger.Info("Redis connection closed")
	return nil
}

// GetRedisClient returns the underlying Redis client
func (r *RedisClient) GetRedisClient() *redis.Client {
	return r.client
} 