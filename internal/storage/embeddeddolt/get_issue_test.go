//go:build cgo

package embeddeddolt_test

import (
	"errors"
	"testing"

	"github.com/steveyegge/beads/internal/storage"
	"github.com/steveyegge/beads/internal/types"
)

func TestGetIssue(t *testing.T) {
	skipUnlessEmbeddedDolt(t)

	t.Run("round_trip", func(t *testing.T) {
		te := newTestEnv(t, "gi")
		ctx := t.Context()

		issue := &types.Issue{
			ID:          "gi-test1",
			Title:       "Round trip test",
			Description: "A test description",
			Status:      types.StatusOpen,
			Priority:    1,
			IssueType:   types.TypeBug,
			Assignee:    "alice",
		}
		if err := te.store.CreateIssue(ctx, issue, "tester"); err != nil {
			t.Fatalf("CreateIssue: %v", err)
		}

		got, err := te.store.GetIssue(ctx, "gi-test1")
		if err != nil {
			t.Fatalf("GetIssue: %v", err)
		}
		if got.ID != "gi-test1" {
			t.Errorf("ID: got %q, want %q", got.ID, "gi-test1")
		}
		if got.Title != "Round trip test" {
			t.Errorf("Title: got %q, want %q", got.Title, "Round trip test")
		}
		if got.Description != "A test description" {
			t.Errorf("Description: got %q, want %q", got.Description, "A test description")
		}
		if got.Priority != 1 {
			t.Errorf("Priority: got %d, want 1", got.Priority)
		}
		if got.IssueType != types.TypeBug {
			t.Errorf("IssueType: got %q, want %q", got.IssueType, types.TypeBug)
		}
		if got.Assignee != "alice" {
			t.Errorf("Assignee: got %q, want %q", got.Assignee, "alice")
		}
	})

	t.Run("not_found", func(t *testing.T) {
		te := newTestEnv(t, "nf")
		ctx := t.Context()

		_, err := te.store.GetIssue(ctx, "nf-nonexistent")
		if err == nil {
			t.Fatal("expected error for non-existent issue")
		}
		if !errors.Is(err, storage.ErrNotFound) {
			t.Errorf("expected ErrNotFound, got: %v", err)
		}
	})

	t.Run("closed_by_session_derived_from_events", func(t *testing.T) {
		// Phase 6 of bd-edi PR1: ClosedBySession is derived from the most
		// recent 'closed' events row, not the issues.closed_by_session column.
		te := newTestEnv(t, "cs")
		ctx := t.Context()

		issue := &types.Issue{
			ID:        "cs-closed",
			Title:     "Closed with session",
			Status:    types.StatusOpen,
			Priority:  2,
			IssueType: types.TypeTask,
		}
		if err := te.store.CreateIssue(ctx, issue, "tester"); err != nil {
			t.Fatalf("CreateIssue: %v", err)
		}
		if err := te.store.CloseIssue(ctx, "cs-closed", "done", "tester", "session-from-close"); err != nil {
			t.Fatalf("CloseIssue: %v", err)
		}

		got, err := te.store.GetIssue(ctx, "cs-closed")
		if err != nil {
			t.Fatalf("GetIssue: %v", err)
		}
		if got.Status != types.StatusClosed {
			t.Fatalf("Status: got %q, want closed", got.Status)
		}
		if got.ClosedBySession != "session-from-close" {
			t.Errorf("ClosedBySession: got %q, want %q (must be derived from events.session, not issues.closed_by_session column)", got.ClosedBySession, "session-from-close")
		}
	})

	t.Run("closed_by_session_open_issue_skips_lookup", func(t *testing.T) {
		// Open issues should never surface ClosedBySession; the derivation is
		// gated on Status == closed in GetIssueInTx.
		te := newTestEnv(t, "co")
		ctx := t.Context()

		issue := &types.Issue{
			ID:        "co-open",
			Title:     "Open issue",
			Status:    types.StatusOpen,
			Priority:  2,
			IssueType: types.TypeTask,
		}
		if err := te.store.CreateIssue(ctx, issue, "tester"); err != nil {
			t.Fatalf("CreateIssue: %v", err)
		}

		got, err := te.store.GetIssue(ctx, "co-open")
		if err != nil {
			t.Fatalf("GetIssue: %v", err)
		}
		if got.ClosedBySession != "" {
			t.Errorf("ClosedBySession on open issue: got %q, want empty", got.ClosedBySession)
		}
	})

	t.Run("closed_by_session_no_event_row_returns_empty", func(t *testing.T) {
		// Reaper-bypass / doltTransaction.CloseIssue (bd-3pc) pattern: an issue
		// is marked closed via a raw UPDATE that bypasses the event-recording
		// helper. No 'closed' event row exists. Per the launch brief §3
		// absence contract, the derived ClosedBySession surfaces as empty.
		te := newTestEnv(t, "cn")
		ctx := t.Context()

		issue := &types.Issue{
			ID:        "cn-reaper",
			Title:     "Reaper-style close",
			Status:    types.StatusOpen,
			Priority:  2,
			IssueType: types.TypeTask,
		}
		if err := te.store.CreateIssue(ctx, issue, "tester"); err != nil {
			t.Fatalf("CreateIssue: %v", err)
		}

		// Bypass CloseIssue: directly mark closed without recording an event.
		te.exec(t, ctx,
			"UPDATE issues SET status = 'closed', closed_at = NOW() WHERE id = ?",
			"cn-reaper")

		got, err := te.store.GetIssue(ctx, "cn-reaper")
		if err != nil {
			t.Fatalf("GetIssue: %v", err)
		}
		if got.Status != types.StatusClosed {
			t.Fatalf("Status: got %q, want closed", got.Status)
		}
		if got.ClosedBySession != "" {
			t.Errorf("ClosedBySession with no close event: got %q, want empty", got.ClosedBySession)
		}
	})

	t.Run("closed_by_session_null_event_session_returns_empty", func(t *testing.T) {
		// Pre-PR1 close pattern: the events row exists but session is NULL
		// (closed before the bd-edi events.session column was wired up, OR
		// closed by a path that doesn't carry session attribution). Per the
		// launch brief §3 absence contract, derived ClosedBySession is empty
		// when the events.session column is NULL.
		te := newTestEnv(t, "cz")
		ctx := t.Context()

		issue := &types.Issue{
			ID:        "cz-pre",
			Title:     "Closed pre-PR1",
			Status:    types.StatusOpen,
			Priority:  2,
			IssueType: types.TypeTask,
		}
		if err := te.store.CreateIssue(ctx, issue, "tester"); err != nil {
			t.Fatalf("CreateIssue: %v", err)
		}
		if err := te.store.CloseIssue(ctx, "cz-pre", "done", "tester", "session-will-be-nulled"); err != nil {
			t.Fatalf("CloseIssue: %v", err)
		}

		// Simulate pre-PR1 state: NULL out the session on the close event row.
		te.exec(t, ctx,
			"UPDATE events SET session = NULL WHERE issue_id = ? AND event_type = 'closed'",
			"cz-pre")

		got, err := te.store.GetIssue(ctx, "cz-pre")
		if err != nil {
			t.Fatalf("GetIssue: %v", err)
		}
		if got.ClosedBySession != "" {
			t.Errorf("ClosedBySession with NULL events.session: got %q, want empty", got.ClosedBySession)
		}
	})

	t.Run("includes_labels", func(t *testing.T) {
		te := newTestEnv(t, "il")
		ctx := t.Context()

		issue := &types.Issue{
			ID:        "il-labeled",
			Title:     "Labeled issue",
			Status:    types.StatusOpen,
			Priority:  2,
			IssueType: types.TypeTask,
		}
		if err := te.store.CreateIssue(ctx, issue, "tester"); err != nil {
			t.Fatalf("CreateIssue: %v", err)
		}
		if err := te.store.AddLabel(ctx, "il-labeled", "bug", "tester", ""); err != nil {
			t.Fatalf("AddLabel: %v", err)
		}
		if err := te.store.AddLabel(ctx, "il-labeled", "urgent", "tester", ""); err != nil {
			t.Fatalf("AddLabel: %v", err)
		}
		if err := te.store.Commit(ctx, "add labels"); err != nil {
			t.Fatalf("Commit: %v", err)
		}

		got, err := te.store.GetIssue(ctx, "il-labeled")
		if err != nil {
			t.Fatalf("GetIssue: %v", err)
		}
		if len(got.Labels) != 2 {
			t.Fatalf("Labels: got %d, want 2", len(got.Labels))
		}
		// Labels should be sorted
		if got.Labels[0] != "bug" || got.Labels[1] != "urgent" {
			t.Errorf("Labels: got %v, want [bug urgent]", got.Labels)
		}
	})
}

