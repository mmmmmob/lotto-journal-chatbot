-- TO CREATE ALL TABLES ON DATABASE
-- ENUM TYPES
CREATE TYPE "account_status" AS ENUM ('active', 'pending', 'suspended');

CREATE TYPE "verification_type" AS ENUM ('email_verification', 'password_reset');

CREATE TYPE "provider_service" AS ENUM ('google', 'facebook', 'apple', 'local');

CREATE TYPE "lottery_type" AS ENUM ('N3', 'N6');

CREATE TYPE "prize_type" AS ENUM (
    'n6_first',
    'n6_second',
    'n6_third',
    'n6_fourth',
    'n6_fifth',
    'n6_last2',
    'n6_last3f',
    'n6_last3b',
    'n6_near_first',
    'n3_straight_three',
    'n3_shuffle',
    'n3_straight_two',
    'n3_special'
);

-- TABLE WITH CASCADE AND DEFAULT VALUE
-- MASTER TABLES (NO RELATION DEPENDENCY ON OTHER TABLES)
CREATE TABLE
    "users" (
        "id" uuid PRIMARY KEY DEFAULT gen_random_uuid (),
        "username" varchar UNIQUE NOT NULL,
        "email" varchar UNIQUE NOT NULL,
        "password_hash" varchar,
        "status" account_status DEFAULT 'active',
        "created_at" timestamp DEFAULT now (),
        "updated_at" timestamp DEFAULT now ()
    );

CREATE TABLE
    "draws" (
        "id" uuid PRIMARY KEY DEFAULT gen_random_uuid (),
        "draw_date" date UNIQUE NOT NULL,
        "is_verified" boolean DEFAULT false,
        "created_at" timestamp DEFAULT now (),
        "updated_at" timestamp DEFAULT now ()
    );

-- HIGHEST LEVEL WHICH RELATED TO MARSTER
CREATE TABLE
    "files" (
        "id" uuid PRIMARY KEY DEFAULT gen_random_uuid (),
        "owner_id" uuid REFERENCES "users" ("id") ON DELETE SET NULL,
        "file_path" varchar UNIQUE NOT NULL,
        "file_type" varchar,
        "created_at" timestamp DEFAULT now ()
    );

-- CHILDREN TABLES
CREATE TABLE
    "user_verifications" (
        "id" uuid PRIMARY KEY DEFAULT gen_random_uuid (),
        "user_id" uuid REFERENCES "users" ("id") ON DELETE CASCADE,
        "otp_code" varchar(6) NOT NULL,
        "type" verification_type NOT NULL,
        "is_used" boolean DEFAULT false,
        "expired_at" timestamp NOT NULL,
        "created_at" timestamp DEFAULT now ()
    );

CREATE TABLE
    "user_profiles" (
        "user_id" uuid PRIMARY KEY REFERENCES "users" ("id") ON DELETE CASCADE,
        "first_name" varchar,
        "last_name" varchar,
        "avatar_file_id" uuid REFERENCES "files" ("id") ON DELETE SET NULL,
        "updated_at" timestamp DEFAULT now ()
    );

CREATE TABLE
    "user_auth_methods" (
        "id" uuid PRIMARY KEY DEFAULT gen_random_uuid (),
        "user_id" uuid REFERENCES "users" ("id") ON DELETE CASCADE,
        "provider" provider_service NOT NULL,
        "provider_user_id" varchar NOT NULL,
        "provider_email" varchar,
        "created_at" timestamp DEFAULT now (),
        "updated_at" timestamp DEFAULT now ()
    );

CREATE TABLE
    "tickets" (
        "id" uuid PRIMARY KEY DEFAULT gen_random_uuid (),
        "owner_id" uuid REFERENCES "users" ("id") ON DELETE CASCADE,
        "draw_id" uuid REFERENCES "draws" ("id"),
        "type" lottery_type,
        "number" varchar(6),
        "quantity" int4 DEFAULT 1,
        "tickets_file_id" uuid REFERENCES "files" ("id") ON DELETE SET NULL,
        "is_checked" boolean DEFAULT false,
        "created_at" timestamp DEFAULT now (),
        "updated_at" timestamp DEFAULT now ()
    );

CREATE TABLE
    "draw_results" (
        "id" uuid PRIMARY KEY DEFAULT gen_random_uuid (),
        "draw_id" uuid REFERENCES "draws" ("id") ON DELETE CASCADE,
        "prize_category" prize_type,
        "winning_number" varchar(6),
        "prize_amount" int4
    );

CREATE TABLE
    "user_winnings" (
        "id" uuid PRIMARY KEY DEFAULT gen_random_uuid (),
        "ticket_id" uuid REFERENCES "tickets" ("id") ON DELETE CASCADE,
        "draw_result_id" uuid REFERENCES "draw_results" ("id") ON DELETE CASCADE,
        "prize_money" int4,
        "created_at" timestamp DEFAULT now ()
    );

-- UNIQUE INDEX
CREATE UNIQUE INDEX ON "user_auth_methods" ("provider", "provider_user_id");

CREATE UNIQUE INDEX ON "draw_results" ("draw_id", "prize_category", "winning_number");

CREATE UNIQUE INDEX ON "user_winnings" ("ticket_id", "draw_result_id");