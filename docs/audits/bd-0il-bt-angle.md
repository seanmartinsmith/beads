# bt TUI Angle - Feature Surfacing Analysis

*Generated: 2026-04-13 for bd-0il*

## Executive Summary

bt already surfaces issues, labels, dependencies, comments, status, priority, and history - the bread and butter of a project dashboard. The biggest untapped opportunities fall into three buckets:

1. **Gates and blocking visibility** - The entire gate system (human, timer, gh:run, gh:pr, bead) is invisible in bt. This is the highest-value gap because gates are the *actual* blocking mechanism, and bt is the visual layer where blocking should be obvious at a glance.

2. **Cross-project coordination in global mode** - bt's global mode (UNION ALL across shared server databases) is uniquely positioned to surface cross-project relationships that are invisible in single-project bd usage: ship/provides capabilities, bead gates between projects, cross-prefix dependencies.

3. **Workflow intelligence** - bt already has a sophisticated analysis engine (PageRank, betweenness, k-core, etc.). The formula/molecule/swarm stack gives bt structured workflow data to analyze that goes beyond ad-hoc dependency graphs.

bt should NOT try to replicate bd's full command surface. bt's value is *visual intelligence and orchestration* - showing patterns, blockers, and cross-project state that are hard to see from individual `bd` commands.

## Priority Recommendations

### P0: Must Have (fill critical visibility gaps)

**1. Gate status in issue detail and list views** (bt-c69c)
- Show gate type icon/badge on blocked issues in list view
- Detail panel shows gate type, await_id, timeout, and resolution status
- Global mode: cross-project bead gates highlighted as external blockers
- Data source: `bd gate list --json` or direct SQL against issues table (await_type, await_id fields)

**2. Human flag indicator**
- Show `human` label as a distinct visual badge (not just another label)
- Separate from gate-blocked - make the advisory vs. blocking distinction visible
- In list: small icon on human-flagged issues
- In detail: "Awaiting human input" banner with respond/dismiss actions (shell to bd)

**3. Rich markdown rendering in detail panel**
- bt already has `pkg/ui/markdown.go` - ensure description, notes, and design fields render properly
- Checklists, code blocks, headers should render in the detail viewport
- This turns bt from a status dashboard into a working document viewer

### P1: Should Have (high-value additions)

**4. Query/BQL integration with bd query syntax**
- bt already has BQL (`pkg/bql/`). Bridge to bd's query language so users can type compound filters in the BQL modal that map to bd query expressions
- Example: typing `status=open AND priority<=2 AND updated>7d` in bt's BQL modal

**5. Swarm visualization**
- Show swarm analysis (wave fronts, max parallelism, ready fronts) as a view mode
- Leverages bt's existing graph view but colored by wave number
- Data source: `bd swarm validate <epic-id> --json` or `bd swarm status --json`
- Particularly valuable in global mode for seeing swarm state across projects

**6. Ship/provides capability map (global mode)**
- New panel or overlay in global mode showing cross-project capability graph
- Which projects export capabilities, which consume them, which are unresolved
- Data source: labels with `export:`, `provides:`, and `external:` prefixes across databases
- This is the cross-project dependency graph that no single `bd` invocation shows

**7. Stale and overdue issue highlighting**
- Visual indicators on issues that are overdue (due_date < now) or stale (not updated in N days)
- bt already has the data (DueDate, UpdatedAt fields in model.Issue)
- Color or badge in list view: overdue = red clock icon, stale = dimmed/faded

**8. Operational state indicators**
- Surface `dimension:value` labels (patrol:active, health:degraded, etc.) as status badges
- Parse state labels into structured indicators in the detail view
- In global mode: dashboard panel showing state dimensions across projects

### P2: Nice to Have (future enhancements)

**9. Formula/molecule browser**
- Browse available formulas from search paths (.beads/formulas/, ~/.beads/formulas/)
- View molecule instances: active vs. complete, progress percentage
- Data source: `bd formula list --json`, `bd mol show <id> --json`

