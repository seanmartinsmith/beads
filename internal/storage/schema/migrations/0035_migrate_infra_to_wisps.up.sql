-- Move infra-type rows from issues into wisps using the intersection of
-- columns common to both tables. Using SELECT * crashes on upgraded DBs
-- where the two tables have drifted (e.g. 0033 added wisp_type and 0034
-- added spec_id to issues only). Mirrors the fix in #2168 / commit
-- 4891d870f for the Go predecessor migration 007.
SET SESSION group_concat_max_len = 32768;
SET @cols = (
    SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY ORDINAL_POSITION SEPARATOR ',')
    FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'wisps'
      AND COLUMN_NAME IN (
          SELECT COLUMN_NAME FROM INFORMATION_SCHEMA.COLUMNS
          WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'issues'
      )
);
SET @sql = IF(
    @cols IS NOT NULL AND @cols <> '',
    CONCAT(
        'INSERT IGNORE INTO wisps (', @cols, ') ',
        'SELECT ', @cols, ' FROM issues ',
        'WHERE issue_type IN (''agent'', ''rig'', ''role'', ''message'')'
    ),
    'SELECT 1'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql = IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.TABLES
        WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'wisps') > 0,
    'UPDATE wisps SET ephemeral = 1 WHERE issue_type IN (''agent'', ''rig'', ''role'', ''message'')',
    'SELECT 1'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql = IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.TABLES
        WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'wisp_labels') > 0,
    'INSERT IGNORE INTO wisp_labels (issue_id, label) SELECT l.issue_id, l.label FROM labels l JOIN issues i ON i.id = l.issue_id WHERE i.issue_type IN (''agent'', ''rig'', ''role'', ''message'')',
    'SELECT 1'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql = IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.TABLES
        WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'wisp_dependencies') > 0,
    'INSERT IGNORE INTO wisp_dependencies (issue_id, depends_on_id, type, created_at, created_by, metadata, thread_id) SELECT d.issue_id, d.depends_on_id, d.type, d.created_at, d.created_by, d.metadata, d.thread_id FROM dependencies d JOIN issues i ON i.id = d.issue_id WHERE i.issue_type IN (''agent'', ''rig'', ''role'', ''message'')',
    'SELECT 1'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql = IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.TABLES
        WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'wisp_events') > 0,
    'INSERT IGNORE INTO wisp_events (id, issue_id, event_type, actor, old_value, new_value, comment, created_at) SELECT e.id, e.issue_id, e.event_type, e.actor, e.old_value, e.new_value, e.comment, e.created_at FROM events e JOIN issues i ON i.id = e.issue_id WHERE i.issue_type IN (''agent'', ''rig'', ''role'', ''message'')',
    'SELECT 1'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql = IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.TABLES
        WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'wisp_comments') > 0,
    'INSERT IGNORE INTO wisp_comments (id, issue_id, author, text, created_at) SELECT c.id, c.issue_id, c.author, c.text, c.created_at FROM comments c JOIN issues i ON i.id = c.issue_id WHERE i.issue_type IN (''agent'', ''rig'', ''role'', ''message'')',
    'SELECT 1'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql = IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.TABLES
        WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'wisp_comments') > 0,
    'DELETE c FROM comments c JOIN issues i ON i.id = c.issue_id WHERE i.issue_type IN (''agent'', ''rig'', ''role'', ''message'')',
    'SELECT 1'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql = IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.TABLES
        WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'wisp_events') > 0,
    'DELETE e FROM events e JOIN issues i ON i.id = e.issue_id WHERE i.issue_type IN (''agent'', ''rig'', ''role'', ''message'')',
    'SELECT 1'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql = IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.TABLES
        WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'wisp_dependencies') > 0,
    'DELETE d FROM dependencies d JOIN issues i ON i.id = d.issue_id WHERE i.issue_type IN (''agent'', ''rig'', ''role'', ''message'')',
    'SELECT 1'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql = IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.TABLES
        WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'wisp_labels') > 0,
    'DELETE l FROM labels l JOIN issues i ON i.id = l.issue_id WHERE i.issue_type IN (''agent'', ''rig'', ''role'', ''message'')',
    'SELECT 1'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql = IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.TABLES
        WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'wisps') > 0,
    'DELETE FROM issues WHERE issue_type IN (''agent'', ''rig'', ''role'', ''message'')',
    'SELECT 1'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
