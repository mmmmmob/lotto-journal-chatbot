<!-- AI-CONTEXT
entities_active: [Go, Fiber-v2, PostgreSQL, GORM, go-migrate, pnpm, Turborepo, air, LINE-Messaging-API]
entities_deprecated: [Next.js]
entities_proposed: []
last_updated: 2026-04-30
-->

---

# Entity Register — Lotto Journal

Last updated: 2026-04-30

---

## Active Entities

| Entity             | Type        | Status | Since   | ADR     | Notes                                                                     |
| ------------------ | ----------- | ------ | ------- | ------- | ------------------------------------------------------------------------- |
| Go                 | tech        | active | 2026-04 | —       | Primary backend language; module: `lotto-journal/api`                     |
| Fiber v2           | dependency  | active | 2026-04 | —       | HTTP framework for Go (`github.com/gofiber/fiber/v2`)                     |
| PostgreSQL         | tech        | active | 2026-04 | —       | Primary database                                                          |
| GORM               | dependency  | active | 2026-04 | —       | ORM for Go + PostgreSQL driver                                            |
| go-migrate         | dependency  | active | 2026-04 | —       | Database migrations (`apps/api/migrations/`)                              |
| pnpm               | tech        | active | 2026-04 | —       | Package manager; monorepo root                                            |
| Turborepo          | tech        | active | 2026-04 | —       | Monorepo build orchestration                                              |
| air                | dependency  | active | 2026-04 | —       | Hot reload for Go during development                                      |
| LINE Messaging API | integration | active | 2026-04 | ADR-001 | Primary user interaction channel — webhook receiver + push message sender |

---

---

## Deprecated / Removed Entities

| Entity  | Type | Status     | Since   | Until   | ADR     | Replaced By                                  |
| ------- | ---- | ---------- | ------- | ------- | ------- | -------------------------------------------- |
| Next.js | tech | deprecated | 2026-04 | 2026-04 | ADR-001 | LINE Messaging API (no web UI for end users) |

---

## DB Schema Notes (post-ADR-001)

**Tables staying unchanged:** `draws`, `tickets`, `draw_results`, `user_winnings`, `files`

**`users` table — to be redesigned in migration 000002:**

- REMOVE: `username`, `email`, `password_hash`
- ADD: `line_user_id varchar UNIQUE NOT NULL`
- KEEP: `status`, `created_at`, `updated_at`

**Tables to be dropped in migration 000002:**

- `user_auth_methods`
- `user_verifications`

**Enums to be dropped in migration 000002:**

- `provider_service`
- `verification_type`
