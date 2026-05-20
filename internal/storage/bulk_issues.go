package storage

import (
	"context"

	"github.com/steveyegge/beads/internal/types"
)

// BulkIssueStore provides extended issue CRUD beyond the base Storage interface.
type BulkIssueStore interface {
	CreateIssuesWithFullOptions(ctx context.Context, issues []*types.Issue, actor string, opts BatchCreateOptions) error
	DeleteIssues(ctx context.Context, ids []string, cascade bool, force bool, dryRun bool) (*types.DeleteIssuesResult, error)
	DeleteIssuesBySourceRepo(ctx context.Context, sourceRepo string) (int, error)
	UpdateIssueID(ctx context.Context, oldID, newID string, issue *types.Issue, actor, session string) error
	// ClaimIssue atomically claims an issue and records a 'claimed' event.
	// session may be empty when session attribution is not configured; when
	// non-empty it is recorded on the claim event (events.session). This is
	// the architectural-win path that closes the bd ready --claim gap.
	ClaimIssue(ctx context.Context, id string, actor, session string) error
	ClaimReadyIssue(ctx context.Context, filter types.WorkFilter, actor, session string) (*types.Issue, error)
	PromoteFromEphemeral(ctx context.Context, id string, actor string) error
	GetNextChildID(ctx context.Context, parentID string) (string, error)
}
