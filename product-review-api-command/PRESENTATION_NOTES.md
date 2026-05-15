# Product Review API Comparison - Presentation Notes

## 1. Presentation Goal

Compare two implementations of the same backend API:

1. Legacy layered MVC version
2. Refactored command-pattern version

The comparison demonstrates:

- how architecture changes affect maintainability
- how Java version and runtime model can be upgraded independently from business behavior
- how to keep API contracts stable while refactoring internals
- how to prepare a fair performance comparison under simulated slow database conditions

---

## 2. Business Domain

Product Review & Rating System with these main capabilities:

- create product review after completed purchase
- reject duplicate reviews
- reject reviews for invalid orders
- return product review summary
- return seller review analytics
- list product reviews
- list customer reviews
- list recent reviews
- return daily review analytics for a selected month and year

Core examples:

- product-level aggregation
- seller-level aggregation
- external client validation
- repository usage
- future migration and refactoring discussion

---

## 3. Project A - Legacy Service-Oriented Version

### Tech Stack

- Java 17
- Spring Boot 2.7
- Spring MVC
- Spring Data JPA
- H2 in-memory database
- Lombok
- JUnit 5
- Mockito

### Architecture

```text
Controller
  -> ProductReviewService
      -> Repository
      -> Client
Entity
DTO
Exception Handler
```

### Characteristics

- one service owns all use cases
- controller calls service methods directly
- readable for small scope
- service grows as endpoints increase
- business orchestration and reusable helpers stay mixed in one class

### Main Public Service Methods

- `createReview`
- `getProductReviewSummary`
- `getSellerReviewAnalytics`
- `getProductReviews`
- `getCustomerReviews`
- `getRecentReviews`
- `getDailyProductReviewAnalytics`

---

## 4. Project B - Command Pattern Version

### Tech Stack

- Java 21
- Spring Boot 3.5
- Spring MVC
- Spring Data JPA
- H2 in-memory database
- Lombok
- JUnit 5
- Mockito
- virtual threads enabled

### Architecture

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

### Command Infrastructure

- `Command<REQUEST, RESPONSE>`
- `CommandExecutor`
- one request object per command
- default `validate()` method in `Command`
- commands override `validate()` only when needed

### Commands Implemented

- `SaveProductReviewCommand`
- `GetProductReviewSummaryCommand`
- `GetSellerReviewAnalyticsCommand`
- `GetListProductReviewsCommand`
- `GetListCustomerReviewsCommand`
- `GetListRecentReviewsCommand`
- `GetDailyProductReviewAnalyticsCommand`

### Characteristics

- each use case has its own class
- controller depends only on `CommandExecutor`
- reusable helpers stay in `ProductReviewService`
- validation and execution flow are explicit
- command names communicate behavior using verbs:
  - `Save...`
  - `Get...`
  - `GetList...`

---

## 5. API Surface Kept Stable

Both projects expose the same endpoints:

```http
POST /api/v1/products/{productId}/reviews
GET  /api/v1/products/{productId}/review-summary
GET  /api/v1/sellers/{sellerId}/review-analytics
GET  /api/v1/products/{productId}/reviews
GET  /api/v1/customers/{customerId}/reviews
GET  /api/v1/reviews/recent?limit=10
GET  /api/v1/products/{productId}/review-analytics/daily?month=5&year=2026
```

Key presentation point:

> Internal architecture changed, but external API behavior remained stable.

---

## 6. Important Functional Additions During Development

### Initial Features

- create review
- duplicate review protection
- invalid-order protection
- product summary aggregation
- seller analytics aggregation

### Added Later

- list reviews by product
- list reviews by customer
- list recent reviews
- daily product analytics by month and year

### Daily Analytics Design

- request params: `month`, `year`
- multiple repository calls allowed
- grouping and summary calculation intentionally done in Java service/command code
- demonstrates mapping logic and aggregation ownership outside SQL

---

## 7. Shared Behaviors Across Both Projects

### Mock External Clients

- `OrderClient`
  - validates order ownership
  - simulates another microservice

- `ProductClient`
  - returns product IDs for a seller
  - supports seller analytics

### Seed Data

- 5 products
- 2 sellers
- 34 seeded product reviews

### Shared Constraints

- rating must be `1..5`
- comment max `500` chars
- unique review per:
  - product
  - customer
  - order

---

## 8. Testing Strategy

### Service-Oriented Project

- service-layer unit tests
- all public service methods covered

