package mongo

import (
	"context"
	"time"

	"lumora/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func NewMongoClient(cfg *config.Config, logger *zap.Logger) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(cfg.URI)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Info("Connecting to MongoDB...", zap.String("uri", cfg.URI))
	if err := client.Connect(ctx); err != nil {
		return nil, err
	}

	// Test the connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	logger.Info("Connected to MongoDB successfully!")
	return client, nil
}

// Call this when shutting down your app
func CloseMongoClient(client *mongo.Client, logger *zap.Logger) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Disconnect(ctx); err != nil {
		logger.Error("Error disconnecting MongoDB", zap.Error(err))
	} else {
		logger.Info("MongoDB connection closed")
	}
}