**10. Merge slot status indicator**
- In global mode: show which projects have active merge slots and who holds them
- Small but critical for multi-agent coordination visibility
- Data source: `bd merge-slot check --json` per project

**11. KV store viewer**
- Read-only display of key-value pairs for a project
- Useful for seeing cross-session state agents have stored
- Data source: `bd kv list --json`

**12. Duplicate detection surface**
- Surface `bd find-duplicates --json` results as a modal or alert
- Could run in background and notify when potential duplicates found
- Integrates with bt's existing alert system

**13. Wisp visibility toggle**
- Show/hide ephemeral issues (wisps) in the list view
- Currently wisps may or may not appear depending on how the data loads
- Toggle in footer or filter to show/hide ephemeral issues

**14. Version control timeline**
- Surface `bd history <id> --json` as a timeline in the detail panel
- Shows how an issue evolved across Dolt commits
- Richer than bt's current git-based history correlation

**15. Integration sync status**
- In global mode: show which projects have external tracker sync (GitHub, Jira, Linear, etc.)
- Badge per project showing sync health/last sync time
- Data source: `bd github status --json` etc. per project

---

## Feature Analysis

### Gates (bd gate)

**TUI Surfacing:** High
**Interaction:** Status indicator in list view (icon per gate type), detail panel section showing gate metadata (type, await_id, timeout, waiters). In global mode: cross-project gate map showing bead gates that span projects.
**Complexity:** Medium - requires either direct SQL access to gate fields (await_type, await_id on issues table) or shelling to `bd gate list --json`. The issue model in bt (`model.Issue`) doesn't currently have gate fields, so model extension is needed.
**Recommendation:** P0. This is the single biggest visibility gap. Users staring at blocked issues in bt have no idea *why* they're blocked or *what* they're waiting for. Gate type determines the action: human gates need a person, timer gates need patience, gh:run gates need CI to pass.

### Human flags (bd human)

**TUI Surfacing:** High
**Interaction:** Distinct badge/icon in list view (separate from gate indicator). Detail panel shows "needs human input" with a respond action. The advisory-vs-blocking distinction must be visually obvious - human label is a yellow flag, human gate is a red block.
**Complexity:** Simple for the label indicator (bt already reads labels). Medium for the respond action (shell to `bd human respond`).
**Recommendation:** P0. Connected to bt-c69c. Critical because agents and users currently confuse advisory human flags with blocking gates. Making the distinction visual eliminates the confusion.

### Rich markdown (description/notes/design fields)

**TUI Surfacing:** High
**Interaction:** Rendered markdown in the detail panel viewport. bt already has `pkg/ui/markdown.go` and uses glamour for rendering. The gap is ensuring all structured fields (design, acceptance_criteria, notes) render with proper markdown treatment, not just description.
**Complexity:** Simple - bt already has the rendering pipeline. May need minor changes to ensure all fields are included in the detail view.
**Recommendation:** P0. Low effort, high impact. Turns bt from "see issue metadata" to "read the full issue with all its structured content."

### Molecules/Formulas (bd mol, bd formula, bd cook)

**TUI Surfacing:** Medium
**Interaction:** Formula browser as a modal or dedicated view. Molecule progress as a progress bar or completion indicator on parent issues. Active molecule list showing running workflows. Could integrate with bt's existing tree view to show molecule structure.
**Complexity:** Complex - new data loading (formulas are TOML/YAML files, molecules are structured issue subgraphs). Would need either a new datasource or shelling to bd for data.
**Recommendation:** P2 for formula browser. P1 for molecule progress indicators (which are just issue subgraphs that bt can already analyze). The molecule progress is achievable because molecules ARE issues with parent-child deps - bt just needs to recognize the pattern.

### Query language (bd query)

**TUI Surfacing:** High
**Interaction:** BQL modal already exists in bt. The opportunity is bridging bt's BQL syntax to bd's query language, or supporting bd query expressions directly in the BQL input. This gives users the full query power without leaving bt.
**Complexity:** Medium - bt has its own BQL (`pkg/bql/`) with its own parser. Two options: (1) extend BQL to support bd query syntax, or (2) add a "bd query passthrough" mode that shells to bd and displays results. Option 2 is simpler.
**Recommendation:** P1. Users who know bd's query syntax should be able to use it in bt. Users who don't should benefit from bt's BQL.

