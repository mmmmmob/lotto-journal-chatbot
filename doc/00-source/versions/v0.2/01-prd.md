# Product Requirements Document — Lotto Journal (LINE-based)

**Version:** v0.2
**Date:** 2026-04-30
**Updated:** 2026-05-12
**Status:** Approved — microcopy TBD (welcome message, push message copy)
**Based on:** ADR-001 (Option B — LINE Messaging API pivot, Accepted 2026-04-30)
**Supersedes:** v0.1 setup placeholder

---

## 1. Product Goal

Lotto Journal is a **LINE chatbot service** that lets Thai lottery players:

1. Record the lottery ticket numbers they own by sending them through LINE chat
2. Be automatically notified via LINE when any of their tickets win a prize

Users interact entirely through LINE. There is no web app or mobile app to install.

---

## 2. Target Users

Thai lottery players who:

- Already have a LINE account (LINE is the dominant messaging platform in Thailand)
- Purchase Thai Government Lottery (สลากกินแบ่งรัฐบาล) tickets regularly
- Want a convenient way to track tickets and get notified of winnings without manually
  checking results

---

## 3. User Flows

### 3.1 User Onboarding (First Interaction)

**Trigger:** User adds the Lotto Journal LINE Official Account as a friend and sends any message.

**Steps:**

1. LINE delivers a `follow` event or a `message` event to the backend webhook
2. Backend extracts `source.userId` (LINE user ID) from the event
3. Backend checks if this `line_user_id` exists in the `users` table
4. If new user → create a `users` record with:
   - `line_user_id` = value from event
   - `status` = `active`
5. Backend sends a **welcome reply message** explaining how to use the service

**Message format for welcome reply:** `<NEEDS_CLARIFICATION: exact welcome message copy>` (microcopy to be finalized later)

---

### 3.2 Ticket Submission

**Trigger:** User sends a message to the LINE chatbot with their lottery ticket number(s).

**Message input format:**

- Plain number — user sends `123456` or `456`
- User can send multiple numbers in one message, separated by spaces or commas — e.g., `123456, 654321` or `456 789`
- Each input can indicate the ticket quantity by appending `xN` — e.g., `123456x2` means 2 tickets of number `123456` (accept with or without spaces before `x`)
- Command (implemented): user can send `โพย` to list tickets already recorded for the upcoming draw
- Other options will be implemented later as features expand (e.g., send ticket photos or calculate sum of bought tickets for financial planning), but MVP relies on users sending plain text messages with ticket numbers and optional quantity

**Steps:**

1. User sends a message containing a lottery number
2. Backend receives webhook `message` event
3. Backend parses the text to extract the lottery number
4. System determines ticket type from digit count:
   - 6-digit number → `L6` (6-digit lottery)
   - 3-digit number → `N3` (3-digit lottery)
5. System identifies or creates the relevant `draws` record for the upcoming draw date
6. System creates a `tickets` record:
   - `owner_id` = resolved `users.id` from `line_user_id`
   - `draw_id` = resolved `draws.id`
   - `type` = `L6` or `N3`
   - `number` = the parsed number
   - `quantity` = 1 by default, or set from `xN` when provided (e.g. `123456x2` → `quantity=2`)
   - `is_checked` = false
7. Backend sends a **reply message** confirming the ticket was recorded

**Error cases:**

- Invalid number format (not 3 or 6 digits) → reply with error message
- Unrecognized message → reply with help instructions

---

### 3.3 Automatic Win Notification (Cronjob Flow)

**Trigger:** Scheduled cronjob fires on lottery draw days, after official results are published.

**Draw day schedule:** Mostly on 1st and 16th of every month (Thai Government Lottery). But sometimes there are exceptions (e.g., holidays). The cronjob schedule should be configurable to accommodate this.

**Cronjob trigger time:** Results are usually published not earlier than 16:00 Bangkok time (GMT +0700). The cronjob can be scheduled to run at 16:00 on draw days to ensure results are available.

**Steps:**

1. Cronjob fires at the scheduled time on a draw day
2. Cronjob calls the Thai Government Lottery results API:
   - **API endpoint:** https://www.glo.or.th/api/lottery/getLatestLottery (POST request without parameters)
   - **Response format:** Read on ./trunk/glo_result.json for expected structure; typically includes draw date, prize categories, and winning numbers
   - Retry logic: if API fails, retry up to 5 times before alerting
