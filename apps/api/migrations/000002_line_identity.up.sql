-- MIGRATION 000002 — LINE IDENTITY REDESIGN (UP)
-- Transforms the schema from email/password auth to LINE Messaging API identity
-- Reference: doc/06-extensions/T-004-migration-002-design.md

-- =============================================================
-- STEP 1: DROP AUTH TABLES (children before parents, FK order)
-- The unique index on user_auth_methods is dropped automatically with its table.
-- =============================================================

DROP TABLE IF EXISTS "user_auth_methods";

DROP TABLE IF EXISTS "user_verifications";

DROP TABLE IF EXISTS "user_profiles";

-- =============================================================
-- STEP 2: DROP AUTH ENUMS (only safe after dependent tables gone)
-- =============================================================

DROP TYPE IF EXISTS "provider_service";

DROP TYPE IF EXISTS "verification_type";

-- =============================================================
-- STEP 3: MODIFY users TABLE
--   REMOVE: username, email, password_hash
--   ADD:    line_user_id varchar UNIQUE NOT NULL
-- Note: no production data; table is expected to be empty.
-- =============================================================

ALTER TABLE "users"
    DROP COLUMN "username",
    DROP COLUMN "email",
    DROP COLUMN "password_hash";

ALTER TABLE "users"
    ADD COLUMN "line_user_id" varchar UNIQUE NOT NULL;

-- =============================================================
-- STEP 4: RECREATE account_status ENUM
-- PostgreSQL has no DROP VALUE; drop-and-recreate is required.
--   REMOVED: pending  (no email auth flow in LINE-based product)
--   ADDED:   inactive (user unfollowed the LINE OA)
--   KEPT:    active, suspended
-- =============================================================

ALTER TABLE "users"
    ALTER COLUMN "status" DROP DEFAULT;

ALTER TABLE "users"
    ALTER COLUMN "status" TYPE text;

DROP TYPE "account_status";

CREATE TYPE "account_status" AS ENUM ('active', 'inactive', 'suspended');

ALTER TABLE "users"
    ALTER COLUMN "status" TYPE account_status USING "status"::account_status;

ALTER TABLE "users"
    ALTER COLUMN "status" SET DEFAULT 'active';

-- =============================================================
-- STEP 5: RENAME lottery_type ENUM VALUE  (N6 → L6)
-- =============================================================

ALTER TYPE "lottery_type" RENAME VALUE 'N6' TO 'L6';

-- =============================================================
-- STEP 6: RENAME prize_type ENUM VALUES  (n6_* → l6_*, 9 values)
-- =============================================================

ALTER TYPE "prize_type" RENAME VALUE 'n6_first'      TO 'l6_first';
ALTER TYPE "prize_type" RENAME VALUE 'n6_second'     TO 'l6_second';
ALTER TYPE "prize_type" RENAME VALUE 'n6_third'      TO 'l6_third';
ALTER TYPE "prize_type" RENAME VALUE 'n6_fourth'     TO 'l6_fourth';
ALTER TYPE "prize_type" RENAME VALUE 'n6_fifth'      TO 'l6_fifth';
ALTER TYPE "prize_type" RENAME VALUE 'n6_last2'      TO 'l6_last2';
ALTER TYPE "prize_type" RENAME VALUE 'n6_last3f'     TO 'l6_last3f';
ALTER TYPE "prize_type" RENAME VALUE 'n6_last3b'     TO 'l6_last3b';
ALTER TYPE "prize_type" RENAME VALUE 'n6_near_first' TO 'l6_near_first';

-- =============================================================
-- STEP 7: ADD user_id TO user_winnings  [FOUND-IN-PASSING]
-- Column was in the DBML design and referenced in PRD §3.3 but
-- accidentally omitted from migration 000001.
-- =============================================================

ALTER TABLE "user_winnings"
    ADD COLUMN "user_id" uuid REFERENCES "users" ("id") ON DELETE CASCADE;
