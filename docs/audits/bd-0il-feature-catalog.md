# bd Feature Catalog - Agent Knowledge Audit

*Generated: 2026-04-13 for bd-0il*

## Summary

- Total subcommands: 103 (top-level, excluding `help`)
- KNOWN: 22 (agents are taught this) - 21%
- PARTIAL: 14 (mentioned but incomplete) - 14%
- UNKNOWN: 67 (agents have no awareness) - 65%

## Critical Confusion Points

### 1. `bd human` is NOT a blocking mechanism

PRIME.md teaches `bd human <id>` as "flag issue for human decision." Beads-onboard reinforces this. But `bd human` is advisory only - it adds a `human` label to the issue. It does NOT block the issue from appearing in `bd ready`. The actual blocking mechanism is `bd gate` with a `human` type gate, which can only be created through the formula/molecule system. Agents who think `bd human` blocks workflow are wrong, and nothing in their training corrects this.

The `bd human` command also has a dual personality that is never explained to agents: when called with no arguments, it displays a help menu for human users (the ~15 essential commands). When called with subcommands like `list`, `respond`, `dismiss`, it manages the human-flagged issues. This is confusing and undocumented in PRIME.md.

### 2. Rich markdown fields are untapped

Beads support full markdown in `description`, `notes`, and `design` fields - including checklists, session setup instructions, code blocks, and suggested task breakdowns. Nothing in PRIME.md or beads-onboard teaches agents to use markdown effectively in these fields. Agents write flat text descriptions when they could be writing structured, searchable, recoverable markdown.

### 3. The gate/formula/molecule stack is invisible

PRIME.md mentions `bd formula list` and `bd mol pour` in exactly two lines under "Structured Workflows" with zero explanation of what they do, when to use them, or how the three concepts relate. The entire async coordination system (gates), workflow templating system (formulas), and work execution system (molecules) is effectively invisible to agents.

### 4. Query language is unknown

`bd query` provides compound boolean filters with date-relative expressions - far more powerful than `bd list` flags. Agents don't know it exists and overuse `bd list` with limited filtering or pipe through jq.

### 5. Cross-project dependency (`bd ship`) is unknown

The entire `ship`/`export:`/`provides:` capability publishing system for cross-project coordination is untaught.

## Feature Catalog

---

### Category: Working With Issues

#### bd assign
**What it does:** Shorthand for `bd update <id> --assignee <name>`. Assigns an issue to a named user. Convenience wrapper to avoid typing the full update command.

**Agent knowledge:** PARTIAL
**Gap:** PRIME.md teaches `bd update <id> --assignee=username` but not the `bd assign` shorthand. Agents use the longer form.

---

#### bd children
**What it does:** Lists all child beads of a specified parent. Convenience alias for `bd list --parent <id> --status all`. Unlike plain `bd list`, this includes closed issues by default since the primary use case is inspecting all work under an epic or parent.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned anywhere. Agents who need to see an epic's children would use `bd list --parent <id>` but might miss closed children since `bd list` excludes them by default.

---

#### bd close
**What it does:** Closes one or more issues. Supports batch closing (`bd close <id1> <id2> ...`), closing with a reason (`--reason`), and suggesting next work (`--suggest-next`). If no ID is provided, closes the last touched issue.

**Agent knowledge:** KNOWN
**Gap:** Fully documented in PRIME.md including batch close and `--suggest-next`.

---

#### bd comment
**What it does:** Shorthand for `bd comments add <id> "text"`. Adds a comment to an issue without the subcommand syntax.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned in PRIME.md. Agents who want to comment must discover this or use `bd comments add`.

---

#### bd comments
**What it does:** View or manage comments on an issue. Subcommands: add, list. Comments provide threaded discussion on issues separate from the description/notes fields.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned in PRIME.md or beads-onboard. The beads plugin has a `beads:comments` skill but it's passive - agents don't know to use comments for discussion, context capture, or decision logging on issues.

---

#### bd create
**What it does:** Creates a new issue with title, description, type, priority, labels, and optional structured fields (design, acceptance criteria, notes). Supports batch creation from markdown or graph JSON. Supports `--validate` for quality checking, `--due` and `--defer` for time-based scheduling.

**Agent knowledge:** KNOWN
**Gap:** Core command well-documented in PRIME.md. Quality flags (`--validate`, `--acceptance`, `--design`, `--notes`) are also covered. Due/defer date flags from CLAUDE.md's quick reference are documented but the markdown-in-fields capability is untaught.

---

#### bd create-form
**What it does:** Interactive terminal form for creating issues with a TUI interface. Fields for title, description, type, priority, labels, and more. Designed for human users rather than agents.

**Agent knowledge:** UNKNOWN
**Gap:** Interactive command - not useful for agents (would block). Not a meaningful gap since agents should use `bd create` with flags.

---

#### bd delete
**What it does:** Permanently deletes one or more issues and cleans up all references. Removes dependency links in both directions, updates text references to `[deleted:ID]` in connected issues, and removes from the database. Destructive and irreversible.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned in PRIME.md. Agents have no way to know they can delete issues. This is arguably fine for safety - accidental deletion would be bad. But there's no guidance on when delete vs close is appropriate.

---

