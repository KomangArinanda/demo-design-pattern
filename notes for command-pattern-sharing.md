# Command Pattern Sharing Session

Demo code base -> https://github.com/KomangArinanda/demo-design-pattern
May 2026
By Komang Arinanda

## Slide 1 — Problem First: Why Command Pattern?

### Context
- Backend projects often start simple with controller → service → repository
- As features grow, service classes/functions can become large
- Business logic from different use cases can become mixed together
- Multiple developers may often edit the same service file

### Goal of this session
- Show how Command Pattern can make business actions more explicit
- Compare normal service layer with command/usecase style
- See how the idea can be applied in Java Spring Boot and Go

### Not covered
- Creating a new base project from zero
- Full clean architecture implementation
- Full production-ready boilerplate

---

## Slide 2 — What Is Command Pattern?

### Definition
- Command Pattern wraps a request/action into a dedicated command
- Each command represents one business use case
- The command owns the flow needed to complete that action

### Pattern category
- Command Pattern is a **Behavioral Design Pattern**
- It focuses on organizing behavior and execution flow

### Example command names
- `CreateReviewCommand`
- `ApproveBudgetCommand`
- `CancelOrderCommand`
- `GenerateReportCommand`

### Simple idea
- Instead of one big service handling many actions
- We create smaller command/usecase units

---

## Slide 3 — Service Layer vs Command Pattern

### Traditional service layer
- Controller calls service
- Service contains many business methods
- Service may call repository, client, validator, mapper, etc.
- Good for simple flows
- Can become crowded when features grow

### Command pattern style
- Controller calls command/usecase
- One command handles one action
- Dependencies are injected into that command
- Business flow becomes easier to locate
- Each use case can be tested separately

### Rule of thumb
- Simple CRUD: service layer may be enough
- Complex business flow: command/usecase can help

---

## Slide 4 — SOLID, Unit Test, and Team Collaboration

### Related SOLID principles

**Single Responsibility Principle**
- One command has one clear reason to change
- Business logic is grouped by use case

**Open/Closed Principle**
- New behavior can be added with a new command
- Less need to modify existing large service classes

**Dependency Inversion Principle**
- Command can depend on interfaces
- Repository/client can be mocked more easily

### Impact on unit testing
- Test one command/usecase at a time
- Mock only dependencies needed by that command
- Test scenarios become smaller and more focused

### Impact on team collaboration
- Different features can live in different command files
- Less chance multiple developers edit the same large service
- Can reduce merge conflict risk

---

## Slide 5 — Demo Plan & Final Message

### Demo 1: Java Spring Boot — Existing Service Layer
- Show controller → service → repository/client
- Identify where the business logic is located
- Show the pain point when service grows

### Demo 2: Java Spring Boot — Command Pattern
- Move one service method into a command
- Show command input, execute method, and output
- Compare with the previous service implementation

### Demo 3: Go — Usecase / Command Style
- Show the same idea in Go
- Use struct/function with `Execute(ctx, input)`
- Inject repository/client dependency through constructor

### Final message
- Command Pattern is not mandatory for every feature
- It is useful when business logic becomes complex
- It improves separation, testability, and collaboration
- Use it consistently, not randomly
