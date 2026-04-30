# Source Documents — Lotto Journal

## ⚠️ Status

**Formal source docs have not been created yet.**

The content below describes the project goal as understood at setup time (2026-04-30).
All planning documents reference `v0.1` until formal source docs are written.

Create `v0.2` after the architecture pivot decision (ADR-001) is finalized.

---

## Current Version

- **`v0.1`** — Setup baseline. No formal PRD.
- Files: `doc/00-source/versions/v0.1/00-setup-placeholder.md`

## Previous Versions

_(none)_

---

## What We Know (Setup Brief)

### Project Goal

A service that lets users record their lottery tickets and be automatically notified
when they win on official draw dates.

### User Interaction Layer — Under Decision (see ADR-001)

| Option            | Description                                                                                    |
| ----------------- | ---------------------------------------------------------------------------------------------- |
| **Current code**  | Web app (Next.js) + Go REST API + email/password auth                                          |
| **Planned pivot** | LINE Messaging API — users send ticket numbers via LINE chat; system notifies winners via LINE |

### Backend

- Language: **Go** (Fiber framework)
- Database: **PostgreSQL** (GORM ORM, go-migrate for migrations)
- Monorepo: **pnpm + Turborepo**

### Core Automated Flow

1. User submits lottery ticket numbers (6-digit or 3-digit)
2. Cronjob runs on every lottery draw day (1st and 16th of each month)
3. Cronjob fetches official results from the Thai Government Lottery Office API
4. System compares results against all user tickets in the database
5. Winners are notified automatically (via LINE if pivot is chosen)

### Existing Code (as of 2026-04-30)

| Component               | Status   | Notes                                                                  |
| ----------------------- | -------- | ---------------------------------------------------------------------- |
| `apps/api`              | Partial  | Fiber setup, basic signup handler (incomplete), DB connection          |
| `apps/web`              | Skeleton | Next.js boilerplate only — under review for removal                    |
| DB schema               | Complete | Full schema: users, tickets, draws, draw_results, user_winnings, files |
| Migrations              | Complete | `000001_init_schema.up.sql`                                            |
| `trunk/db_diagram.dbml` | Complete | Full ER diagram                                                        |

---

## Source Version Policy

Create a **new source version** when any of the following change:

- Overall MVP scope
- Primary user interaction channel (e.g., web → LINE)
- Core data model
- Core architecture pattern
- Key acceptance criteria

Use an **extension doc** (`doc/06-extensions/`) for:

- Additional details that don't change the "product truth"
- Feature ideas still under exploration
- Design decisions scoped to a single component

---

## Version History

| Version | Date       | Summary                                                               |
| ------- | ---------- | --------------------------------------------------------------------- |
| v0.1    | 2026-04-30 | Setup baseline — no formal PRD. Architecture pivot pending (ADR-001). |
