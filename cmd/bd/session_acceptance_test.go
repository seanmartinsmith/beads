//go:build cgo

package main

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/steveyegge/beads/internal/configfile"
	"github.com/steveyegge/beads/internal/storage/embeddeddolt"
)

// queryEventSession returns events.session for the most recent event row of
// the given (issue_id, event_type), or empty if none exists. Used to verify
// that session attribution reaches the events table.
//
// Tries the events table first, then wisp_events.
func queryEventSession(t *testing.T, beadsDir, issueID, eventType string) string {
	t.Helper()
	dataDir := filepath.Join(beadsDir, "embeddeddolt")
	cfg, _ := configfile.Load(beadsDir)
	database := ""
	if cfg != nil {
		database = cfg.GetDoltDatabase()
	}
	db, cleanup, err := embeddeddolt.OpenSQL(t.Context(), dataDir, database, "main")
	if err != nil {
		t.Fatalf("OpenSQL: %v", err)
	}
	defer cleanup()
	var session string
	err = db.QueryRowContext(context.Background(),
		"SELECT COALESCE(session, '') FROM events WHERE issue_id = ? AND event_type = ? ORDER BY created_at DESC, id DESC LIMIT 1",
		issueID, eventType).Scan(&session)
	if err != nil {
		// Fall back to wisp_events for promoted wisps / wisp-only issues.
		err = db.QueryRowContext(context.Background(),
			"SELECT COALESCE(session, '') FROM wisp_events WHERE issue_id = ? AND event_type = ? ORDER BY created_at DESC, id DESC LIMIT 1",
			issueID, eventType).Scan(&session)
		if err != nil {
			t.Fatalf("query session for issue=%s event_type=%s: %v", issueID, eventType, err)
		}
	}
	return session
}

// TestSessionAttribution_CloseRoundTrip verifies the resolver chain's
// precedence ordering survives end-to-end: when multiple sources are set,
// the correct source wins ALL the way to events.session, not just at the
// resolveSession() return point.
//
// TestResolveSession asserts the resolver returns the right value in
// isolation. This test adds the contract that the resolver's output reaches
// the storage layer through a real cobra command. bd close is the simplest
// session-aware path.
//
// Refs: bd-edi / ADR-0003 / bd-bkzj.
func TestSessionAttribution_CloseRoundTrip(t *testing.T) {
	if os.Getenv("BEADS_TEST_EMBEDDED_DOLT") != "1" {
		t.Skip("set BEADS_TEST_EMBEDDED_DOLT=1 to run embedded dolt integration tests")
	}
	t.Parallel()

	bd := buildEmbeddedBD(t)
	dir, beadsDir, _ := bdInit(t, bd, "--prefix", "sa")

	closeAndQuery := func(t *testing.T, issueID string, args []string, extraEnv ...string) string {
		t.Helper()
		fullArgs := append([]string{"close", issueID}, args...)
		cmd := exec.Command(bd, fullArgs...)
		cmd.Dir = dir
		env := bdEnv(dir)
		env = append(env, extraEnv...)
		cmd.Env = env
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("bd close %s failed: %v\n%s", strings.Join(fullArgs, " "), err, out)
		}
		return queryEventSession(t, beadsDir, issueID, "closed")
	}

	t.Run("flag_beats_env", func(t *testing.T) {
		issue := bdCreate(t, bd, dir, "flag beats env", "--type", "task")
		got := closeAndQuery(t, issue.ID,
			[]string{"--session", "flag-wins"},
			"BD_CORE_CAPTURE_SESSION=true",
			"BEADS_SESSION_ID=beads-loses",
			"CLAUDE_CODE_SESSION_ID=claude-code-loses",
			"CLAUDE_SESSION_ID=claude-loses",
		)
		if got != "flag-wins" {
			t.Errorf("expected events.session=flag-wins, got %q", got)
		}
	})

	t.Run("beads_env_beats_claude_code", func(t *testing.T) {
		issue := bdCreate(t, bd, dir, "beads beats claude_code", "--type", "task")
		got := closeAndQuery(t, issue.ID,
			nil,
			"BD_CORE_CAPTURE_SESSION=true",
			"BEADS_SESSION_ID=beads-wins",
			"CLAUDE_CODE_SESSION_ID=claude-code-loses",
		)
		if got != "beads-wins" {
			t.Errorf("expected events.session=beads-wins, got %q", got)
		}
	})

	t.Run("claude_code_beats_claude_legacy", func(t *testing.T) {
		issue := bdCreate(t, bd, dir, "claude_code beats claude_legacy", "--type", "task")
		got := closeAndQuery(t, issue.ID,
			nil,
			"BD_CORE_CAPTURE_SESSION=true",
			"CLAUDE_CODE_SESSION_ID=claude-code-wins",
			"CLAUDE_SESSION_ID=claude-legacy-loses",
		)
		if got != "claude-code-wins" {
			t.Errorf("expected events.session=claude-code-wins, got %q", got)
		}
	})

	t.Run("default_ignores_env_vars", func(t *testing.T) {
		issue := bdCreate(t, bd, dir, "default ignores env", "--type", "task")
		got := closeAndQuery(t, issue.ID,
			nil,
			// No BD_CORE_CAPTURE_SESSION.
			"BEADS_SESSION_ID=beads-ignored",
			"CLAUDE_CODE_SESSION_ID=claude-code-ignored",
			"CLAUDE_SESSION_ID=claude-legacy-ignored",
		)
		if got != "" {
			t.Errorf("expected events.session='' (default ignores env vars), got %q", got)
		}
	})

	t.Run("flag_always_honored_without_opt_in", func(t *testing.T) {
		issue := bdCreate(t, bd, dir, "flag without opt-in", "--type", "task")
		got := closeAndQuery(t, issue.ID,
			[]string{"--session", "flag-no-opt-in"},
			"BEADS_SESSION_ID=ignored",
			"CLAUDE_CODE_SESSION_ID=ignored",
			"CLAUDE_SESSION_ID=ignored",
		)
		if got != "flag-no-opt-in" {
			t.Errorf("expected events.session=flag-no-opt-in, got %q", got)
		}
	})
}

// TestSessionAttribution_ReadyClaim is the architectural-win acceptance
// test: bd ready --claim must capture session in the claim event. This was
// the canonical gap that bd-edi closes (the #3578 path didn't thread
// session at all on the reference branch — fresh-impl extends ClaimIssue's
// signature to fix it).
func TestSessionAttribution_ReadyClaim(t *testing.T) {
	if os.Getenv("BEADS_TEST_EMBEDDED_DOLT") != "1" {
		t.Skip("set BEADS_TEST_EMBEDDED_DOLT=1 to run embedded dolt integration tests")
	}
	t.Parallel()

	bd := buildEmbeddedBD(t)
	dir, beadsDir, _ := bdInit(t, bd, "--prefix", "rc")

	issue := bdCreate(t, bd, dir, "claim-me", "--type", "task")

	cmd := exec.Command(bd, "ready", "--claim", "--session", "ready-claim-session", "--json")
	cmd.Dir = dir
	cmd.Env = bdEnv(dir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("bd ready --claim failed: %v\n%s", err, out)
	}

	got := queryEventSession(t, beadsDir, issue.ID, "claimed")
	if got != "ready-claim-session" {
		t.Errorf("expected events.session=ready-claim-session on claimed event, got %q", got)
	}
}
