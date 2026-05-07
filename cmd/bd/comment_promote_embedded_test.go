//go:build cgo

package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// bdComment runs "bd comments add" with the given args and returns stdout.
func bdComment(t *testing.T, bd, dir string, args ...string) string {
	t.Helper()
	fullArgs := append([]string{"comments", "add"}, args...)
	cmd := exec.Command(bd, fullArgs...)
	cmd.Dir = dir
	cmd.Env = bdEnv(dir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("bd comments add %s failed: %v\n%s", strings.Join(args, " "), err, out)
	}
	return string(out)
}

// bdCommentList runs "bd comments list" and returns stdout.
func bdCommentList(t *testing.T, bd, dir, issueID string) string {
	t.Helper()
	cmd := exec.Command(bd, "comments", "list", issueID)
	cmd.Dir = dir
	cmd.Env = bdEnv(dir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("bd comments list %s failed: %v\n%s", issueID, err, out)
	}
	return string(out)
}

func TestEmbeddedComments(t *testing.T) {
	if os.Getenv("BEADS_TEST_EMBEDDED_DOLT") != "1" {
		t.Skip("set BEADS_TEST_EMBEDDED_DOLT=1 to run embedded dolt integration tests")
	}
	t.Parallel()

	bd := buildEmbeddedBD(t)
	dir, beadsDir, _ := bdInit(t, bd, "--prefix", "cm")

	t.Run("add_and_list_comment", func(t *testing.T) {
		issue := bdCreate(t, bd, dir, "Comment target", "--type", "task")

		store := openStore(t, beadsDir, "cm")
		comment, err := store.AddIssueComment(t.Context(), issue.ID, "tester", "Hello world")
		if err != nil {
			t.Fatalf("AddIssueComment: %v", err)
		}
		if comment.Text != "Hello world" {
			t.Errorf("expected comment text 'Hello world', got %q", comment.Text)
		}
		if comment.Author != "tester" {
			t.Errorf("expected author 'tester', got %q", comment.Author)
		}
		if comment.ID == "" {
			t.Error("expected comment ID to be set")
		}

		// Verify via GetIssueComments.
		comments, err := store.GetIssueComments(t.Context(), issue.ID)
		if err != nil {
			t.Fatalf("GetIssueComments: %v", err)
		}
		found := false
		for _, c := range comments {
			if c.Text == "Hello world" {
				found = true
				break
			}
		}
		if !found {
			t.Error("expected to find 'Hello world' comment in GetIssueComments")
		}
	})

	t.Run("add_comment_event", func(t *testing.T) {
		issue := bdCreate(t, bd, dir, "Event comment target", "--type", "task")

		store := openStore(t, beadsDir, "cm")
		if err := store.AddComment(t.Context(), issue.ID, "actor", "", "A comment event"); err != nil {
			t.Fatalf("AddComment: %v", err)
		}
	})

	t.Run("add_comment_nonexistent_issue", func(t *testing.T) {
		store := openStore(t, beadsDir, "cm")
		_, err := store.AddIssueComment(t.Context(), "cm-nonexistent999", "tester", "nope")
		if err == nil {
			t.Error("expected error for nonexistent issue")
		}
	})
}

