package issueops

import (
	"context"
	"database/sql"

	"github.com/steveyegge/beads/internal/types"
)

// TestRecordFullEventInTable_CapturesSession_SignatureCheck verifies at compile
// time that RecordFullEventInTable accepts a session parameter between actor and
// oldValue. There is no Dolt test infrastructure in this package, so the
// round-trip assertion (write + read-back) lives in the integration test suite
// in internal/storage/dolt/events_test.go (Phase 3). This compile assertion is
// the TDD first-step: it must fail to compile with the old 8-parameter
// signature, and compile cleanly once session is added as the 9th parameter.
//
// Parameter order: (ctx, tx, table, issueID, eventType, actor, session, oldValue, newValue)
var _ func(context.Context, *sql.Tx, string, string, types.EventType, string, string, string, string) error = RecordFullEventInTable