### Swarm (bd swarm)

**TUI Surfacing:** High
**Interaction:** Wave visualization in the graph view: issues colored by their wave number (ready front), with wave 0 highlighted as immediately actionable. Max parallelism indicator in the insights panel. Swarm status as a view mode for epic-focused work.
**Complexity:** Medium - bt already has graph analysis and visualization. Swarm data (waves, ready fronts, parallelism) maps naturally onto bt's existing graph infrastructure. Data from `bd swarm validate --json` provides the wave assignments.
**Recommendation:** P1. Swarm is about parallel coordination - exactly what a visual dashboard should make obvious. Seeing "wave 1: 3 issues ready, wave 2: blocked on wave 1" is more intuitive than reading JSON output.

### Merge slot (bd merge-slot)

**TUI Surfacing:** Medium
**Interaction:** Status indicator in global mode footer or header: "merge slot: available" / "merge slot: held by agent-X (2 waiters)". Clicking/expanding shows waiter queue.
**Complexity:** Simple - single status check per project. Data from `bd merge-slot check --json`.
**Recommendation:** P2. Only relevant in multi-agent scenarios. When it matters, it matters a lot (prevents merge conflicts), but most single-user sessions don't need it. Could be an opt-in indicator.

### Ship/cross-project deps (bd ship)

**TUI Surfacing:** High (in global mode)
**Interaction:** Cross-project capability graph as a new view or overlay in global mode. Nodes are projects, edges are capability dependencies. Unresolved capabilities highlighted in red. Resolved capabilities shown as green connections.
**Complexity:** Complex - requires scanning labels across all databases for export:/provides:/external: patterns. No single bd command produces this cross-project view; bt would need to synthesize it from per-database label data.
**Recommendation:** P1 for global mode. This is a view that ONLY bt can provide - no single `bd` invocation shows the full cross-project capability graph. This is bt's unique value proposition as the orchestration layer.

### State management (bd set-state, bd state)

**TUI Surfacing:** Medium
**Interaction:** State dimensions parsed from labels and displayed as structured badges in the detail panel. In global mode: state dashboard showing all projects' operational state (health, patrol status, mode).
**Complexity:** Simple - labels are already loaded. Parsing `dimension:value` patterns is string splitting. Display is badge rendering bt already does for labels.
**Recommendation:** P1. Low implementation effort for meaningful operational visibility. Especially valuable in global mode where seeing "project X: health=degraded" across all projects is immediate situational awareness.

### KV store (bd kv)

**TUI Surfacing:** Low
**Interaction:** Read-only key-value list in the detail panel or a dedicated modal. Useful for debugging what agents have persisted.
**Complexity:** Simple - `bd kv list --json` per project.
**Recommendation:** P2. Niche use case but trivial to implement. Could be a "project settings" or "project state" section in a future project detail view.

### Graph visualization (bd graph)

**TUI Surfacing:** Medium (bt already has this)
**Interaction:** bt already has ViewGraph with full graph visualization. The opportunity is importing bd's `--html` and `--dot` export formats, and adding bd's layer-based execution order view as an alternative layout in bt's graph.
**Complexity:** Low - bt already has the infrastructure. Just importing additional layout algorithms.
**Recommendation:** P2. bt's existing graph view is already richer than bd's terminal output. The HTML export from bd (`bd graph --html`) is a separate artifact - bt's TUI graph serves the same purpose natively.

### Find duplicates (bd find-duplicates)

**TUI Surfacing:** Medium
**Interaction:** Background analysis that surfaces duplicate candidates as alerts. User can review pairs and mark as duplicate, dismiss, or merge. Integrates with bt's existing alert system.
**Complexity:** Medium - requires running `bd find-duplicates --json` (potentially slow for AI mode) and displaying results as an actionable list.
**Recommendation:** P2. Useful for project hygiene but not daily workflow. Could be a periodic background check that surfaces alerts when duplicates are found.

