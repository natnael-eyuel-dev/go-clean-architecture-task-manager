# Go Gin Task Manager API (Clean Architecture)

A clean-architecture iteration of the Task Manager API focused on maintainability, testability, and scalability through strict dependency boundaries.

## Phase Position

- **A2SV Go Phase:** Task 7 (Architecture Maturity)
- **Previous Project:** `go-gin-task-manager-api-jwt-rbac`
- **Next Project:** `go-gin-task-manager-api-testing`
- **Program Milestone:** Architecture refactor before dedicated testing and CI hardening

## Features

- Clean Architecture layering
- JWT authentication and authorization
- Role-based access control
- Task CRUD operations with MongoDB persistence
- API documentation in `docs/`

## Tech Stack

- Go
- Gin
- MongoDB
- JWT

## Project Structure

```text
.
├── Delivery/        # HTTP handlers and route entry points
├── Domain/          # Entities and interface contracts
├── Infrastructure/  # Cross-cutting concerns (JWT, hashing, config)
├── Repositories/    # Concrete data implementations
├── Usecases/        # Business rules and orchestration
├── docs/
├── .gitignore
├── go.mod
├── go.sum
└── Readme.md
```

## Run

```bash
go mod tidy
go run Delivery/main.go
```

## Documentation

- API documentation: `docs/api_documentation.md`

## Key Design Decisions

- **Dependency Rule:** Outer layers depend inward via domain interfaces
- **Use Case Validation:** Business constraints enforced in use case layer
- **Framework Isolation:** Delivery details separated from core domain logic

## Learning Outcomes

- Applying clean architecture in real backend services
- Isolating business logic from transport and infra concerns
- Building a production-ready API foundation in Go
