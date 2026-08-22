package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/FARINATTAR/GoTrack/internal/models"
	"github.com/FARINATTAR/GoTrack/internal/store"
)

// OrderHandler handles HTTP requests for order operations
type OrderHandler struct {
	store *store.OrderStore
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(s *store.OrderStore) *OrderHandler {
	return &OrderHandler{store: s}
}

// CreateOrder handles POST /api/orders
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req models.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// basic validation
	if req.CustomerName == "" {
		writeError(w, http.StatusBadRequest, "customer_name is required")
		return
	}
	if len(req.Items) == 0 {
		writeError(w, http.StatusBadRequest, "at least one item is required")
		return
	}
	if req.DeliveryAddr == "" {
		writeError(w, http.StatusBadRequest, "delivery_address is required")
		return
	}

	order := h.store.Create(req)
	writeJSON(w, http.StatusCreated, order)
}

// GetOrder handles GET /api/orders/{id}
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	order, err := h.store.Get(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, order)
}

// ListOrders handles GET /api/orders
func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	orders := h.store.List()
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"orders": orders,
		"count":  len(orders),
	})
}

// UpdateStatus handles PUT /api/orders/{id}/status
func (h *OrderHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var req models.UpdateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	order, err := h.store.UpdateStatus(id, req.Status)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, order)
}

// helper to write JSON responses
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// helper to write error responses
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
