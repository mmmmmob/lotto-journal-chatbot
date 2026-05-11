# Versioning Policy (SemVer)

This project uses **Semantic Versioning** with the format `MAJOR.MINOR.PATCH`.

## Scope

- The primary release version is tracked at the repository root in `package.json` (`version`).
- Git tags are the release marker of record.

## Rules

- **PATCH**: backward-compatible fixes and internal improvements.
  - Examples: parser bugfixes, logging improvements, race-condition fixes.
- **MINOR**: backward-compatible new functionality.
  - Examples: new bot commands, new cronjob workflows, new endpoints.
- **MAJOR**: breaking changes.
  - Examples: API contract changes, behavior changes requiring migration, incompatible schema/protocol changes.

## Pre-1.0 policy

While version is `0.x.y`:

- Breaking changes may still occur.
- Treat **MINOR** bumps as potentially impactful and communicate clearly in release notes.

## Tag format

- Use annotated tags: `vMAJOR.MINOR.PATCH` (example: `v0.2.0`).

## Source of truth for release notes

- `doc/03-log/work-log-index.md` (session-level evidence)
- `doc/02-task/task-board.md` (task completion evidence)
