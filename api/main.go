package main

import (
	"context"
	"net/http"
	"time"

	"lumora/internal/auth"
	"lumora/internal/chat"
	"lumora/internal/config"
	"lumora/internal/logger"
	"lumora/internal/mongo"
	"lumora/internal/server"
	"lumora/internal/user"

	"go.uber.org/fx"
)

// Handler is the Vercel serverless entry point
func Handler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var mux *http.ServeMux

	app := fx.New(
		fx.Provide(
			config.NewConfig,
			logger.NewLogger,
			mongo.NewMongoClient,
			server.NewMux,
			user.NewUserRepository,
			user.NewUserService,
			user.NewUserHandler,
			auth.NewGoogleAuth,
			chat.NewChatRepository,
			chat.NewChatService,
		),
		fx.Invoke(
			func(m *http.ServeMux) {
				mux = m
			},
			chat.NewChatHandler,
		),
	)

	startCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := app.Start(startCtx); err != nil {
		http.Error(w, "Failed to start app: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer func() {
		stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = app.Stop(stopCtx)
	}()

	// Serve the HTTP request through the mux
	mux.ServeHTTP(w, r)
}