3. System creates or finds the `draws` record for today's draw date
4. System inserts results into `draw_results` table:
   - One row per `(draw_id, prize_category, winning_number)` — unique constraint prevents duplicates
5. System marks `draws.is_verified = true`
6. System queries all `tickets` where `draw_id = current_draw_id AND is_checked = false`
7. For each ticket, compare `tickets.number` against relevant `draw_results` rows:
   - L6 ticket → compare against all L6 prize categories
   - N3 ticket → compare against all N3 prize categories
8. For each matching ticket:
   - Insert a `user_winnings` record: `(user_id, ticket_id, draw_result_id, prize_money)`
   - Unique constraint `(ticket_id, draw_result_id)` prevents duplicate win records
9. System sends a **LINE push message** to each winning user
   - Content: draw date, ticket number, prize category, prize amount
   - **Push message format:** TBA (will decide the exact copy and formatting later)
10. Mark all processed tickets as `is_checked = true`

**Non-win case:**

- For tickets that do not match any winning numbers, send a push message indicating the ticket did not win (optional, but can increase user engagement and trust in the system)
- **Non-winning message format:** TBA (will decide the exact copy and formatting later)

---

### 3.4 On-demand Status Query (Optional / Post-MVP)

Implement in a post-MVP phase. This allows users to query ticket status at any time, not only on draw days.

Possible flow: user sends a command like "ผล" (results) to get recent draw results or their ticket status.

---

## 4. System Architecture

```
User (LINE app)
     |
     | HTTPS
     v
LINE Platform
     |
     | Webhook POST (X-Line-Signature header)
     v
+---------------------------+
|  Lotto Journal Backend    |
|  Go + Fiber               |
|  apps/api                 |
|                           |
|  /webhook  (LINE handler) |
|  /health                  |
+---------------------------+
     |                    ^
     | GORM / SQL         | LINE Push Message API
     v                    |
+------------+     +-----------------------------+
| PostgreSQL |     | LINE Messaging API          |
|            |     | - Webhook delivery          |
+------------+     | - Push Message (notify)     |
     ^             +-----------------------------+
     |
+---------------------------+
| Cronjob (Go scheduler)    |
|  - runs on draw days      |
|  - fetch lottery results  |
|  - compare & notify       |
+---------------------------+
     |
     v (POST)
Thai Government Lottery API
+-------------------------------------------------------------------------------------+
| [/api/lottery/getLatestLottery](https://www.glo.or.th/api/lottery/getLatestLottery)|
+-------------------------------------------------------------------------------------+
```

**Component list (`apps/api/internal/`):**

| Component            | Package                           | Responsibility                       |
| -------------------- | --------------------------------- | ------------------------------------ |
| LINE webhook handler | `handler/line_handler.go`         | Receive + verify LINE events         |
| Message parser       | `handler/` or `service/`          | Parse ticket number from LINE text   |
| User service         | `service/user_service.go`         | LINE user identity, create/find user |
| Ticket service       | `service/ticket_service.go`       | Ticket CRUD                          |
| Draw service         | `service/draw_service.go`         | Draw date management                 |
| Result checker       | `service/result_service.go`       | Compare tickets vs draw results      |
| Notification service | `service/notification_service.go` | LINE push message sender             |
| Cronjob              | `cronjob/` or `cmd/cronjob/`      | Orchestrate draw-day flow            |
| Lottery API client   | `client/lottery_client.go`        | HTTP client for Thai Gov Lottery API |

---

## 5. Data Model

### 5.1 Migration 000002 — User Identity Redesign

This migration implements the changes required by ADR-001 (Option B).

**`users` table — MODIFIED:**

| Column                           | Change | Value                             |
| -------------------------------- | ------ | --------------------------------- |
| `line_user_id`                   | ADD    | `varchar UNIQUE NOT NULL`         |
| `username`                       | REMOVE | —                                 |
| `email`                          | REMOVE | —                                 |
| `password_hash`                  | REMOVE | —                                 |
| `status`                         | KEEP   | `account_status DEFAULT 'active'` |
| `id`, `created_at`, `updated_at` | KEEP   | unchanged                         |

**Tables DROPPED:**

- `user_auth_methods`
- `user_verifications`
- `user_profiles`

**Enums DROPPED:**

- `provider_service`
- `verification_type`

### 5.2 Tables (post-migration snapshot)

