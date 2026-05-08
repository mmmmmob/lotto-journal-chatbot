-- MIGRATION 000003 — WEBHOOK EVENTS IDEMPOTENCY (UP)
-- Creates the webhook_events table to track processed LINE webhook event IDs.
-- Reference: doc/02-task/task-board.md T-002
--
-- LINE may re-deliver webhook events. This table ensures each webhookEventId
-- is processed exactly once. Inserting with ON CONFLICT DO NOTHING lets the
-- application detect duplicates atomically without a separate SELECT.

CREATE TABLE "webhook_events" (
    "event_id"     varchar   PRIMARY KEY,
    "processed_at" timestamp NOT NULL DEFAULT now()
);
