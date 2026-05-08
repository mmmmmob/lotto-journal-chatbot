# apps/api — Go API Service

The backend service for Lotto Journal. Built with Go + Fiber.

> **Setup and running instructions are in the [root README](../../README.md).**
> This document covers API-specific development reference.

---

## Make targets

All commands run from `apps/api/`.

### App

| Command      | What it does                    |
| ------------ | ------------------------------- |
| `make run`   | Start API with hot reload (air) |
| `make build` | Build binary to `dist/`         |
| `make clean` | Remove `tmp/` and `dist/`       |

### Database

| Command         | What it does                                                                |
| --------------- | --------------------------------------------------------------------------- |
| `make db-start` | Start PostgreSQL container                                                  |
| `make db-stop`  | Stop PostgreSQL container (pause — data kept, restart with `make db-start`) |

### Migrations

| Command                      | What it does                                        |
| ---------------------------- | --------------------------------------------------- |
| `make migrate-up`            | Apply all pending migrations                        |
| `make migrate-up-one`        | Apply the next 1 migration only                     |
| `make migrate-down`          | Roll back the last migration                        |
| `make migrate-down-all`      | Roll back all migrations                            |
| `make migrate-version`       | Show current schema version                         |
| `make migrate-force N=<ver>` | Force-set version (use to recover from dirty state) |

---

## Project structure

```
apps/api/
├── app/
│   └── main.go              # Entry point
├── internal/
│   ├── config/              # Env config loader
│   ├── database/            # DB connection
│   ├── handler/             # HTTP handlers
│   ├── models/              # GORM models
│   ├── repository/          # DB access layer
│   └── service/             # Business logic
├── migrations/              # SQL migration files
├── middlewares/             # Fiber middlewares
├── Makefile
└── go.mod
```

---

## Adding a new migration

1. Create two files in `migrations/` following the naming convention:
   ```
   000004_<description>.up.sql
   000004_<description>.down.sql
   ```
2. Write the `up` SQL (schema change) and the `down` SQL (full reversal).
3. Apply with `make migrate-up-one` and verify with `make migrate-version`.
4. Always test the `down` migration too: `make migrate-down` then `make migrate-up-one`.

### Migration history

| Version | File                    | Description                                                                                      |
| ------- | ----------------------- | ------------------------------------------------------------------------------------------------ |
| 000001  | `000001_init_schema`    | Initial schema — all tables, enums, indexes                                                      |
| 000002  | `000002_line_identity`  | LINE identity redesign — replace email/password with `line_user_id`; rename `N6→L6`, `n6_*→l6_*` |
| 000003  | `000003_webhook_events` | Idempotency table — store processed LINE `webhookEventId` values (ON CONFLICT DO NOTHING)        |
