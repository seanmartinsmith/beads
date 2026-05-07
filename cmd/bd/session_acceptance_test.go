//go:build cgo

package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestSessionAttribution_PrecedenceAndOptIn is the acceptance-level proof that
// the resolver chain's precedence ordering survives end-to-end: when multiple
// sources are set, the correct source wins ALL the way to events.session, not
// just at the resolveSession() return point.
//
// TestResolveSession in session_test.go already asserts the resolver returns
// the right value in isolation. These sub-tests add the contract that the
// resolver's output reaches the storage layer through a real cobra command.
//
// bd close is the test vehicle — it's the simplest session-aware path
// (single operation, no setup beyond a fresh issue) and the path bd-fwb's
// pre-faa391c03 fix lives on. Choosing close also keeps the precedence proof
// orthogonal to TestSessionAttribution_Close (which exercises each source
// individually) — together they triangulate end-to-end correctness.
//
// Refs: bd-edi (PR1 Phase 7), gastownhall/beads#3583, faa391c03.
func TestSessionAttribution_PrecedenceAndOptIn(t *testing.T) {
	if os.Getenv("BEADS_TEST_EMBEDDED_DOLT") != "1" {
		t.Skip("set BEADS_TEST_EMBEDDED_DOLT=1 to run embedded dolt integration tests")
	}
	t.Parallel()

	bd := buildEmbeddedBD(t)
	dir, beadsDir, _ := bdInit(t, bd, "--prefix", "sp")

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
		return queryEventSessionSQL(t, beadsDir, issueID, "closed")
	}

	t.Run("flag_beats_env", func(t *testing.T) {
		// --session always wins regardless of opt-in or any env var.
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
		// BEADS_SESSION_ID (priority 2) beats CLAUDE_CODE_SESSION_ID
		// (priority 3) when both are set under opt-in.
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
		// CLAUDE_CODE_SESSION_ID (priority 3) beats CLAUDE_SESSION_ID
		// (priority 4) when both are set under opt-in. Mirrors the Claude
		// Code 2.1.132+ rollout: the new env var supersedes the legacy
		// one for Claude Code contexts but the legacy one is retained
		// for upstream tooling.
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
		// Without core.capture-session=true (the default) all three env
		// vars MUST be ignored. This is the "no unattended logging"
		// anchor for environments where CLAUDE_CODE_SESSION_ID is
		// auto-populated globally.
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
		// --session is the explicit-consent path; it must work even when
		// core.capture-session is false. This is the contract that
		// distinguishes the flag from env vars and gives upstream tooling
		// (Gas Town automation, Steve's scripts) a deterministic surface.
		issue := bdCreate(t, bd, dir, "flag without opt-in", "--type", "task")
		got := closeAndQuery(t, issue.ID,
			[]string{"--session", "flag-no-opt-in"},
			// No BD_CORE_CAPTURE_SESSION.
			"BEADS_SESSION_ID=ignored",
			"CLAUDE_CODE_SESSION_ID=ignored",
			"CLAUDE_SESSION_ID=ignored",
		)
		if got != "flag-no-opt-in" {
			t.Errorf("expected events.session=flag-no-opt-in, got %q", got)
		}
	})
}
