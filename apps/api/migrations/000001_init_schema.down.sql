-- TO ROLLBACK ALL TABLES AND RESTART MIGRATION PROCESS
-- DROP INDICES
DROP INDEX IF EXISTS "user_auth_methods_provider_provider_user_id_idx";

DROP INDEX IF EXISTS "draw_results_draw_id_prize_category_winning_number_idx";

DROP INDEX IF EXISTS "user_winnings_ticket_id_draw_result_id_idx";

-- DROP TABLE FROM CHILDREN TO MASTERS
DROP TABLE IF EXISTS "user_winnings" CASCADE;

DROP TABLE IF EXISTS "draw_results" CASCADE;

DROP TABLE IF EXISTS "tickets" CASCADE;

DROP TABLE IF EXISTS "user_auth_methods" CASCADE;

DROP TABLE IF EXISTS "user_profiles" CASCADE;

DROP TABLE IF EXISTS "user_verifications" CASCADE;

DROP TABLE IF EXISTS "files" CASCADE;

DROP TABLE IF EXISTS "draws" CASCADE;

DROP TABLE IF EXISTS "users" CASCADE;

-- DROP ENUM TYPES
DROP TYPE IF EXISTS "prize_type";

DROP TYPE IF EXISTS "lottery_type";

DROP TYPE IF EXISTS "provider_service";

DROP TYPE IF EXISTS "verification_type";

DROP TYPE IF EXISTS "account_status";