func TestEmbeddedPromote(t *testing.T) {
	if os.Getenv("BEADS_TEST_EMBEDDED_DOLT") != "1" {
		t.Skip("set BEADS_TEST_EMBEDDED_DOLT=1 to run embedded dolt integration tests")
	}
	t.Parallel()

	bd := buildEmbeddedBD(t)
	dir, beadsDir, _ := bdInit(t, bd, "--prefix", "pm")

	t.Run("promote_wisp", func(t *testing.T) {
		// Create an ephemeral issue (routes to wisps table).
		issue := bdCreate(t, bd, dir, "Promote me", "--ephemeral")

		store := openStore(t, beadsDir, "pm")

		// Verify it's a wisp before promote.
		got, err := store.GetIssue(t.Context(), issue.ID)
		if err != nil {
			t.Fatalf("GetIssue before promote: %v", err)
		}
		if !got.Ephemeral {
			t.Skip("issue is not ephemeral; cannot test promote")
		}

		// Promote.
		if err := store.PromoteFromEphemeral(t.Context(), issue.ID, "tester"); err != nil {
			t.Fatalf("PromoteFromEphemeral: %v", err)
		}

		// Verify it's now permanent.
		got, err = store.GetIssue(t.Context(), issue.ID)
		if err != nil {
			t.Fatalf("GetIssue after promote: %v", err)
		}
		if got.Ephemeral {
			t.Error("expected issue to be non-ephemeral after promote")
		}
	})

	t.Run("promote_nonexistent_wisp", func(t *testing.T) {
		store := openStore(t, beadsDir, "pm")
		err := store.PromoteFromEphemeral(t.Context(), "pm-nonexistent999", "tester")
		if err == nil {
			t.Error("expected error for nonexistent wisp")
		}
	})
}

// TestSessionAttribution_Comment verifies that the session-aware comment-event
// path (Storage.AddComment → AddCommentEventInTx → INSERT INTO events) writes
// the supplied session value to events.session. Two sub-tests cover the two
// surfaces that drive AddComment in PR1's scope:
//
//  1. Storage-method directly — proves the lowest-level contract.
//  2. bd promote cobra path — proves the bd-8kt fix (resolveSession() threading
//     into the compact/promote AddComment call sites) actually round-trips.
//
// Note: bd comment (the user-facing comment cobra command) does NOT go through
// Storage.AddComment. It goes through AddIssueComment, which writes to the
// comments table only and is not session-aware by API design. That's a
// separate structural gap (closer to outcome C than bd-fwb shape) and is not
// in PR1's scope to close.
//
// Refs: bd-edi (PR1 Phase 7), bd-8kt (the call-site escapees this test
// validates the fix for), gastownhall/beads#3583.
func TestSessionAttribution_Comment(t *testing.T) {
	if os.Getenv("BEADS_TEST_EMBEDDED_DOLT") != "1" {
		t.Skip("set BEADS_TEST_EMBEDDED_DOLT=1 to run embedded dolt integration tests")
	}
	t.Parallel()

	bd := buildEmbeddedBD(t)
	dir, beadsDir, _ := bdInit(t, bd, "--prefix", "sc")

	t.Run("storage_method_round_trips_session", func(t *testing.T) {
		// Directly invoke the session-aware storage method against a fresh
		// permanent issue. Asserts events.session for the resulting commented
		// event row. This is the lowest-level contract: callers that thread
		// session through must see it land.
		issue := bdCreate(t, bd, dir, "AddComment storage round-trip", "--type", "task")
		store := openStore(t, beadsDir, "sc")
		if err := store.AddComment(t.Context(), issue.ID, "tester", "storage-sess", "round-trip event"); err != nil {
			t.Fatalf("AddComment: %v", err)
		}
		got := queryEventSessionSQL(t, beadsDir, issue.ID, "commented")
		if got != "storage-sess" {
			t.Errorf("expected events.session=storage-sess, got %q", got)
		}
	})

	t.Run("bd_promote_cobra_round_trips_session_via_flag", func(t *testing.T) {
		// bd promote runs PromoteFromEphemeral and then writes a 'Promoted
		// from wisp to permanent bead' comment via Storage.AddComment. Pre-
		// bd-8kt that call site hardcoded "" for session; post-bd-8kt it
		// passes resolveSession(). This sub-test exercises the now-fixed
		// path and asserts the comment event carries --session.
		wisp := bdCreate(t, bd, dir, "Promote with session", "--ephemeral")
		cmd := exec.Command(bd, "promote", wisp.ID, "--session", "promote-sess")
		cmd.Dir = dir
		cmd.Env = bdEnv(dir)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("bd promote --session failed: %v\n%s", err, out)
		}
		got := queryEventSessionSQL(t, beadsDir, wisp.ID, "commented")
		if got != "promote-sess" {
			t.Errorf("expected events.session=promote-sess, got %q", got)
		}
	})
}

