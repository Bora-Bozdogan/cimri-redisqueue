# Cimri Redis Queue

This service is an HTTP gateway that validates and enqueues incoming product update events from merchants into Redis queues (high, medium, low priority) based on a computed score, and exposes Prometheus metrics. Built as part of a summer internship at Cimri.

## High-Level Flow
- **HTTP server (Fiber):** starts a Fiber app and exposes endpoints for enqueueing and metrics.
- **Config:** loads Redis and server parameters from config (YAML/env).
- **Redis client:** connects to Redis and pushes requests into priority queues.
- **Services layer:** unpacks requests, determines priority, enqueues into Redis, and updates Prometheus metrics.

## HTTP Endpoints
### POST /enqueue
- Accepts a JSON payload, validates the score, and enqueues the request into one of three Redis lists (high, med, low).

### GET /metrics
- Exposes Prometheus metrics via promhttp from the same Fiber app.

## Architecture
```
Merchant
  -> [Fiber /enqueue]
      -> Body parse & Validate (score between 100–1000)
      -> Select queue: high / med / low
      -> Redis.LPUSH(queue, request)
      -> Update Prometheus metrics
```

- **Main wiring (main.go):** builds the Fiber app, sets up metrics registry, initializes Redis client, services, and handlers.
- **Redis client (worker_service_client.go):** wraps go-redis client with enqueue logic for different priority queues.
- **Service layer (service.go):** delegates to client and updates metrics.
- **Handler (handlers.go):** validates score, routes request to correct queue, increments metrics.

## Service Logic
- **UnpackRequest(body):** parses a QueMessage JSON into Request and Score.
- **Score validation:** valid if 100 <= score <= 1000.
- **Queue selection:**
  - score >= 800 → enqueue to high
  - 500 <= score < 800 → enqueue to med
  - 100 <= score < 500 → enqueue to low

## Metrics (Prometheus gauges)
- `queue_requests_made` — all /enqueue hits.
- `queue_valid_requests_made` — requests successfully enqueued.

## API
### POST /enqueue
**Headers**
```
Content-Type: application/json
```
**Request body (example used in tests)**
```json
{
  "Message": {
    "ApiKey": "amazon-key",
    "ProductName": "iPhone 16",
    "ProductDescription": "Latest Apple flagship phone",
    "ProductImage": "https://example.com/iphone16.jpg",
    "StoreName": "Amazon",
    "Price": 1475,
    "Stock": 50,
    "PopularityScore": 5,
    "UrgencyScore": 5
  },
  "Score": 650
}
```

**Validation**
- Score must be between 100–1000, otherwise 400.

**Responses**
- `200 OK` — "Enqueued request score {score} on queue {queue}".
- `400 Bad Request` — invalid/missing score or enqueue error.

### GET /metrics
- Prometheus metrics endpoint exposed via the same Fiber app using promhttp.

## Development
### Requirements
- Go toolchain
- Redis (for running the app locally)

### Run
Set up `internal/config/config.yaml` (or environment variables) with Redis and server parameters.  
Start the app:
```
go run ./...
```
`main.go` initializes Redis client and services; the app listens on `server_params.listen_port`.

### Test
The test suite wires the same Fiber routes (`/enqueue`, `/metrics`) and mocks the Redis client to simulate enqueueing. It covers scenarios for invalid scores and each queue branch.
```
go test ./...
```

## Notable Files
- `main.go` — app wiring, metrics, Redis client, routes.
- `internal/handlers/handlers.go` — HTTP handler for /enqueue.
- `internal/service/service.go` — delegation to Redis client + metrics updates.
- `internal/client/worker_service_client.go` — Redis LPUSH logic.
- `internal/metrics/metrics.go` — Prometheus gauges.
- `internal/models/request.go` — product update request model.
- `internal/models/que_message.go` — wrapper with score.
- `*_test.go` — handler and flow tests.

## Notes / Future Work
- Metrics could be refactored to counters instead of gauges for monotonic tracking.

## Why It Mattered at Cimri
This gateway centralized scoring-based prioritization of product updates before workers processed them. It protected the system by discarding invalid inputs, ordered work by importance, and provided observability through metrics.

## Tech
Built with Fiber, go-redis, Prometheus client, and Go testing.
