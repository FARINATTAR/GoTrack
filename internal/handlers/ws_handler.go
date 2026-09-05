package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/FARINATTAR/GoTrack/internal/store"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // allow all origins for now
	},
}

// WebSocketHandler handles websocket connections for live tracking
type WebSocketHandler struct {
	store *store.OrderStore
}

// NewWebSocketHandler creates a new websocket handler
func NewWebSocketHandler(s *store.OrderStore) *WebSocketHandler {
	return &WebSocketHandler{store: s}
}

// HandleTrack handles GET /ws/orders/{id}/track
// clients connect here to get live status updates for an order
func (h *WebSocketHandler) HandleTrack(w http.ResponseWriter, r *http.Request) {
	orderID := r.PathValue("id")

	// check if order exists before upgrading
	order, err := h.store.Get(orderID)
	if err != nil {
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}

	// upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("websocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	// send current status first
	initialMsg, _ := json.Marshal(map[string]interface{}{
		"type":   "current_status",
		"order":  order,
	})
	conn.WriteMessage(websocket.TextMessage, initialMsg)

	log.Printf("client connected to track order %s", orderID)

	// subscribe to updates for this order
	updateCh := h.store.Subscribe(orderID)
	defer h.store.Unsubscribe(orderID, updateCh)

	// also listen for client disconnect
	done := make(chan struct{})
	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				close(done)
				return
			}
		}
	}()

	// push updates to client until they disconnect or order is delivered
	for {
		select {
		case update := <-updateCh:
			msg, _ := json.Marshal(map[string]interface{}{
				"type":  "status_update",
				"order": update,
			})
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Printf("failed to send update: %v", err)
				return
			}
			log.Printf("sent status update for order %s: %s", orderID, update.Status)

		case <-done:
			log.Printf("client disconnected from order %s", orderID)
			return
		}
	}
}
