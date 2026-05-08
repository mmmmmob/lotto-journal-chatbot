<!-- AI-CONTEXT
active: T-003(todo) T-012(todo) T-013(todo) T-014(todo,blocked:T-013) T-015(todo,blocked:T-014) T-016(todo) T-017(todo)
blocked: T-014(needs T-013) T-015(needs T-014)
done: T-000 T-001 T-005 T-008 T-004 T-007 T-006 T-002 T-010 T-011
future: T-009(liff-planning post-MVP)
priority_next: T-003
src: v0.2
updated: 2026-05-08
-->

---

# Task Board — Lotto Journal

Last updated: 2026-05-08 (session 7)

## Rules

- Every task must have a source reference
- Status values: `todo` `design_validate` `in_progress` `review` `done` `blocked`
- If a task changes scope, create an extension doc or new source version
- Tasks found unplanned: tag `[FOUND-IN-PASSING]`

## Definition of Ready (before moving to `in_progress`)

- [ ] Clear source reference (`doc/00-source/...` or ADR)
- [ ] Scope defined: what is included and what is not
- [ ] No unresolved dependencies
- [ ] design_validate passed (or confirmed "scope clear, no changes needed")

## Definition of Done (before moving to `review`)

- [ ] Work matches the scope defined at design_validate
- [ ] Compliance scan passed (no Level 1 violations pending)
- [ ] Validation evidence exists (test pass / manual check / screenshot)
- [ ] `work-status.md` and `work-log-index.md` updated

---

## Current Tasks

| ID    | Task                                                     | Type        | Source Reference                                  | Priority | Status | Notes                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       |
| ----- | -------------------------------------------------------- | ----------- | ------------------------------------------------- | -------- | ------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| T-003 | Design cronjob: lottery result fetch + comparison flow   | chore       | doc/00-source/versions/v0.2/01-prd.md §§3.3, §6.2 | High     | todo   | API: POST https://www.glo.or.th/api/lottery/getLatestLottery. Response format: see trunk/glo_result.json. Retry=5. Schedule configurable. Non-win push = YES.                                                                                                                                                                                                                                                                                                                                                               |
| T-012 | Feature: list upcoming draw tickets (summary on demand)  | feature     | doc/00-source/versions/v0.2/01-prd.md §3          | Medium   | todo   | User sends keyword → bot replies with all tickets for the upcoming draw (number, type, qty, total count). Add ListByOwnerAndDraw to ticket repo, ListForUpcomingDraw to ticket service, keyword routing in line_handler.go. Keyword TBD (e.g. "ดูตั๋ว" / "สรุป").                                                                                                                                                                                                                                                           |
| T-013 | Infra prep: Dockerfile + fly.toml + env secrets mapping  | chore/infra | —                                                 | High     | todo   | Multi-stage Dockerfile for Go API. fly.toml (app name, region, health check path). Document and apply env var assignment: Fly.io secrets vs GitHub Actions secrets (see Env Map section below). Prerequisite for T-014.                                                                                                                                                                                                                                                                                                     |
| T-014 | First production deploy to Fly.io + Neon wiring          | chore/infra | —                                                 | High     | todo   | Blocked by T-013. Steps: fly launch, fly secrets set (DATABASE_URL + LINE secrets), run migrations against Neon, update LINE Developer Console webhook URL to Fly.io app URL, verify GET /health, smoke-test bot end-to-end.                                                                                                                                                                                                                                                                                                |
| T-015 | GitHub Actions CI/CD pipeline                            | chore/infra | —                                                 | Medium   | todo   | Blocked by T-014. .github/workflows/deploy.yml: build + go vet + go test on every PR; auto-deploy to Fly.io on push to main via flyctl deploy --remote-only. Only one GitHub Actions secret needed: FLY_API_TOKEN.                                                                                                                                                                                                                                                                                                          |
| T-016 | Bug: ticket parsing breaks when x has surrounding spaces | bug         | apps/api/internal/service/ticket_service.go       | Medium   | todo   | Two broken cases observed in production test: (1) "144333 x2" — spaceXRe normalisation silently fails; "2" appears as lone invalid token instead of qty. (2) "122222 x 3" — space on BOTH sides of x; number saves correctly but "3" becomes invalid token. Fix: change spaceXRe to allow optional whitespace after x — `(\d+)\s+x\s*(\d+)`. Also investigate whether LINE sends non-ASCII space or non-ASCII x (×) character, which would defeat the current regex entirely. Add unit tests for both cases before closing. |
| T-017 | Improvement: atomic draws upsert via GORM clause.OnConflict | improvement | doc/01-plan/work-status.md (Risks and Notes)      | Low      | todo   | Replace `FirstOrCreate` in `repository/draw_repository.go` with `db.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "draw_date"}}, DoUpdates: clause.Assignments(map[string]interface{}{"draw_date": gorm.Expr("draws.draw_date")})}).Create(&draw)`. Eliminates the SELECT+INSERT race condition. No raw SQL needed. Non-blocking for MVP (≤100 users) — do before scaling. |

