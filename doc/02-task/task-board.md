<!-- AI-CONTEXT
active: T-001(todo) T-002(todo) T-003(todo) T-004(todo) T-005(todo)
blocked: none
done: T-000
priority_next: T-001
src: v0.1
updated: 2026-04-30
-->

---

# Task Board — Lotto Journal

Last updated: 2026-04-30

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

| ID    | Task                                                          | Type  | Source Reference                 | Priority | Status | Notes                                                                                                     |
| ----- | ------------------------------------------------------------- | ----- | -------------------------------- | -------- | ------ | --------------------------------------------------------------------------------------------------------- |
| T-001 | Decide architecture pivot: web app vs LINE Messaging API      | chore | doc/00-source/README.md, ADR-001 | High     | todo   | Must be resolved before M1 — finalize ADR-001                                                             |
| T-002 | Design LINE Messaging API webhook handler                     | chore | [NEEDS SOURCE VALIDATION]        | High     | todo   | Depends on T-001 (LINE path); design LINE channel access, webhook signature verification, message parsing |
| T-003 | Design cronjob: lottery result fetch + comparison flow        | chore | [NEEDS SOURCE VALIDATION]        | High     | todo   | Confirm Thai Gov Lottery API format first; depends on T-001                                               |
| T-004 | Redesign user identity model (line_user_id vs email/password) | chore | [NEEDS SOURCE VALIDATION]        | High     | todo   | Current schema uses email/password; LINE pivot requires line_user_id — may need DB migration              |
| T-005 | Write formal source docs (PRD v0.2)                           | chore | doc/00-source/README.md          | High     | todo   | After T-001 resolves; creates the v0.2 source version                                                     |

---

## Blocked Tasks

| ID  | Task | Reason | Waiting On | Notes                      |
| --- | ---- | ------ | ---------- | -------------------------- |
| —   | —    | —      | —          | No blocked tasks currently |

---

## Completed Tasks

| ID    | Task                                        | Closed     | Evidence                                               |
| ----- | ------------------------------------------- | ---------- | ------------------------------------------------------ |
| T-000 | Documentation setup: doc/ structure created | 2026-04-30 | All required files created; bootstrap checklist passed |
