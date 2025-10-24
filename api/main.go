package handler

import (
	"net/http"

	"lumora/auth"
	"lumora/chat"
	"lumora/config"
	"lumora/logger"
	"lumora/mongo"
	"lumora/server"
	"lumora/user"
)

var handler http.Handler

// init runs automatically when the package is loaded
func init() {
	// Load config and logger
	cfg := config.NewConfig()
	log, err := logger.NewLogger()
	if err != nil {
		panic("Failed to create NewLogger client: " + err.Error())
	}

	// Mongo client
	mongoClient, err := mongo.NewMongoClient(nil, cfg, log)
	if err != nil {
		panic("Failed to create mongo client: " + err.Error())
	}
	// User and Auth handlers
	userRepo := user.NewUserRepository(mongoClient)
	userService := user.NewUserService(userRepo)
	userHandler := user.NewUserHandler(userService, log)
	googleAuth := auth.NewGoogleAuth(userService, log)

	// Chat handler
	chatRepo := chat.NewChatRepository(mongoClient)
	chatService := chat.NewChatService(chatRepo)
	chatHandler := chat.NewChatHandler(chatService, log)

	// Create mux and register routes
	mux := server.NewMux(userHandler, googleAuth, chatHandler)
	chatHandler.RegisterRoutes(mux)

	// Assign to global handler
	handler = mux
}

// Handler is what Vercel calls for each HTTP request
func Handler(w http.ResponseWriter, r *http.Request) {
	if handler == nil {
		http.Error(w, "Server not initialized", http.StatusInternalServerError)
		return
	}
	handler.ServeHTTP(w, r)
}
