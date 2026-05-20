-- Add session column to events for per-event session attribution.
-- VARCHAR(255) matches the existing actor column shape.
-- Idempotent via INFORMATION_SCHEMA checks (Dolt-style; matches 0023/0027/0033/0034).
-- wisp_events.session lives in the ignored tree (see ignored/0005) per post-#3918
-- convention: all wisp_* schema changes ship via the ignored tree.
-- See ADR-0003 (bd-edi) for design rationale.

SET @needs_add = (
    SELECT IF(COUNT(*) = 0, 1, 0)
    FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'events'
      AND COLUMN_NAME = 'session'
);
SET @sql = IF(@needs_add = 1,
    'ALTER TABLE events ADD COLUMN session VARCHAR(255) NULL',
    'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
