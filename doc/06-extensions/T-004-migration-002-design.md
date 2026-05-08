# Migration 000002 — Design Document

**Task:** T-004
**Status:** done (migration 000002 implemented and verified — 2026-04-30)
**Source reference:** `doc/00-source/versions/v0.2/01-prd.md` §5.1, §5.3
**Ground truth (pre-migration state):** `apps/api/migrations/000001_init_schema.up.sql`
**Target schema:** `trunk/db_diagram.dbml` (post-migration 000002)
**Date:** 2026-04-30

---

## Purpose

Migration 000002 implements the ADR-001 pivot from email/password authentication to
LINE Messaging API user identity. It redesigns the `users` table, removes auth infrastructure
that no longer applies, and renames enums to match the v0.2 naming convention.

---

## Change Summary

### 1. `users` table — MODIFIED

| Column          | Action | New Definition                    | Notes                          |
| --------------- | ------ | --------------------------------- | ------------------------------ |
| `id`            | KEEP   | uuid PK                           | unchanged                      |
| `line_user_id`  | ADD    | `varchar UNIQUE NOT NULL`         | LINE platform user ID          |
| `status`        | KEEP   | `account_status DEFAULT 'active'` | enum values change — see below |
| `created_at`    | KEEP   | `timestamp DEFAULT now()`         | unchanged                      |
| `updated_at`    | KEEP   | `timestamp DEFAULT now()`         | unchanged                      |
| `username`      | DROP   | —                                 | no username in LINE flow       |
| `email`         | DROP   | —                                 | no email in LINE flow          |
| `password_hash` | DROP   | —                                 | no password in LINE flow       |

### 2. `account_status` enum — RECREATED

PostgreSQL does not support `DROP VALUE` on an enum. Since there is no production data,
the migration will drop and recreate this type.

| Value       | Action | Notes                                |
| ----------- | ------ | ------------------------------------ |
| `active`    | KEEP   | user is following the LINE OA        |
| `inactive`  | ADD    | user unfollowed the LINE OA          |
| `pending`   | DROP   | no longer applicable (no email auth) |
| `suspended` | KEEP   | admin action — misbehaving or bot    |

**Migration approach (no production data):**

1. `ALTER TABLE users ALTER COLUMN status DROP DEFAULT`
2. `ALTER TABLE users ALTER COLUMN status TYPE text`
3. `DROP TYPE account_status`
4. `CREATE TYPE account_status AS ENUM ('active', 'inactive', 'suspended')`
5. `ALTER TABLE users ALTER COLUMN status TYPE account_status USING status::account_status`
6. `ALTER TABLE users ALTER COLUMN status SET DEFAULT 'active'`

### 3. `lottery_type` enum — RENAME VALUE

| Old Value | New Value | Notes      |
| --------- | --------- | ---------- |
| `N6`      | `L6`      | naming fix |
| `N3`      | `N3`      | unchanged  |

SQL: `ALTER TYPE lottery_type RENAME VALUE 'N6' TO 'L6';`

### 4. `prize_type` enum — RENAME VALUES (9 values)

All `n6_*` prefix values renamed to `l6_*` to match the new `L6` naming.

| Old Value       | New Value       |
| --------------- | --------------- |
| `n6_first`      | `l6_first`      |
| `n6_second`     | `l6_second`     |
| `n6_third`      | `l6_third`      |
| `n6_fourth`     | `l6_fourth`     |
| `n6_fifth`      | `l6_fifth`      |
| `n6_last2`      | `l6_last2`      |
| `n6_last3f`     | `l6_last3f`     |
| `n6_last3b`     | `l6_last3b`     |
| `n6_near_first` | `l6_near_first` |

`n3_*` values are unchanged.

### 5. Tables — DROPPED

These tables supported email/OAuth authentication, which is no longer used.

| Table                | Reason                                |
| -------------------- | ------------------------------------- |
| `user_auth_methods`  | OAuth providers no longer used        |
| `user_verifications` | Email/OTP verification no longer used |
| `user_profiles`      | Decided: drop in migration 000002     |

Drop order (respect FK constraints):

1. `user_auth_methods` (FK → users.id)
2. `user_verifications` (FK → users.id)
3. `user_profiles` (FK → users.id, files.id)

### 6. Enums — DROPPED

| Enum                | Reason                                      |
| ------------------- | ------------------------------------------- |
| `provider_service`  | Used only by `user_auth_methods` (dropped)  |
| `verification_type` | Used only by `user_verifications` (dropped) |

Drop after their dependent tables are dropped.

---

## Found-in-Passing: `user_winnings.user_id` Missing

**Issue:** Migration 000001 SQL created `user_winnings` without a `user_id` column.
The `trunk/db_diagram.dbml` (pre-edit) showed `user_id`, and PRD §3.3 references it
as part of the win record `(user_id, ticket_id, draw_result_id, prize_money)`.

**Resolution:** Add `user_id` to `user_winnings` in migration 000002.

```
ALTER TABLE user_winnings
  ADD COLUMN user_id uuid REFERENCES users(id) ON DELETE CASCADE;
```

Note: since no production data exists, this column can be added as non-nullable in
practice, but adding it nullable first then constraining is safer idiomatically.
For this project with no data, add as direct FK reference.

---

## Tables Unchanged

| Table           | Notes                                           |
| --------------- | ----------------------------------------------- |
| `draws`         | No changes                                      |
| `tickets`       | No changes; `owner_id` still refs `users.id`    |
| `draw_results`  | No changes                                      |
| `files`         | Kept for future photo upload feature (post-MVP) |
| `user_winnings` | `user_id` column ADDED (see Found-in-Passing)   |

---

## Down Migration Notes

The down migration must restore the schema to the exact 000001 state.
Key reversal steps:

1. Remove `user_id` from `user_winnings`
2. Restore `users` table columns (`username`, `email`, `password_hash`); remove `line_user_id`
3. Restore `account_status` enum to (`active`, `pending`, `suspended`)
4. Rename `lottery_type` value `L6` → `N6`
5. Rename 9 `prize_type` values `l6_*` → `n6_*`
6. Recreate dropped tables: `user_auth_methods`, `user_verifications`, `user_profiles`
7. Recreate dropped enums: `provider_service`, `verification_type`
8. Restore indexes dropped as part of table recreation

---

## Go Model Changes (for T-007 implementation)

The `apps/api/internal/models/user.go` GORM model must be updated:

| Field          | Action | Notes                                             |
| -------------- | ------ | ------------------------------------------------- |
| `LineUserID`   | ADD    | `gorm:"type:varchar;uniqueIndex;not null"`        |
| `Username`     | DROP   | —                                                 |
| `Email`        | DROP   | —                                                 |
| `PasswordHash` | DROP   | —                                                 |
| `DeletedAt`    | DROP   | Not in SQL; `status=inactive` handles soft-delete |
| `Status`       | KEEP   | update comment to reflect new enum                |

A new `UserWinning` model will also need `UserID` added.

---

## Validation Checklist (before T-007 starts)

- [x] New `users` schema defined (line_user_id, status, timestamps)
- [x] `account_status` enum values confirmed: `active`, `inactive`, `suspended`
- [x] Tables to drop confirmed: user_auth_methods, user_verifications, user_profiles
- [x] Enums to drop confirmed: provider_service, verification_type
- [x] Enum renames confirmed: N6→L6, n6*\*→l6*\* (×9)
- [x] Found-in-passing fix confirmed: user_winnings.user_id to be added
- [x] Down migration approach noted
- [x] DBML updated to reflect post-000002 target state
