package main

import (
	"os"
	"testing"

	"github.com/steveyegge/beads/internal/config"
)

// TestResolveSession covers the two-layer opt-in semantics:
//   - --session flag is always honored, regardless of config.
//   - BEADS_SESSION_ID > CLAUDE_SESSION_ID, but only when core.capture-session=true.
//   - Default (no flag, no opt-in) returns "" even with env vars set.
func TestResolveSession(t *testing.T) {
	tests := []struct {
		name           string
		flag           string
		beadsSession   string
		claudeSession  string
		captureSession bool
		want           string
	}{
		{
			name: "default returns empty",
			want: "",
		},
		{
			name:           "opt-in off ignores env vars",
			beadsSession:   "from-beads-env",
			claudeSession:  "from-claude-env",
			captureSession: false,
			want:           "",
		},
		{
			name:           "BEADS_SESSION_ID wins over CLAUDE_SESSION_ID when opted in",
			beadsSession:   "from-beads-env",
			claudeSession:  "from-claude-env",
			captureSession: true,
			want:           "from-beads-env",
		},
		{
			name:           "CLAUDE_SESSION_ID is fallback when BEADS_SESSION_ID empty",
			claudeSession:  "from-claude-env",
			captureSession: true,
			want:           "from-claude-env",
		},
		{
			name:           "flag wins over env vars even when opted in",
			flag:           "from-flag",
			beadsSession:   "from-beads-env",
			claudeSession:  "from-claude-env",
			captureSession: true,
			want:           "from-flag",
		},
		{
			name:           "flag wins when opt-in is off",
			flag:           "from-flag",
			beadsSession:   "from-beads-env",
			claudeSession:  "from-claude-env",
			captureSession: false,
			want:           "from-flag",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withCleanSessionEnv(t)
			t.Setenv("BEADS_TEST_IGNORE_REPO_CONFIG", "1")
			config.ResetForTesting()
			t.Cleanup(config.ResetForTesting)
			if err := config.Initialize(); err != nil {
				t.Fatalf("config.Initialize: %v", err)
			}
			config.Set("core.capture-session", tt.captureSession)

			origSession := session
			session = tt.flag
			t.Cleanup(func() { session = origSession })

			if tt.beadsSession != "" {
				t.Setenv("BEADS_SESSION_ID", tt.beadsSession)
			}
			if tt.claudeSession != "" {
				t.Setenv("CLAUDE_SESSION_ID", tt.claudeSession)
			}

			got := resolveSession()
			if got != tt.want {
				t.Errorf("resolveSession() = %q, want %q", got, tt.want)
			}
		})
	}
}

// withCleanSessionEnv unsets session env vars for the duration of the test
// so that ambient values from the developer's shell do not bleed into cases
// that exercise the unset path.
func withCleanSessionEnv(t *testing.T) {
	t.Helper()
	for _, key := range []string{"BEADS_SESSION_ID", "CLAUDE_SESSION_ID"} {
		orig, ok := os.LookupEnv(key)
		os.Unsetenv(key)
		if ok {
			t.Cleanup(func() { os.Setenv(key, orig) })
		} else {
			t.Cleanup(func() { os.Unsetenv(key) })
		}
	}
}
