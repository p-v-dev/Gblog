# Gblog API

API REST para gerenciamento de artigos de um blog, construída com **Go**, **Gin**, **GORM** e **PostgreSQL**, seguindo **Clean Architecture**.

---

## Funcionalidades

- **Posts**: criar (draft), editar, publicar, deletar (soft delete)
- **Tags**: many2many com posts, find-or-create por nome
- **Comentários**: criar, listar, deletar por post
- **Usuários**: criar, login (bcrypt + JWT), atualizar, desativar
- **Auth**: JWT Bearer token em rotas protegidas
- **Markdown**: renderizado server-side via goldmark
- **Swagger**: documentação automática dos endpoints

---

## Estrutura do Projeto

```
Gblog/
├── cmd/
│   ├── main.go                  # Entry point: env, DB, use cases, rotas, CORS
│   └── docs/                    # Swagger docs (gerado por swag init)
├── internal/
│   ├── blogPost/                # Posts + tags many2many
│   │   ├── entity.go
│   │   ├── dto.go
│   │   ├── repository.go
│   │   ├── usecase.go
│   │   └── http/handlers.go
│   ├── comment/                 # Comentários
│   │   ├── entity.go
│   │   ├── dto.go
│   │   ├── repository.go
│   │   ├── usecase.go
│   │   └── http/handlers.go
│   ├── user/                    # Usuários + auth
│   │   ├── entity.go
│   │   ├── dto.go
│   │   ├── repository.go
│   │   ├── usecase.go
│   │   └── http/handlers.go
│   ├── tag/                     # Tags
│   │   ├── entity.go
│   │   ├── dto.go
│   │   ├── repository.go
│   │   └── usecase.go
│   └── infra/
│       └── database.go          # DB, JWT, auth middleware
├── pkg/
│   └── blogstatus.go            # Enum: draft, published, archived
├── Dockerfile
├── docker-compose.yml
├── .env                         # (não versionado)
├── go.mod
└── go.sum
```

---

## Tecnologias

| Tecnologia | Versão | Descrição |
|------------|--------|-----------|
| Go | 1.26.3 | Linguagem |
| Gin | v1.12.0 | HTTP framework |
| GORM | v1.31.1 | ORM |
| PostgreSQL | — | Banco de dados |
| Swag | v1.16.6 | Swagger docs |
| goldmark | — | Markdown → HTML |
| golang-jwt | v5 | JWT HS256 |
| gin-contrib/cors | — | CORS middleware |

---

## Pré-requisitos

- [Go](https://go.dev/dl/) 1.26.3+
- [PostgreSQL](https://www.postgresql.org/)
- [Swag CLI](https://github.com/swaggo/swag) (opcional, para gerar docs)
- [Docker](https://www.docker.com/) (opcional)

---

## Variáveis de Ambiente

| Variável | Descrição | Exemplo |
|----------|-----------|---------|
| `API_PORT` | Porta do servidor | `8080` |
| `DB_HOST` | Host do PostgreSQL | `localhost` |
| `DB_USER` | Usuário do banco | `postgres` |
| `DB_PASSWORD` | Senha do banco | `senha123` |
| `DB_NAME` | Nome do banco | `gblog` |
| `DB_PORT` | Porta do PostgreSQL | `5432` |
| `JWT_SECRET` | Chave secreta JWT | `minha-chave-super-secreta` |

---

## Execução

```bash
# Local
go run cmd/main.go

# Docker
docker compose up
```

A API estará em `http://localhost:{API_PORT}`. Docs Swagger em `/swagger/index.html`.

---

## Endpoints da API

Base URL: `/api/v1`

### Públicos
| Método | Rota | Descrição |
|--------|------|-----------|
| POST | `/auth/token` | Login → JWT |
| GET | `/posts` | Listar posts (`?limit=&offset=`) |
| GET | `/posts/slug/:slug` | Buscar post por slug |
| GET | `/posts/:id/comments` | Listar comentários |
| GET | `/tags` | Listar tags |
| GET | `/users/:id` | Ver usuário |
| POST | `/users` | Criar usuário |

### Protegidos (Bearer token)
| Método | Rota | Descrição |
|--------|------|-----------|
| POST | `/posts` | Criar rascunho |
| PUT | `/posts/:id` | Editar post |
| PATCH | `/posts/:id/publish` | Publicar |
| DELETE | `/posts/:id` | Soft delete |
| POST | `/posts/:id/comments` | Comentar |
| DELETE | `/comments/:id` | Deletar comentário |
| POST | `/tags` | Criar tag |
| PUT | `/users/:id` | Atualizar perfil |
| DELETE | `/users/:id` | Desativar conta |

---

## Modelos de Dados

### BlogPost

| Campo | Tipo | Descrição |
|-------|------|-----------|
| `ID` | `uuid` | Chave primária |
| `Title` | `string` | Título (255 max) |
| `Content` | `text` | Markdown |
| `Status` | `BlogStatus` | `draft`, `published`, `archived` |
| `Slug` | `string` | URL única (com timestamp) |
| `IsActive` | `bool` | Soft delete |
| `UserID` | `uuid` | Autor |
| `Tags` | `[]Tag` | Many2many via `post_tags` |
| `CreatedAt` | `time` | |
| `UpdatedAt` | `time` | |
| `DeletedAt` | `time` | |

### User

| Campo | Tipo | Descrição |
|-------|------|-----------|
| `ID` | `uuid` | Chave primária |
| `Name` | `string` | Nome |
| `Email` | `string` | Login (unique) |
| `Password` | `string` | Bcrypt hash |
| `IsActive` | `bool` | Soft delete |

### Tag

| Campo | Tipo | Descrição |
|-------|------|-----------|
| `ID` | `uuid` | Chave primária |
| `Name` | `string` | Unique |

### Comment

| Campo | Tipo | Descrição |
|-------|------|-----------|
| `ID` | `uuid` | Chave primária |
| `Content` | `text` | |
| `PostID` | `uuid` | FK → blog_posts |
| `UserID` | `uuid` | Autor |
| `IsActive` | `bool` | Soft delete |

---

## Regras de Negócio

### Criação
- Título obrigatório
- Status inicial: `draft`
- Slug gerado automaticamente (título + timestamp)

### Publicação
- Post deve existir e estar ativo
- Não pode já estar publicado
- Conteúdo ≥ 10 caracteres

### Edição
- Post deve existir e estar ativo
- Owner check (UserID do token)

### Exclusão
- Soft delete: `IsActive = false`
- Owner check

### Tags
- Find-or-create por nome
- Many2many via tabela `post_tags`

---

## Licença

MIT