---

## Env Map — Where Each Secret Lives

> Reference for T-013 and T-014.

| Secret / Env Var          | Value source                       | Stored in              | Why                                                              |
| ------------------------- | ---------------------------------- | ---------------------- | ---------------------------------------------------------------- |
| DATABASE_URL              | Neon dashboard → Connection string | Fly.io secrets         | Runtime DB connection; never in code or GitHub                   |
| LINE_CHANNEL_SECRET       | LINE Developer Console             | Fly.io secrets         | Runtime webhook signature verification                           |
| LINE_CHANNEL_ACCESS_TOKEN | LINE Developer Console             | Fly.io secrets         | Runtime push/reply API calls                                     |
| APP_ENV                   | Hardcoded value: production        | fly.toml [env] section | Non-secret; safe to commit                                       |
| FLY_API_TOKEN             | Fly.io dashboard → Access Tokens   | GitHub Actions secret  | Only needed by CI/CD to run flyctl deploy; never touches the app |

> **Neon itself** stores no secrets — it IS the database. Copy the connection string from the
> Neon dashboard and paste it as the DATABASE_URL Fly.io secret. Done.

---

## Future Tasks (post-MVP)

| ID    | Task                                     | Type  | Source Reference                                        | Priority | Status | Notes                                                                                                                                |
| ----- | ---------------------------------------- | ----- | ------------------------------------------------------- | -------- | ------ | ------------------------------------------------------------------------------------------------------------------------------------ |
| T-009 | Plan LIFF (LINE Front-end Framework) app | chore | doc/00-source/versions/v0.2/01-prd.md §8 (Out of Scope) | Low      | todo   | LIFF web app to complement the chatbot. Lives in apps/liff. Monorepo kept intentionally for this. Design when post-MVP phase begins. |

---

## Blocked Tasks

| ID    | Task                                            | Reason                          | Waiting On | Notes                                   |
| ----- | ----------------------------------------------- | ------------------------------- | ---------- | --------------------------------------- |
| T-014 | First production deploy to Fly.io + Neon wiring | Dockerfile + fly.toml not ready | T-013      | Unblock by completing T-013             |
| T-015 | GitHub Actions CI/CD pipeline                   | No live Fly.io app yet          | T-014      | Needs FLY_API_TOKEN from live app first |

---

## Completed Tasks

| ID    | Task                                                         | Closed     | Evidence                                                                                                            |
| ----- | ------------------------------------------------------------ | ---------- | ------------------------------------------------------------------------------------------------------------------- |
| T-011 | Implement GET /health endpoint                               | 2026-05-08 | Build passes; DB ping via db.DB().Ping(); 200 ok / 503 degraded JSON response                                       |
| T-010 | Add middleware: recover, requestid, enhanced logger, timeout | 2026-05-08 | Build passes; recover+requestid global; log upgraded (status+req_id); 25s timeout on /webhook; Fiber v2→v3 (v3.2.0) |
| T-002 | Design + implement LINE webhook handler                      | 2026-05-07 | Build passes; all event types handled; idempotency via webhook_events table                                         |
| T-000 | Documentation setup: doc/ structure created                  | 2026-04-30 | All required files created; bootstrap checklist passed                                                              |
| T-001 | Decide architecture pivot: web app vs LINE Messaging API     | 2026-04-30 | ADR-001 accepted — Option B chosen                                                                                  |
| T-005 | Write formal source docs (PRD v0.2)                          | 2026-04-30 | doc/00-source/versions/v0.2/01-prd.md created                                                                       |
| T-008 | Commit trunk/glo_result.json                                 | 2026-04-30 | Committed by owner; file now in repo                                                                                |
| T-004 | Design user identity model (line_user_id)                    | 2026-04-30 | Design doc + DBML updated; owner approved                                                                           |
| T-007 | Write migration 000002                                       | 2026-04-30 | SQL up/down written; Go model + code updated; build ✓                                                               |
| T-006 | Remove apps/web from monorepo                                | 2026-04-30 | apps/web deleted; turbo.json + pnpm-workspace cleaned                                                               |
