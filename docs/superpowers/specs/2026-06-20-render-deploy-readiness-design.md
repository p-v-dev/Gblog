# Render Deploy Readiness — Health Check & Graceful Shutdown

**Date:** 2026-06-20
**Status:** Approved

## Motivation

First deploy to Render free tier (Docker + Neon PostgreSQL). The app is functionally complete but lacks two features Render's platform expects:

1. A health check endpoint so Render knows the app is alive
2. Graceful shutdown so ongoing requests aren't dropped during restarts

## Changes

All changes are confined to `cmd/main.go`. No new files, no new packages.

### 1. Health Check Endpoint

- **Route:** `GET /health` (public, no auth)
- **Response:** `200 {"status":"ok"}`
- **Scope:** App-only — confirms the HTTP server is accepting connections
- **DB check omitted** intentionally to avoid hammering Neon free tier on Render's periodic pings
- Registered before `/api/v1` group, outside any auth middleware

### 2. Graceful Shutdown

Replace `r.Run()` with explicit `http.Server` lifecycle:

```
main:
  create *http.Server{Addr, Handler: r}
  go srv.ListenAndServe()
  signal.Notify(quit, SIGTERM, SIGINT)
  <-quit
  log "Desligando servidor..."
  ctx, cancel = context.WithTimeout(10s)
  srv.Shutdown(ctx)
  sqlDB.Close()   # close GORM connection pool
```

- `SIGTERM` covers Render's stop signal
- `SIGINT` covers local Ctrl+C
- 10-second timeout drains in-flight requests before hard exit
- GORM's underlying `*sql.DB` is closed after `Shutdown` returns
- If `ListenAndServe` errors before shutdown, `log.Fatalf`

### 3. API_PORT Validation

- Add explicit check at startup: if `API_PORT` is empty, `log.Fatal` with clear message
- No fallback chain — user explicitly wants fail-fast if env is misconfigured

## CI

No changes. Current GitHub Actions workflow (`go vet` + `go build`) is sufficient.

## Render Dashboard Configuration

| Variable | Value |
|----------|-------|
| `API_PORT` | `8080` |
| `JWT_SECRET` | (generate a random key) |
| `DATABASE_URL` | Neon connection string with sslmode=require |

- Health check path: `/health`
- Health check type: HTTP

## Files Changed

| File | Change |
|------|--------|
| `cmd/main.go` | Add `GET /health` route, replace `r.Run()` with graceful shutdown, add `API_PORT` validation |
