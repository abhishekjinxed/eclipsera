package mongo

import (
	"context"
	"time"

	"lumora/internal/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewMongoClient(lc fx.Lifecycle, cfg *config.Config, logger *zap.Logger) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(cfg.URI)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Connecting to MongoDB...", zap.String("uri", cfg.URI))
			ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()
			if err := client.Connect(ctx); err != nil {
				return err
			}
			if err := client.Ping(ctx, nil); err != nil {
				return err
			}
			logger.Info("Connected to MongoDB")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Disconnecting MongoDB...")
			return client.Disconnect(ctx)
		},
	})

	return client, nil
}
