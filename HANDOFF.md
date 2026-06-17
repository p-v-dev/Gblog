# Handoff — Gblog

## Estado atual

Fases 1–4 em andamento, Phase 4 parcial (rich content pendente). API funcional com auth JWT, posts, tags, comentários.

## Rotas (`/api/v1`)

### Públicas
| Método | Rota | Descrição |
|--------|------|-----------|
| POST | `/auth/token` | Login (email + senha) → JWT |
| GET | `/posts` | Listar posts (`?limit=&offset=`) |
| GET | `/posts/:slug` | Buscar post por slug |
| GET | `/posts/:id/comments` | Listar comentários do post |
| GET | `/tags` | Listar tags |
| GET | `/users/:id` | Ver usuário |
| POST | `/users` | Criar usuário |

### Protegidas (Bearer token)
| Método | Rota | Descrição |
|--------|------|-----------|
| POST | `/posts` | Criar post (draft) |
| PUT | `/posts/:id` | Editar post (owner) |
| PATCH | `/posts/:id/publish` | Publicar post |
| DELETE | `/posts/:id` | Soft delete (owner) |
| POST | `/posts/:id/comments` | Comentar |
| DELETE | `/comments/:id` | Deletar comentário (owner) |
| POST | `/tags` | Criar tag |
| PUT | `/users/:id` | Atualizar perfil |
| DELETE | `/users/:id` | Desativar conta |

## Módulos

```
internal/
├── blogPost/     — Posts, tags (many2many via post_tags)
├── comment/      — Comentários por post
├── user/         — Usuários, login com bcrypt
├── tag/          — Tags (find-or-create por nome)
└── infra/        — DB, JWT, auth middleware
```

## Dependências

- `github.com/golang-jwt/jwt/v5` — JWT HS256
- `golang.org/x/crypto/bcrypt` — hash de senha
- `gorm.io/gorm` — ORM
- `github.com/gin-gonic/gin` — HTTP framework
- `github.com/swaggo/swag` — Swagger docs

## Padrões importantes

- **Owner check**: handlers leem `user_id` do context (`c.Get("user_id")`, setado pelo middleware) com fallback pro header `X-User-ID`
- **Auth**: `POST /auth/token` com `{"email","password"}` → token. Requer `JWT_SECRET` no `.env`
- **Tags**: find-or-create por nome. `POST /posts` aceita `{"tags":["golang"]}`, `PUT /posts/:id` também
- **Soft delete**: `IsActive = false` em posts e comments, `DeletedAt` do GORM

## Pendente

Ver `FUTURE.md`:
- Rich content (markdown, images) — Phase 4
- Swagger annotations nos endpoints de comment
- Validar post existe antes de comentar
- Docker, Nginx, CI/CD — Phase 5

## Comandos

```bash
go run cmd/main.go                    # Rodar
swag init -g cmd/main.go              # Gerar Swagger
scarf mod -n <name>                   # Scaffold novo módulo
```

## Notas

- `POST /tags` foi movido pra rota protegida — se precisar público, mover de volta
- `currentUserID()` duplicado entre `blogPost/http/handlers.go` e `comment/http/handlers.go` — extrair pra infra se incomodar
