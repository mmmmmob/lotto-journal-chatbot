# LINE Webhook — How It's Wired and How It Works

> A step-by-step walkthrough of `POST /webhook` from the moment LINE sends a request
> to the moment the user receives a reply.
>
> **Code lives in:** `apps/api/`

---

## The Big Picture

```
LINE Platform
    │
    │  POST /webhook
    │  Header: X-Line-Signature: <hmac>
    │  Body:   { "events": [...] }
    ▼
┌─────────────────────────────────────────────────────┐
│  Fiber app  (app/main.go)                           │
│                                                     │
│  Global middleware (every request):                 │
│    recoverer → requestid → Logging                  │
│                                                     │
│  Route: POST /webhook                               │
│    timeout.New(25 s) → lineHandler.Handle           │
└───────────────────┬─────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────────────────┐
│  LineHandler  (internal/handler/line_handler.go)    │
│                                                     │
│  1. Verify signature                                │
│  2. Parse events                                    │
│  3. Deduplicate (webhook_events table)              │
│  4. Route → follow / unfollow / message             │
└───────────────────┬─────────────────────────────────┘
                    │
          ┌─────────┼──────────┐
          ▼         ▼          ▼
       follow   unfollow   message (text)
          │         │          │
          │         │    ┌─────┴──────────────┐
          │         │    │ isTicketListCmd?   │
          │         │    └──────┬─────────────┘
          │         │        yes│          no│
          │         │           ▼            ▼
          ▼         ▼     TicketService  TicketService
     UserService  User    ListTickets    SubmitTickets
     FindOrCreate Deactivate    │            │
          │                     ▼            ▼
    UserRepository         DrawRepository + TicketRepository
    (users table)          (draws + tickets tables)
          │                   │
          └─────────┬─────────┘
                    ▼
             PostgreSQL DB
                    │
                    ▼ (reply via LINE SDK)
           messaging_api.ReplyMessage
                    │
                    ▼
             LINE Platform
                    │
                    ▼
              User's phone
```

---

## Step 1 — Register the Route (`app/main.go`)

All requests pass through a global middleware chain before reaching any handler:

```go
app.Use(recoverer.New(recoverer.Config{EnableStackTrace: true})) // panic → 500
app.Use(requestid.New())                                         // assign X-Request-ID
app.Use(middlewares.Logging)                                     // log after response
```