#### bd edit
**What it does:** Opens an issue field in `$EDITOR` for editing. By default edits the description. Supports flags for other fields.

**Agent knowledge:** KNOWN (as a warning)
**Gap:** PRIME.md explicitly warns: "WARNING: Do NOT use `bd edit` - it opens $EDITOR (vim/nano) which blocks agents." This is correct and sufficient.

---

#### bd gate
**What it does:** Manages async coordination gates - wait conditions that block workflow steps. Gate types: human (manual close), timer (expires after timeout), gh:run (GitHub workflow), gh:pr (PR merge), bead (cross-rig bead close). Subcommands: list, check, resolve, show, add-waiter, discover. Gates are created automatically when a formula step has a gate field and must be closed for blocked steps to proceed.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned in PRIME.md at all. Not covered in beads-onboard. This is the actual async coordination mechanism - the real blocking primitive - and agents have zero awareness. The `bd human` advisory flag is taught; the `bd gate` blocking mechanism is not.

---

#### bd label
**What it does:** Manages issue labels with subcommands: add, remove, list, list-all, propagate (parent to children). Labels are the primary taxonomy mechanism in beads.

**Agent knowledge:** PARTIAL
**Gap:** PRIME.md doesn't mention `bd label` as a command. It teaches `bd update <id> --add-label` for adding labels. The `label propagate` command (push parent labels to children) and `label list-all` (see all labels in use across the database) are unknown.

---

#### bd link
**What it does:** Shorthand for `bd dep add <id1> <id2>`. Links two issues with a dependency. By default creates a "blocks" dependency. Supports `--type` for different relationship types.

**Agent knowledge:** PARTIAL
**Gap:** PRIME.md teaches `bd dep add` but not the `bd link` shorthand. The shorthand is cleaner for simple linking.

---

#### bd list
**What it does:** Lists issues with extensive filtering options. Supports status, type, priority, assignee, label, parent, date-based filters (--due-before, --due-after, --defer-before, --defer-after, --overdue, --deferred), and sorting. The workhorse query command.

**Agent knowledge:** KNOWN
**Gap:** Core filtering is taught. CLAUDE.md's quick reference covers time-based query flags. However, the full flag inventory is large and agents may not know about all filtering capabilities.

---

#### bd merge-slot
**What it does:** Exclusive access primitive for serialized conflict resolution. Only one agent can hold the merge slot at a time, preventing race conditions during merge operations. Subcommands: create, check, acquire, release. Uses status and metadata fields to track holder and waiters queue.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned anywhere. Critical for multi-agent workflows where parallel agents need to serialize merge operations. Agents in swarm or parallel work scenarios have no awareness this exists.

---

#### bd note
**What it does:** Shorthand for `bd update <id> --append-notes "text"`. Appends text to an issue's notes field.

**Agent knowledge:** PARTIAL
**Gap:** PRIME.md teaches `--notes` on create and `--notes` on update, but not the `bd note` convenience command for appending.

---

#### bd priority
**What it does:** Shorthand for `bd update <id> --priority <n>`. Sets the priority of an issue.

**Agent knowledge:** PARTIAL
**Gap:** PRIME.md uses `bd update` syntax for priority changes. The shorthand is untaught but the capability is known.

---

#### bd promote
**What it does:** Promotes a wisp (ephemeral issue) to a permanent bead. Copies the issue from the wisps table (dolt_ignored) to the permanent issues table (Dolt-versioned), preserving labels, dependencies, events, and comments. The original ID is preserved so all links keep working.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned in PRIME.md or beads-onboard. Agents who create wisps have no way to know they can promote them to permanent beads if the work turns out to be more significant than initially expected. Part of the invisible molecule system.

---

#### bd q
**What it does:** Quick capture - creates an issue and outputs only the ID. Designed for scripting: `ISSUE=$(bd q "New feature")`. Minimal output for pipeline use.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned anywhere. Useful for agents in scripting contexts or when they need to capture the ID programmatically without parsing verbose output. The `--json` flag on `bd create` achieves similar results but `bd q` is cleaner.

---

#### bd query
**What it does:** Full query language with compound filters, boolean operators (AND, OR, NOT), comparison operators, date-relative expressions, and grouping with parentheses. Supports fields: status, priority, type, assignee, owner, label, title, description, notes, created, updated, closed, id, spec, pinned, ephemeral, template, parent, mol_type. Far more powerful than `bd list` flags.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned in PRIME.md or beads-onboard. Agents use `bd list` with flags or pipe through jq. Examples like `bd query "status=open AND priority<=2 AND updated>7d"` are dramatically more efficient than the alternatives agents currently use.

---

#### bd reopen
**What it does:** Reopens closed issues by setting status to open and clearing closed_at timestamp. More explicit than `bd update --status open` and emits a Reopened event for audit trail.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned in PRIME.md. Agents who need to reopen an issue would try `bd update <id> --status open` which works but doesn't emit the proper event.

---

#### bd search
**What it does:** Searches issues across title and ID. ID-like queries use fast exact/prefix matching. Text queries search titles. Supports `--desc-contains` for description search. Excludes closed issues by default; use `--status all` to include them.

**Agent knowledge:** KNOWN
**Gap:** Mentioned in PRIME.md as `bd search <query>`. The `--desc-contains` flag for description search is not taught.

