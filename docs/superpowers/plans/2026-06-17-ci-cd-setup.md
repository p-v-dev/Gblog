# CI/CD Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** GitHub Actions CI (vet + build) em push na main + deploy automático no Render via Docker.

**Architecture:** Workflow YAML único em `.github/workflows/ci.yml`. Render builda direto do Dockerfile existente — sem push de imagem, sem código de deploy.

**Tech Stack:** GitHub Actions (ubuntu-latest, Go 1.26), Render (Docker deploy)

## Global Constraints

- Go 1.26 (`actions/setup-go@v5` com `go-version: '1.26'`)
- Apenas push na branch `main` dispara o workflow
- `go vet ./...` + `go build ./...`
- Sem secrets no GitHub Actions
- Sem mudanças no código Go

---

### Task 1: Workflow CI

**Files:**
- Create: `.github/workflows/ci.yml`

- [ ] **Step 1: Criar diretório**

```bash
New-Item -ItemType Directory -Path ".github\workflows" -Force
```

- [ ] **Step 2: Criar ci.yml**

```yaml
name: CI

on:
  push:
    branches: [main]

jobs:
  ci:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.26'

      - name: Download deps
        run: go mod download

      - name: Vet
        run: go vet ./...

      - name: Build
        run: go build ./...
```

- [ ] **Step 3: Verificar arquivo existe**

```powershell
Test-Path -LiteralPath ".github\workflows\ci.yml"
```
Expected: `True`

- [ ] **Step 4: Commit**

```bash
git add .github/workflows/ci.yml
git commit -m "ci: add GitHub Actions workflow (vet + build on main push)"
```

---

### Task 2: Deploy no Render (manual, fora do repositório)

> Ação manual no dashboard do Render, não código.

- [ ] **Step 1:** Criar conta em render.com (se não tiver)
- [ ] **Step 2:** Dashboard → New Web Service → conectar repositório GitHub
- [ ] **Step 3:** Configurar:
  - **Name:** `gblog-api`
  - **Runtime:** `Docker`
  - **Branch:** `main`
  - **Root Directory:** (vazio — Dockerfile na raiz)
- [ ] **Step 4:** Adicionar env vars:
  - `DB_HOST`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_PORT`
  - `API_PORT` — definir com o mesmo valor do `PORT` que Render injeta (ex: `10000`)
  - `JWT_SECRET`
- [ ] **Step 5:** Criar Web Service e verificar deploy
- [ ] **Step 6:** Testar endpoint:
  ```bash
  curl https://<render-url>.onrender.com/api/v1/posts
  ```
  Expected: `200` com JSON (array vazio ou posts)
