package simulator

import (
	"log"
	"math/rand"
	"time"

	"github.com/FARINATTAR/GoTrack/internal/store"
)

// Simulator automatically progresses orders through their lifecycle
// this simulates real-world order processing using goroutines
type Simulator struct {
	store *store.OrderStore
}

// NewSimulator creates a new order lifecycle simulator
func NewSimulator(s *store.OrderStore) *Simulator {
	return &Simulator{store: s}
}

// Run starts the simulator — checks for active orders every 10 seconds
// and randomly progresses them to the next status
func (s *Simulator) Run() {
	log.Println("order lifecycle simulator started")

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		activeOrders := s.store.GetActiveOrders()

		for _, order := range activeOrders {
			// 40% chance to progress each active order on each tick
			// this makes it feel more organic/realistic
			if rand.Float64() < 0.4 {
				nextStatus, ok := order.Status.NextStatus()
				if !ok {
					continue
				}

				// simulate processing in a separate goroutine
				go func(id string, next string) {
					_, err := s.store.UpdateStatus(id, nextStatus)
					if err != nil {
						log.Printf("simulator: failed to update %s: %v", id, err)
						return
					}
					log.Printf("simulator: order %s → %s", id, next)
				}(order.ID, string(nextStatus))
			}
		}
	}
}