---

#### bd set-state
**What it does:** Atomically sets operational state on an issue. Creates an event bead recording the state change (source of truth), removes any existing label for the dimension, and adds the new dimension:value label. State labels follow `<dimension>:<value>` convention (e.g., patrol:active, health:healthy). The `--reason` flag provides context for the event.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned anywhere. Part of the operational state management system for long-running agents and monitors. Agents managing patrol or health-check workflows have no awareness of this atomic state transition mechanism.

---

#### bd show
**What it does:** Shows detailed information about one or more issues including dependencies, labels, comments, and all structured fields. Supports `--current` for the most recently touched issue.

**Agent knowledge:** KNOWN
**Gap:** Well-documented in PRIME.md.

---

#### bd state
**What it does:** Queries the current value of a state dimension from an issue's labels. Extracts the value for a given dimension (e.g., `bd state witness-abc patrol` returns "active"). Subcommand `list` shows all state dimensions on an issue.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned anywhere. Companion to `bd set-state` - both are invisible.

---

#### bd tag
**What it does:** Shorthand for `bd update <id> --add-label <label>`. Adds a label to an issue.

**Agent knowledge:** PARTIAL
**Gap:** PRIME.md teaches `--add-label` on update. The `bd tag` shorthand is untaught but the capability is known.

---

#### bd todo
**What it does:** Lightweight task convenience wrapper. `bd todo add "Title"` maps to `bd create "Title" -t task -p 2`. `bd todo` lists open tasks. `bd todo done <id>` closes a task. Designed for quick capture without specifying type/priority.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned in PRIME.md. Agents always use the full `bd create` command with explicit type and priority flags. The `bd todo` shorthand would be useful for lightweight task tracking.

---

#### bd update
**What it does:** Updates one or more issues. Supports title, description, notes, design, status, priority, assignee, labels, and more. Supports `--claim` as shorthand for claiming work. If no ID provided, updates the last touched issue.

**Agent knowledge:** KNOWN
**Gap:** Well-documented in PRIME.md including `--claim` and various field flags.

---

### Category: Views & Reports

#### bd count
**What it does:** Counts issues matching specified filters. Supports `--by-*` flags to group counts by different attributes (type, status, priority, etc.). Useful for dashboards and reporting without fetching full issue data.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned anywhere. Agents who need counts would run `bd list` and count results, or pipe through `wc -l`.

---

#### bd diff
**What it does:** Shows changes between two Dolt commits or branches in the issue database. Refs can be commit hashes or branch names. Useful for understanding what changed in the issue database over time.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned in PRIME.md. Part of the version control capabilities that agents don't know about.

---

#### bd find-duplicates
**What it does:** Semantic duplicate detection using either mechanical (token-based Jaccard similarity) or AI (LLM-based) methods. The mechanical approach is fast and free; the AI approach uses Claude for semantic comparison with mechanical pre-filtering. Configurable threshold and status filters.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned anywhere. Valuable for project hygiene - agents could proactively check for duplicates before creating new issues or during preflight checks.

---

#### bd history
**What it does:** Shows the complete version history of an issue including all Dolt commits where the issue was modified. Provides a full audit trail of changes over time.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned in PRIME.md. Useful for understanding how an issue evolved, especially after compaction or when debugging why a field changed.

---

#### bd lint
**What it does:** Checks issues for missing recommended sections based on issue type. By default lints all open issues. Can target specific issue IDs. Section requirements vary by type (bugs need reproduction steps, features need acceptance criteria, etc.).

**Agent knowledge:** KNOWN
**Gap:** Documented in PRIME.md under Quality Tools and Lifecycle & Hygiene.

---

#### bd stale
**What it does:** Shows issues that haven't been updated recently. Identifies in-progress issues with no recent activity (may be abandoned), forgotten open issues, and issues that might be outdated.

**Agent knowledge:** KNOWN
**Gap:** Documented in PRIME.md under Lifecycle & Hygiene.

---

#### bd status
**What it does:** Shows a snapshot of the issue database state and statistics. Provides counts by state (open, in_progress, blocked, closed), ready work count, extended statistics (pinned issues, average lead time), and recent activity from git history.

**Agent knowledge:** PARTIAL
**Gap:** PRIME.md mentions `bd stats` but the actual command is `bd status`. The CLAUDE.md project instructions reference `bd stats` which is an alias. The full output (lead time, recent activity) is not described.

---

#### bd statuses
**What it does:** Lists all valid issue statuses and their categories. Shows built-in statuses (open, in_progress, blocked, etc.) and any custom statuses configured via `status.custom`.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned anywhere. Agents know the standard statuses from context but don't know this discovery command exists or that custom statuses are possible.

---

#### bd types
**What it does:** Lists all valid issue types that can be used with `bd create --type`. Shows core types (bug, task, feature, chore, epic, decision) and any custom types configured via `types.custom`.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned anywhere. Agents know the standard types from context but don't know about the `decision` type or that custom types are configurable.

---

### Category: Dependencies & Structure

#### bd dep
**What it does:** Manages dependencies between issues. Core subcommands: add, remove, list, tree. Supports dependency types: blocks, related, parent-child, discovered-from. The `tree` subcommand visualizes the full dependency tree. Also supports `bd dep <blocker> --blocks <blocked>` syntax.

