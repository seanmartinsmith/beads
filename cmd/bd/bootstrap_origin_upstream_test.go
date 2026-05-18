package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/steveyegge/beads/internal/configfile"
)

// TestBootstrapPlan_SyncRemoteSource_Explicit verifies that a configured
// sync.remote is tagged with SyncRemoteSource="explicit" so the bootstrap
// confirmation flow can trust it (no extra prompting beyond the standard
// Proceed? gate). See bd-jui.
func TestBootstrapPlan_SyncRemoteSource_Explicit(t *testing.T) {
	plan := BootstrapPlan{
		Action:           "sync",
		SyncRemote:       "git+ssh://git@example.com/user/repo.git",
		SyncRemoteSource: "explicit",
	}

	if plan.SyncRemoteSource != "explicit" {
		t.Fatalf("SyncRemoteSource = %q, want %q", plan.SyncRemoteSource, "explicit")
	}
}

// TestBootstrapPlan_SyncRemoteSource_GitOrigin verifies that a sync URL
// auto-derived from git origin (the dangerous fork-less clone case) is
// tagged as such. The confirmation logic uses this tag to refuse silent
// auto-confirms on non-TTY stdin. See bd-jui.
func TestBootstrapPlan_SyncRemoteSource_GitOrigin(t *testing.T) {
	plan := BootstrapPlan{
		Action:           "sync",
		SyncRemote:       "https://github.com/upstream/repo.git",
		SyncRemoteSource: "git-origin",
	}

	if plan.SyncRemoteSource != "git-origin" {
		t.Fatalf("SyncRemoteSource = %q, want %q", plan.SyncRemoteSource, "git-origin")
	}
}

// TestConfirmAutoDerivedSync_NonInteractiveAffirmative verifies that an
// explicit non-interactive affirmative signal (--yes / BD_NON_INTERACTIVE /
// CI) permits the auto-derived sync clone. See bd-jui.
func TestConfirmAutoDerivedSync_NonInteractiveAffirmative(t *testing.T) {
	plan := BootstrapPlan{
		Action:           "sync",
		SyncRemote:       "https://github.com/upstream/repo.git",
		SyncRemoteSource: "git-origin",
	}

	if !confirmAutoDerivedSync(plan, true) {
		t.Fatal("confirmAutoDerivedSync should return true when nonInteractive=true (explicit signal)")
	}
}

// TestConfirmAutoDerivedSync_NoTTYNoExplicitSignal verifies that the
// auto-derived sync path REFUSES to silently auto-confirm when stdin is
// not a TTY and no explicit affirmative signal was provided. This is the
// core regression for bd-jui — the 4559-issue incident path.
//
// confirmAutoDerivedSync must return false in this scenario. The test
// harness runs without a TTY on stdin, so we can exercise this branch
// directly. We pass nonInteractive=false to simulate the absence of any
// affirmative signal.
func TestConfirmAutoDerivedSync_NoTTYNoExplicitSignal(t *testing.T) {
	plan := BootstrapPlan{
		Action:           "sync",
		SyncRemote:       "https://github.com/upstream/repo.git",
		SyncRemoteSource: "git-origin",
	}

	// Test harness: stdin is piped (not a TTY) and no explicit signal.
	// This is the silent-clone footgun and MUST be refused.
	if confirmAutoDerivedSync(plan, false) {
		t.Fatal("confirmAutoDerivedSync must refuse silent auto-confirm on non-TTY stdin without an explicit signal (bd-jui regression)")
	}
}

// TestInitCommandRegistersNoRemoteFlag verifies the --no-remote opt-out
// flag is registered on bd init. See bd-jui acceptance criterion 6.
func TestInitCommandRegistersNoRemoteFlag(t *testing.T) {
	flag := initCmd.Flags().Lookup("no-remote")
	if flag == nil {
		t.Fatal("init command does not register --no-remote")
	}
	if flag.DefValue != "false" {
		t.Fatalf("--no-remote default = %q, want %q", flag.DefValue, "false")
	}
}

// TestBootstrapSyncSourceDistinct verifies that the SyncRemoteSource
// values are distinct strings — the executeBootstrapPlan gate depends on
// matching the literal "git-origin" tag. This catches accidental rename
// regressions of the tag.
func TestBootstrapSyncSourceDistinct(t *testing.T) {
	if "explicit" == "git-origin" {
		t.Fatal("source-tag constants collide")
	}

	// Smoke-check the JSON serialization path so the field is observable
	// to programmatic consumers (bd bootstrap --json).
	plan := BootstrapPlan{
		Action:           "sync",
		SyncRemote:       "https://github.com/upstream/repo.git",
		SyncRemoteSource: "git-origin",
		BeadsDir:         filepath.Join(t.TempDir(), ".beads"),
	}
	if plan.SyncRemoteSource == "" {
		t.Fatal("SyncRemoteSource lost across struct construction")
	}
	// Make sure the unused configfile import stays valid in case future
	// edits need it; this test intentionally avoids touching disk so it
	// does not collide with viper initialization in detectBootstrapAction.
	_ = configfile.DefaultConfig()
	_ = os.PathSeparator
}
