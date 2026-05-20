package issueops

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/steveyegge/beads/internal/storage"
	"github.com/steveyegge/beads/internal/types"
)

// GetIssueInTx retrieves a single issue by ID within an existing transaction,
// including its labels. Automatically routes to the wisps/wisp_labels tables
// if the ID is an active wisp. Returns storage.ErrNotFound (wrapped) if the
// issue does not exist in either table.
//
// For closed issues, ClosedBySession is derived from the most recent matching
// events row (event_type='closed') rather than the issues.closed_by_session
// column. The column itself is retained for backward compat but is no longer
// the source of truth (bd-edi PR1 / ADR-0003). Closed issues without a
// matching events row (e.g. reaper bypass) fall back to the column value, so
// pre-PR1 data remains readable.
func GetIssueInTx(ctx context.Context, tx *sql.Tx, id string) (*types.Issue, error) {
	issue, err := getIssueFromTableInTx(ctx, tx, "issues", "labels", "events", id)
	if err == nil {
		return issue, nil
	}
	if !errors.Is(err, storage.ErrNotFound) {
		return nil, err
	}

	issue, err = getIssueFromTableInTx(ctx, tx, "wisps", "wisp_labels", "wisp_events", id)
	if err == nil {
		return issue, nil
	}
	if errors.Is(err, storage.ErrNotFound) {
		return nil, fmt.Errorf("%w: issue %s", storage.ErrNotFound, id)
	}
	return nil, err
}

func getIssueFromTableInTx(ctx context.Context, tx *sql.Tx, issueTable, labelTable, eventTable, id string) (*types.Issue, error) {
	//nolint:gosec // G201: issueTable is a hardcoded literal supplied by GetIssueInTx ("issues" or "wisps")
	row := tx.QueryRowContext(ctx, fmt.Sprintf(`SELECT %s FROM %s WHERE id = ?`, IssueSelectColumns, issueTable), id)
	issue, err := ScanIssueFrom(row)
	if err == sql.ErrNoRows || isTableNotExistError(err) {
		return nil, storage.ErrNotFound
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
		derived, ok, err := lookupClosedBySession(ctx, tx, eventTable, id)
		if err != nil {
			return nil, fmt.Errorf("get derived closed_by_session: %w", err)
		}
		if ok {
			issue.ClosedBySession = derived
		}
		// If no events row exists (ok=false), retain the column-scanned value
		// to handle pre-PR1 data and reaper-bypass writes.
	}

	return issue, nil
}

// lookupClosedBySession returns the session attribution for the most recent
// 'closed' event for issueID in eventTable, plus a flag indicating whether a
// matching event row exists. ok=false means no row; ok=true with empty string
// means the row exists but session was unset (NULL).
func lookupClosedBySession(ctx context.Context, tx *sql.Tx, eventTable, issueID string) (string, bool, error) {
	var derived sql.NullString
	//nolint:gosec // G201: eventTable is a hardcoded literal supplied by callers ("events" or "wisp_events")
	err := tx.QueryRowContext(ctx, fmt.Sprintf(
		`SELECT session FROM %s WHERE issue_id = ? AND event_type = 'closed' ORDER BY created_at DESC, id DESC LIMIT 1`,
		eventTable,
	), issueID).Scan(&derived)
	if err == sql.ErrNoRows || isTableNotExistError(err) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	if !derived.Valid {
		return "", true, nil
	}
	return derived.String, true, nil
}
