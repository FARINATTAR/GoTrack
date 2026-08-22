package store

import (
	"fmt"
	"sync"
	"time"

	"github.com/FARINATTAR/GoTrack/internal/models"
)

// OrderStore is a thread-safe in-memory store for orders
type OrderStore struct {
	mu      sync.RWMutex
	orders  map[string]*models.Order
	counter int

	// subscribers for real-time updates
	subMu       sync.RWMutex
	subscribers map[string][]chan models.Order
}

// NewOrderStore creates a new in-memory order store
func NewOrderStore() *OrderStore {
	return &OrderStore{
		orders:      make(map[string]*models.Order),
		subscribers: make(map[string][]chan models.Order),
	}
}

// Create adds a new order to the store
func (s *OrderStore) Create(req models.CreateOrderRequest) *models.Order {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.counter++
	id := fmt.Sprintf("ORD-%04d", s.counter)

	var total float64
	for _, item := range req.Items {
		total += item.Price * float64(item.Quantity)
	}

	order := &models.Order{
		ID:            id,
		CustomerName:  req.CustomerName,
		Items:         req.Items,
		Status:        models.StatusPlaced,
		TotalAmount:   total,
		DeliveryAddr:  req.DeliveryAddr,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		EstimatedTime: 30, // default 30 mins
	}

	s.orders[id] = order
	return order
}

// Get retrieves an order by ID
func (s *OrderStore) Get(id string) (*models.Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, exists := s.orders[id]
	if !exists {
		return nil, fmt.Errorf("order %s not found", id)
	}
	return order, nil
}

// List returns all orders
func (s *OrderStore) List() []*models.Order {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*models.Order, 0, len(s.orders))
	for _, order := range s.orders {
		result = append(result, order)
	}
	return result
}

// UpdateStatus changes the status of an order with validation
func (s *OrderStore) UpdateStatus(id string, newStatus models.OrderStatus) (*models.Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, exists := s.orders[id]
	if !exists {
		return nil, fmt.Errorf("order %s not found", id)
	}

	if !order.Status.CanTransitionTo(newStatus) {
		return nil, fmt.Errorf("cannot transition from %s to %s", order.Status, newStatus)
	}

	order.Status = newStatus
	order.UpdatedAt = time.Now()

	// notify all subscribers for this order
	s.notifySubscribers(id, *order)

	return order, nil
}

// GetActiveOrders returns all orders that haven't been delivered or cancelled
func (s *OrderStore) GetActiveOrders() []*models.Order {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var active []*models.Order
	for _, order := range s.orders {
		if order.Status != models.StatusDelivered && order.Status != models.StatusCancelled {
			active = append(active, order)
		}
	}
	return active
}

// Subscribe creates a channel that receives updates for a specific order
func (s *OrderStore) Subscribe(orderID string) chan models.Order {
	s.subMu.Lock()
	defer s.subMu.Unlock()

	ch := make(chan models.Order, 10)
	s.subscribers[orderID] = append(s.subscribers[orderID], ch)
	return ch
}

// Unsubscribe removes a subscriber channel for a specific order
func (s *OrderStore) Unsubscribe(orderID string, ch chan models.Order) {
	s.subMu.Lock()
	defer s.subMu.Unlock()

	subs := s.subscribers[orderID]
	for i, sub := range subs {
		if sub == ch {
			s.subscribers[orderID] = append(subs[:i], subs[i+1:]...)
			close(ch)
			break
		}
	}
}

// notifySubscribers sends order updates to all subscribers
func (s *OrderStore) notifySubscribers(orderID string, order models.Order) {
	s.subMu.RLock()
	defer s.subMu.RUnlock()

	for _, ch := range s.subscribers[orderID] {
		select {
		case ch <- order:
		default:
			// skip if subscriber is too slow
		}
	}
}