**Agent knowledge:** KNOWN
**Gap:** Well-documented in PRIME.md. The `--blocks` shorthand syntax and `tree` subcommand are in the CLAUDE.md project instructions.

---

#### bd duplicate
**What it does:** Marks an issue as a duplicate of a canonical issue. The duplicate is automatically closed with a reference to the canonical. Syntax: `bd duplicate <id> --of <canonical>`.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned in PRIME.md. Agents who find duplicates have no taught mechanism for marking them - they might close with a note or create a related dependency.

---

#### bd duplicates
**What it does:** Finds issues with identical content (title, description, design, acceptance criteria). Groups by content hash and suggests merge targets chosen by reference count and lexicographic ID. Supports `--auto-merge` for automatic deduplication and `--dry-run` for preview.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned anywhere. Different from `bd find-duplicates` (semantic similarity) - this finds exact content matches.

---

#### bd epic
**What it does:** Epic management commands. Subcommands: `close-eligible` (close epics where all children are complete), `status` (show epic completion status). Automates the common pattern of checking whether an epic can be closed.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned in PRIME.md. Agents manually check epic children and close epics. The `close-eligible` automation is unknown.

---

#### bd graph
**What it does:** Dependency visualization with multiple output formats. Default DAG with columns and box-drawing edges (terminal-native), `--box` for ASCII boxes with layers, `--compact` for tree format, `--dot` for Graphviz DOT output, `--html` for interactive D3.js visualization. Shows execution order where Layer 0 has no dependencies. Includes `check` subcommand for graph integrity verification.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned in PRIME.md. CLAUDE.md project instructions mention `bd dep tree` but not the richer `bd graph` visualization. The HTML interactive output would be useful for sharing with human users.

---

#### bd supersede
**What it does:** Marks an issue as superseded by a newer version. The superseded issue is automatically closed with a reference to the replacement. Designed for design docs, specs, and evolving artifacts.

**Agent knowledge:** KNOWN
**Gap:** Documented in PRIME.md under Lifecycle & Hygiene.

---

#### bd swarm
**What it does:** Structured parallel work on epics. Creates swarm molecules from epics with DAG-based coordination. Subcommands: create, list, status, validate. The `validate` subcommand checks whether an epic's structure is suitable for swarming before committing to it.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned in PRIME.md or beads-onboard. This is the primary mechanism for coordinating parallel agent work on complex epics and agents have zero awareness. The validate-before-swarm pattern is especially valuable.

---

### Category: Sync & Data

#### bd backup
**What it does:** Database backup for off-machine recovery. Subcommands: init (set up backup destination - filesystem or DoltHub), sync (push to destination), restore (restore from backup), remove (remove destination), status (show backup status).

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned in PRIME.md. Agents have no awareness of backup capabilities. Data loss recovery depends on Dolt push/pull which is taught, but dedicated backup to external destinations is not.

---

#### bd branch
**What it does:** Lists all Dolt branches or creates a new branch. Requires the Dolt storage backend. Enables branching the issue database itself for experimental changes.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned in PRIME.md. Part of the Dolt version control system - agents know about `bd dolt push/pull` but not branching.

---

#### bd export
**What it does:** Exports all issues to JSONL format. Each line is a complete JSON object with labels, dependencies, and comments. Used for migration, interoperability, and git-tracked snapshots.

**Agent knowledge:** KNOWN
**Gap:** Documented in CLAUDE.md project instructions as part of the commit workflow (`bd export -o .beads/issues.jsonl`).

---

#### bd federation
**What it does:** Peer-to-peer sync between Dolt-backed databases. Enables synchronized issue tracking across multiple workspaces. Subcommands: add-peer, list-peers, remove-peer, status, sync. Requires the Dolt storage backend.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned anywhere. Multi-workspace federation is an advanced feature for distributed teams.

---

#### bd import
**What it does:** Imports issues from a JSONL file with upsert semantics - new issues are created, existing issues are updated. Automatically detects memory records (lines with `_type: memory`) and imports them as persistent memories. Default source is `.beads/issues.jsonl`.

**Agent knowledge:** PARTIAL
**Gap:** CLAUDE.md project instructions mention `bd import .beads/issues.jsonl` for post-pull sync. The memory record detection and upsert semantics are not documented.

---

#### bd restore
**What it does:** Restores full history of a compacted issue from Dolt version history. When issues are compacted, descriptions and notes are truncated. This command queries Dolt's history tables to find pre-compaction content. Read-only - does not modify the database.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned in PRIME.md. Agents who encounter compacted issues have no way to know they can recover the full content.

---

#### bd vc
**What it does:** Git-like version control for issue data. Subcommands: commit, merge, status. Provides branching and merging capabilities for the Dolt issue database beyond simple push/pull.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned in PRIME.md. Agents know `bd dolt push/pull` but not the full version control capabilities (commit, merge, branch).

---

### Category: Setup & Configuration

#### bd bootstrap
**What it does:** Non-destructive database setup for fresh clones and recovery. Auto-detects the right action: clones from remote, clones from git, restores from backup, or imports from JSONL. Unlike `bd init --force`, never deletes existing issues.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned in PRIME.md. Agents in fresh-clone scenarios would try `bd init` which could be destructive. `bd bootstrap` is the safe alternative.

