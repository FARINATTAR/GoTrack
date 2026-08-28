package main

import (
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/FARINATTAR/GoTrack/internal/handlers"
    "github.com/FARINATTAR/GoTrack/internal/store"
)

func main() {
    orderStore := store.NewOrderStore()
    mux := http.NewServeMux()
    orderHandler := handlers.NewOrderHandler(orderStore)

    mux.HandleFunc("POST /api/orders", orderHandler.CreateOrder)
    mux.HandleFunc("GET /api/orders/{id}", orderHandler.GetOrder)
    mux.HandleFunc("GET /api/orders", orderHandler.ListOrders)
    mux.HandleFunc("PUT /api/orders/{id}/status", orderHandler.UpdateStatus)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    fmt.Printf("GoTrack server running on http://localhost:%s\n", port)
    log.Fatal(http.ListenAndServe(":"+port, mux))
}