The `/webhook` route then wraps the handler in a **25 s timeout**. If the handler does
not return within 25 seconds, Fiber abandons the context and immediately sends `408`.
See [`apps/api/README.md`](../apps/api/README.md#why-does-webhook-need-a-timeout) for the
full rationale (short version: LINE redelivers if no `2xx` arrives, so a hung handler
causes a retry storm on an already-stalled server).

```go
app.Get("/health", healthHandler.Handle)  // no timeout — DB ping is sub-millisecond

app.Post("/webhook", timeout.New(lineHandler.Handle, timeout.Config{
    Timeout: 25 * time.Second,
}))
```

`lineHandler` is built by wiring every layer together in `main.go`:

```go
// Repositories — own all DB access
userRepo    := repository.NewUserRepository(db)
drawRepo    := repository.NewDrawRepository(db)
ticketRepo  := repository.NewTicketRepository(db)
webhookRepo := repository.NewWebhookEventRepository(db)

// Services — own all business logic
userSvc   := service.NewUserService(userRepo)
ticketSvc := service.NewTicketService(ticketRepo, drawRepo)

// Handler — owns HTTP concerns only
lineHandler := handler.NewLineHandler(
    cfg.LineChannelSecret,      // for signature verification
    bot,                        // LINE messaging client (for replies)
    userSvc,
    ticketSvc,
    webhookRepo,                // for deduplication
)
```

**Why this matters:** Each layer only knows about the layer directly below it.
The handler never touches the DB. The service never touches HTTP. This is called
**layered architecture** — dependencies always flow downward.

---

## Step 2 — Receive and Verify the Request (`Handle`)

```go
func (h *LineHandler) Handle(c fiber.Ctx) error {
    req := &http.Request{
        Method: "POST",
        Header: http.Header{
            "X-Line-Signature": []string{c.Get("X-Line-Signature")},
            "Content-Type":     []string{"application/json"},
        },
        Body: io.NopCloser(bytes.NewReader(c.Body())),
    }

    cb, err := webhook.ParseRequest(h.channelSecret, req)
    ...
}
```

**Why the synthetic `*http.Request`?**

Fiber uses `fasthttp` internally — a different HTTP engine from Go's standard `net/http`.
The LINE SDK only speaks `net/http`. So the handler builds a fake `*http.Request`
from Fiber's raw context (`c.Body()`, `c.Get(...)`) and hands it to the SDK.

**What `webhook.ParseRequest` does internally:**

1. Reads the raw body bytes
2. Computes `HMAC-SHA256(body, channelSecret)`
3. Base64-encodes the result
4. Compares it against the `X-Line-Signature` header
5. If they match → parse the JSON body into typed event structs
6. If they don't → return `webhook.ErrInvalidSignature`

LINE signs every webhook delivery with your channel secret. Any request that fails
signature check is rejected with `400 Bad Request` — this prevents forged events.

---

## Step 3 — Deduplicate the Event (`dispatch`)

LINE's platform has an **at-least-once delivery guarantee** — it may send the same
event more than once if it doesn't receive a `200 OK` in time.

```go
func (h *LineHandler) dispatch(event webhook.EventInterface) {
    eventID := webhookEventID(event)    // e.g. "01H..." (ULID format)

    isNew, err := h.webhookRepo.MarkProcessed(eventID)
    if !isNew {
        log.Printf("[webhook] duplicate event %s — skipping", eventID)
        return  // already handled, do nothing
    }

    switch e := event.(type) { ... }
}
```

`MarkProcessed` runs this SQL atomically:

```sql
INSERT INTO webhook_events (event_id, processed_at)
VALUES ($1, NOW())
ON CONFLICT DO NOTHING;
```

- **First delivery** → INSERT succeeds → `RowsAffected = 1` → `isNew = true` → process
- **Re-delivery** → INSERT hits the PK conflict → silently skipped → `RowsAffected = 0` → `isNew = false` → skip

The whole check-and-record is one atomic database operation — no race condition possible.

---

## Step 4a — `follow` Event

**Trigger:** User adds the LINE Official Account as a friend (or unblocks it).

```go
func (h *LineHandler) handleFollow(e webhook.FollowEvent) {
    lineUserID := sourceUserID(e.Source)        // extract "Uxxxxx" from event

    _, isNew, err := h.userSvc.FindOrCreate(lineUserID)
    // → SQL: SELECT ... WHERE line_user_id = ?
    //         INSERT INTO users (...) ON CONFLICT DO NOTHING

    h.replyText(e.ReplyToken, welcomeMessage)   // send reply via LINE SDK
}
```

**`FindOrCreate` under the hood (GORM `FirstOrCreate`):**

```sql
-- 1. Try to find
SELECT * FROM users WHERE line_user_id = $1 LIMIT 1;

-- 2. If not found, create
INSERT INTO users (line_user_id, status) VALUES ($1, 'active');
```

`RowsAffected` tells the service whether this is a brand new user (`1`) or a
returning user who re-followed (`0`). Both cases get the welcome reply.

**Reply path:**

```go
h.bot.ReplyMessage(&messaging_api.ReplyMessageRequest{
    ReplyToken: e.ReplyToken,       // one-time token from LINE, expires in 30s
    Messages: []messaging_api.MessageInterface{
        messaging_api.TextMessage{Text: welcome},
    },
})
```

The `ReplyToken` is embedded in the event — it's LINE's way of linking a reply to
a specific incoming message. It can only be used once and expires after 30 seconds.

---

## Step 4b — `unfollow` Event

**Trigger:** User removes the LINE Official Account (blocks or unfriends).

```go
func (h *LineHandler) handleUnfollow(e webhook.UnfollowEvent) {
    lineUserID := sourceUserID(e.Source)
    h.userSvc.Deactivate(lineUserID)
    // → SQL: UPDATE users SET status = 'inactive' WHERE line_user_id = $1
    // No reply — LINE does not allow replying to unfollow events (no ReplyToken)
}
```

The user is **soft-deleted** (status set to `inactive`), not removed from the DB.
Their ticket history is preserved. If they follow again, `FindOrCreate` finds the
existing row and the `follow` handler can re-activate them if needed.

---

## Step 4c — `message` Event

**Trigger:** User sends a text message.

The handler first checks whether the message is a recognised command keyword, then routes accordingly.

### 4c-0: Route by keyword

```go
if isTicketListCmd(textMsg.Text) {
    // → list tickets for upcoming draw
} else {
    // → parse and submit tickets
}
```

`isTicketListCmd` is a simple equality check: `text == "โพย"`.
Anything that is not a recognised command is treated as a ticket submission attempt.

---

### 4c — List tickets command (`โพย`)

```go
userTickets, err := h.ticketSvc.ListTickets(user.ID)
h.replyText(e.ReplyToken, buildTicketListReply(userTickets))
```

`ListTickets` resolves the upcoming draw (same `FindOrCreate` logic as submission)
then fetches all tickets for that draw and user:

```sql
SELECT * FROM tickets WHERE owner_id = $1 AND draw_id = $2;
```

`buildTicketListReply` formats the result:

| Situation          | Reply                                                |
| ------------------ | ---------------------------------------------------- |
| Has tickets        | `"สลากที่คุณบันทึกไว้ในงวดนี้ 📝\n  • 123456 x2 (L6)"` |
| No tickets yet     | `"คุณยังไม่ได้บันทึกสลากในงวดนี้"`                    |

---

### 4c — Ticket submission (any other text)

**Trigger:** User sends a text message, e.g. `"123456 x2, 789"`.

### 4c-1: Extract user

```go
lineUserID := sourceUserID(e.Source)
user, _, err := h.userSvc.FindOrCreate(lineUserID)
// Handles edge case: message arrives before the follow event
```

### 4c-2: Parse the message text

```go
saved, invalid, err := h.ticketSvc.SubmitTickets(user.ID, textMsg.Text)
```

`ParseTicketInput` inside `TicketService` runs these steps:

| Input text         | After normalise    | After merge x    | Tokens                | Result             |
| ------------------ | ------------------ | ---------------- | --------------------- | ------------------ |
| `"123456 x2, 789"` | `"123456 x2  789"` | `"123456x2 789"` | `["123456x2", "789"]` | L6×2, N3×1         |
| `"12345"`          | `"12345"`          | `"12345"`        | `["12345"]`           | invalid (5 digits) |
| `"hello"`          | `"hello"`          | `"hello"`        | `[]`                  | unrecognised       |

Rules:

- **6 digits** → `L6` lottery type
- **3 digits** → `N3` lottery type
- **Anything else** → added to `invalid` list

### 4c-3: Resolve the upcoming draw

```go
draw, err := s.drawRepo.FindOrCreate(NextDrawDate(time.Now()))
```

`NextDrawDate` calculates the nearest 1st or 16th of the month in Bangkok time
(the Thai Government Lottery draw days). The draw record is found or created:

```sql
SELECT * FROM draws WHERE draw_date = $1;
-- if not found:
INSERT INTO draws (draw_date) VALUES ($1);
```

### 4c-4: Save tickets

```go
for _, pt := range parsed {
    ticket := &models.Ticket{
        OwnerID:  ownerID,
        DrawID:   draw.ID,
        Type:     pt.Type,    // "L6" or "N3"
        Number:   pt.Number,
        Quantity: pt.Quantity,
    }
    s.ticketRepo.Create(ticket)
    // → INSERT INTO tickets (...) VALUES (...)
}
```

One row per parsed ticket entry. `quantity` stores how many physical tickets
the user holds for that number (from the `x2` syntax).

### 4c-5: Reply

```go
h.replyText(e.ReplyToken, buildReply(saved, invalid))
```

`buildReply` produces different messages depending on the outcome:

| Situation        | Reply                                                        |
| ---------------- | ------------------------------------------------------------ |
| All valid        | `"บันทึกสลากเรียบร้อย ✅\n  • 123456 x2 (L6)\n  • 789 (N3)"` |
| Some invalid     | Success list + `"เลขที่ไม่ถูกต้อง (ข้ามไป): 12345"`          |
| All invalid      | `"ไม่พบเลขสลากที่ถูกต้อง ❌ ..."`                            |
| No digits at all | Help instructions                                            |

---

## Step 5 — Always Return `200 OK`

```go
return c.SendStatus(fiber.StatusOK)
```

This is critical. LINE's platform **redelivers** the webhook if it does not receive a `200 OK`
(the retry count and interval are not publicly disclosed — see
[LINE docs](https://developers.line.biz/en/docs/messaging-api/receiving-messages/#redeliver-a-webhook-that-failed-to-be-received)).
The handler always returns `200` — even if an individual event fails internally (logged but
not bubbled up). The 25 s timeout on the route guarantees a response is always emitted;
the deduplication in Step 3 prevents a redelivered event from being processed twice.

---

## Data Flow Summary

```
LINE sends POST /webhook
    │
    ├── [Handle] verify X-Line-Signature → reject 400 if invalid
    │
    ├── for each event in cb.Events:
    │       │
    │       ├── [dispatch] MarkProcessed(eventId)
    │       │       ├── new  → continue
    │       │       └── dupe → skip, return
    │       │
    │       ├── follow event
    │       │       ├── UserService.FindOrCreate(lineUserId)
    │       │       └── ReplyMessage(welcomeText)
    │       │
    │       ├── unfollow event
    │       │       └── UserService.Deactivate(lineUserId)
    │       │
    │               └── message event (text)
    │                       ├── UserService.FindOrCreate(lineUserId)
    │                       ├── isTicketListCmd?
    │                       │       ├── yes → TicketService.ListTickets(userId)
    │                       │       │               ├── DrawRepo.FindOrCreate(nextDrawDate)
    │                       │       │               └── TicketRepo.List(drawId, userId)
    │                       │       │           ReplyMessage(ticketListText)
    │                       │       └── no  → TicketService.SubmitTickets(userId, text)
    │                       │                       ├── ParseTicketInput(text)
    │                       │                       ├── DrawRepo.FindOrCreate(nextDrawDate)
    │                       │                       └── TicketRepo.Create(ticket) × N
    │                       │           ReplyMessage(confirmationText)
    │
    └── return 200 OK
```

---

## Files Involved

| File                                              | Role                                                                         |
| ------------------------------------------------- | ---------------------------------------------------------------------------- |
| `app/main.go`                                     | Wires all layers; registers global middleware + `POST /webhook` with timeout |
| `middlewares/log.go`                              | Logs every request: method, path, status, duration, req_id                   |
| `internal/handler/line_handler.go`                | Receives request, verifies, routes events, sends replies                     |
| `internal/service/user_service.go`                | Find-or-create user; deactivate on unfollow                                  |
| `internal/service/ticket_service.go`              | Parse message text; save tickets                                             |
| `internal/service/draw_service.go`                | Calculate next draw date                                                     |
| `internal/repository/user_repository.go`          | `users` table — FindOrCreate, UpdateStatus                                   |
| `internal/repository/draw_repository.go`          | `draws` table — FindOrCreate                                                 |
| `internal/repository/ticket_repository.go`        | `tickets` table — `Create`, `List`                                           |
| `internal/repository/webhook_event_repository.go` | `webhook_events` table — MarkProcessed                                       |
| `internal/config/config.go`                       | Loads `LINE_CHANNEL_SECRET`, `LINE_CHANNEL_ACCESS_TOKEN`                     |
| `migrations/000003_webhook_events.up.sql`         | Creates the idempotency table                                                |
