package migrations

import (
	"database/sql"
	"testing"
)

// createEventTableStubs creates minimal events + wisp_events tables (without
// session column) so migration tests can exercise the ADD COLUMN path.
// The shared test schema (testmain_test.go: initMigrationSharedSchema) creates
// only the issues table, so each test that needs event tables must seed them.
// Schema mirrors production shape minimally — id/issue_id/event_type/actor —
// matching the columns the migration cares about (actor for VARCHAR(255)
// shape parity, no session column to force the ADD path).
func createEventTableStubs(t *testing.T, db *sql.DB) {
	t.Helper()
	for _, ddl := range []string{
		`CREATE TABLE events (
			id CHAR(36) NOT NULL PRIMARY KEY,
			issue_id VARCHAR(255) NOT NULL,
			event_type VARCHAR(32) NOT NULL,
			actor VARCHAR(255) NOT NULL
		)`,
		`CREATE TABLE wisp_events (
			id CHAR(36) NOT NULL PRIMARY KEY,
			issue_id VARCHAR(255) NOT NULL,
			event_type VARCHAR(32) NOT NULL,
			actor VARCHAR(255) NOT NULL
		)`,
	} {
		if _, err := db.Exec(ddl); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}
}

// TestMigrateAddSessionColumns_FreshDB verifies the migration adds session to
// both events and wisp_events on a database that doesn't have the column yet.
func TestMigrateAddSessionColumns_FreshDB(t *testing.T) {
	db := openTestDoltBranch(t)
	createEventTableStubs(t, db)

	if err := MigrateAddSessionColumns(db); err != nil {
		t.Fatalf("first run: %v", err)
	}

	for _, table := range []string{"events", "wisp_events"} {
		has, err := columnExists(db, table, "session")
		if err != nil {
			t.Fatalf("check %s.session: %v", table, err)
		}
		if !has {
			t.Errorf("expected %s.session to exist after migration", table)
		}
	}
}

// TestMigrateAddSessionColumns_Idempotent verifies that running the migration
// twice does not error — the column-existence guard makes the second pass a
// no-op. Tables are seeded so the first pass actually exercises ADD COLUMN
// (otherwise the test would only prove the table-missing skip path twice).
func TestMigrateAddSessionColumns_Idempotent(t *testing.T) {
	db := openTestDoltBranch(t)
	createEventTableStubs(t, db)

	if err := MigrateAddSessionColumns(db); err != nil {
		t.Fatalf("first run: %v", err)
	}
	if err := MigrateAddSessionColumns(db); err != nil {
		t.Fatalf("second run (must be no-op): %v", err)
	}

	// Confirm the column is still present after the no-op pass.
	for _, table := range []string{"events", "wisp_events"} {
		has, err := columnExists(db, table, "session")
		if err != nil {
			t.Fatalf("check %s.session: %v", table, err)
		}
		if !has {
			t.Errorf("expected %s.session to still exist after idempotent re-run", table)
		}
	}
}

// TestMigrateAddSessionColumns_PreExistingColumnNoOp simulates environments
// where the session column was already added out-of-band before the migration
// runs. The motivating real-world case is Gas Town's
// internal/doltserver/wisps_migrate.go which independently CREATEd wisp_events
// on its shared Dolt server; we cover both events and wisp_events here so the
// loop's column-already-present branch is exercised on each leg.
func TestMigrateAddSessionColumns_PreExistingColumnNoOp(t *testing.T) {
	db := openTestDoltBranch(t)

	// Pre-create both tables WITH the session column already present, to
	// exercise the columnExists-returns-true skip branch on both loop iters.
	for _, ddl := range []string{
		`CREATE TABLE events (
			id CHAR(36) NOT NULL PRIMARY KEY,
			issue_id VARCHAR(255) NOT NULL,
			event_type VARCHAR(32) NOT NULL,
			actor VARCHAR(255) NOT NULL,
			session VARCHAR(255) NULL
		)`,
		`CREATE TABLE wisp_events (
			id CHAR(36) NOT NULL PRIMARY KEY,
			issue_id VARCHAR(255) NOT NULL,
			event_type VARCHAR(32) NOT NULL,
			actor VARCHAR(255) NOT NULL,
			session VARCHAR(255) NULL
		)`,
	} {
		if _, err := db.Exec(ddl); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}

	if err := MigrateAddSessionColumns(db); err != nil {
		t.Fatalf("migration must tolerate pre-existing column: %v", err)
	}
}
