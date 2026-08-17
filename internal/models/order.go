package models

import (
	"time"
)

// OrderStatus represents the current state of an order
type OrderStatus string

const (
	StatusPlaced         OrderStatus = "placed"
	StatusConfirmed      OrderStatus = "confirmed"
	StatusPreparing      OrderStatus = "preparing"
	StatusOutForDelivery OrderStatus = "out_for_delivery"
	StatusDelivered      OrderStatus = "delivered"
	StatusCancelled      OrderStatus = "cancelled"
)

// Order represents a food delivery order
type Order struct {
	ID            string      `json:"id"`
	CustomerName  string      `json:"customer_name"`
	Items         []OrderItem `json:"items"`
	Status        OrderStatus `json:"status"`
	TotalAmount   float64     `json:"total_amount"`
	DeliveryAddr  string      `json:"delivery_address"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
	EstimatedTime int         `json:"estimated_time_mins"`
}

// OrderItem represents a single item in the order
type OrderItem struct {
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

// CreateOrderRequest is the payload for creating a new order
type CreateOrderRequest struct {
	CustomerName string      `json:"customer_name"`
	Items        []OrderItem `json:"items"`
	DeliveryAddr string      `json:"delivery_address"`
}

// UpdateStatusRequest is the payload for updating order status
type UpdateStatusRequest struct {
	Status OrderStatus `json:"status"`
}

// StatusTransitions defines valid status transitions
var StatusTransitions = map[OrderStatus][]OrderStatus{
	StatusPlaced:         {StatusConfirmed, StatusCancelled},
	StatusConfirmed:      {StatusPreparing, StatusCancelled},
	StatusPreparing:      {StatusOutForDelivery, StatusCancelled},
	StatusOutForDelivery: {StatusDelivered},
	StatusDelivered:      {},
	StatusCancelled:      {},
}

// CanTransitionTo checks if moving from current status to the next one is valid
func (s OrderStatus) CanTransitionTo(next OrderStatus) bool {
	allowed, exists := StatusTransitions[s]
	if !exists {
		return false
	}
	for _, status := range allowed {
		if status == next {
			return true
		}
	}
	return false
}

// NextStatus returns the next status in the order lifecycle
func (s OrderStatus) NextStatus() (OrderStatus, bool) {
	switch s {
	case StatusPlaced:
		return StatusConfirmed, true
	case StatusConfirmed:
		return StatusPreparing, true
	case StatusPreparing:
		return StatusOutForDelivery, true
	case StatusOutForDelivery:
		return StatusDelivered, true
	default:
		return "", false
	}
}