func TestGetLabels(t *testing.T) {
	skipUnlessEmbeddedDolt(t)

	t.Run("empty", func(t *testing.T) {
		te := newTestEnv(t, "gl")
		ctx := t.Context()

		issue := &types.Issue{
			ID:        "gl-nolabels",
			Title:     "No labels",
			Status:    types.StatusOpen,
			Priority:  2,
			IssueType: types.TypeTask,
		}
		if err := te.store.CreateIssue(ctx, issue, "tester"); err != nil {
			t.Fatalf("CreateIssue: %v", err)
		}

		labels, err := te.store.GetLabels(ctx, "gl-nolabels")
		if err != nil {
			t.Fatalf("GetLabels: %v", err)
		}
		if len(labels) != 0 {
			t.Errorf("expected empty labels, got %v", labels)
		}
	})

	t.Run("sorted", func(t *testing.T) {
		te := newTestEnv(t, "gs")
		ctx := t.Context()

		issue := &types.Issue{
			ID:        "gs-sorted",
			Title:     "Sorted labels",
			Status:    types.StatusOpen,
			Priority:  2,
			IssueType: types.TypeTask,
		}
		if err := te.store.CreateIssue(ctx, issue, "tester"); err != nil {
			t.Fatalf("CreateIssue: %v", err)
		}
		for _, l := range []string{"zebra", "alpha", "middle"} {
			if err := te.store.AddLabel(ctx, "gs-sorted", l, "tester", ""); err != nil {
				t.Fatalf("AddLabel(%s): %v", l, err)
			}
		}
		if err := te.store.Commit(ctx, "add labels"); err != nil {
			t.Fatalf("Commit: %v", err)
		}

		labels, err := te.store.GetLabels(ctx, "gs-sorted")
		if err != nil {
			t.Fatalf("GetLabels: %v", err)
		}
		want := []string{"alpha", "middle", "zebra"}
		if len(labels) != len(want) {
			t.Fatalf("got %v, want %v", labels, want)
		}
		for i, l := range labels {
			if l != want[i] {
				t.Errorf("labels[%d]: got %q, want %q", i, l, want[i])
			}
		}
	})
}

