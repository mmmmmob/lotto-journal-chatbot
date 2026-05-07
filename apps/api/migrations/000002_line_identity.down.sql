-- MIGRATION 000002 — LINE IDENTITY REDESIGN (DOWN)
-- Restores the schema to the exact 000001 state
-- Reference: doc/06-extensions/T-004-migration-002-design.md

-- =============================================================
-- STEP 1: REMOVE user_id FROM user_winnings
-- =============================================================

ALTER TABLE "user_winnings"
    DROP COLUMN IF EXISTS "user_id";

-- =============================================================
-- STEP 2: RENAME lottery_type ENUM VALUE BACK  (L6 → N6)
-- =============================================================

ALTER TYPE "lottery_type" RENAME VALUE 'L6' TO 'N6';

-- =============================================================
-- STEP 3: RENAME prize_type ENUM VALUES BACK  (l6_* → n6_*, 9 values)
-- =============================================================

ALTER TYPE "prize_type" RENAME VALUE 'l6_first'      TO 'n6_first';
ALTER TYPE "prize_type" RENAME VALUE 'l6_second'     TO 'n6_second';
ALTER TYPE "prize_type" RENAME VALUE 'l6_third'      TO 'n6_third';
ALTER TYPE "prize_type" RENAME VALUE 'l6_fourth'     TO 'n6_fourth';
ALTER TYPE "prize_type" RENAME VALUE 'l6_fifth'      TO 'n6_fifth';
ALTER TYPE "prize_type" RENAME VALUE 'l6_last2'      TO 'n6_last2';
ALTER TYPE "prize_type" RENAME VALUE 'l6_last3f'     TO 'n6_last3f';
ALTER TYPE "prize_type" RENAME VALUE 'l6_last3b'     TO 'n6_last3b';
ALTER TYPE "prize_type" RENAME VALUE 'l6_near_first' TO 'n6_near_first';

-- =============================================================
-- STEP 4: RESTORE account_status ENUM TO 000001 STATE
--   RESTORE: pending
--   REMOVE:  inactive
-- =============================================================

ALTER TABLE "users"
    ALTER COLUMN "status" DROP DEFAULT;

ALTER TABLE "users"
    ALTER COLUMN "status" TYPE text;

DROP TYPE "account_status";

CREATE TYPE "account_status" AS ENUM ('active', 'pending', 'suspended');

ALTER TABLE "users"
    ALTER COLUMN "status" TYPE account_status USING "status"::account_status;

ALTER TABLE "users"
    ALTER COLUMN "status" SET DEFAULT 'active';

-- =============================================================
-- STEP 5: RESTORE users TABLE COLUMNS
--   REMOVE: line_user_id
--   ADD:    username, email, password_hash
-- Note: DEFAULT '' used for NOT NULL columns to handle any dev rows;
--       defaults are dropped immediately after — not for production use.
-- =============================================================

ALTER TABLE "users"
    DROP COLUMN IF EXISTS "line_user_id";

ALTER TABLE "users"
    ADD COLUMN "username"      varchar UNIQUE NOT NULL DEFAULT '',
    ADD COLUMN "email"         varchar UNIQUE NOT NULL DEFAULT '',
    ADD COLUMN "password_hash" varchar;

ALTER TABLE "users"
    ALTER COLUMN "username" DROP DEFAULT,
    ALTER COLUMN "email"    DROP DEFAULT;

-- =============================================================
-- STEP 6: RECREATE DROPPED ENUMS
-- =============================================================

CREATE TYPE "verification_type" AS ENUM ('email_verification', 'password_reset');

CREATE TYPE "provider_service" AS ENUM ('google', 'facebook', 'apple', 'local');

-- =============================================================
-- STEP 7: RECREATE DROPPED TABLES (parents before children)
-- =============================================================

CREATE TABLE
    "user_verifications" (
        "id"         uuid PRIMARY KEY DEFAULT gen_random_uuid (),
        "user_id"    uuid REFERENCES "users" ("id") ON DELETE CASCADE,
        "otp_code"   varchar(6) NOT NULL,
        "type"       verification_type NOT NULL,
        "is_used"    boolean DEFAULT false,
        "expired_at" timestamp NOT NULL,
        "created_at" timestamp DEFAULT now ()
    );

CREATE TABLE
    "user_profiles" (
        "user_id"        uuid PRIMARY KEY REFERENCES "users" ("id") ON DELETE CASCADE,
        "first_name"     varchar,
        "last_name"      varchar,
        "avatar_file_id" uuid REFERENCES "files" ("id") ON DELETE SET NULL,
        "updated_at"     timestamp DEFAULT now ()
    );

CREATE TABLE
    "user_auth_methods" (
        "id"               uuid PRIMARY KEY DEFAULT gen_random_uuid (),
        "user_id"          uuid REFERENCES "users" ("id") ON DELETE CASCADE,
        "provider"         provider_service NOT NULL,
        "provider_user_id" varchar NOT NULL,
        "provider_email"   varchar,
        "created_at"       timestamp DEFAULT now (),
        "updated_at"       timestamp DEFAULT now ()
    );

CREATE UNIQUE INDEX ON "user_auth_methods" ("provider", "provider_user_id");
