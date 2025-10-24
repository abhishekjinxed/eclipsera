package handler

import (
	"context"
	"net/http"
	"os"
	"time"

	"lumora/auth"
	"lumora/chat"
	"lumora/config"
	"lumora/logger"
	"lumora/mongo"
	"lumora/server"
	"lumora/user"

	"go.uber.org/fx"
)

var handler http.Handler // shared between Vercel + local

func main() {
	app := fx.New(
		fx.Provide(
			config.NewConfig,
			logger.NewLogger,
			mongo.NewMongoClient,
			server.NewMux, // returns *http.ServeMux
			user.NewUserRepository,
			user.NewUserService,
			user.NewUserHandler,
			auth.NewGoogleAuth,
			chat.NewChatRepository,
			chat.NewChatService,
			chat.NewChatHandler,
		),
		fx.Invoke(func(mux *http.ServeMux) {
			handler = mux
		}),
	)

	// Start the app
	startCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := app.Start(startCtx); err != nil {
		panic(err)
	}

	// 🧩 Detect if running locally
	if os.Getenv("VERCEL") == "" {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		println("✅ Server running locally on http://localhost:" + port)
		if err := http.ListenAndServe(":"+port, handler); err != nil {
			panic(err)
		}
	} else {
		// On Vercel, do nothing — it uses Handler() below
		<-app.Done()
	}

	stopCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := app.Stop(stopCtx); err != nil {
		panic(err)
	}
}

// ✅ This is what Vercel calls for each HTTP request
func Handler(w http.ResponseWriter, r *http.Request) {
	handler.ServeHTTP(w, r)
}
