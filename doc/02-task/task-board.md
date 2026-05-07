<!-- AI-CONTEXT
active: T-002(todo) T-003(todo) T-006(todo)
blocked: none
done: T-000 T-001 T-005 T-008 T-004 T-007
priority_next: T-006 (independent) or T-002 (design)
src: v0.2
updated: 2026-04-30
-->

---

# Task Board — Lotto Journal

Last updated: 2026-04-30 (session 4)

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

| ID        | Task                                                                                   | Type  | Source Reference                                                    | Priority | Status | Notes                                                                                                                                                                                    |
| --------- | -------------------------------------------------------------------------------------- | ----- | ------------------------------------------------------------------- | -------- | ------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| T-002     | Design LINE Messaging API webhook handler                                              | chore | doc/00-source/versions/v0.2/01-prd.md §3.1–3.2, §6.1                | High     | todo   | Input format defined: plain number(s), comma/space separated, optional xN quantity. Idempotency key = webhookEventId. follow=create user, unfollow=mark status inactive.                 |
| T-003     | Design cronjob: lottery result fetch + comparison flow                                 | chore | doc/00-source/versions/v0.2/01-prd.md §3.3, §6.2                    | High     | todo   | API: POST https://www.glo.or.th/api/lottery/getLatestLottery. Response format: see trunk/glo_result.json (T-008 must be done first). Retry=5. Schedule configurable. Non-win push = YES. |
| T-004     | Design user identity model: users table with line_user_id                              | chore | doc/00-source/versions/v0.2/01-prd.md §5.1                          | High     | done   | Design doc: doc/06-extensions/T-004-migration-002-design.md. DBML updated. Owner approved.                                                                                               |
| T-006     | Remove apps/web (Next.js) from monorepo                                                | chore | doc/07-decisions/ADR-001-line-messaging-pivot.md                    | Medium   | todo   | Delete apps/web directory; update pnpm-workspace.yaml and turbo.json                                                                                                                     |
| T-007     | Write migration 000002: redesign users table + drop unused tables/enums + rename enums | chore | doc/00-source/versions/v0.2/01-prd.md §5.1, §5.3 + T-004 design doc | High     | done   | 000002_line_identity.up/down.sql written. User model updated. auth_handler/service/dto removed. main.go cleaned up. Build passes.                                                        |
| ~~T-008~~ | ~~Commit trunk/glo_result.json~~                                                       | chore | —                                                                   | —        | done   | Committed by owner before session 4.                                                                                                                                                     |

---

## Blocked Tasks

| ID  | Task | Reason | Waiting On | Notes                      |
| --- | ---- | ------ | ---------- | -------------------------- |
| —   | —    | —      | —          | No blocked tasks currently |

---

## Completed Tasks

| ID    | Task                                                     | Closed     | Evidence                                               |
| ----- | -------------------------------------------------------- | ---------- | ------------------------------------------------------ |
| T-000 | Documentation setup: doc/ structure created              | 2026-04-30 | All required files created; bootstrap checklist passed |
| T-001 | Decide architecture pivot: web app vs LINE Messaging API | 2026-04-30 | ADR-001 accepted — Option B chosen                     |
| T-005 | Write formal source docs (PRD v0.2)                      | 2026-04-30 | doc/00-source/versions/v0.2/01-prd.md created          |
| T-008 | Commit trunk/glo_result.json                             | 2026-04-30 | Committed by owner; file now in repo                   |
| T-004 | Design user identity model (line_user_id)                | 2026-04-30 | Design doc + DBML updated; owner approved              |
| T-007 | Write migration 000002                                   | 2026-04-30 | SQL up/down written; Go model + code updated; build ✓  |
