package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/FARINATTAR/GoTrack/internal/handlers"
	"github.com/FARINATTAR/GoTrack/internal/middleware"
	"github.com/FARINATTAR/GoTrack/internal/simulator"
	"github.com/FARINATTAR/GoTrack/internal/store"
)

func main() {
	// initialize the in-memory order store
	orderStore := store.NewOrderStore()

	// start the order lifecycle simulator in background
	sim := simulator.NewSimulator(orderStore)
	go sim.Run()

	// setup routes
	mux := http.NewServeMux()

	orderHandler := handlers.NewOrderHandler(orderStore)
	wsHandler := handlers.NewWebSocketHandler(orderStore)

	mux.HandleFunc("POST /api/orders", orderHandler.CreateOrder)
	mux.HandleFunc("GET /api/orders/{id}", orderHandler.GetOrder)
	mux.HandleFunc("GET /api/orders", orderHandler.ListOrders)
	mux.HandleFunc("PUT /api/orders/{id}/status", orderHandler.UpdateStatus)
	mux.HandleFunc("GET /ws/orders/{id}/track", wsHandler.HandleTrack)

	// wrap with middleware
	handler := middleware.Logger(middleware.CORS(mux))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 GoTrack server running on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
