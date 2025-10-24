package server

import (
	"net/http"

	"lumora/internal/auth"
	"lumora/internal/chat"
	"lumora/internal/user"
)

// NewMux builds the HTTP request multiplexer (router)
func NewMux(
	userHandler *user.Handler,
	googleAuth *auth.GoogleAuth,
	chatHandler *chat.Handler,
) *http.ServeMux {
	mux := http.NewServeMux()

	// ✅ Register routes for user, auth, and chat modules
	userHandler.RegisterRoutes(mux)
	googleAuth.RegisterRoutes(mux)
	chatHandler.RegisterRoutes(mux)

	// ✅ Health check
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	return mux
}