| Table           | Notes                                                                |
| --------------- | -------------------------------------------------------------------- |
| `draws`         | Unchanged                                                            |
| `tickets`       | Unchanged; `owner_id` still references `users.id`                    |
| `draw_results`  | Unchanged                                                            |
| `user_winnings` | Includes `user_id` FK (added in migration 000002)                   |
| `files`         | Reserved for a future MVP extension: ticket photo upload + OCR transcription |

### 5.3 Enums (post-migration snapshot)

| Enum             | Values                                                                                                                                                                                  |
| ---------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `account_status` | `active`, `inactive`, `suspended`                                                                                                                                                       |
| `lottery_type`   | `N3`, `L6`                                                                                                                                                                              |
| `prize_type`     | `l6_first`, `l6_second`, `l6_third`, `l6_fourth`, `l6_fifth`, `l6_last2`, `l6_last3f`, `l6_last3b`, `l6_near_first`, `n3_straight_three`, `n3_shuffle`, `n3_straight_two`, `n3_special` |

---

## 6. External Integrations

### 6.1 LINE Messaging API

| Item                   | Value                                                           |
| ---------------------- | --------------------------------------------------------------- |
| Platform               | LINE Messaging API v2                                           |
| Webhook endpoint       | `POST /webhook`                                                 |
| Signature verification | `X-Line-Signature` header using HMAC-SHA256 with channel secret |
| Push message API       | `POST https://api.line.me/v2/bot/message/push`                  |
| SDK                    | line-bot-sdk-go (official Go SDK)                               |
| ENV: channel secret    | `LINE_CHANNEL_SECRET`                                           |
| ENV: access token      | `LINE_CHANNEL_ACCESS_TOKEN`                                     |

(Recheck exact API endpoints and request/response formats in the official LINE Messaging API documentation to ensure accuracy: https://developers.line.biz/en/docs/messaging-api/overview/)

**Webhook event types to handle:**

- `message` (type: text) — user sends a ticket number
- `follow` — user adds the LINE Official Account as a friend
- `unfollow` — mark user as inactive (soft-delete behavior)

**Idempotency:** LINE may re-deliver webhook events. The backend must use the event's
`webhookEventId` (or `message.id`) to detect and skip duplicates.

### 6.2 Thai Government Lottery Results API

| Item            | Value                                                                               |
| --------------- | ----------------------------------------------------------------------------------- |
| Purpose         | Fetch official draw results on draw days                                            |
| Endpoint        | [/api/lottery/getLatestLottery](https://www.glo.or.th/api/lottery/getLatestLottery) |
| Authentication  | open                                                                                |
| Response format | JSON (example: `./trunk/glo_result.json`)                                            |
| Draw dates      | 1st and 16th of each month (subject to change if marked as holiday)                 |
| Rate limits     | TBC                                                                                 |
| Fallback        | Manual fallback if API remains unavailable after configured retries (TBD)            |

---

## 7. Non-functional Requirements

| Requirement | Detail                                                                                                |
| ----------- | ----------------------------------------------------------------------------------------------------- |
| Idempotency | LINE webhook events must be deduplicated; `user_winnings` unique index prevents duplicate win records |
| Reliability | Cronjob must retry on external API failure; log failures clearly                                      |
| Security    | Webhook signature verified on every request; no secrets hardcoded; no PII in logs                     |
| Scalability | TBC; MVP target is up to 100 users                                                                       |
| Monitoring  | MVP uses backend terminal logs only                                                                    |
| Timezone    | All draw-day scheduling uses Bangkok time (GMT +0700)                                                 |

---

## 8. Out of Scope (MVP)

- Ticket photo upload / OCR for bulk entry
- Ticket resale or marketplace
- Admin dashboard or back-office web UI
- Historical statistics or analytics for users
- Web UI for end users (apps/web removed per ADR-001)
- Multiple notification channels simultaneously
- On-demand status query from users (flow 3.4 — deferred post-MVP)
- **Localization & Multi-Language Support (EN/TH):** Out of scope for MVP (Thai language only). Future expansion plans include:
  * Fetching user locale preference dynamically via LINE Profile API (`profile.Language` field returns `"th"` or `"en"`).
  * Supporting custom command overrides (e.g., `EN` / `TH` settings commands).
  * Storing language preferences in a new `users.language` database column, and loading corresponding localization dictionaries dynamically.
- **LIFF (LINE Front-end Framework) web app** — intentionally deferred, not abandoned.
  LIFF would run inside LINE's in-app browser and complement the chatbot (e.g. ticket history
  view, result display, settings). When added it will live in `apps/liff` alongside `apps/api`
  in the existing monorepo. See T-009.
