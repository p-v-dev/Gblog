# Gblog — Blog Post Backend API

A study project to build a clean, production-style blog backend in Go. The goal is to keep the codebase well-structured and disciplined so it can serve as a solid foundation for learning DevOps and infrastructure tooling around it (Docker, Nginx, GitHub Actions, monitoring, etc.).

The [Scarf CLI](https://github.com/p-v-dev/scarf) was built specifically to enforce consistent module scaffolding and prevent the codebase from drifting into a messy state as it grows.

## Stack

- **Language:** Go 1.26.3
- **Framework:** Gin v1.12.0
- **ORM:** GORM v1.31.1
- **Database:** PostgreSQL
- **Swagger:** swaggo/swag v1.16.6
- **Env loader:** godotenv
- **Scaffolding:** [Scarf](https://github.com/p-v-dev/scarf) CLI (`scarf mod -n <name>`)

## Project Structure

```
Gblog/
├── cmd/
│   ├── main.go              # Entry point: env, DB, use cases, routes
│   └── docs/                # Auto-generated Swagger docs
├── internal/
│   ├── <module>/            # Vertical slice per domain (scarf generates this)
│   │   ├── entity.go        # Domain entity (GORM model + business rules)
│   │   ├── dto.go           # Input/output DTOs for use cases
│   │   ├── repository.go    # Interface + GORM implementation
│   │   ├── usecase.go       # Use case interfaces + implementations
│   │   └── http/
│   │       └── handlers.go  # Gin HTTP handlers (Swagger annotations)
│   └── infra/
│       └── database.go      # DB connection, env loading, auto-migrate
├── pkg/
│   └── blogStatus.go        # Shared domain types/enums
├── .env                     # DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT, API_PORT
└── AGENTS.md
```

## Creating a New Module

Use Scarf to scaffold a new domain module:

```
scarf mod -n <name>
```

This generates `internal/<name>/` with `entity.go`, `dto.go`, `repository.go`, and `usecase.go`. After scaffolding you must manually create `internal/<name>/http/handlers.go` and wire everything in `cmd/main.go`.

## Patterns & Conventions

### Module package naming
The package name matches the directory name (e.g. `package blogPost`, `package user`). The `blogPost` module uses camelCase because the directory name has no underscore — keep the package name consistent with the directory.

### Entity (`entity.go`)
- Struct embeds `gorm.Model` for ID, timestamps, soft delete
- Uses `pkg.BlogStatus` for status fields (`draft`, `published`, `archived`)
- Uses GORM tags on fields (`gorm:"type:varchar(255);not null"`)
- Soft delete via `IsActive bool` field

### DTO (`dto.go`)
- Separate input/output structs per use case
- JSON tags for request binding

### Repository (`repository.go`)
- Define an **interface** first (contract) with `context.Context` in signatures
- Private **implementation** struct (e.g. `blogPostRepositoryORM`)
- **Constructor** returns the interface: `func New<Name>Repository(db *gorm.DB) <Name>Repository`
- Methods: `Create`, `GetBySlug`, `FindByID`, `FetchAll`, `Update`, `Delete`
- `FetchAll` uses `Where("is_active = ?", true)`, `Limit`, `Offset`, `Order("created_at DESC")`
- `Delete` sets `IsActive = false` then calls GORM's `Delete` (soft delete)

### Use Case (`usecase.go`)
- Each **action** gets its own interface + private implementation (e.g. `CreatePostUseCase`, `PublishPostUseCase`)
- Interface: `Execute(ctx, ...params) error`
- Constructor: `func New<Action>UseCase(repo <Name>Repository) <Action>UseCase`
- Business rules go here, not in handlers or entities
- Error messages are plain strings with `errors.New(...)` (Portuguese)

### Handler (`internal/<module>/http/handlers.go`)
- Package is named `http` (not `http` aliased when imported)
- Single handler struct holds all use cases for the module
- Constructor receives each use case individually
- Methods: `Create`, `Publish`, `Update`, `Delete`
- ID param is `c.Param("id")`, parsed with `strconv.Atoi`, validated `> 0`
- Use `c.ShouldBindJSON` for request body
- Response pattern: `c.JSON(http.StatusXXX, gin.H{"message": "...", "error": "..."})`
- Swagger annotations on each handler method

### Main wiring (`cmd/main.go`)
- Load env → connect DB → auto-migrate → create repo → create use cases → create handler → register routes → run
- Route prefix: `/api/v1`
- Swagger at `/swagger/*any`

## Key Commands

```bash
# Run
go run cmd/main.go

# Generate Swagger docs (after changing annotations)
swag init -g cmd/main.go

# Scaffold a new module
scarf mod -n <name>

# Tidy dependencies
go mod tidy
```

## Business Rules

### Creation
- Title is required
- Default status is `draft`

### Publishing
- Post must exist and be active
- Post must not already be published
- Content must be at least 10 characters

### Editing
- Post must exist and be active
- Inactive/deleted posts cannot be edited

### Deletion
- Logical (soft delete) — sets `IsActive = false`, GORM fills `DeletedAt`

## .env
- Never read or modify `.env` files
- Required vars: `DB_HOST`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_PORT`, `API_PORT`
- A `.env.example` may exist; use it as reference if you need env vars

## Future Infrastructure (planned studies)

- **Docker** — containerize the app
- **Nginx** — reverse proxy in front of the API
- **GitHub Actions** — CI/CD pipeline (lint, test, build, deploy)
- **Design systems** — service boundaries, observability, reliability patterns

Every decision today should keep these future layers in mind — no shortcuts that would force a rewrite later.

## General Rules
- Write **Portuguese** for user-facing messages (errors, logs)
- Write **English** for code identifiers, comments, and documentation
- Do not read or expose `.env` file contents
- Keep layers separated: handler never calls DB directly, use case never handles HTTP
- Never compromise Clean Architecture boundaries — even for "quick fixes"
- Scarf is the only way to create new modules — never copy-paste a module folder manually
