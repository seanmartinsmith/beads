package issueops

import (
	"context"
	"database/sql"

	"github.com/steveyegge/beads/internal/types"
)

// Compile-time assertion that RecordFullEventInTable carries the session
// parameter in the expected position.
//
// Intended parameter order: (ctx, tx, table, issueID, eventType, actor, session, oldValue, newValue)
var _ func(context.Context, *sql.Tx, string, string, types.EventType, string, string, string, string) error = RecordFullEventInTable
