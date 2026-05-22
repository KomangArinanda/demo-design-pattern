# Product Review & Rating API - Go Executor Version

Go implementation of the product review demo API, modeled after the Java command-pattern project with use cases grouped under `internal/usecase`.

## Tech Stack

- Go 1.26.1
- Gorilla Mux
- In-memory repository
- Standard library HTTP server
- Testify
- Shared `Execute(ctx, input any)` executor contract
- Internal `appctx.Response` envelope for success/error handling

## Architecture

```text
Handler
  -> ProductReview
      -> Executor
          -> Repository
          -> Client
          -> Shared ProductReviewService helpers
```

Each executor implements the same `Execute(ctx, input any)` method and returns a shared `appctx.Response`. The HTTP handler keeps the public JSON shape of the API stable by writing success data directly and mapping errors to `{ "message": "..." }`.

## Run

```bash
go run ./cmd
```

Default URL:

```text
http://localhost:7082
```

Optional configuration:

```bash
APP_PORT=7082 DB_SIMULATED_LATENCY_MS=200 go run ./cmd
```

## Test

Run all tests:

```bash
go test ./...
```

If your local Go environment has multiple installed versions, force this project to use Go 1.26.1:

```bash
GOTOOLCHAIN=local GOROOT=/Users/komangarinanda/.goenv/versions/1.26.1 /Users/komangarinanda/.goenv/versions/1.26.1/bin/go test ./...
```

## Endpoints

```http
POST /api/v1/products/{productId}/reviews
GET  /api/v1/products/{productId}/review-summary
GET  /api/v1/sellers/{sellerId}/review-analytics
GET  /api/v1/products/{productId}/reviews
GET  /api/v1/customers/{customerId}/reviews
GET  /api/v1/reviews/recent?limit=10
GET  /api/v1/products/{productId}/review-analytics/daily?month=5&year=2026
```

## Simulated Database Latency

The repository waits before each operation to mimic blocking database time:

```bash
DB_SIMULATED_LATENCY_MS=200
```

## Load Test

```bash
k6 run load-test/review-summary-1000-rps.js
```

Phases:

```text
10s warm-up at 100 req/s
10s measured load at 1000 req/s
```
