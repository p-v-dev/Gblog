# Roadmap

## Phase 1 — Blog reader endpoints
- [x] `GET /posts` — list published posts (uses `FetchAll` in repo, already exists)
- [x] `GET /posts/:slug` — view single post by slug (uses `GetBySlug` in repo, already exists)
- [x] `GET /users/:id` — view author info
- [x] Pagination validation (min/max limit, clamp offset)
- [x] Response DTOs (strip internal fields like `is_active`)
- [x] Swagger annotations on new endpoints

After this the API is consumable as a real blog.

## Phase 2 — Complete user management
- [x] `PUT /users/:id` — update profile
- [x] `DELETE /users/:id` — deactivate account
- [x] Validate `UserID` exists on post creation (app-level check)
- [x] Input validation on update (at least one field required)
- [x] Swagger annotations on new user endpoints (PUT/DELETE)
- [x] Auto-migrate user table (`db.AutoMigrate(&user.User{})`)

## Phase 3 — Ownership and auth
- [x] Owner-only post edit/delete (compare `UserID` from request)
- [x] Auth tokens (JWT or basic token)
- [x] Protected routes (middleware extracts user from token)
- [x] JWT secret in env (required at startup, fails if unset)
- [x] Password-based login instead of raw `user_id`

## Phase 4 — Content enrichment
- [x] Categories / tags
- [x] Comments
- [x] Markdown support (rendered server-side via goldmark)
- [ ] Upload de imagens e GIFs
- [x] Swagger annotations on comment endpoints
- [x] Validate post exists before creating comment

## Phase 5 — Infrastructure
- [x] Docker (Dockerfile + docker-compose com api e postgres)
- [ ] Nginx reverse proxy
- [ ] CI/CD (GitHub Actions)