### Comments (bd comments)

**TUI Surfacing:** Already partially surfaced
**Interaction:** bt's model.Issue already has a Comments field. Ensure comments render in the detail panel with author/timestamp. Add ability to add comments via shell-out to `bd comment <id> "text"`.
**Complexity:** Simple for read, Medium for write (needs text input modal + shell execution).
**Recommendation:** P1 for ensuring read display is complete. P2 for write capability (part of the broader "writes via shell-out to bd" work).

### Epic management (bd epic)

**TUI Surfacing:** Medium
**Interaction:** Epic completion indicator on epic-type issues in list view (e.g., "3/7 children complete"). "Close eligible" badge when all children are done. Could be a button/action in the detail panel.
**Complexity:** Simple - bt already analyzes parent-child relationships. Computing "N of M children closed" is straightforward.
**Recommendation:** P1. Low effort, high visibility. Seeing epic progress at a glance is exactly what a dashboard should provide.

### Promote (bd promote)

**TUI Surfacing:** Low
**Interaction:** Action on wisp issues in the detail panel: "Promote to permanent bead." Shells to `bd promote <id>`.
**Complexity:** Simple - single action, single bd command.
**Recommendation:** P2. Only relevant when wisps are visible and the user decides one should be permanent. Depends on wisp visibility toggle (also P2).

### Integration sync (bd github, bd jira, bd linear, etc.)

**TUI Surfacing:** Medium (in global mode)
**Interaction:** Per-project badge showing external tracker integration status. Sync health indicator (last sync time, error count). In global mode: overview of which projects sync where.
**Complexity:** Medium - requires per-project `bd <tracker> status --json` queries. Multiple integration types to support.
**Recommendation:** P2. Valuable for mixed environments but not core to the TUI experience. Better as a "project health" view in global mode.

### Worktree (bd worktree)

**TUI Surfacing:** Low
**Interaction:** Informational indicator showing if current project has active worktrees. Not actionable from TUI.
**Complexity:** Simple - read-only metadata display.
**Recommendation:** P2. Worktree management is inherently a git/terminal operation, not a TUI operation. bt might note "3 active worktrees" in project info but shouldn't manage them.

### Dolt version control (bd vc, bd branch, bd diff)

**TUI Surfacing:** Low
**Interaction:** Time travel is already in bt (ModalTimeTravelInput). Branch listing and diff viewing would be additional panels.
**Complexity:** Medium - bt already connects to Dolt. Branch listing is a simple query. Diff viewing requires rendering table diffs.
**Recommendation:** P2. bt's time travel covers the most common use case. Full branch management is a power-user need better served by bd commands directly.

### Batch operations (bd batch)

**TUI Surfacing:** N/A (backend optimization)
**Interaction:** Not visible to users. bt should use batch internally when performing multiple writes to reduce Dolt overhead.
**Complexity:** Simple - bt shells to bd for writes. Could batch multiple updates into a single `bd batch` invocation.
**Recommendation:** Internal optimization, not a surfacing concern. File as a separate bt performance bead.

### Defer/undefer (bd defer, bd undefer)

