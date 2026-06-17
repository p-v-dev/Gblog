# Docker Setup — Gblog API

## Objetivo

Containerizar a API Go + PostgreSQL para rodar em qualquer ambiente com Docker, sem dependências locais além do próprio Docker.

## Escopo

- Dockerfile multi-stage para a API Go
- docker-compose.yml com serviços `api` e `db`
- .dockerignore
- Sem mudanças no código Go (`godotenv.Load()` falha silenciosamente sem `.env`)
- Nginx, CI/CD, upload de imagens — fases futuras

## Arquivos

```
Gblog/
├── Dockerfile
├── docker-compose.yml
└── .dockerignore
```

## Dockerfile — Multi-stage

Stage 1 (`builder`): `golang:1.26-alpine` compila o binário com `CGO_ENABLED=0`.

Stage 2 (`runtime`): `alpine:3.21` com `ca-certificates` e `tzdata`. Copia só o binário.

`EXPOSE 8080` — porta default para documentação da imagem. O compose mapeia a porta real em runtime.

## docker-compose.yml

### Serviço `db`

- Imagem: `postgres:16-alpine`
- Health check com `pg_isready` (5s intervalo, 5 tentativas)
- Volume nomeado `pgdata` para persistência
- Porta `5432` exposta (opcional, apenas para ferramentas externas)

### Serviço `api`

- Build: `Dockerfile` do diretório atual
- `depends_on` com `condition: service_healthy` — só sobe após DB pronto
- Variáveis de ambiente mapeadas diretamente (sem volume `.env`)
- Porta configurável via `API_PORT`

### Variáveis

Todas as variáveis têm default via `${VAR:-default}` para dev local. `JWT_SECRET` não tem default — obrigatório.

`DB_HOST` fixo como `db` (nome do serviço). `DB_PORT` fixo como `5432` (porta interna do PostgreSQL).

### Volumes

`pgdata` — volume nomeado para dados do PostgreSQL.

## .dockerignore

```
.git/
.env
```

Exclui git e `.env` local. `cmd/docs/` NÃO é excluído — o binário importa os Swagger docs gerados.

## Teste

Sem Docker na máquina atual. O repositório será testado em notebook com Fedora (ou Windows dual boot) executando:

```bash
docker compose up --build
```

## Fora do escopo (futuro)

- Nginx reverse proxy
- CI/CD (GitHub Actions)
- Upload de imagens/GIFs
- Health check na API (endpoint `/health`)
- Docker Compose profile para produção vs dev
