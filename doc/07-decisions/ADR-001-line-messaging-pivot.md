# ADR-001: Architecture Pivot — Web App vs LINE Messaging API

**Date:** 2026-04-30
**Status:** Accepted
**Proposed by:** AI session — 2026-04-30
**Accepted by:** Project owner — 2026-04-30
**Source Reference:** `doc/00-source/versions/v0.2/01-prd.md`
**Related Tasks:** T-001

---

## Context

Lotto Journal was initially built as a monorepo web application:

- `apps/web` — Next.js frontend (currently a skeleton with no business logic)
- `apps/api` — Go + Fiber REST API

The existing backend has:

- A partial user authentication system (email/password + OAuth providers)
- A complete PostgreSQL schema for tickets, draws, draw_results, user_winnings
- Only one implemented endpoint: `POST /signup` (incomplete — uses hardcoded password)

The intended product goal is: users record lottery ticket numbers and are automatically
notified when they win. The team is evaluating whether the user interaction layer should
remain a web app or be replaced with LINE Messaging API.

This decision is architecturally significant because it affects:

- User identity model (email/password vs LINE user ID)
- How tickets are submitted (web form vs LINE chat message)
- How winners are notified (web UI / email vs LINE push message)
- Whether `apps/web` is kept, repurposed, or removed
- Whether the `users` table and auth system need to be redesigned

---

## Options Considered

### Option A — Keep Web App

The product continues as a web application. Users log in via browser, submit tickets
through a web form, and check results on the site (or receive email notifications).

**Pros:**

- Current code direction (no pivot cost)
- Richer UI possibilities (history views, search, filters)
- No dependency on LINE platform

**Cons:**

- Requires significant frontend build-out (`apps/web` is currently a skeleton)
- Higher friction for users — requires opening a browser/app
- Push notifications on web are less reliable than LINE messages for Thai users
- More user acquisition friction (registration, password management)

---

### Option B — Pivot to LINE Messaging API ✅ CHOSEN

Replace the web frontend with LINE as the primary user interaction channel.
Users interact entirely through the LINE chat interface:

- Send ticket numbers as LINE messages
- Receive winner notifications as LINE push messages

The `apps/web` Next.js app would be removed or repurposed for admin/ops only.
The Go backend would handle LINE webhook events instead of REST API calls from a web client.

A cronjob would still run on draw days to fetch results and trigger comparisons.

**Pros:**

- LINE has very high penetration among Thai users — no install friction
- Push messages via LINE are highly visible and reliable
- Simpler UX: no registration form, no password (LINE identity is pre-authenticated)
- Reduces frontend scope significantly — no web UI to build for end users
- The lottery data model (tickets, draws, draw_results, user_winnings) is already solid
  and survives the pivot

**Cons:**

- LINE platform dependency — subject to LINE API policy changes and rate limits
- Webhook-driven architecture requires careful idempotency design (LINE may re-deliver events)
- Current user identity model (email/password + OAuth) must be redesigned around `line_user_id`
- The `users` table needs to change: remove email/password/OAuth columns, add `line_user_id`
- `user_auth_methods` and `user_verifications` tables may become obsolete
- Requires LINE Developer account + LINE Official Account setup

---

## Decision

**Option B — LINE Messaging API pivot** was chosen.

**Rationale:**

- LINE penetration in Thailand is extremely high; the target users (Thai lottery players)
  are already on LINE, eliminating registration and installation friction entirely
- The web app is currently a skeleton — the cost of removing it is minimal
- The core lottery data model is already well-designed and survives the pivot unchanged
- LINE push messages are a more reliable and visible notification channel for this use case

---

## Consequences

**Immediate actions required:**

1. **Remove `apps/web`** (T-006) — Next.js app is no longer the user-facing product
2. **Redesign user identity** (T-004, T-007) — replace `users` table email/password schema
   with `line_user_id`; drop `user_auth_methods` and `user_verifications` tables
3. **Implement LINE webhook handler** (T-002, T-008) — handle webhook events from LINE platform
4. **Implement cronjob** (T-003) — fetch results + compare + push notifications

**What stays unchanged:**

- Go + Fiber backend framework
- PostgreSQL database
- `draws`, `tickets`, `draw_results`, `user_winnings`, `files` tables (unchanged)
- `lottery_type` and `prize_type` enums

**What is deprecated / removed:**

- `apps/web` (Next.js) — to be removed
- `user_auth_methods` table — removed in migration 000002
- `user_verifications` table — removed in migration 000002
- `provider_service` enum — removed in migration 000002
- `verification_type` enum — removed in migration 000002

---

## Review Trigger

Revisit this ADR if:

- LINE changes its API terms in a way that breaks the chosen design
- The user base grows beyond LINE's messaging limits
- A significant portion of target users are not on LINE
