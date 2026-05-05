package migrations

import (
	"database/sql"
	"fmt"
)

// MigrateAddSessionColumns adds a session VARCHAR(255) NULL column to the
// events and wisp_events tables if not already present. Idempotent; safe to
// re-run any number of times.
//
// Reason: schema migration 0033 adds these columns on fresh DBs, but
// pre-existing embedded DBs and shared Dolt servers may not have run that
// migration yet. Specifically, Gas Town's internal/doltserver/wisps_migrate.go
// independently CREATEd wisp_events without a session column; the compat layer
// ensures the column exists before any code path attempts to write or read it.
//
// VARCHAR(255) matches the existing actor column shape on both tables.
//
// See gastownhall/beads#3583 and bd-xdc for design rationale.
func MigrateAddSessionColumns(db *sql.DB) error {
	for _, table := range []string{"events", "wisp_events"} {
		tableOK, err := TableExists(db, table)
		if err != nil {
			return fmt.Errorf("check %s table existence: %w", table, err)
		}
		if !tableOK {
			// Table not present at all — schema migration 0033 will create it
			// with the column; nothing to do here.
			continue
		}

		has, err := columnExists(db, table, "session")
		if err != nil {
			return fmt.Errorf("check %s.session: %w", table, err)
		}
		if has {
			continue
		}

		// VARCHAR(255) parallels the existing `actor` column shape on both tables.
		//nolint:gosec // G201: table is from hardcoded list, not user input
		if _, err := db.Exec(fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN session VARCHAR(255) NULL", table)); err != nil {
			return fmt.Errorf("add %s.session: %w", table, err)
		}
	}
	return nil
}