---

#### bd config
**What it does:** Manages configuration settings for integrations and preferences. Stores settings per-project in the beads database. Supports namespaces for integrations (jira.*, linear.*, github.*, etc.) and general preferences (validation.on-create, etc.).

**Agent knowledge:** PARTIAL
**Gap:** PRIME.md mentions `bd config set validation.on-create warn` but doesn't teach the config system broadly. Integration configuration namespaces are not documented.

---

#### bd context
**What it does:** Shows effective backend identity including repository paths, backend configuration, and sync settings. Reads directly from config files without requiring the database to be open - useful for diagnostics in degraded states.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned in PRIME.md. Useful for debugging when the database isn't working properly.

---

#### bd dolt
**What it does:** Manages the Dolt database server and settings. Subcommands: start, stop, status (server lifecycle), show, set, test (configuration), commit, push, pull (version control), remote add/list/remove (remote management), clean-databases, killall (maintenance).

**Agent knowledge:** PARTIAL
**Gap:** PRIME.md teaches `bd dolt push` and `bd dolt pull` for sync. Server lifecycle management (start/stop/status), configuration (set/show/test), and remote management are not taught.

---

#### bd forget
**What it does:** Removes a persistent memory by its key. Companion to `bd remember` and `bd memories`.

**Agent knowledge:** KNOWN
**Gap:** Documented in PRIME.md under Memory section.

---

#### bd hooks
**What it does:** Installs, uninstalls, or lists git hooks for beads integration. Provides pre-commit, post-merge, pre-push, post-checkout, and prepare-commit-msg hooks. The prepare-commit-msg hook adds agent identity trailers for forensics.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned in PRIME.md. CLAUDE.md project instructions reference git hooks in `examples/git-hooks/` but not the `bd hooks` management command.

---

#### bd human
**What it does:** Has two distinct behaviors: (1) With no arguments or the bare command, displays a focused help menu showing the ~15 essential commands for human users. (2) With subcommands (list, respond, dismiss, stats), manages human-flagged issues. The `bd human <id>` pattern (from older conventions) adds a `human` label to flag an issue for human attention. This is ADVISORY ONLY - it does NOT block the issue from appearing in `bd ready`.

**Agent knowledge:** PARTIAL
**Gap:** PRIME.md teaches `bd human <id>` for flagging and `bd human <id> --respond` for recording decisions. However, it does not clarify that this is advisory-only and does not block workflow. Agents may incorrectly believe `bd human` prevents an issue from being picked up by `bd ready`. The actual blocking mechanism is `bd gate` with a human-type gate, created through formulas - which agents don't know about.

---

#### bd info
**What it does:** Shows database information including absolute path, statistics (issue count), schema information (`--schema`), and what's new in recent versions (`--whats-new`).

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned in PRIME.md. Useful for debugging and orientation but agents use `bd doctor` instead.

---

#### bd init
**What it does:** Initializes bd in the current directory by creating `.beads/` directory and Dolt database. Supports custom prefix, backend specification, and agents-file configuration. The foundation command for any new beads project.

**Agent knowledge:** PARTIAL
**Gap:** CLAUDE.md project instructions show `bd init --prefix test`. PRIME.md doesn't cover init since it assumes the project already has beads. Advanced flags (--agents-file, --database) are untaught.

---

#### bd kv
**What it does:** Key-value store persisting across sessions. Subcommands: set, get, clear, list. Stores flags, environment variables, or user-defined data in the beads database. Values survive session compaction and agent rotation.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned anywhere. A general-purpose persistence mechanism that agents could use for cross-session state (feature flags, environment markers, workflow state) but don't know exists.

---

#### bd memories
**What it does:** Lists all persistent memories, or searches by keyword. The read companion to `bd remember`.

**Agent knowledge:** KNOWN
**Gap:** Documented in PRIME.md under Memory section.

---

#### bd onboard
**What it does:** Displays a minimal snippet to add to the agent instructions file for bd integration. The same snippet that `bd init` generates. Points to `bd prime` for full workflow context.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned in PRIME.md. The beads-onboard skill exists but it's a separate interactive process. The `bd onboard` command itself is a simple output command.

---

#### bd prime
**What it does:** Outputs essential beads workflow context in AI-optimized markdown. Auto-detects MCP mode (brief) vs CLI mode (full). Designed for Claude Code hooks (SessionStart, PreCompact) to inject workflow context. The primary training mechanism for agents.

**Agent knowledge:** KNOWN
**Gap:** Well-understood as the source of agent context. PRIME.md is the output of this command.

---

#### bd quickstart
**What it does:** Displays a quick start guide showing common bd workflows and patterns. Designed for new users.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned in PRIME.md. The beads plugin has a deprecated `beads:quickstart` skill.

---

#### bd recall
**What it does:** Retrieves the full content of a memory by its key. More targeted than `bd memories` which lists/searches all memories.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned in PRIME.md. Agents use `bd memories <keyword>` for search but don't know about exact key retrieval.

---

#### bd remember
**What it does:** Stores a persistent memory that survives sessions and account rotations. Memories are injected at prime time (`bd prime`) so they're available in every session. Supports `--key` for explicit key naming.