### Command-Pattern Project

- tests moved to one test class per command
- command-level orchestration covered
- shared test factory used for fixtures

### Why This Matters

- tests follow the architecture
- refactor changes test ownership, not required behavior
- command tests are more localized by use case

---

## 9. Simulated Slow Database

### Why We Added It

H2 in-memory database is too fast to show the behavior of real I/O-bound workloads.

### What We Added

- `DatabaseLatencySimulator`
- configurable property:

```yaml
app:
  database:
    simulated-latency-ms: 200
```

### What It Simulates

- blocking wait before repository calls
- similar application-thread behavior to waiting on a slow JDBC/database response

### Important Limitation

It simulates wait time, not:

- connection pool contention
- lock contention
- network failures
- database CPU work
- slow query plans

---

## 10. Performance Test Setup

### Tool

- k6

### Target Endpoint

```http
GET /api/v1/products/{productId}/review-summary
```

### Load Profile

```text
30s warm-up at 100 req/s
30s measured load at 1000 req/s
```

### Why Warm-Up Is Needed

Java applications need time for:

- class loading
- JIT compilation
- bean initialization effects
- connection pool warm-up
- first-query effects

Only the steady-state phase should be compared.

---

## 11. Java Runtime Comparison Angle

### Legacy Version

- Java 17
- traditional platform-thread request handling

### Command Version

- Java 21
- virtual threads enabled

### Why This Matters

The workload is intentionally blocking:

- each request waits on simulated database latency
- many concurrent requests are mostly waiting rather than computing

### Expected Discussion Point

- virtual threads do not make a single slow query faster
- virtual threads can reduce the cost of handling many blocking requests concurrently
- architecture refactor and runtime upgrade solve different problems

---

## 12. Architecture Comparison

| Topic | Service-Oriented Version | Command-Pattern Version |
|---|---|---|
| Entry from controller | direct service call | command executor |
| Use-case ownership | one service class | one command per use case |
| Reuse | helper methods inside service | shared helper service + isolated commands |
| Extensibility | service grows wider | add a new command |
| Test organization | per service | per command |
| Readability at small size | very high | moderate |
| Readability at larger size | degrades over time | remains more localized |

---

## 13. Trade-Offs

### Service-Oriented Version

Pros:

- simpler to understand initially
- fewer classes
- faster to build for small systems

Cons:

- service can become a large orchestration hub
- harder to isolate growing use cases
- shared and use-case-specific logic mix together

### Command-Pattern Version

Pros:

- stronger use-case isolation
- easier local reasoning
- cleaner extension path
- better test ownership per use case

Cons:

- more files and indirection
- more ceremony for simple endpoints
- naming discipline becomes important

---

## 14. Key Lessons

1. Keep external API contracts stable during refactoring.
2. Refactor only when the structure starts to justify it.
3. Use commands to isolate use cases, not to remove all shared services.
4. Runtime upgrades and architecture refactors are complementary, not interchangeable.
5. Fair performance testing requires warm-up, repeatability, and identical workloads.
6. Simulated latency is useful for teaching blocking behavior, but it is not a full database benchmark.

---

## 15. Suggested Slide Flow

1. Problem statement
2. Business domain and API scope
3. Original architecture
4. Why the service started to grow
5. Refactored command architecture
6. Side-by-side code structure comparison
7. Testing strategy before and after
8. Simulated slow database design
9. Performance test design with warm-up
10. Java 17 vs Java 21 virtual threads discussion
11. Trade-offs
12. Closing lessons

---

## 16. Recommended Demo Sequence

1. Show same endpoint list in both projects
2. Show original `ProductReviewService`
3. Show command controller using `CommandExecutor`
4. Open one command such as:
   - `SaveProductReviewCommand`
   - `GetDailyProductReviewAnalyticsCommand`
5. Show shared helper service
6. Show command-level tests
7. Show k6 script
8. Run or discuss comparison result

---

## 17. One-Sentence Summary

> We started with a conventional Spring MVC service architecture, expanded the domain until the service began to widen, refactored the same API into command-based use cases, upgraded the runtime from Java 17 to Java 21 with virtual threads, and built a fair load-testing setup to compare both approaches under simulated slow-database conditions.


# Benefit
1. colaborasi antar team lebih mudah minim conflict git
2. lbeih traceable baca code business logic di one place
3. Lebih predictable ketika melakukan code change, hampir tidak impacted ke fitur lain
2. lebih mudah baca / bikin unit test