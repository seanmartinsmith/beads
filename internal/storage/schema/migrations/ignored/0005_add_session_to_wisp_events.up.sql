-- Mirror migration 0043 on the ignored tree so freshly-recreated wisp_events
-- (via the ignored migration tree) carry the session column. Idempotent.
SET @needs_add = IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.TABLES
        WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'wisp_events') > 0
    AND
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
        WHERE TABLE_SCHEMA = DATABASE()
          AND TABLE_NAME = 'wisp_events'
          AND COLUMN_NAME = 'session') = 0,
    1, 0
);
SET @sql = IF(@needs_add = 1,
    'ALTER TABLE wisp_events ADD COLUMN session VARCHAR(255) NULL',
    'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
