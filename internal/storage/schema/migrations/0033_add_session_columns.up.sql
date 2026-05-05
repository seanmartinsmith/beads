-- Add session column to events and wisp_events for per-event session attribution.
-- VARCHAR(255) matches the existing actor column shape on both tables.
-- Bare ADD COLUMN (no IF NOT EXISTS) matches existing style (see 0023, 0027).
-- Normal re-runs are prevented by schema_migrations.version tracking.
-- Pre-existing-column environments (e.g., Gas Town's independently-created
-- wisp_events) are handled by compat migration 018_add_session_columns.go.
-- See gastownhall/beads#3583 and bd-xdc for design rationale.

ALTER TABLE events ADD COLUMN session VARCHAR(255) NULL;

ALTER TABLE wisp_events ADD COLUMN session VARCHAR(255) NULL;
