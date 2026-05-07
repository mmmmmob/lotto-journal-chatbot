# AI Decision Protocol — Lotto Journal

This document defines how AI should make decisions when facing ambiguous, conflicting,
or information-poor situations. Follow this protocol in every session.

Source: adapted from `core/11-ai-decision-protocol-template.md`

---

## 1. Three-Question Check (Before Any Work)

Answer these 3 questions before starting any task:

1. Do I know the expected outcome clearly?
2. Is this work consistent with the source docs?
3. If I make a mistake, can this work be reversed?

→ If all 3 are "yes": proceed.
→ If any is "no": use the Decision Tree below.

---

## 2. Decision Tree by Scenario

### Scenario A — Task is unclear; don't know what to do

```
Is there a similar completed task in the task-board?
  → Yes: follow that pattern; note "based on T-XXX"
  → No: is there a relevant section in source docs?
         → Yes: interpret from source docs; note which section
         → No: STOP → change task status to [BLOCKED]
                     → state clearly what information is needed
                     → update work-status
```

### Scenario B — Source docs conflict with existing code

```
Check which is newer (doc version vs git log)
→ Do not silently resolve — always surface the conflict
→ Create an extension doc recording the conflict
→ Update work-status to [NEEDS HUMAN DECISION: describe conflict]
→ Wait for human decision
```

### Scenario C — Work in progress exceeds original scope

```
→ Do only the part within the original scope
→ Create a new task T-XXX in task-board for the excess scope
→ If source reference exists: include it
→ If not: tag [NEEDS SOURCE VALIDATION]
→ Log the additional scope found in work-log
```

### Scenario D — Found a bug unrelated to the current task

```
→ Do not fix it silently
→ Create task T-XXX in task-board tagged [FOUND-IN-PASSING]
→ Log it in work-log
→ Return to the current task
```

### Scenario E — Two requirements conflict with each other

```
→ Do not silently pick one side
→ Create an extension doc explaining the conflict, referencing both requirements
→ Change task status to [BLOCKED: CONFLICT] with references to both points
→ Update work-status; wait for human decision
```

### Scenario F — No information in any document

```
→ Do not guess or invent data
→ Use placeholder: <NEEDS_CLARIFICATION: [describe what information is needed]>
→ Log the gap in work-status
→ Continue with other parts where information is sufficient
```

### Scenario G — Context window nearly full mid-task

```
→ Stop and save current progress to work-log before context is lost
→ Update work-status to reflect the current checkpoint
→ Change task status to [IN_PROGRESS: checkpoint saved — <summary of what was done>]
→ Next session reads work-status to resume
```

### Scenario H — Gap between task board and source docs

Triggered when: task references an old source version / in_progress task has no source
reference / source doc was updated but task board wasn't updated

```
→ Do not start implementation until the gap is clearly described
→ List all gaps found: task ID + what is missing or misaligned
→ Small gap (task missing reference): update task now; log in work-log
→ Large gap (source changed, multiple tasks affected):
     change work-status to [NEEDS HUMAN DECISION: source/task gap — <affected task list>]
→ Wait for human confirmation before continuing
```

### Scenario I — Found code without documentation

```
→ Do not modify the code immediately
→ Perform reverse-document per protocol in doc/06-extensions/
→ Create task [REVERSE-DOC] in task-board
→ Wait for human review before modifying
```

### Scenario J — Found [ENTITY:deprecated] or [ENTITY:superseded] tag

```
→ Do not use that entity without checking first
→ Open doc/07-decisions/entity-register.md
→ Check current status, replaced_by (if any), and related ADR
→ If entity is replaced: use the new entity; log in work-log
→ If uncertain: mark task [BLOCKED: deprecated entity]; wait for human decision
```

### Scenario K — Unsure where to store new information

```
Ask in order:
1. Is this an architectural decision?     → ADR (doc/07-decisions/)
2. Is this a new entity or status change? → Entity Register
3. Is this a cross-project pattern?       → cross-project-memory (ask user first)
4. Is this session progress or a detail?  → work-log-index
5. Is this a new or changed task?         → task-board

If multiple apply: store in all relevant places.
If none fit: save to work-log and note "location uncertain".
```

---

## 3. Escalation Levels

| Level                        | Situation                           | Action                                                  | Scenarios        |
| ---------------------------- | ----------------------------------- | ------------------------------------------------------- | ---------------- |
| **Level 1 — Log & Continue** | Low impact, reversible              | Log decision in work-log; continue                      | C, D, G, I, K    |
| **Level 2 — Block & Flag**   | Medium impact, uncertain            | Mark task blocked; update work-status; wait for human   | A, B, E, F, H, J |
| **Level 3 — Stop Session**   | High impact, irreversible, or risky | Stop immediately; document clearly why; do not continue | See below        |

**Always Level 3, regardless of scenario:**

- Risk of data loss
- Involves security or credentials
- Affects production environment
- Directly conflicts with a stated requirement

---

## 4. Responsibility Boundary

| Decision type                | AI may do                                 |
| ---------------------------- | ----------------------------------------- |
| Code style / formatting      | Decide independently                      |
| Refactoring (same logic)     | Decide independently                      |
| Adding a new dependency      | Propose in work-log; wait for approval    |
| Interpreting ambiguous scope | Propose + note risk                       |
| Changing architecture        | Create ADR draft; wait for human approval |
| Changing requirements        | **Strictly prohibited — humans only**     |
| Production operations        | **Strictly prohibited**                   |

---

## 5. Core Principle

> **"Do less, document more."**
> Stopping clearly and documenting why is always better than silently making a wrong decision.

- When unsure between Level 1 and Level 2: choose Level 2
- When unsure between Level 2 and Level 3: choose Level 3
- A clear placeholder is always better than a plausible guess

---

## 6. Status Tags

| Tag                               | Meaning                                                   |
| --------------------------------- | --------------------------------------------------------- |
| `[BLOCKED]`                       | Stopped; waiting for information or decision              |
| `[BLOCKED: CONFLICT]`             | Requirement conflict; waiting for resolution              |
| `[NEEDS HUMAN DECISION]`          | Human decision required before continuing                 |
| `[NEEDS SOURCE VALIDATION]`       | Task has no source doc support; needs verification        |
| `[FOUND-IN-PASSING]`              | Found while working on another task; not yet assigned     |
| `[IN_PROGRESS: checkpoint saved]` | Partially done; checkpoint recorded                       |
| `[REVERSE-DOC]`                   | Undocumented code found; waiting for human review         |
| `[ENTITY:deprecated:X]`           | Entity X is deprecated — check entity-register before use |
| `[ENTITY:superseded:X→Y]`         | Entity X replaced by Y                                    |
| `[ENTITY:proposed:X]`             | Entity X awaiting ADR                                     |
| `<NEEDS_CLARIFICATION: ...>`      | Placeholder for missing information                       |
