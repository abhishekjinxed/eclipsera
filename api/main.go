package main

import (
	"context"
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

func main() {
	app := fx.New(
		fx.Provide(
			config.NewConfig,
			logger.NewLogger,
			mongo.NewMongoClient,
			server.NewMux,
			server.NewServer,
			user.NewUserRepository,
			user.NewUserService,
			user.NewUserHandler,
			auth.NewGoogleAuth, // ✅ Make sure this is here
			chat.NewChatRepository,
			chat.NewChatService,
		),
		fx.Invoke(func(*server.Server) {}, chat.NewChatHandler),
	)

	startCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := app.Start(startCtx); err != nil {
		panic(err)
	}

	<-app.Done()

	stopCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := app.Stop(stopCtx); err != nil {
		panic(err)
	}
}
