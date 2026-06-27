-- MIGRATION 000004 — WIDEN WINNING NUMBER FOR N3 JACKPOT (DOWN)
-- Alters draw_results.winning_number back from varchar(12) to varchar(6).

ALTER TABLE "draw_results" ALTER COLUMN "winning_number" TYPE varchar(6);
