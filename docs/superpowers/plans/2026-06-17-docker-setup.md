# Docker Setup Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Containerizar API Go + PostgreSQL via Docker multi-stage e docker-compose.

**Architecture:** Multi-stage Dockerfile (build em `golang:1.26-alpine`, runtime em `alpine:3.21`), docker-compose com serviços `api` e `db` linkados por rede interna, health check no PostgreSQL.

**Tech Stack:** Docker, docker-compose, Go 1.26, PostgreSQL 16, Alpine 3.21

**Anotações para futuro:** Nginx reverse proxy, health check endpoint `/health`, perfis dev/prod no compose, CI/CD.

## Global Constraints

- Go 1.26.3 (`golang:1.26-alpine` no build stage)
- PostgreSQL 16 (`postgres:16-alpine`)
- Runtime: `alpine:3.21` com `ca-certificates` e `tzdata`
- `CGO_ENABLED=0` — binário estático
- Sem mudanças no código Go — `godotenv.Load()` falha silenciosamente sem `.env`
- Porta default da API: 8080
- `JWT_SECRET` obrigatório (sem default)
- Docker não disponível na máquina atual — verificação será feita em notebook Fedora/Windows

---

### Task 1: .dockerignore

**Files:**
- Create: `.dockerignore`

- [ ] **Step 1: Criar .dockerignore**

```dockerignore
.git/
.env
```

- [ ] **Step 2: Verificar arquivo existe**

```bash
Test-Path -LiteralPath ".dockerignore"
```
Expected: `True`

- [ ] **Step 3: Commit**

```bash
git add .dockerignore
git commit -m "chore: add .dockerignore"
```

---

### Task 2: Dockerfile

**Files:**
- Create: `Dockerfile`

**Interfaces:**
- Produces: imagem Docker com binário em `/app/bin`, porta `EXPOSE 8080`

- [ ] **Step 1: Criar Dockerfile**

```dockerfile
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/bin ./cmd/main.go

FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /app/bin /app/bin
EXPOSE 8080
CMD ["/app/bin"]
```

- [ ] **Step 2: Verificar arquivo existe**

```bash
Test-Path -LiteralPath "Dockerfile"
```
Expected: `True`

- [ ] **Step 3: Commit**

```bash
git add Dockerfile
git commit -m "feat: add multi-stage Dockerfile"
```

---

### Task 3: docker-compose.yml

**Files:**
- Create: `docker-compose.yml`

**Interfaces:**
- Produces: serviços `api` (build local) e `db` (postgres:16-alpine)

- [ ] **Step 1: Criar docker-compose.yml**

```yaml
services:
  db:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: ${DB_USER:-gblog}
      POSTGRES_PASSWORD: ${DB_PASSWORD:-gblog}
      POSTGRES_DB: ${DB_NAME:-gblog}
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER:-gblog}"]
      interval: 5s
      retries: 5

  api:
    build: .
    ports:
      - "${API_PORT:-8080}:${API_PORT:-8080}"
    environment:
      DB_HOST: db
      DB_USER: ${DB_USER:-gblog}
      DB_PASSWORD: ${DB_PASSWORD:-gblog}
      DB_NAME: ${DB_NAME:-gblog}
      DB_PORT: "5432"
      API_PORT: ${API_PORT:-8080}
      JWT_SECRET: ${JWT_SECRET}
    depends_on:
      db:
        condition: service_healthy

volumes:
  pgdata:
```

- [ ] **Step 2: Verificar arquivo existe**

```bash
Test-Path -LiteralPath "docker-compose.yml"
```
Expected: `True`

- [ ] **Step 3: Commit**

```bash
git add docker-compose.yml
git commit -m "feat: add docker-compose with api and postgres"
```

---

### Task 4: Verificação em máquina com Docker

> Executar no notebook Fedora/Windows com Docker instalado.

- [ ] **Step 1: Build e start**

```bash
# criar .env com JWT_SECRET (mínimo pra teste)
echo "JWT_SECRET=chave-de-teste" > .env

# build e up
docker compose up --build -d
```

Expected: containers `gblog-db-1` e `gblog-api-1` running sem erros.

- [ ] **Step 2: Verificar health check do DB**

```bash
docker compose ps
```

Expected: db com status `healthy`.

- [ ] **Step 3: Testar endpoint da API**

```bash
curl -s http://localhost:8080/api/v1/posts | head -20
```

Expected: JSON com lista de posts (pode ser vazia `[]`) — sem erro de conexão.

- [ ] **Step 4: Testar Swagger docs**

```bash
curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/swagger/index.html
```

Expected: `200`.

- [ ] **Step 5: Logs sem erros**

```bash
docker compose logs api
```

Expected: sem stack traces ou panic. Apenas logs do Gin (`[GIN] ...`).

- [ ] **Step 6: Cleanup**

```bash
docker compose down -v
```
