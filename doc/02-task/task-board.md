<!-- AI-CONTEXT
active: none
blocked: none
done: T-000 T-001 T-005 T-008 T-004 T-007 T-006 T-002 T-010 T-011 T-012 T-013 T-014 T-016 T-018 T-015 T-017 T-019 T-003 T-023 T-022 T-021
future: T-009(liff-planning post-MVP), T-020(photo-ocr-openai-r2 post-MVP)
priority_next: none
src: v0.3
updated: 2026-06-29
-->

---

# Task Board — Lotto Journal

Last updated: 2026-06-29

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

---

## Env Map — Where Each Secret Lives

> Reference for T-013 and T-014.

| Secret / Env Var          | Value source                       | Stored in              | Why                                                              |
| ------------------------- | ---------------------------------- | ---------------------- | ---------------------------------------------------------------- |
| DB_DSN                    | Neon dashboard → Connection string | Fly.io secrets         | Runtime DB connection key used by current Go config loader        |
| LINE_CHANNEL_SECRET       | Production LINE channel            | Fly.io secrets         | Runtime webhook signature verification — production channel only |
| LINE_CHANNEL_ACCESS_TOKEN | Production LINE channel            | Fly.io secrets         | Runtime push/reply API calls — production channel only           |
| APP_ENV                   | Hardcoded value: production        | fly.toml [env] section | Non-secret; safe to commit                                       |
| FLY_API_TOKEN             | Fly.io dashboard → Access Tokens   | GitHub Actions secret  | Only needed by CI/CD to run flyctl deploy; never touches the app |

> **LINE channels:** Keep two separate channels under the same LINE provider.
> Dev channel webhook = Cloudflare tunnel URL (local). Production channel webhook = Fly.io URL.
> Each channel has its own secret + access token. Never mix them.
>
> **Neon itself** stores no secrets — it IS the database. Copy the connection string from the
> Neon dashboard and paste it as the DB_DSN Fly.io secret. Done.

---

## Future Tasks (post-MVP)

| ID    | Task                                     | Type  | Source Reference                                        | Priority | Status | Notes                                                                                                                                |
| ----- | ---------------------------------------- | ----- | ------------------------------------------------------- | -------- | ------ | ------------------------------------------------------------------------------------------------------------------------------------ |
| T-009 | Plan LIFF (LINE Front-end Framework) app | chore | doc/00-source/versions/v0.3/01-prd.md §8 (Out of Scope) | Low      | todo   | LIFF web app to complement the chatbot. Lives in apps/liff. Monorepo kept intentionally for this. Design when post-MVP phase begins. |
| T-020 | Photo ticket OCR via OpenAI + Cloudflare R2 (confirm-before-save) | feat  | doc/06-extensions/T-020-photo-ocr-openai-r2-proposal.md | Medium   | todo   | Post-MVP candidate. Single-image flow first; quantity-only confirm when OCR is correct, fallback to `numberxquantity` correction flow. Prioritized before T-009 by latest decision draft. |


---

## Blocked Tasks

None currently.

---

## Completed Tasks

