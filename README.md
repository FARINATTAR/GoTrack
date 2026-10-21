# 🚀 GoTrack

> **High-Performance Real-Time Order Tracking & Concurrency Engine in Go**  
> Built with Go 1.22+, WebSockets, Channels, and Goroutines — simulating production-scale logistics systems like Swiggy and Zomato.

[![Go Version](https://img.shields.io/badge/Go-1.22%2B-00ADD8?style=flat&logo=go)](https://golang.org)
[![WebSocket](https://img.shields.io/badge/WebSocket-Gorilla-orange?style=flat)](https://github.com/gorilla/websocket)
[![Architecture](https://img.shields.io/badge/Architecture-Clean%20Standard%20Layout-blue)](https://github.com/golang-standards/project-layout)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

---

## 📌 Overview

**GoTrack** is a lightweight, concurrent backend microservice designed for low-latency food delivery and logistics order lifecycle management. It leverages Go's native concurrency primitives (goroutines, channels, and `sync.RWMutex`) alongside duplex WebSockets to deliver instantaneous state updates to connected clients without high polling overhead.

It includes an automated background **Lifecycle Simulation Engine** that drives orders through real-world fulfillment stages (`placed` → `confirmed` → `preparing` → `out_for_delivery` → `delivered`) with non-blocking fan-out notifications.

---

## ⚡ Core Features

- **Blazing Fast REST API**: Uses Go 1.22's enhanced `net/http` router with path-value extraction and zero third-party routing overhead.
- **Real-Time WebSocket Pub/Sub**: Full-duplex connection per order tracking session. Supports dynamic broadcast subscriber registration and graceful teardown on client disconnect.
- **Background Goroutine Simulator**: Autonomous ticker-driven worker with stochastic state transitions simulating driver allocation and kitchen preparation.
- **Strict State Machine**: Enforces valid business rules (e.g., an order cannot transition from `delivered` to `preparing`).
- **Thread-Safe In-Memory Store**: Optimized read-write locks (`sync.RWMutex`) supporting concurrent reads and safe atomic mutations.
- **Modular Production Layout**: Follows standard Go idioms (`cmd/server`, `internal/handlers`, `internal/store`, `internal/models`, `internal/simulator`, `internal/middleware`).

---

## 🏗️ Architecture & Data Flow

```text
               +---------------------------------------------------+
               |               Client Applications                 |
               |       (Mobile App / Web Dashboard / Driver)       |
               +---------------------------------------------------+
                               |                   |
                        REST / HTTP           WebSocket
                      (CRUD & Admin)        (/ws/.../track)
                               |                   |
                               v                   v
               +---------------------------------------------------+
               |                 Middleware Layer                  |
               |           (Structured Logger & CORS)              |
               +---------------------------------------------------+
                               |                   |
                               v                   v
                     [ Order Handler ]     [ WebSocket Handler ]
                               |                   |
                               +---------+---------+
                                         |
                                         v
                         +-------------------------------+
                         |   Thread-Safe In-Memory Store |
                         |   - sync.RWMutex              |
                         |   - State Machine Validator   |
                         |   - Subscriber Channel Map    |
                         +-------------------------------+
                                         ^
                                         | Tick / Events
                         +-------------------------------+
                         |  Background Lifecycle Engine  |
                         |  (Goroutines + time.Ticker)   |
                         +-------------------------------+
```

---

## 🔄 Order Lifecycle State Machine

```text
[ Placed ] ------> [ Confirmed ] ------> [ Preparing ] ------> [ Out for Delivery ] ------> [ Delivered ]
     |                   |                     |
     +-------------------+---------------------+------> [ Cancelled ]
```

---

## 📡 API Reference

### Base URL: `http://localhost:8080`

### 1. Create a New Order
- **Endpoint**: `POST /api/orders`
- **Request Body**:
  ```json
  {
    "customer_name": "Farin Attar",
    "delivery_address": "Kothrud, Pune, Maharashtra",
    "items": [
      { "name": "Butter Chicken", "quantity": 1, "price": 380.00 },
      { "name": "Garlic Naan", "quantity": 2, "price": 60.00 }
    ]
  }
  ```
- **Response** (`201 Created`):
  ```json
  {
    "id": "ord_8f1b2c4e",
    "customer_name": "Farin Attar",
    "delivery_address": "Kothrud, Pune, Maharashtra",
    "items": [
      { "name": "Butter Chicken", "quantity": 1, "price": 380.00 },
      { "name": "Garlic Naan", "quantity": 2, "price": 60.00 }
    ],
    "total_amount": 500.00,
    "status": "placed",
    "estimated_time_mins": 35,
    "created_at": "2026-10-23T12:00:00Z",
    "updated_at": "2026-10-23T12:00:00Z"
  }
  ```

---

### 2. Fetch Order Details
- **Endpoint**: `GET /api/orders/{id}`
- **Response** (`200 OK`):
  ```json
  {
    "id": "ord_8f1b2c4e",
    "status": "preparing",
    "total_amount": 500.00,
    "estimated_time_mins": 35
  }
  ```

---

### 3. List All Orders
- **Endpoint**: `GET /api/orders`
- **Response** (`200 OK`): Array of all recorded orders.

---

### 4. Manually Update Status
- **Endpoint**: `PUT /api/orders/{id}/status`
- **Request Body**:
  ```json
  {
    "status": "out_for_delivery"
  }
  ```
- **Response** (`200 OK`): Updated order entity.

---

### 5. Real-Time Tracking (WebSocket)
- **Endpoint**: `ws://localhost:8080/ws/orders/{id}/track`
- **Event Flow**:
  1. **On Connection**: Client immediately receives the current state:
     ```json
     {
       "type": "current_status",
       "order": { "id": "ord_8f1b2c4e", "status": "placed", ... }
     }
     ```
  2. **On Background Progression**: Channel broadcasts status transition events:
     ```json
     {
       "type": "status_update",
       "order": { "id": "ord_8f1b2c4e", "status": "confirmed", ... }
     }
     ```

---

## 🛠️ Project Structure

```bash
GoTrack/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point, router, and graceful startup
├── internal/
│   ├── handlers/
│   │   ├── order_handler.go     # REST controller for CRUD & state manipulation
│   │   └── ws_handler.go        # Gorilla WebSocket upgrader & broadcast consumer
│   ├── middleware/
│   │   └── middleware.go        # HTTP request logging and CORS headers
│   ├── models/
│   │   └── order.go             # Domain entities, validation rules, state machine
│   ├── simulator/
│   │   └── simulator.go         # Async order lifecycle progression via Goroutines
│   └── store/
│       └── order_store.go       # Concurrent sync.RWMutex store & pub/sub channel registry
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

---

## 🚀 Getting Started

### Prerequisites
- [Go 1.22+](https://go.dev/dl/) installed on your machine.
- Git.

### 1. Clone & Setup
```bash
git clone https://github.com/FARINATTAR/GoTrack.git
cd GoTrack
go mod download
```

### 2. Run the Server
```bash
go run cmd/server/main.go
```
The server will boot on port `8080`:
```text
🚀 GoTrack server running on http://localhost:8080
order lifecycle simulator started
```

---

## 🧪 Testing the Service

### Quick Test via cURL
```bash
# 1. Place an order
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{"customer_name":"Farin","delivery_address":"Pune","items":[{"name":"Burger","quantity":2,"price":150}]}'

# 2. Check all orders
curl http://localhost:8080/api/orders
```

### Live WebSocket Tracking in Browser Console
Open any browser DevTools console (F12) and run:
```javascript
const ws = new WebSocket("ws://localhost:8080/ws/orders/YOUR_ORDER_ID/track");

ws.onopen = () => console.log("🟢 Connected to live tracker");
ws.onmessage = (event) => console.log("📦 Update received:", JSON.parse(event.data));
ws.onclose = () => console.log("🔴 Tracker connection closed");
```

---

## 🎯 Concurrency & Design Decisions

1. **Why Go Channels for Pub/Sub?**  
   Avoided heavy external dependencies (Redis / Kafka) for core event dissemination. Go's native buffered channels allow zero-allocation, sub-millisecond in-process broadcast across goroutines.
2. **Preventing Goroutine Leaks:**  
   Subscriber channels are monitored with `select` blocks and explicitly pruned via `defer store.Unsubscribe(...)` upon WebSocket disconnects or client terminations.
3. **Data Integrity:**  
   Enforced strict `sync.RWMutex` locks separating write contention from high-throughput read operations.

---

## 👤 Author

**Farin Attar**  
- GitHub: [@FARINATTAR](https://github.com/FARINATTAR)  
- Portfolio / LinkedIn: [linkedin.com/in/farinattar](https://linkedin.com/in/farinattar)

---

## 📄 License
This project is licensed under the MIT License — see the [LICENSE](LICENSE) file for details.
