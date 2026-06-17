# CI/CD — Gblog API

## Objetivo

CI via GitHub Actions (validação em push na main) + CD via Render (deploy automático do Dockerfile).

## Stack

- GitHub Actions (ubuntu-latest, Go 1.26)
- Render Web Service (Docker deploy)

## Escopo

- Workflow `.github/workflows/ci.yml` com `go vet` + `go build`
- Deploy delegado ao Render (nenhum código novo)
- Nginx não faz parte deste ciclo

## Fora do escopo

- Testes (`go test`) — não existem ainda no projeto; adicionar quando houver testes
- Linter externo (golangci-lint) — `go vet` já cobre o essencial
- Build/push de imagem Docker no Actions — Render builda direto do Dockerfile
- Secrets no GitHub Actions — nenhum necessário

## Arquivos

```
Gblog/
├── .github/
│   └── workflows/
│       └── ci.yml
```

## Workflow CI

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

- `go vet` — análise estática sem dependência extra
- `go build` — valida que compila
- Sem segredos ou variáveis de ambiente

## Deploy (Render — configuração manual)

1. Dashboard Render → New Web Service → conectar GitHub repo
2. Runtime: **Docker** (aponta para o Dockerfile existente)
3. Branch: `main`
4. Variáveis de ambiente obrigatórias:
   - `DB_HOST`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_PORT`
   - `API_PORT` (definir com o mesmo valor da variável `PORT` que o Render injeta automaticamente, ex: `10000`)
   - `JWT_SECRET`
5. Render faz auto-deploy em cada push na branch configurada

O PostgreSQL pode ser um Render Postgres ou serviço externo.
