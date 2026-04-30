<!-- AI-CONTEXT
entities_active: [Go, Fiber-v2, PostgreSQL, GORM, go-migrate, pnpm, Turborepo, air, Next.js]
entities_deprecated: []
entities_proposed: [LINE-Messaging-API]
last_updated: 2026-04-30
-->

---

# Entity Register — Lotto Journal

Last updated: 2026-04-30

---

## Active Entities

| Entity     | Type       | Status | Since   | ADR     | Notes                                                                |
| ---------- | ---------- | ------ | ------- | ------- | -------------------------------------------------------------------- |
| Go         | tech       | active | 2026-04 | —       | Primary backend language; module: `lotto-journal/api`                |
| Fiber v2   | dependency | active | 2026-04 | —       | HTTP framework for Go (`github.com/gofiber/fiber/v2`)                |
| PostgreSQL | tech       | active | 2026-04 | —       | Primary database                                                     |
| GORM       | dependency | active | 2026-04 | —       | ORM for Go + PostgreSQL driver                                       |
| go-migrate | dependency | active | 2026-04 | —       | Database migrations (`apps/api/migrations/`)                         |
| pnpm       | tech       | active | 2026-04 | —       | Package manager; monorepo root                                       |
| Turborepo  | tech       | active | 2026-04 | —       | Monorepo build orchestration                                         |
| air        | dependency | active | 2026-04 | —       | Hot reload for Go during development                                 |
| Next.js    | tech       | active | 2026-04 | ADR-001 | `apps/web` — under review; may be deprecated if LINE pivot is chosen |

---

## Proposed Entities (pending ADR decision)

| Entity             | Type        | Status   | ADR     | Notes                                                 |
| ------------------ | ----------- | -------- | ------- | ----------------------------------------------------- |
| LINE Messaging API | integration | proposed | ADR-001 | Will become active if Option B (LINE pivot) is chosen |

---

## Deprecated / Removed Entities

_(none yet)_

| Entity | Type | Status | Since | Until | ADR | Replaced By |
| ------ | ---- | ------ | ----- | ----- | --- | ----------- |
| —      | —    | —      | —     | —     | —   | —           |

---

## Notes on User Identity

The existing `users` table in the database uses email/password + OAuth (`provider_service`
enum: google, facebook, apple, local) and OTP verification (`user_verifications`).

If ADR-001 resolves to Option B (LINE pivot), the following entities become candidates
for deprecation:

- email/password auth pattern
- OAuth provider integration
- `user_verifications` table (OTP flow)

And the following would be added:

- LINE user identity (`line_user_id`) as the primary user identifier
