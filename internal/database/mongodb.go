package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"nutrient_be/internal/config"
	"nutrient_be/internal/pkg/logger"
)

// MongoDB represents a MongoDB connection
type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
	logger   logger.Logger
}

// NewMongoDB creates a new MongoDB connection
func NewMongoDB(cfg *config.DatabaseConfig, log logger.Logger) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ConnectTimeout*time.Second)
	defer cancel()

	// Set client options
	clientOptions := options.Client().
		ApplyURI(cfg.URI).
		SetMaxPoolSize(cfg.MaxPoolSize).
		SetMinPoolSize(cfg.MinPoolSize)

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the database to verify connection
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	database := client.Database(cfg.Database)

	log.InfoLegacy("Successfully connected to MongoDB",
		logger.String("database", cfg.Database),
		logger.String("uri", cfg.URI))

	return &MongoDB{
		Client:   client,
		Database: database,
		logger:   log,
	}, nil
}

// Close closes the MongoDB connection
func (m *MongoDB) Close(ctx context.Context) error {
	if err := m.Client.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to disconnect from MongoDB: %w", err)
	}
	m.logger.InfoLegacy("MongoDB connection closed")
	return nil
}

// Ping checks if the database is accessible
func (m *MongoDB) Ping(ctx context.Context) error {
	return m.Client.Ping(ctx, readpref.Primary())
}

// GetCollection returns a MongoDB collection
func (m *MongoDB) GetCollection(name string) *mongo.Collection {
	return m.Database.Collection(name)
}
