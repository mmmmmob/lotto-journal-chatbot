-- MIGRATION 000005 — AUDIT LOGS FOR OUTBOUND NOTIFICATIONS (UP)
CREATE TYPE "notification_status" AS ENUM ('success', 'failed');
CREATE TYPE "notification_type" AS ENUM ('welcome', 'ticket_submitted', 'ticket_list', 'draw_result');

CREATE TABLE "notification_logs" (
    "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    "user_id" uuid NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE,
    "line_user_id" varchar NOT NULL,
    "notification_type" notification_type NOT NULL,
    "draw_id" uuid REFERENCES "draws" ("id") ON DELETE SET NULL,
    "status" notification_status NOT NULL,
    "error_message" text,
    "created_at" timestamp NOT NULL DEFAULT now()
);

CREATE INDEX ON "notification_logs" ("draw_id");
CREATE INDEX ON "notification_logs" ("user_id");
