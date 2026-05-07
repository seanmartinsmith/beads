package issueops

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/steveyegge/beads/internal/storage"
	"github.com/steveyegge/beads/internal/types"
)

// GetIssueInTx retrieves a single issue by ID within an existing transaction,
// including its labels. Automatically routes to the wisps/wisp_labels tables
// if the ID is an active wisp. Returns storage.ErrNotFound (wrapped) if the
// issue does not exist in either table.
//
// For closed issues, ClosedBySession is derived from the most recent
// matching events row (event_type='closed') rather than the
// issues.closed_by_session column. The column itself is retained for
// backward compat but is no longer the source of truth (Phase 6 of bd-edi
// PR1). Closed issues without a matching events row (e.g. reaper bypass)
// surface ClosedBySession = "".
func GetIssueInTx(ctx context.Context, tx *sql.Tx, id string) (*types.Issue, error) {
	isWisp := IsActiveWispInTx(ctx, tx, id)
	issueTable, labelTable, eventTable, _ := WispTableRouting(isWisp)

	//nolint:gosec // G201: issueTable is from WispTableRouting ("issues" or "wisps")
	row := tx.QueryRowContext(ctx, fmt.Sprintf(`SELECT %s FROM %s WHERE id = ?`, IssueSelectColumns, issueTable), id)
	issue, err := ScanIssueFrom(row)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("%w: issue %s", storage.ErrNotFound, id)
	}
	if err != nil {
		return nil, fmt.Errorf("get issue: %w", err)
	}

	// Fetch labels in the same transaction to avoid MaxOpenConns=1 deadlock.
	labels, err := GetLabelsInTx(ctx, tx, labelTable, id)
	if err != nil {
		return nil, fmt.Errorf("get issue labels: %w", err)
	}
	issue.Labels = labels

	if issue.Status == types.StatusClosed {
		derived, err := lookupClosedBySession(ctx, tx, eventTable, id)
		if err != nil {
			return nil, fmt.Errorf("get derived closed_by_session: %w", err)
		}
		issue.ClosedBySession = derived
	}

	return issue, nil
}

// lookupClosedBySession returns the session attribution for the most recent
// 'closed' event for issueID in eventTable. Empty string if no matching
// event row exists (e.g. issue closed by a path that bypasses event recording,
// such as the Gas Town reaper or the doltTransaction txn-scoped close path
// tracked at bd-3pc).
func lookupClosedBySession(ctx context.Context, tx *sql.Tx, eventTable, issueID string) (string, error) {
	var derived sql.NullString
	//nolint:gosec // G201: eventTable is from WispTableRouting ("events" or "wisp_events")
	err := tx.QueryRowContext(ctx, fmt.Sprintf(
		`SELECT session FROM %s WHERE issue_id = ? AND event_type = 'closed' ORDER BY id DESC LIMIT 1`,
		eventTable,
	), issueID).Scan(&derived)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	if !derived.Valid {
		return "", nil
	}
	return derived.String, nil
}
