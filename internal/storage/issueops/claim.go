package issueops

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/steveyegge/beads/internal/storage"
	"github.com/steveyegge/beads/internal/types"
)

// ClaimResult holds the result of a ClaimIssueInTx call.
type ClaimResult struct {
	OldIssue *types.Issue
	IsWisp   bool
}

// ClaimIssueInTx atomically claims an issue using compare-and-swap semantics.
// It sets the assignee to actor and status to "in_progress" only if the issue
// is currently open and unassigned or already assigned to the same actor.
// Returns storage.ErrAlreadyClaimed if already claimed by a different user.
// Idempotent: re-claiming an in_progress issue by the same actor is a no-op
// success (supports agent retry workflows).
// Routes to the correct table (issues/wisps) automatically.
// The caller is responsible for Dolt versioning (DOLT_ADD/COMMIT) if needed.
//
// The session parameter records the Claude Code session that claimed this
// issue in claimed_by_session, mirroring the per-event attribution pattern
// established by closed_by_session. Overwrites on re-claim (last-writer-wins);
// the audit log remains the source of truth for full multi-session history.
//
//nolint:gosec // G201: table names come from WispTableRouting (hardcoded constants)
func ClaimIssueInTx(ctx context.Context, tx *sql.Tx, id string, actor, session string) (*ClaimResult, error) {
	isWisp := IsActiveWispInTx(ctx, tx, id)
	issueTable, _, eventTable, _ := WispTableRouting(isWisp)

	// Read old issue inside the transaction for event recording.
	oldIssue, err := GetIssueInTx(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get issue for claim: %w", err)
	}

	now := time.Now().UTC()

	// Conditional UPDATE: only succeeds while the issue is still claimable.
	// Also set started_at on first transition to in_progress (GH#2796); preserve
	// any existing value so re-claims don't overwrite the original start time.
	// claimed_by_session overwrites on every claim — last-writer-wins semantics.
	var (
		result sql.Result
	)
	if oldIssue.StartedAt == nil {
		result, err = tx.ExecContext(ctx, fmt.Sprintf(`
			UPDATE %s
			SET assignee = ?, status = 'in_progress', updated_at = ?, started_at = ?, claimed_by_session = ?
			WHERE id = ? AND status = 'open' AND (assignee = '' OR assignee IS NULL OR assignee = ?)
		`, issueTable), actor, now, now, session, id, actor)
	} else {
		result, err = tx.ExecContext(ctx, fmt.Sprintf(`
			UPDATE %s
			SET assignee = ?, status = 'in_progress', updated_at = ?, claimed_by_session = ?
			WHERE id = ? AND status = 'open' AND (assignee = '' OR assignee IS NULL OR assignee = ?)
		`, issueTable), actor, now, session, id, actor)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to claim issue: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		// Query current state inside the same transaction for consistency.
		var currentAssignee sql.NullString
		var currentStatus types.Status
		err := tx.QueryRowContext(ctx, fmt.Sprintf(
			`SELECT assignee, status FROM %s WHERE id = ?`, issueTable), id).Scan(&currentAssignee, &currentStatus)
		if err != nil {
			return nil, fmt.Errorf("failed to get current claim state: %w", err)
		}
		assignee := ""
		if currentAssignee.Valid {
			assignee = currentAssignee.String
		}
		// Idempotent: if already claimed in_progress by the same actor, treat as success.
		// This supports agent retry workflows where claim may be called multiple
		// times after transient failures (GH#8).
		if assignee == actor && currentStatus == types.StatusInProgress {
			return &ClaimResult{OldIssue: oldIssue, IsWisp: isWisp}, nil
		}
		if assignee != "" && assignee != actor {
			return nil, fmt.Errorf("%w by %s", storage.ErrAlreadyClaimed, assignee)
		}
		return nil, fmt.Errorf("%w: status %s", storage.ErrNotClaimable, currentStatus)
	}

	// Record the claim event.
	oldData, _ := json.Marshal(oldIssue)
	newUpdates := map[string]interface{}{
		"assignee": actor,
		"status":   "in_progress",
	}
	newData, _ := json.Marshal(newUpdates)

	if err := RecordFullEventInTable(ctx, tx, eventTable, id, "claimed", actor, string(oldData), string(newData)); err != nil {
		return nil, fmt.Errorf("failed to record claim event: %w", err)
	}

	return &ClaimResult{OldIssue: oldIssue, IsWisp: isWisp}, nil
}

// ClaimReadyIssueInTx claims the first currently ready issue matching filter in
// the same transaction that computes readiness. It returns nil when no matching
// ready issue can be claimed.
func ClaimReadyIssueInTx(
	ctx context.Context,
	tx *sql.Tx,
	filter types.WorkFilter,
	actor string,
	computeBlockedFn func(ctx context.Context, tx *sql.Tx, includeWisps bool) ([]string, error),
) (*types.Issue, error) {
	claimFilter := filter
	claimFilter.Status = types.StatusOpen
	claimFilter.Unassigned = true
	claimFilter.Assignee = nil
	claimFilter.Limit = 0

	readyIssues, err := GetReadyWorkInTx(ctx, tx, claimFilter, computeBlockedFn)
	if err != nil {
		return nil, err
	}
	for _, issue := range readyIssues {
		// Session attribution gap: ClaimReadyIssue path doesn't yet propagate
		// CLAUDE_SESSION_ID. bd ready --claim won't stamp claimed_by_session.
		// bd update --claim still does. Follow-up: thread session from cmd/bd/ready.go.
		if _, err := ClaimIssueInTx(ctx, tx, issue.ID, actor, ""); err != nil {
			if errors.Is(err, storage.ErrAlreadyClaimed) || errors.Is(err, storage.ErrNotClaimable) {
				continue
			}
			return nil, err
		}
		claimed, err := GetIssueInTx(ctx, tx, issue.ID)
		if err != nil {
			return nil, fmt.Errorf("get claimed issue: %w", err)
		}
		return claimed, nil
	}
	return nil, nil
}
