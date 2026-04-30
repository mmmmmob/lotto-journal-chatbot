# ADR-001: Architecture Pivot — Web App vs LINE Messaging API

**Date:** 2026-04-30
**Status:** Proposed
**Proposed by:** AI session — 2026-04-30
**Source Reference:** `doc/00-source/versions/v0.1/00-setup-placeholder.md`
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

### Option B — Pivot to LINE Messaging API

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

**⚠️ PENDING — Human decision required.**

This ADR is in Proposed status. The team must decide between Option A and Option B.

Once decided, update this ADR's status to Accepted and record the rationale below.

**Considerations to weigh:**

- Target user demographic (Thai users → LINE penetration is very high)
- Timeline and resource for frontend development (web app requires significant UI work)
- Acceptable platform dependency risk
- Current state of `apps/web` (it's a skeleton — reverting or removing is low-cost now)

---

## Consequences (to be filled in after decision)

**If Option A (keep web app):**

- Continue building `apps/web` as a Next.js SPA
- Complete the auth system (signup, login, OAuth)
- Add ticket submission and result viewing pages
- Design notification delivery (email or web push)
- Entity register: keep all current entities active; no LINE integration

**If Option B (LINE Messaging API pivot):**

- Remove or repurpose `apps/web`
- Redesign `users` table: replace email/password with `line_user_id`
- Deprecate `user_auth_methods` and `user_verifications` tables
- Add LINE webhook handler to `apps/api`
- Design message parsing: extract ticket numbers from LINE text messages
- Design notification sender: LINE push message to winner's LINE user ID
- Create new migration for user identity change
- Entity register: add LINE Messaging API (integration); deprecate Next.js, OAuth providers

---

## Review Trigger

Revisit this ADR if:

- LINE changes its API terms in a way that breaks the chosen design
- The user base grows beyond LINE's messaging limits
- A significant portion of target users are not on LINE