func TestAddLabel(t *testing.T) {
	skipUnlessEmbeddedDolt(t)

	t.Run("idempotent", func(t *testing.T) {
		te := newTestEnv(t, "al")
		ctx := t.Context()

		issue := &types.Issue{
			ID:        "al-idem",
			Title:     "Idempotent label",
			Status:    types.StatusOpen,
			Priority:  2,
			IssueType: types.TypeTask,
		}
		if err := te.store.CreateIssue(ctx, issue, "tester"); err != nil {
			t.Fatalf("CreateIssue: %v", err)
		}
		// Add same label twice — should not error.
		if err := te.store.AddLabel(ctx, "al-idem", "dup", "tester", ""); err != nil {
			t.Fatalf("AddLabel (first): %v", err)
		}
		if err := te.store.AddLabel(ctx, "al-idem", "dup", "tester", ""); err != nil {
			t.Fatalf("AddLabel (second): %v", err)
		}
		if err := te.store.Commit(ctx, "add labels"); err != nil {
			t.Fatalf("Commit: %v", err)
		}

		labels, err := te.store.GetLabels(ctx, "al-idem")
		if err != nil {
			t.Fatalf("GetLabels: %v", err)
		}
		if len(labels) != 1 {
			t.Errorf("expected 1 label, got %v", labels)
		}
	})
}
