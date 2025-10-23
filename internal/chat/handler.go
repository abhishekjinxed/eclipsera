package chat

import (
	"net/http"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Handler struct {
	service Service
	hub     *Hub
	logger  *zap.Logger
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func NewChatHandler(mux *http.ServeMux, service Service, logger *zap.Logger) *Handler {
	hub := NewHub()
	h := &Handler{service: service, hub: hub, logger: logger}
	mux.HandleFunc("/ws/chat", h.HandleConnections)
	return h
}

func (h *Handler) HandleConnections(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "Missing user_id", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error("WebSocket upgrade failed", zap.Error(err))
		return
	}

	client := &Client{
		UserID:  userID,
		Conn:    conn,
		Service: h.service,
		Hub:     h.hub,
		Ctx:     r.Context(),
		//logger:  h.logger,
	}

	h.hub.AddClient(userID, client)
	go client.ReadPump()
}
