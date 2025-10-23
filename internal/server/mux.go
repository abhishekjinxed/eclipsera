package server

import (
	"net/http"

	"lumora/internal/auth"
	"lumora/internal/user"
)

func NewMux(userHandler *user.Handler, googleAuth *auth.GoogleAuth) *http.ServeMux {
	mux := http.NewServeMux()

	// User and Google auth handlers register themselves
	userHandler.RegisterRoutes(mux)
	googleAuth.RegisterRoutes(mux)
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	return mux
}