| ID    | Task                                                         | Closed     | Evidence                                                                                                            |
| ----- | ------------------------------------------------------------ | ---------- | ------------------------------------------------------------------------------------------------------------------- |
| T-021 | Multi-language & Localization support (EN/TH)                | 2026-06-29 | Persisted user language preference in `users.language`; automatically detected profile language from LINE Profile API on follow; handled manual switcher commands (`ไทย`/`english`); localized text messages and win/loss push notifications; enabled one-tap Quick Replies navigation buttons. |
| T-022 | Implement win notification via LINE push message             | 2026-06-28 | Implemented NotificationService to check and group winnings per active user, format push texts, handle 3-attempt backoff retries, and write success/failed audits to notification_logs table. |
| T-023 | Swagger Documentation and Mockery Setup                      | 2026-06-27 | Extracted interface layers; configured Mockery v3 with `.mockery.yml`; added Swagger specs and Fiber v3 swaggo middleware; set up dev-only access; automated compilation reloads in `air`; updated package scripts/READMEs |
| T-003 | Design and implement cronjob: lottery result fetch + comparison flow | 2026-06-27 | Implemented LotteryClient for GLO APIs with retries/dedup; database-first draw resolver in DrawService with low-frequency background sync and emergency fallback; ResultService win-checking logic for L6 and N3; Bangkok time CronScheduler running startup sync, 3 AM sync, and 16:00 draw day checking; migration 000004 applied; all build/tests pass |
| T-019 | UX: loading indicator + personalized follow welcome [FOUND-IN-PASSING] | 2026-05-11 | Added `ShowLoadingAnimation` call in `handleMessage` (5s, clamped to LINE constraints); follow welcome now fetches profile display name via `GetProfile` and personalizes greeting; added tests for welcome message builder; `pnpm test:api` and `pnpm build` pass |
| T-017 | Improvement: atomic draws upsert via GORM clause.OnConflict | 2026-05-11 | Replaced `FirstOrCreate` in `internal/repository/draw_repository.go` with atomic `Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "draw_date"}}, DoUpdates: clause.AssignmentColumns([]string{"draw_date"})}).Create(&draw)`; eliminates SELECT+INSERT race; `pnpm test:api` and `pnpm build` pass |
| T-015 | GitHub Actions CI/CD pipeline                               | 2026-05-11 | `.github/workflows/deploy.yml` implemented; owner added repository secret `FLY_API_TOKEN`; first GitHub Actions run confirmed green (PR checks + deploy on `main`) |
| T-018 | Improve list command parsing for spaced/Unicode input [FOUND-IN-PASSING] | 2026-05-11 | Updated `isTicketListCmd` to normalize internal/Unicode spaces (including zero-width chars), so variants like `โ พย` and `โ\u200Bพย` match; added `internal/handler/line_handler_test.go`; `pnpm test:api` and `pnpm build` pass |
| T-016 | Bug: ticket parsing breaks when x has surrounding spaces     | 2026-05-11 | Fixed `spaceXRe` + replacement to `${1}x${2}`; added Unicode normalization for non-ASCII spaces and `×/ｘ`; added `internal/service/ticket_service_test.go` covering `144333 x2`, `122222 x 3`, and Unicode variants; `go test ./...` and `pnpm build` pass |
| T-014 | First production deploy to Fly.io + Neon wiring             | 2026-05-11 | Owner confirmed Fly deploy complete with 1 machine; Neon schema migrations applied; production LINE webhook wired; smoke-test ticket message stored in Neon DB |
| T-013 | Infra prep: Dockerfile + fly.toml + env secrets mapping     | 2026-05-09 | `Dockerfile`, `fly.toml`, `.dockerignore` created; env mapping fixed to use `DB_DSN`; `pnpm build` passes |
| T-012 | Feature: list upcoming draw tickets (summary on demand)      | 2026-05-08 | Build passes; keyword "โพย" routes to ListTickets; empty state handled; TicketRepository.List + TicketService.ListTickets + buildTicketListReply implemented |
| T-011 | Implement GET /health endpoint                               | 2026-05-08 | Build passes; DB ping via db.DB().Ping(); 200 ok / 503 degraded JSON response |
| T-010 | Add middleware: recover, requestid, enhanced logger, timeout | 2026-05-08 | Build passes; recover+requestid global; log upgraded (status+req_id); 25s timeout on /webhook; Fiber v2→v3 (v3.2.0) |
| T-002 | Design + implement LINE webhook handler                      | 2026-05-07 | Build passes; all event types handled; idempotency via webhook_events table                                         |
| T-000 | Documentation setup: doc/ structure created                  | 2026-04-30 | All required files created; bootstrap checklist passed                                                              |
| T-001 | Decide architecture pivot: web app vs LINE Messaging API     | 2026-04-30 | ADR-001 accepted — Option B chosen                                                                                  |
| T-005 | Write formal source docs (PRD v0.2)                          | 2026-04-30 | doc/00-source/versions/v0.2/01-prd.md created                                                                       |
| T-008 | Commit trunk/glo_result.json                                 | 2026-04-30 | Committed by owner; file now in repo                                                                                |
| T-004 | Design user identity model (line_user_id)                    | 2026-04-30 | Design doc + DBML updated; owner approved                                                                           |
| T-007 | Write migration 000002                                       | 2026-04-30 | SQL up/down written; Go model + code updated; build ✓                                                               |
| T-006 | Remove apps/web from monorepo                                | 2026-04-30 | apps/web deleted; turbo.json + pnpm-workspace cleaned                                                               |
