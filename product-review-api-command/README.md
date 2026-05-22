# Product Review & Rating API - Command Pattern Version

Demo backend API for product reviews and ratings, refactored from a service-heavy MVC flow into a command-based architecture.

## Tech Stack

- Java 21
- Spring Boot 3.5.14
- Spring MVC
- Spring Data JPA
- H2 in-memory database
- Lombok
- JUnit 5
- Mockito
- Virtual threads enabled

## Architecture

```text
Controller
  -> CommandExecutor
      -> Command
          -> Repository
          -> Client
          -> Shared ProductReviewService helpers
Entity
DTO
Exception Handler
```

The controller depends only on `CommandExecutor`. Each use case has its own command and command request object.

## Command Layout

```text
com.example.review.command
  ├── Command.java
  ├── CommandExecutor.java
  ├── SaveProductReviewCommand.java
  ├── GetProductReviewSummaryCommand.java
  ├── GetSellerReviewAnalyticsCommand.java
  ├── GetListProductReviewsCommand.java
  ├── GetListCustomerReviewsCommand.java
  ├── GetListRecentReviewsCommand.java
  └── GetDailyProductReviewAnalyticsCommand.java

com.example.review.command.request
  ├── SaveProductReviewCommandRequest.java
  ├── GetProductReviewSummaryCommandRequest.java
  ├── GetSellerReviewAnalyticsCommandRequest.java
  ├── GetListProductReviewsCommandRequest.java
  ├── GetListCustomerReviewsCommandRequest.java
  ├── GetListRecentReviewsCommandRequest.java
  └── GetDailyProductReviewAnalyticsCommandRequest.java
```

## Runtime Configuration

Virtual threads are enabled in `application.yml`:

```yaml
spring:
  threads:
    virtual:
      enabled: true
  main:
    keep-alive: true
```

## Run

Use Java 21 when running the project:

```bash
JAVA_HOME=$(/usr/libexec/java_home -v 21) PATH="$JAVA_HOME/bin:$PATH" mvn spring-boot:run
```

Application URL:

```text
http://localhost:8080
```

H2 console:

```text
http://localhost:8080/h2-console
```

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

## Example Create Review Request

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

## Tests

Run the command-level unit tests with Java 21:

```bash
JAVA_HOME=$(/usr/libexec/java_home -v 21) PATH="$JAVA_HOME/bin:$PATH" mvn test
```

Current test coverage is organized per command class under:

```text
src/test/java/com/example/review/command
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
BASE_URL=http://localhost:7081 PRODUCT_ID=PROD-001 k6 run load-test/review-summary-1000-rps.js
```