// TestSessionAttribution_PromotePreservesSession verifies that the session
// column on existing event rows survives the wisp_events → events boundary
// crossing during promote. This is the test that would have caught the gap
// fixed by 0d78608fe ('feat(session): preserve session column on promote/
// demote round-trip'): pre-fix, the INSERT...SELECT bulk copy didn't include
// the session column in either projection or column list, silently zeroing
// out attribution on every event row for the promoted issue.
//
// The Phase 3 closure audit missed this because it grepped INSERT...VALUES
// patterns and didn't catch INSERT...SELECT round-trips. The Phase 4
// verification grep caught it. This test would have caught it earlier.
//
// Demote is not exercised here because DemoteToWisp is dolt-only — the
// embedded backend supports PromoteFromEphemeral but not its inverse, so
// demote testing requires a dolt-backed test which is out of scope for the
// embedded test suite. Phase 4.1 fix at 0d78608fe applies to both
// directions; promote-only coverage validates the symmetric contract
// pattern. Demote-side coverage is tracked as a follow-up.
//
// Refs: bd-edi (PR1 Phase 7), 0d78608fe (the Phase 4.1 fix this would
// have caught).
func TestSessionAttribution_PromotePreservesSession(t *testing.T) {
	if os.Getenv("BEADS_TEST_EMBEDDED_DOLT") != "1" {
		t.Skip("set BEADS_TEST_EMBEDDED_DOLT=1 to run embedded dolt integration tests")
	}
	t.Parallel()

	bd := buildEmbeddedBD(t)
	dir, beadsDir, _ := bdInit(t, bd, "--prefix", "spr")

	t.Run("label_event_session_survives_promote", func(t *testing.T) {
		// Create a wisp, label it with --session (writing a label_added
		// event to wisp_events with session attribution), then promote.
		// The promote operation must preserve session as it copies the
		// event row from wisp_events to events.
		wisp := bdCreate(t, bd, dir, "promote preserves session", "--ephemeral")

		// Pre-promote: add a label with session attribution. This writes
		// to wisp_events because the issue is still ephemeral.
		labelCmd := exec.Command(bd, "label", "add", wisp.ID, "preserve-test",
			"--session", "wisp-era-session")
		labelCmd.Dir = dir
		labelCmd.Env = bdEnv(dir)
		if out, err := labelCmd.CombinedOutput(); err != nil {
			t.Fatalf("bd label add failed: %v\n%s", err, out)
		}

		// Sanity: the event lives in wisp_events with the expected session.
		// queryEventSessionSQL falls back to wisp_events when events is
		// empty, so it picks up the wisp-era row directly.
		preGot := queryEventSessionSQL(t, beadsDir, wisp.ID, "label_added")
		if preGot != "wisp-era-session" {
			t.Fatalf("pre-promote: expected wisp_events.session=wisp-era-session, got %q (test setup wrong, not a promote bug)", preGot)
		}

		// Promote: triggers the INSERT...SELECT from wisp_events to events.
		promoteCmd := exec.Command(bd, "promote", wisp.ID)
		promoteCmd.Dir = dir
		promoteCmd.Env = bdEnv(dir)
		if out, err := promoteCmd.CombinedOutput(); err != nil {
			t.Fatalf("bd promote failed: %v\n%s", err, out)
		}

		// Post-promote: the label_added row now lives in events; assert
		// session was preserved across the boundary. The Phase 4.1 fix
		// (0d78608fe) added session to both the SELECT projection and the
		// INSERT column list; without it, this assertion would see "".
		postGot := queryEventSessionSQL(t, beadsDir, wisp.ID, "label_added")
		if postGot != "wisp-era-session" {
			t.Errorf("post-promote: expected events.session=wisp-era-session (preserved across boundary), got %q", postGot)
		}
	})
}