**Agent knowledge:** KNOWN
**Gap:** Documented in PRIME.md under Memory section.

---

#### bd setup
**What it does:** Installs integration files for AI editors and coding assistants. Built-in recipes: cursor, claude, gemini, aider, factory, codex, mux, opencode, junie, windsurf, cody, kilocode. Supports `--project` flag for workspace-level installation.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned in PRIME.md. Handles the initial editor integration that beads-onboard partially covers. Agents don't know about the multi-editor support.

---

#### bd where
**What it does:** Shows the active beads database location including redirect information. Useful for debugging when using redirects to understand which workspace is actually being used.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned in PRIME.md. Diagnostic command for multi-workspace setups.

---

### Category: Maintenance

#### bd batch
**What it does:** Runs multiple write operations in a single database transaction. Commands read from stdin or file. All operations execute in a single dolt transaction - any error rolls back the whole batch. Reduces write amplification compared to invoking `bd` many times in a loop.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned anywhere. Significant performance optimization for agents that create/update many issues in sequence. Could dramatically reduce Dolt write overhead in bulk operations.

---

#### bd compact
**What it does:** Squashes Dolt commits older than N days into a single commit. Preserves recent commits within the retention window via cherry-pick. Reduces storage overhead from auto-commit history. For semantic issue compaction (summarizing closed issues), use `bd admin compact` instead.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned in PRIME.md. Part of database maintenance that humans typically manage.

---

#### bd doctor
**What it does:** Sanity checks the beads installation. Checks `.beads/` directory, database version, migration status, schema compatibility, ID format, CLI version currency. Supports `--fix` for automatic repairs and `--check=conventions` for convention drift detection.

**Agent knowledge:** KNOWN
**Gap:** Documented in PRIME.md under Project Health. The `--check=conventions` flag is noted.

---

#### bd flatten
**What it does:** Nuclear option: squashes ALL Dolt commit history into a single commit. Uses the Tim Sehn recipe: create new branch, soft-reset to initial commit, commit everything as single snapshot, swap main branch, run GC. Irreversible.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned in PRIME.md. Extreme maintenance operation - appropriate that agents don't know about it.

---

#### bd gc
**What it does:** Full lifecycle garbage collection running three phases: (1) DECAY - delete closed issues older than N days (default 90), (2) COMPACT - squash old Dolt commits, (3) GC - run Dolt garbage collection. Each phase can be skipped individually. Supports `--dry-run`.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned in PRIME.md. Database maintenance that would typically be triggered by humans or scheduled automation.

---

#### bd migrate
**What it does:** Database migration and data transformation. Without subcommand, checks and updates database metadata to current version. Subcommands: hooks (plan git hook migration), issues (move issues between repositories), sync (set up sync.branch workflow).

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned in PRIME.md. Administrative command for version upgrades and repository reorganization.

---

#### bd preflight
**What it does:** Pre-PR checklist: lint, stale, orphans, and other common issues. Catches problems before pushing to CI.

**Agent knowledge:** KNOWN
**Gap:** Documented in PRIME.md under Lifecycle & Hygiene.

---

#### bd purge
**What it does:** Permanently deletes closed ephemeral beads (wisps, transient molecules) and their associated data. Reclaims storage from accumulated closed ephemeral issues. Skips pinned beads (protected).

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned anywhere. Cleanup for the molecule/wisp system that agents don't know about because they don't know about molecules.

---

#### bd rename-prefix
**What it does:** Renames the issue prefix for all issues in the database. Updates all IDs and text references across all fields. Use cases: shortening prefixes, rebranding, consolidating after corruption, migrating to team standards.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned in PRIME.md. The beads plugin has a `beads:rename-prefix` skill but PRIME.md doesn't teach it.

---

#### bd rules
**What it does:** Audits and compacts Claude rules. Subcommands: audit (scan for contradictions and merge opportunities), compact (merge related rules into composites). Meta-tool for maintaining the quality of the instruction files that govern agent behavior.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned in PRIME.md. Niche meta-tooling for maintaining CLAUDE.md/AGENTS.md quality.

---

#### bd sql
**What it does:** Executes raw SQL queries against the underlying Dolt database. Useful for debugging, maintenance, and working around bugs in higher-level commands. Supports any valid SQL including SELECT, UPDATE, DELETE.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned in PRIME.md. Power-user escape hatch. Potentially dangerous (can corrupt data) but valuable for debugging.

---

