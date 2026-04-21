package issueops

import (
	"context"
	"database/sql"
	"encoding/json"
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
// currently has no assignee. Returns storage.ErrAlreadyClaimed if already
// claimed by a different user. Idempotent: re-claiming by the same actor is
// a no-op success (supports agent retry workflows).
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

	// Conditional UPDATE: only succeeds if assignee is currently empty.
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
			WHERE id = ? AND (assignee = '' OR assignee IS NULL)
		`, issueTable), actor, now, now, session, id)
	} else {
		result, err = tx.ExecContext(ctx, fmt.Sprintf(`
			UPDATE %s
			SET assignee = ?, status = 'in_progress', updated_at = ?, claimed_by_session = ?
			WHERE id = ? AND (assignee = '' OR assignee IS NULL)
		`, issueTable), actor, now, session, id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to claim issue: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		// Query current assignee inside the same transaction for consistency.
		var currentAssignee string
		err := tx.QueryRowContext(ctx, fmt.Sprintf(
			`SELECT assignee FROM %s WHERE id = ?`, issueTable), id).Scan(&currentAssignee)
		if err != nil {
			return nil, fmt.Errorf("failed to get current assignee: %w", err)
		}
		// Idempotent: if already claimed by the same actor, treat as success.
		// This supports agent retry workflows where claim may be called multiple
		// times after transient failures (GH#8).
		if currentAssignee == actor {
			return &ClaimResult{OldIssue: oldIssue, IsWisp: isWisp}, nil
		}
		return nil, fmt.Errorf("%w by %s", storage.ErrAlreadyClaimed, currentAssignee)
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
