# Handoff — Gblog

## Estado atual

API funcional com auth JWT, posts (draft → publish), tags (many2many), comentários, users. Fases 1–4 completas (exceto upload de imagens). Docker configurado.

## Rotas (`/api/v1`)

### Públicas
| Método | Rota | Descrição |
|--------|------|-----------|
| POST | `/auth/token` | Login (email + senha) → JWT |
| GET | `/posts` | Listar posts (`?limit=&offset=`) |
| GET | `/posts/slug/:slug` | Buscar post por slug |
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
├── blogPost/     — Posts CRUD + publish, tags many2many
├── comment/      — Comentários por post
├── user/         — Usuários, login com bcrypt + JWT
├── tag/          — Tags (find-or-create por nome, unique)
└── infra/        — DB, JWT, auth middleware
```

## Dependências

- `github.com/golang-jwt/jwt/v5` — JWT HS256
- `golang.org/x/crypto/bcrypt` — hash de senha
- `gorm.io/gorm` — ORM
- `github.com/gin-gonic/gin` — HTTP framework
- `github.com/gin-contrib/cors` — CORS middleware
- `github.com/swaggo/swag` — Swagger docs
- `github.com/yuin/goldmark` — Markdown → HTML

## Padrões importantes

- **Owner check**: handlers leem `user_id` do context (`c.Get("user_id")`, setado pelo middleware) com fallback pro header `X-User-ID`
- **Auth**: `POST /auth/token` com `{"email","password"}` → token. Requer `JWT_SECRET` no `.env`
- **Tags**: find-or-create por nome. `POST /posts` aceita `{"tags":["golang"]}`, `PUT /posts/:id` também
- **Soft delete**: `IsActive = false` em posts e comments, `DeletedAt` do GORM
- **CORS**: configurado com `AllowOrigins: *`, permite header `Authorization` para rotas protegidas
- **Slug**: gerado automaticamente do título com sufixo timestamp (`meu-titulo-1718765432100`) para evitar colisão
- **Create response**: retorna o post criado (com ID e slug populados) em vez de mensagem fixa

## Pendente

Ver `FUTURE.md`:
- Upload de imagens — Phase 4
- CI/CD (GitHub Actions)
- Nginx reverse proxy — Phase 5

## Comandos

```bash
go run cmd/main.go                    # Rodar
swag init -g cmd/main.go              # Gerar Swagger
scarf mod -n <name>                   # Scaffold novo módulo
docker compose up                     # Rodar com Docker
```

## Notas

- `POST /tags` foi movido pra rota protegida — se precisar público, mover de volta
- `currentUserID()` duplicado entre `blogPost/http/handlers.go` e `comment/http/handlers.go` — extrair pra infra se incomodar
- `.env` requer `JWT_SECRET` além das vars de DB — o servidor falha no startup se não estiver setado
