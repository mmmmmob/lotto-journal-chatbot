# Coding Standards — Lotto Journal

## Working Principle

Build the primary user flow before secondary features.
Before starting any task, answer:

- Who is the user for this flow?
- Which endpoint or message type handles it?
- What comes before and after in the user journey?
- Where does this task connect to the system's core workflow?

---

## Project Tech Stack

| Layer            | Technology                  | Notes                                   |
| ---------------- | --------------------------- | --------------------------------------- |
| Backend API      | Go + Fiber v2               | `apps/api`                              |
| ORM              | GORM                        | PostgreSQL driver                       |
| Migrations       | go-migrate                  | `apps/api/migrations/`                  |
| Database         | PostgreSQL                  |                                         |
| Frontend         | Next.js (TypeScript)        | `apps/web` — under review (see ADR-001) |
| Monorepo         | pnpm + Turborepo            | root `pnpm-workspace.yaml`              |
| Hot reload (dev) | air                         | `apps/api/.air.toml`                    |
| LINE integration | LINE Messaging API SDK (Go) | Proposed — pending ADR-001              |

---

## Naming

- Use names that clearly describe intent
- Functions: describe the behavior (e.g., `fetchDrawResults`, `sendWinNotification`)
- Variables: describe the data they hold (e.g., `userTickets`, `drawResult`)
- Go types/structs: PascalCase; match the domain entity name closely
- Database columns: snake_case (matches existing schema)
- Files: snake_case for Go (`auth_handler.go`), kebab-case for TS (`task-board.md`)

---

## File Size and Structure

- Go files: no more than **500 lines** (excluding comments and blank lines)
- If a file grows large, split by concern (e.g., separate handler, service, repository)
- Follow the existing project structure:
  - `internal/handler/` — HTTP / webhook handlers
  - `internal/service/` — business logic
  - `internal/repository/` — database operations
  - `internal/models/` — GORM models
  - `internal/dto/` — request/response data transfer objects

---

## Comments

- Write comments only when needed
- Comments should explain **why** (reason, constraint, invariant, tradeoff), not **what**
- Mark deferred cleanup with the compliance tag format:
  `// REFACTOR-PENDING[C-XX]: <description> — T-XXX`

---

## Go-Specific Rules

- Use GORM's parameterized queries — never concatenate raw SQL strings from user input (C-11 / SQL injection)
- All HTTP handlers must validate input; use DTO structs with struct tags
- Secrets (DB password, LINE channel secret, etc.) must come from environment variables — never hardcoded
- Return meaningful HTTP status codes and consistent JSON error shapes
- Goroutines used for the cronjob must handle panics gracefully

---

## LINE Messaging API Rules (once ADR-001 is accepted)

- Always verify webhook request signature using the channel secret before processing
- Idempotently handle webhook events — LINE may re-deliver events
- Never log full LINE event payloads that contain user profile data (C-11 / sensitive data in logs)
- Use the official LINE Bot SDK for Go

---

## Database and Migration Rules

- All schema changes go through `go-migrate` — never alter tables manually
- Migration files: `NNNNNN_<description>.up.sql` and `NNNNNN_<description>.down.sql`
- Always test `down` migration before merging
- Do not change existing migration files — create new ones for corrections

---

## Change Discipline

- Keep changes within the scope of the current task
- If changing workflow, architecture, database schema, or API contracts:
  update documentation in the same session
- Do not fix bugs unrelated to the current task without creating a separate task
  `[FOUND-IN-PASSING]`

---

## Verification

- Before closing any task, there must be at least one validation:
  test pass, manual smoke test, or documented check
- If something was not tested, state explicitly what was not verified

---

## Git Policy

- Never commit `.env`, credentials, runtime logs, export outputs, or local-only files
- Daily logs (`doc/03-log/YYYY/MM/`) are local only — not committed
- See `.gitignore` for full exclusion list
- `_template/` must remain in `.gitignore`
