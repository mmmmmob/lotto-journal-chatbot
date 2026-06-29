-- Step 1: Temporarily convert the column to text
ALTER TABLE "notification_logs" ALTER COLUMN "notification_type" TYPE text;

-- Step 2: Drop the enum type
DROP TYPE "notification_type";

-- Step 3: Recreate the enum type with only the original values
CREATE TYPE "notification_type" AS ENUM ('welcome', 'ticket_submitted', 'ticket_list', 'draw_result');

-- Step 4: Convert the column back to the enum type (casting existing rows)
-- Note: Any rows with the new enum values must be cleaned up first.
DELETE FROM "notification_logs" WHERE "notification_type" IN ('language_changed', 'help_add', 'help_notify');

ALTER TABLE "notification_logs" ALTER COLUMN "notification_type" TYPE notification_type USING "notification_type"::notification_type;