#### bd upgrade
**What it does:** Checks and manages bd version upgrades. Subcommands: status (check if version changed), review (show what's new), ack (acknowledge current version). Version tracking is automatic.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned in PRIME.md. Agents don't know how to check for or acknowledge version changes.

---

#### bd worktree
**What it does:** Manages git worktrees with proper beads configuration. Worktrees share the same beads database via git common directory discovery - no manual redirect needed. Enables parallel development across multiple working directories.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned in PRIME.md. Detailed docs exist at `docs/WORKTREES.md` but agents aren't pointed there. Important for multi-agent parallel development scenarios.

---

### Category: Integrations & Advanced

#### bd admin
**What it does:** Administrative database maintenance commands. Subcommands: cleanup (delete closed issues by age), compact (summarize old closed issues to save space), reset (remove all beads data). For routine maintenance, `bd doctor --fix` is preferred.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned in PRIME.md. Destructive operations that should require human approval anyway.

---

#### bd ado (Azure DevOps)
**What it does:** Bidirectional sync with Azure DevOps. Subcommands: projects, pull, push, status, sync. Configuration via `bd config` or environment variables (AZURE_DEVOPS_ORG, AZURE_DEVOPS_PROJECT, AZURE_DEVOPS_PAT, etc.).

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned anywhere. Part of the integration suite.

---

#### bd audit
**What it does:** Append-only JSONL interaction log at `.beads/interactions.jsonl`. Subcommands: record (append entry), label (add label referencing existing entry). Designed for auditing agent decisions and dataset generation (SFT/RL fine-tuning). Git-versionable.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned in PRIME.md. The beads plugin has a `beads:audit` skill but agents aren't taught the auditing workflow or when to use it.

---

#### bd blocked
**What it does:** Shows all blocked issues. Supports `--parent` to filter to descendants of a specific epic. Uses dependency-aware semantics to find issues that are actually blocked by open dependencies.

**Agent knowledge:** KNOWN
**Gap:** Documented in PRIME.md under Dependencies & Blocking.

---

#### bd completion
**What it does:** Generates shell autocompletion scripts for bash, fish, zsh, and PowerShell. Standard Cobra completion support.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned in PRIME.md. Shell completion is a human-user convenience, not relevant to agent workflows.

---

#### bd cook
**What it does:** Compiles a formula file into a proto. Two modes: compile-time (default, preserves `{{variable}}` placeholders for modeling/estimation) and run-time (substitutes variables for execution). The middle step of the Rig-Cook-Run lifecycle.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned in PRIME.md. Part of the formula/molecule system that is almost entirely untaught.

---

#### bd defer
**What it does:** Defers issues to put them on ice for later. Deferred issues are deliberately postponed - not blocked by dependencies, just set aside. Deferred issues don't show in `bd ready` but remain visible in `bd list`. Supports `--until` for timed undefer.

**Agent knowledge:** KNOWN
**Gap:** Documented in PRIME.md under Lifecycle & Hygiene.

---

#### bd formula
**What it does:** Manages workflow formulas - YAML/JSON files defining workflows with composition rules. The source layer for molecule templates. Rig-Cook-Run lifecycle: compose formulas, transform to proto, agents execute. Search paths: project (.beads/formulas/), user (~/.beads/formulas/), orchestrator ($GT_ROOT). Subcommands: list, show, convert.

**Agent knowledge:** PARTIAL
**Gap:** PRIME.md mentions `bd formula list` in a single line under "Structured Workflows" with zero explanation. Agents don't know what formulas are, when to use them, how to create them, or how they relate to molecules and gates.

---

#### bd github
**What it does:** Bidirectional sync with GitHub Issues. Subcommands: pull, push, repos, status, sync. Configuration via `bd config` or environment variables (GITHUB_TOKEN, GITHUB_OWNER, GITHUB_REPO, etc.). Supports GitHub Enterprise with custom API URLs.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned in PRIME.md. Agents working in GitHub-hosted projects don't know they can sync beads issues bidirectionally with GitHub Issues.

---

#### bd gitlab
**What it does:** Bidirectional sync with GitLab Issues. Subcommands: pull, push, projects, status, sync. Supports group-level sync with `gitlab.group_id`.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned anywhere. Part of the integration suite.

---

#### bd jira
**What it does:** Bidirectional sync with Jira. Subcommands: pull, push, status, sync. Supports Jira Cloud and Server. Configuration: URL, project, API token, username. Supports multiple projects and prefix-scoped push.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned in PRIME.md. Significant for enterprise environments where Jira is the canonical tracker.

---

#### bd linear
**What it does:** Bidirectional sync with Linear. Subcommands: pull, push, status, sync, teams. Supports multiple teams and project-scoped sync.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned in PRIME.md. Linear is popular among startups and dev-forward teams.

---

#### bd mail
**What it does:** Delegates mail operations to an external mail provider. Bridges the gap when agents type `bd mail` expecting mail functionality. Configured via environment variable or config setting to delegate to external commands (e.g., `gt mail`).

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned anywhere. Niche bridging command for orchestrator-integrated workflows.

---

#### bd mol
**What it does:** Work templates for agent workflows. Protos are template epics that define a DAG of work. Spawning creates molecules (real issues) from protos with variable substitution. Subcommands: pour (instantiate persistent), wisp (instantiate ephemeral), bond (combine protos/molecules), squash (compress to digest), burn (discard), distill (extract proto from ad-hoc epic), current, last-activity, progress, ready, seed, show, stale.

**Agent knowledge:** PARTIAL
**Gap:** PRIME.md mentions `bd mol pour <name>` in one line. Agents don't understand the proto/molecule/wisp/bond/distill concepts, variable substitution, the molecule lifecycle, or any subcommand beyond `pour`. The entire workflow templating system is functionally invisible.

---

#### bd notion
**What it does:** Bidirectional sync with Notion databases. Subcommands: connect (to existing database), init (create dedicated Beads database in Notion), pull, push, status, sync.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned anywhere. Part of the integration suite.

---

#### bd orphans
**What it does:** Identifies orphaned issues - issues referenced in commit messages but remaining open/in_progress. Helps find implemented work that wasn't formally closed. Supports `--fix` for interactive closing and `--label` for filtering.

**Agent knowledge:** KNOWN
**Gap:** Documented in PRIME.md under Lifecycle & Hygiene as part of the hygiene workflow.

---

#### bd ready
**What it does:** Shows work ready to be claimed - open issues with no active blockers. Excludes in_progress, blocked, deferred, and hooked issues. Uses blocker-aware semantics (not just status filtering). Supports `--mol` to filter within a specific molecule. Note: `bd list --ready` is NOT equivalent - it only filters by status=open.

**Agent knowledge:** KNOWN
**Gap:** Well-documented in PRIME.md as the primary work-finding command. The `--mol` filter and the distinction from `bd list --ready` are not taught.

---

#### bd rename
**What it does:** Renames an issue from one ID to another. Updates the primary ID, all references in other issues, dependencies, labels, comments, and events. Useful for giving memorable names to important issues.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned in PRIME.md. Could be useful for agents who want to give semantic names to important issues (e.g., renaming a hash-based ID to `bd-auth-bug`).

---

#### bd repo
**What it does:** Multi-repo hydration. Add/remove/list/sync across multiple beads repositories. Enables unified cross-repo issue tracking from a single database. Configuration in `.beads/config.yaml` under `repos` section.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned in PRIME.md. Important for monorepo or multi-service architectures where issues span repositories.

---

#### bd ship
**What it does:** Publishes a capability for cross-project dependencies. Finds issue with `export:<capability>` label, validates it's closed, adds `provides:<capability>` label. External projects depend on capabilities using `bd dep add <issue> external:<project>:<capability>`.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned anywhere. The entire cross-project dependency resolution system (export/provides labels, external dependencies) is untaught. Agents working across multiple beads projects have no way to coordinate capabilities.

---

#### bd undefer
**What it does:** Restores deferred issues to open status. Issues will appear in `bd ready` if they have no blockers. Supports multiple IDs.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned in PRIME.md. PRIME.md teaches `bd defer` but not the reverse. Agents who defer issues have no taught way to bring them back (they'd try `bd update --status open`).

---

#### bd version
**What it does:** Prints version information for the bd CLI.

**Agent knowledge:** UNKNOWN
**Gap:** Not mentioned in PRIME.md. The beads plugin has a `beads:version` skill. Minor gap - agents rarely need version info.

---

## Gap Analysis by Severity

### Tier 1: Critical Gaps (agents are making wrong decisions)

| Feature | Issue |
|---------|-------|
| **bd human** | Taught as blocking, actually advisory. Agents think flagging = blocking. |
| **bd gate** | The actual blocking mechanism is completely unknown. |
| **Markdown in fields** | Rich description/notes/design fields are used as flat text. |

### Tier 2: Major Gaps (significant capabilities invisible)

| Feature | Issue |
|---------|-------|
| **bd formula / bd mol / bd cook** | Entire workflow templating system ~invisible. Two lines in PRIME.md. |
| **bd query** | Powerful query language unknown; agents overuse `bd list` + jq. |
| **bd swarm** | Parallel work coordination unknown despite being core to multi-agent. |
| **bd merge-slot** | Serialization primitive for multi-agent merge conflicts unknown. |
| **bd ship** | Cross-project capability publishing unknown. |
| **bd set-state / bd state** | Operational state management unknown. |
| **bd batch** | Performance optimization for bulk operations unknown. |
| **bd kv** | General-purpose cross-session persistence unknown. |

### Tier 3: Moderate Gaps (useful but situational)

| Feature | Issue |
|---------|-------|
| **bd graph** | Rich visualization (HTML, DOT) unknown; only `bd dep tree` taught. |
| **bd find-duplicates** | Semantic duplicate detection unknown. |
| **bd bootstrap** | Safe clone recovery unknown; agents might use destructive `bd init`. |
| **bd comments / bd comment** | Issue-level threaded discussion unknown. |
| **bd children** | Convenient epic child listing unknown. |
| **bd epic** | Epic close-eligible automation unknown. |
| **bd promote** | Wisp-to-bead promotion unknown. |
| **bd reopen** | Proper reopen with event tracking unknown. |
| **bd undefer** | Reverse of taught `bd defer` is untaught. |
| **Integration commands** | GitHub, Jira, Linear, GitLab, ADO, Notion sync all unknown. |
| **bd worktree** | Parallel development with shared beads DB unknown. |
| **bd vc / bd branch / bd diff** | Full Dolt version control beyond push/pull unknown. |

### Tier 4: Low Impact (admin/niche, appropriate to leave untaught)

| Feature | Issue |
|---------|-------|
| **bd admin** | Destructive ops - should require human anyway. |
| **bd flatten** | Nuclear history reset - human-only. |
| **bd sql** | Raw SQL escape hatch - dangerous for agents. |
| **bd create-form** | Interactive TUI - blocks agents. |
| **bd completion** | Shell completion - human convenience. |
| **bd rules** | Meta-tooling for rule maintenance. |
| **bd version** | Version info - rarely needed. |
| **bd mail** | Orchestrator bridging - niche. |