**TUI Surfacing:** Medium
**Interaction:** Action on issues in detail panel: "Defer" with optional until-date input, "Undefer" on deferred issues. Deferred issues could have a distinct visual treatment (dimmed, snowflake icon matching bd's conventions).
**Complexity:** Simple - bt already shows deferred status. Adding actions requires shell-out to `bd defer`/`bd undefer`.
**Recommendation:** P1 for visual treatment of deferred issues. P2 for defer/undefer actions (part of broader write capability).

### Audit (bd audit)

**TUI Surfacing:** Low
**Interaction:** Read-only view of interaction log entries. Niche use case.
**Complexity:** Simple - parse JSONL file.
**Recommendation:** P2. Not relevant for most users. Could be a debugging/forensics view for advanced users.

### Backup/restore (bd backup, bd restore)

**TUI Surfacing:** N/A
**Interaction:** Administrative operations not appropriate for TUI surfacing.
**Complexity:** N/A
**Recommendation:** Leave as bd CLI commands. bt should not manage backups.

---

## Global Mode Opportunities

bt's global mode (cross-project UNION ALL) creates unique opportunities that no single `bd` invocation can match:

### 1. Cross-project gate map
Show all bead-type gates across projects: "Project A is waiting on Project B's auth-module capability." This is the only view that shows cross-project blocking relationships in one place.

### 2. Capability graph (ship/provides/external)
Scan `export:`, `provides:`, and `external:<project>:<capability>` labels across all databases. Render as a directed graph where nodes are projects and edges are capability dependencies. Unresolved edges highlighted.

### 3. Cross-project state dashboard
Aggregate `dimension:value` state labels across all projects. Show a matrix: rows = projects, columns = state dimensions (health, patrol, mode), cells = current values with color coding.

### 4. Global swarm coordination
When multiple projects have active swarms, show aggregate parallelism and ready fronts across all projects. Useful for understanding overall system load and coordination state.

### 5. Project health heatmap
Combine stale issue counts, overdue issues, gate count, and velocity per project into a health score. Render as a heatmap in the project picker or a dedicated view.

### 6. Cross-project duplicate detection
In global mode, scan for issues with similar titles across projects. Surface potential cross-project duplicates or redundant work.

---

## New Panel/View Ideas

### Gate Dashboard (P0)
A dedicated view mode (like ViewGraph or ViewBoard) showing all gates in the current project (or across projects in global mode). Columns by gate type: human, timer, gh:run, gh:pr, bead. Each gate shows what it's waiting for, how long it's been waiting, and which issues are blocked behind it.

### Molecule Progress View (P1)
An extension of the tree view that shows active molecules as progress-tracked workflows. Each molecule shows: total steps, completed steps, current wave, blocked steps with gate details. Think of it as a pipeline/workflow visualization.

### Operational State Panel (P1)
A compact panel (could be a sidebar or footer section) showing the current project's operational state dimensions. In global mode, expands to show all projects' state. Think: miniature control room dashboard.

### Cross-project Dependency Graph (P1, global mode only)
A graph view where nodes are projects (not individual issues) and edges are cross-project dependencies (external:, bead gates, ship/provides). Drill down into a project node to see its contributing edges.

### Ready Work Dashboard (P1)
Reimagine the list view as a "ready to work" dashboard that mirrors `bd ready` semantics. Instead of showing all issues and filtering, start from ready issues and explain *why* each is ready (no blockers, no gates, not deferred). For non-ready issues, show *what's blocking* (which gate, which dependency, which deferred date).

### Timeline View (P2)
A horizontal timeline showing issue lifecycle events: created, claimed, gated, unblocked, completed. Useful for retrospectives and understanding flow. Data from `bd history --json` per issue.

---

## Implementation Notes

### Data Access Strategy

bt currently loads issues via:
1. JSONL files (legacy)
2. Direct Dolt SQL queries (single project and global mode)

For gate and molecule data, bt needs access to fields not in the current model:
- **Gate fields**: `await_type`, `await_id`, `timeout`, `waiters` on issues
- **Molecule fields**: `ephemeral`, `template`, `mol_type` on issues
- **State labels**: Already available via labels (just need parsing logic)
- **Ship/provides**: Already available via labels (just need cross-project aggregation)

The most pragmatic approach: extend the SQL queries in `internal/datasource/` to include gate and molecule columns, extend `model.Issue` to carry them, and parse state/capability labels in bt's analysis layer.

For features that require bd CLI execution (writes, formula listing, swarm validation), shell out to `bd` commands - this is already the planned approach for bt writes.

### Connected beads

- **bt-c69c**: TUI human/gate surfacing - directly addressed by P0 recommendations here
- **bt-8f34**: Project registry - the global mode opportunities here (capability graph, state dashboard) depend on or complement the project registry work
- **bd-chc**: Ad-hoc gate creation - bt needs this to enable gate actions from the TUI without requiring formulas
