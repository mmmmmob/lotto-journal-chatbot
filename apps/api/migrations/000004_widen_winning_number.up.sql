-- MIGRATION 000004 — WIDEN WINNING NUMBER FOR N3 JACKPOT (UP)
-- Alters draw_results.winning_number from varchar(6) to varchar(12)
-- to accommodate the 12-digit N3 Special Prize (Jackpot) winning number.

ALTER TABLE "draw_results" ALTER COLUMN "winning_number" TYPE varchar(12);
