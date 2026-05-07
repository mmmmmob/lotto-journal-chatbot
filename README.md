# Lotto Journal

A LINE chatbot service that lets Thai lottery players record their ticket numbers and get
automatically notified via LINE when any of their tickets win a prize.

**Stack:** Go + Fiber · PostgreSQL · LINE Messaging API

---

## Prerequisites

- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- [golang-migrate CLI](https://github.com/golang-migrate/migrate)

```shell
brew install golang-migrate
```

---

## Local Setup

> Run all commands from the **repo root** unless noted otherwise.

### 1. Environment variables

Copy the example env file and fill in your values:

```shell
cp .env.example .env.local
```

Key variables:

| Variable      | Used by        | Example value                                                                   |
| ------------- | -------------- | ------------------------------------------------------------------------------- |
| `DB_USERNAME` | docker-compose | `postgres`                                                                      |
| `DB_PASSWORD` | docker-compose | `yourpassword`                                                                  |
| `DB_NAME`     | docker-compose | `lotto_journal`                                                                 |
| `DB_DSN`      | Go app         | `postgres://postgres:yourpassword@localhost:5432/lotto_journal?sslmode=disable` |
| `PORT`        | Go app         | `:3000`                                                                         |

### 2. Start the database

```shell
pnpm db:start
```

### 3. Run migrations

```shell
pnpm migrate:up
```

### 4. Start the API (with hot reload)

```shell
pnpm dev
```

---

## Scripts reference

All `pnpm` commands run from the **repo root**.
All `make` commands run from `apps/api/` — the `pnpm` shortcuts above call these automatically.

### App

| pnpm (root)  | make (apps/api) | What it does                    |
| ------------ | --------------- | ------------------------------- |
| `pnpm dev`   | `make run`      | Start API with hot reload (air) |
| `pnpm build` | `make build`    | Build binary to `dist/`         |
| —            | `make clean`    | Remove `tmp/` and `dist/`       |

### Database

| pnpm (root)     | make (apps/api) | What it does               |
| --------------- | --------------- | -------------------------- |
| `pnpm db:start` | `make db-start` | Start PostgreSQL container |
| `pnpm db:stop`  | `make db-stop`  | Stop PostgreSQL container  |

### Migrations

| pnpm (root)            | make (apps/api)              | What it does                                 |
| ---------------------- | ---------------------------- | -------------------------------------------- |
| `pnpm migrate:up`      | `make migrate-up`            | Apply all pending migrations                 |
| —                      | `make migrate-up-one`        | Apply the next 1 migration only              |
| `pnpm migrate:down`    | `make migrate-down`          | Roll back the last migration                 |
| —                      | `make migrate-down-all`      | Roll back all migrations                     |
| `pnpm migrate:version` | `make migrate-version`       | Show current schema version                  |
| —                      | `make migrate-force N=<ver>` | Force-set version (recover from dirty state) |

### Migration history

| Version | File                   | Description                                                                                      |
| ------- | ---------------------- | ------------------------------------------------------------------------------------------------ |
| 000001  | `000001_init_schema`   | Initial schema — all tables, enums, indexes                                                      |
| 000002  | `000002_line_identity` | LINE identity redesign — replace email/password with `line_user_id`; rename `N6→L6`, `n6_*→l6_*` |
