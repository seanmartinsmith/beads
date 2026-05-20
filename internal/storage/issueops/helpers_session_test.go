package issueops

import (
	"context"
	"database/sql"

	"github.com/steveyegge/beads/internal/types"
)

// Compile-time assertion that RecordEventInTable carries the session parameter
// in the expected position. Round-trip behavior is asserted in the
// integration tests in internal/storage/dolt.
//
// Intended parameter order: (ctx, tx, table, issueID, eventType, actor, session, newValue)
var _ func(context.Context, *sql.Tx, string, string, types.EventType, string, string, string) error = RecordEventInTable
