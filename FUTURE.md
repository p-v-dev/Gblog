# Roadmap

## Phase 1 — Blog reader endpoints
- [x] `GET /posts` — list published posts
- [x] `GET /posts/slug/:slug` — view single post by slug
- [x] `GET /users/:id` — view author info
- [x] Pagination validation (min/max limit, clamp offset)
- [x] Response DTOs (strip internal fields like `is_active`)
- [x] Swagger annotations on new endpoints

## Phase 2 — Complete user management
- [x] `PUT /users/:id` — update profile
- [x] `DELETE /users/:id` — deactivate account
- [x] Validate `UserID` exists on post creation (app-level check)
- [x] Input validation on update (at least one field required)
- [x] Swagger annotations on new user endpoints (PUT/DELETE)
- [x] Auto-migrate user table

## Phase 3 — Ownership and auth
- [x] Owner-only post edit/delete (compare `UserID` from request)
- [x] Auth tokens (JWT)
- [x] Protected routes (middleware extracts user from token)
- [x] JWT secret in env (required at startup, fails if unset)
- [x] Password-based login instead of raw `user_id`

## Phase 4 — Content enrichment
- [x] Categories / tags (many2many via post_tags)
- [x] Comments
- [x] Markdown support (rendered server-side via goldmark)
- [ ] Upload de imagens e GIFs
- [x] Swagger annotations on comment endpoints
- [x] Validate post exists before creating comment

## Phase 5 — Infrastructure
- [x] Docker (Dockerfile multi-stage + docker-compose com api e postgres)
- [x] CI/CD (GitHub Actions — `go vet` + `go build` em push na main)
- [ ] Nginx reverse proxy
