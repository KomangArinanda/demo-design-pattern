# Product Review & Rating API

Demo backend API built with Java 17, Spring Boot, Spring MVC, Spring Data JPA, and H2.

## Run

```bash
mvn spring-boot:run
```

Application URL: `http://localhost:8080`

H2 console: `http://localhost:8080/h2-console`

- JDBC URL: `jdbc:h2:mem:reviewdb`
- Username: `sa`
- Password: empty

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

## Example create review request

```bash
curl -X POST http://localhost:8080/api/v1/products/PROD-001/reviews \
  -H 'Content-Type: application/json' \
  -d '{
    "customerId": "CUST-001",
    "orderId": "ORD-001",
    "rating": 5,
    "comment": "Good product, fast delivery"
  }'
```

## Load Test

Run a 1000 requests/second k6 test against product review summary:

```bash
k6 run load-test/review-summary-1000-rps.js
```

The script runs in two phases:

```text
10s warm-up at 100 req/s
10s measured load at 1000 req/s
```

Optional overrides:

```bash
BASE_URL=http://localhost:7080 PRODUCT_ID=PROD-001 k6 run load-test/review-summary-1000-rps.js
```
