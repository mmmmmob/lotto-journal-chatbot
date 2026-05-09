# Lotto Journal

A LINE chatbot service that lets Thai lottery players record their ticket numbers and get
automatically notified via LINE when any of their tickets win a prize.

**Stack:** Go + Fiber · PostgreSQL · LINE Messaging API

---

## Bot Commands

All interaction happens inside the LINE chat with the bot.

| Message sent         | What happens                                                              |
| -------------------- | ------------------------------------------------------------------------- |
| `123456`             | Records one L6 ticket for the upcoming draw                               |
| `456`                | Records one N3 ticket for the upcoming draw                               |
| `123456 x3`          | Records 3 copies of the same L6 ticket                                    |
| `123456, 789012`     | Records two tickets in one message (comma or space separated)             |
| `โพย`               | Lists all tickets you have registered for the current upcoming draw       |
| Follow bot           | Creates your account; sends welcome message                               |
| Unfollow bot         | Marks your account inactive; ticket history is preserved                  |

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

| Variable                    | Used by        | Example value                                                                   |
| --------------------------- | -------------- | ------------------------------------------------------------------------------- |
| `DB_USERNAME`               | docker-compose | `postgres`                                                                      |
| `DB_PASSWORD`               | docker-compose | `yourpassword`                                                                  |
| `DB_NAME`                   | docker-compose | `lotto_journal`                                                                 |
| `DB_DSN`                    | Go app         | `postgres://postgres:yourpassword@localhost:5432/lotto_journal?sslmode=disable` |
| `PORT`                      | Go app         | `:3000`                                                                         |
| `LINE_CHANNEL_SECRET`       | Go app         | from LINE Developers console → Basic Settings                                   |
| `LINE_CHANNEL_ACCESS_TOKEN` | Go app         | from LINE Developers console → Messaging API                                    |

### 2. Start the database

```shell
pnpm db:start
```

Wait until the container is healthy before running migrations:

```shell
docker ps --filter name=lotto-db --format "table {{.Names}}\t{{.Status}}"
```

You should see `healthy` in the status before proceeding.

### 3. Run migrations

```shell
pnpm migrate:up
```

### 4. Start the API (with hot reload)

```shell
pnpm dev
```

---

## API Testing (Bruno)

The collection lives in `trunk/bruno/`. Open it in [Bruno](https://www.usebruno.com/) by choosing **Open Collection** and selecting that folder.

### Requests

| Folder | Request            | What it tests                                                                             |
| ------ | ------------------ | ----------------------------------------------------------------------------------------- |
| REST   | Webhook - Follow   | User adds the bot as a friend → creates user record, sends welcome reply                  |
| REST   | Webhook - Message  | User sends ticket numbers (e.g. `123456 x2, 789`) → parses and saves tickets              |
| REST   | Webhook - Unfollow | User removes the bot → marks user `inactive`                                              |
| REST   | Health             | `GET /health` — liveness + DB readiness; `200` ok / `503` degraded                        |
| GLO    | Check Result       | Calls the Thai Government Lottery API directly — useful for inspecting the result payload |

### One-time setup

**1. Set the environment**

In Bruno, select the **dev** environment (top-right dropdown). Then open **Configure** and set the secret variable:

| Variable              | Value                                                             |
| --------------------- | ----------------------------------------------------------------- |
| `line_channel_secret` | Your channel secret from LINE Developers Console → Basic Settings |

This value is stored locally by Bruno and never committed to git.

**2. Make sure the API is running**

```shell
pnpm dev
```

The `endpoint_url` in the dev environment points to `http://localhost:3000` by default. Adjust if your `PORT` is different.

### How the signature works

Every webhook request has a pre-request script that automatically computes the `X-Line-Signature` header before sending:

```js
const crypto = require('crypto');
const secret = bru.getEnvVar('line_channel_secret');
const body = JSON.stringify(req.body);
const signature = crypto.createHmac('sha256', secret).update(body).digest('base64');
req.setHeader('X-Line-Signature', signature);
```

You don't need to compute the signature manually — just send the request.

### Idempotency caveat

Each request body contains a `webhookEventId`. On first send, this ID is written to the `webhook_events` table. Sending the **same request a second time** will be silently skipped (the deduplication is working correctly).

To re-trigger processing, change the `webhookEventId` to any unique value before sending again.

### `userId` and `destination` fields

| Field                    | What it is                                        | Does our code use it?                            |
| ------------------------ | ------------------------------------------------- | ------------------------------------------------ |
| `events[].source.userId` | The LINE user ID of the person who sent the event | **Yes** — used to find or create the user record |
| `destination`            | Your bot's own LINE user ID                       | **No** — ignored by the handler                  |

The sample bodies use a fake `userId` (`U1234567890abcdef…`). Because `FindOrCreate` is idempotent, re-running the same request always operates on the same test user record in your local DB.

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

| pnpm (root)     | make (apps/api) | What it does                                                           |
| --------------- | --------------- | ---------------------------------------------------------------------- |
| `pnpm db:start` | `make db-start` | Start PostgreSQL container                                             |
| `pnpm db:stop`  | `make db-stop`  | Stop PostgreSQL container (pause — data kept, restart with `db:start`) |

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

| Version | File                    | Description                                                                                      |
| ------- | ----------------------- | ------------------------------------------------------------------------------------------------ |
| 000001  | `000001_init_schema`    | Initial schema — all tables, enums, indexes                                                      |
| 000002  | `000002_line_identity`  | LINE identity redesign — replace email/password with `line_user_id`; rename `N6→L6`, `n6_*→l6_*` |
| 000003  | `000003_webhook_events` | Idempotency table — store processed LINE `webhookEventId` values (ON CONFLICT DO NOTHING)        |
