<!-- AI-CONTEXT
active: T-003(todo)
blocked: none
done: T-000 T-001 T-005 T-008 T-004 T-007 T-006 T-002 T-010
future: T-009(liff-planning post-MVP)
priority_next: T-003
src: v0.2
updated: 2026-05-08
-->

---

# Task Board — Lotto Journal

Last updated: 2026-05-08 (session 6)

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

| ID    | Task                                                   | Type  | Source Reference                                  | Priority | Status | Notes                                                                                                                                                         |
| ----- | ------------------------------------------------------ | ----- | ------------------------------------------------- | -------- | ------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| T-003 | Design cronjob: lottery result fetch + comparison flow | chore | doc/00-source/versions/v0.2/01-prd.md §§3.3, §6.2 | High     | todo   | API: POST https://www.glo.or.th/api/lottery/getLatestLottery. Response format: see trunk/glo_result.json. Retry=5. Schedule configurable. Non-win push = YES. |

---

## Future Tasks (post-MVP)

| ID    | Task                                     | Type  | Source Reference                                        | Priority | Status | Notes                                                                                                                                  |
| ----- | ---------------------------------------- | ----- | ------------------------------------------------------- | -------- | ------ | -------------------------------------------------------------------------------------------------------------------------------------- |
| T-009 | Plan LIFF (LINE Front-end Framework) app | chore | doc/00-source/versions/v0.2/01-prd.md §8 (Out of Scope) | Low      | todo   | LIFF web app to complement the chatbot. Lives in `apps/liff`. Monorepo kept intentionally for this. Design when post-MVP phase begins. |

---

## Blocked Tasks

| ID  | Task | Reason | Waiting On | Notes                      |
| --- | ---- | ------ | ---------- | -------------------------- |
| —   | —    | —      | —          | No blocked tasks currently |

---

## Completed Tasks

| ID    | Task                                                         | Closed     | Evidence                                                                                                            |
| ----- | ------------------------------------------------------------ | ---------- | ------------------------------------------------------------------------------------------------------------------- |
| T-010 | Add middleware: recover, requestid, enhanced logger, timeout | 2026-05-08 | Build passes; recover+requestid global; log upgraded (status+req_id); 25s timeout on /webhook; Fiber v2→v3 (v3.2.0) |
| T-002 | Design + implement LINE webhook handler                      | 2026-05-07 | Build passes; all event types handled; idempotency via webhook_events table                                         |
| T-000 | Documentation setup: doc/ structure created                  | 2026-04-30 | All required files created; bootstrap checklist passed                                                              |
| T-001 | Decide architecture pivot: web app vs LINE Messaging API     | 2026-04-30 | ADR-001 accepted — Option B chosen                                                                                  |
| T-005 | Write formal source docs (PRD v0.2)                          | 2026-04-30 | doc/00-source/versions/v0.2/01-prd.md created                                                                       |
| T-008 | Commit trunk/glo_result.json                                 | 2026-04-30 | Committed by owner; file now in repo                                                                                |
| T-004 | Design user identity model (line_user_id)                    | 2026-04-30 | Design doc + DBML updated; owner approved                                                                           |
| T-007 | Write migration 000002                                       | 2026-04-30 | SQL up/down written; Go model + code updated; build ✓                                                               |
| T-006 | Remove apps/web from monorepo                                | 2026-04-30 | apps/web deleted; turbo.json + pnpm-workspace cleaned                                                               |
