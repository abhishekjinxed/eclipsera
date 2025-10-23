package server

import (
	"context"
	"net/http"

	"lumora/internal/config"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Server = http.Server

func NewServer(lc fx.Lifecycle, mux *http.ServeMux, cfg *config.Config, logger *zap.Logger) *Server {
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Starting HTTP server", zap.String("port", cfg.Port))
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Fatal("Server failed", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping HTTP server...")
			return srv.Shutdown(ctx)
		},
	})

	return srv
}
