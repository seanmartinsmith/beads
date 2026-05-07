//go:build cgo

package embeddeddolt_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/steveyegge/beads/internal/storage/embeddeddolt"
	"github.com/steveyegge/beads/internal/types"
)

// skipUnlessEmbeddedDoltBench is the *testing.B parallel of
// skipUnlessEmbeddedDolt, since the package-level helper is *testing.T-typed.
func skipUnlessEmbeddedDoltBench(b *testing.B) {
	b.Helper()
	if os.Getenv("BEADS_TEST_EMBEDDED_DOLT") != "1" {
		b.Skip("set BEADS_TEST_EMBEDDED_DOLT=1 to run embedded dolt benchmarks")
	}
}

// BenchmarkGetIssue_ClosedBySessionDerivation measures the per-call cost of
// GetIssue for a closed issue (which exercises the events-table lookup added
// in Phase 6 of bd-edi PR1) versus an open issue (which skips the lookup
// entirely via the `if Status == closed` gate in GetIssueInTx).
//
// The delta between the Closed and Open sub-benchmarks is the wall-clock
// overhead introduced by approach (c) (in-code enrichment via a single
// indexed lookup against events / wisp_events). Brief §6 commit 4 gate:
// >2x regression versus baseline triggers a halt for architecture review.
//
// To run:
//
//	BEADS_TEST_EMBEDDED_DOLT=1 go test -tags gms_pure_go \
//	  -run XXX -bench BenchmarkGetIssue_ClosedBySessionDerivation \
//	  ./internal/storage/embeddeddolt/ -benchtime=2s
//
// The benchmark seeds 500 issues (250 closed with session attribution, 250
// open) and selects a representative ID for each sub-benchmark. 500 issues
// is enough to make per-call costs stable; smaller sets get drowned by
// fixed per-test overhead, larger sets risk Windows test-server timeouts.
//
// For the brief's 10K issues + 100K events scale gate, see the follow-up
// dolt-server-backed benchmark (deferred — runs only on Linux/CI where the
// shared Dolt server is available).
func benchmarkSeedClosedAndOpen(b *testing.B) (te *testEnv, closedID, openID string) {
	b.Helper()
	ctx := b.Context()
	beadsDir := filepath.Join(b.TempDir(), ".beads")
	store, err := embeddeddolt.Open(ctx, beadsDir, "bench", "main")
	if err != nil {
		b.Fatalf("Open: %v", err)
	}
	b.Cleanup(func() { store.Close() })
	if err := store.SetConfig(ctx, "issue_prefix", "bench"); err != nil {
		b.Fatalf("SetConfig: %v", err)
	}
	if err := store.Commit(ctx, "bench init"); err != nil {
		b.Fatalf("Commit: %v", err)
	}

	te = &testEnv{
		store:    store,
		dataDir:  filepath.Join(beadsDir, "embeddeddolt"),
		database: "bench",
	}

	// Seed 500 issues: half open, half closed-with-session.
	const total = 500
	for i := 0; i < total; i++ {
		issue := &types.Issue{
			ID:        fmt.Sprintf("bench-%04d", i),
			Title:     fmt.Sprintf("Bench issue %d", i),
			Status:    types.StatusOpen,
			Priority:  2,
			IssueType: types.TypeTask,
		}
		if err := store.CreateIssue(ctx, issue, "bench"); err != nil {
			b.Fatalf("CreateIssue %s: %v", issue.ID, err)
		}
		// Close half of them with session attribution. Choosing every other
		// issue gives a deterministic split and lets us pick a closed ID and
		// open ID in the middle of the seed (warm-cache representative).
		if i%2 == 0 {
			if err := store.CloseIssue(ctx, issue.ID, "bench close", "bench", fmt.Sprintf("session-%d", i)); err != nil {
				b.Fatalf("CloseIssue %s: %v", issue.ID, err)
			}
		}
	}

	// Pick representative IDs from the middle of the seed set.
	closedID = "bench-0250" // even → closed
	openID = "bench-0251"   // odd → open
	return te, closedID, openID
}

func BenchmarkGetIssue_ClosedBySessionDerivation(b *testing.B) {
	skipUnlessEmbeddedDoltBench(b)

	te, closedID, openID := benchmarkSeedClosedAndOpen(b)
	ctx := b.Context()

	b.Run("closed_issue_runs_lookup", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			got, err := te.store.GetIssue(ctx, closedID)
			if err != nil {
				b.Fatalf("GetIssue(closed): %v", err)
			}
			if got.ClosedBySession == "" {
				b.Fatalf("expected derived ClosedBySession on closed issue %q, got empty", closedID)
			}
		}
	})

	b.Run("open_issue_skips_lookup", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			got, err := te.store.GetIssue(ctx, openID)
			if err != nil {
				b.Fatalf("GetIssue(open): %v", err)
			}
			if got.Status != types.StatusOpen {
				b.Fatalf("expected open status on %q, got %q", openID, got.Status)
			}
		}
	})
}
