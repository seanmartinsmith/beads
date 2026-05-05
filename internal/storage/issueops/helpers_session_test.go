package issueops

import (
	"context"
	"database/sql"

	"github.com/steveyegge/beads/internal/types"
)

// TestRecordEventInTable_CapturesSession_SignatureCheck verifies at compile
// time that RecordEventInTable accepts the expected number and types of
// parameters. There is no Dolt test infrastructure in this package, so the
// round-trip assertion (write + read-back) lives in the integration test suite
// in internal/storage/dolt/events_test.go (Phase 3). This compile assertion is
// the TDD first-step: it must fail to compile with the old 7-parameter
// signature, and compile cleanly once session is added as the 8th parameter.
//
// What this catches: removal of the session parameter, change of its type
// (e.g., to *string), or any arity change.
//
// What this does NOT catch: actor/session positional transposition, because
// both are string and the type system cannot distinguish two adjacent string
// parameters. Positional correctness is enforced by the SQL INSERT column
// order in helpers.go and by the round-trip tests in Phase 7.
//
// Intended parameter order: (ctx, tx, table, issueID, eventType, actor, session, newValue)
var _ func(context.Context, *sql.Tx, string, string, types.EventType, string, string, string) error = RecordEventInTable
