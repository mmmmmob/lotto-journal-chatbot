<!-- AI-CONTEXT
entities_active: [Go, Fiber-v3, PostgreSQL, GORM, go-migrate, pnpm, Turborepo, air, LINE-Messaging-API, line-bot-sdk-go-v8]
entities_deprecated: [Next.js, Fiber-v2]
entities_proposed: []
last_updated: 2026-05-08
-->

---

# Entity Register — Lotto Journal

Last updated: 2026-05-08

---

## Active Entities

| Entity             | Type        | Status | Since   | ADR     | Notes                                                                                                 |
| ------------------ | ----------- | ------ | ------- | ------- | ----------------------------------------------------------------------------------------------------- |
| Go                 | tech        | active | 2026-04 | —       | Primary backend language; module: `lotto-journal/api`                                                 |
| Fiber v3           | dependency  | active     | 2026-05 | —       | HTTP framework for Go (`github.com/gofiber/fiber/v3`) — upgraded from v2 in session 6            |
| PostgreSQL         | tech        | active | 2026-04 | —       | Primary database                                                                                      |
| GORM               | dependency  | active | 2026-04 | —       | ORM for Go + PostgreSQL driver                                                                        |
| go-migrate         | dependency  | active | 2026-04 | —       | Database migrations (`apps/api/migrations/`)                                                          |
| pnpm               | tech        | active | 2026-04 | —       | Package manager; monorepo root                                                                        |
| Turborepo          | tech        | active | 2026-04 | —       | Monorepo build orchestration                                                                          |
| air                | dependency  | active | 2026-04 | —       | Hot reload for Go during development                                                                  |
| line-bot-sdk-go v8 | dependency  | active | 2026-05 | ADR-001 | Official LINE Bot SDK for Go (`github.com/line/line-bot-sdk-go/v8`) — webhook parsing + messaging API |
| LINE Messaging API | integration | active | 2026-04 | ADR-001 | Primary user interaction channel — webhook receiver + push message sender                             |

---

---

## Deprecated / Removed Entities

| Entity  | Type | Status     | Since   | Until   | ADR     | Replaced By                                  |
| ------- | ---- | ---------- | ------- | ------- | ------- | -------------------------------------------- |
| Next.js | tech | deprecated | 2026-04 | 2026-04 | ADR-001 | LINE Messaging API (no web UI for end users) |
| Fiber v2 | dependency | deprecated | 2026-04 | 2026-05 | — | Upgraded to Fiber v3 (v3.2.0) in session 6 — `timeout.New` race-free fix |

---

## DB Schema Notes (post-migration 000003)

**Tables unchanged:** `draws`, `tickets`, `draw_results`, `user_winnings`, `files`

**`users` table — redesigned in migration 000002 (done):**

- REMOVED: `username`, `email`, `password_hash`
- ADDED: `line_user_id varchar UNIQUE NOT NULL`
- KEPT: `status`, `created_at`, `updated_at`

**Tables dropped in migration 000002 (done):**

- `user_auth_methods`
- `user_verifications`

**Enums dropped in migration 000002 (done):**

- `provider_service`
- `verification_type`

**New in migration 000003 (done):**

- `webhook_events` table — idempotency store for LINE webhook event IDs
