# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build/Run

```bash
# Run locally (loads .env.local then .env)
go run main.go

# Build binary
go build -o telegram-expense-bot

# Unit tests (no external dependencies)
go test ./... -count=1

# Integration tests (requires Podman socket or Docker)
# These files have //go:build integration tag
systemctl --user enable --now podman.socket   # one-time Podman setup
go test -tags=integration ./... -count=1 -timeout 180s

# Coverage
go test ./... -cover
go test -tags=integration ./... -cover -timeout 180s

# Single test
go test ./domain -run TestUserFromEntity -v
```

## Architecture

This is a Telegram bot for expense tracking, written in Go. It uses PostgreSQL for persistence and communicates with users exclusively through the Telegram Bot API.

### Layer stack

Each data type follows an identical layered pattern:

| Layer | Package | Role |
|-------|---------|------|
| **entity** | `entity/` | Database row structs (flat, raw DB types) |
| **dao** | `dao/` | SQL queries via `pgxpool.Pool` + `pgxscan` |
| **repo** | `repo/` | Maps entity ↔ domain, contains business logic |
| **domain** | `domain/` | Rich types used by handlers (uses `go-money`) |
| **handler** | `handler/` | Bot commands and callback query processing |

The wiring order in `main.go` is: config → db pool → DAOs → repos → handlers → bot. The repoes own all data transformation between entity and domain types.

### Message flow

- **Polling mode** (default) or **webhook mode** — configured at startup via `WEBHOOK_ENABLED`
- Incoming updates are dispatched concurrently with `NUM_ROUTINES` goroutines
- `main.go` routes by update type:
  - **Messages**: If it's a command (`/start`, `/help`, `/stats`, `/undo`, `/list`, `/export`), routed to `CommandHandler`. Otherwise, treated as a new expense entry and routed to `StartTransaction`.
  - **Callback queries**: dispatched by `CallbackType` JSON discriminator field (`Category`, `Pagination`, `Undo`, `Cancel`) to `CallbackHandler`

### Inline keyboard flow

Interactive flows (category selection, pagination, undo confirmation) use inline keyboards with callback data encoded as compact JSON. The flow:

1. Handler creates a `domain.Callback` struct with a `Type` discriminator and a `MessageContextId`
2. Struct is marshalled to JSON and embedded in the button's callback data
3. `MessageContext` table stores the original user message text, allowing subsequent callbacks to recover context (e.g., re-parsing the original amount when the user picks a category)

### Key dependencies

- `go-telegram-bot-api/v5` — Telegram Bot API client
- `pgx/v5` + `scany/v2` — PostgreSQL driver and row scanning
- `go-money` — currency-safe money representation (amounts stored as integers in the smallest denomination)
- `zerolog` — structured logging, with optional Telegram hook for error/fatal/panic alerts
- `excelize/v2` — Excel export (`/export` command)
- `caarlos0/env/v6` — env var parsing into config struct
- `joho/godotenv` — `.env` file loading
