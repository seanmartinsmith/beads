---
description: Set up project-specific quality conventions for beads
argument-hint:
---

Onboard a project with beads quality conventions. Analyzes the codebase and
interactively generates project-specific rules that make agents produce
high-quality, searchable, non-redundant beads.

Run this once per project, or re-run after upgrading to get improved guidance.

## Phase 0: State Detection

Before starting the interview, detect the current project state.

**Onboard version:** v3

1. Check for `.beads/` directory:
   - **If missing:** offer to initialize with AskUserQuestion:
     "No beads found. Initialize now? The prefix is used for issue IDs
     (e.g., `<dir-name>-1`, `<dir-name>-2`)."
     - Option 1: "Yes, use `<dir-name>` as prefix" (auto-detected from directory name)
     - Option 2: "Yes, but let me choose a prefix" (prompt for custom prefix)
     - Option 3: "Yes, with team mode" (runs `bd init --team` for multi-user sync)
     - If user selects option 1: run `bd init --prefix <dir-name>`
     - If user selects option 2: ask for prefix, then run `bd init --prefix <custom>`
     - If user selects option 3: run `bd init --team --prefix <dir-name>` (or custom)
     - If user declines: stop. "Run `bd init` when ready, then re-run `/beads:onboard`."
     - Continue to step 2 after init completes.
2. Read `.beads/config.yaml` for the issue prefix.
3. Check for `.beads/conventions/.onboard-state.yaml` and `.beads/conventions/` directory:
   - **No state file AND no conventions dir:** first-time onboard. Proceed normally.
   - **No state file BUT conventions dir exists:** onboarded by older version (pre-v3).
     Read existing convention files to extract previous answers as defaults.
     Tell user: "Conventions exist from an earlier onboard version. v3 adds description
     templates, priority calibration, and lifecycle guidance. Re-running the interview
     with your existing answers as defaults." Offer: "Update with new guidance
     (Recommended)" or "Start completely fresh."
   - **State file exists, version < current:** read saved answers for pre-filling.
     Tell user what's new in this version. Same offer as above.
   - **State file exists, version = current:** "Conventions are current (v3). Want to
     review/update your answers, or start fresh?"
   - In all re-onboard cases, check if convention files were manually edited after
     generation (compare file content against what saved answers would produce). If
     customized, warn: "Your convention files have custom edits. New guidance will be
     merged alongside your customizations."
4. Check for `<!-- BEGIN BEADS INTEGRATION -->` markers in AGENTS.md:
   - If markers found: will replace content between markers (interop with `bd setup`).
   - If no markers but beads content found: will wrap replacement in markers.
   - If no AGENTS.md: will create it with markers.
5. Check if CLAUDE.md exists and whether it references `@AGENTS.md`.
   - If CLAUDE.md has substantial content beyond `@AGENTS.md` (more than 5 non-empty
     lines after removing the `@AGENTS.md` reference), flag it for Phase 4 guidance.
     Store the content summary (line count, what it covers) for the recommendation.
6. Run `bd doctor --agent 2>&1 || bd doctor 2>&1 || true` (try `--agent` first
   for rich agent-facing output, fall back to plain `bd doctor` for older versions).
   Capture output, don't show it yet. Note errors and warnings to surface as
   recommendations after generation.
7. Detect project state for Phase 1 detection strategy:
   - **Existing project:** significant code exists (>5 source files or package manager
     config like package.json, go.mod, Cargo.toml, pyproject.toml).
   - **Planned project:** minimal/no code but planning docs exist (check
     `docs/brainstorms/`, `docs/plans/`, `docs/`, root `.md` files).
   - **Empty project:** neither code nor planning docs.
   Store the detected state for Phase 1.

**Do not block on doctor results.** The hard gate is `.beads/` existence (or willingness
to init).

## Phase 1: Analyze the Project

Detection adapts based on project state from Phase 0. The goal is the same
regardless of state: gather defaults for the interview questions.

### If existing project (has code):

1. **Directory structure**: List top-level directories. Map each to a candidate
   area label (e.g., `cmd/` → `area:cli`, `internal/` → `area:internal`,
   `docs/` → `area:docs`). Include subdirectories that represent distinct areas.
2. **Languages and frameworks**: Check file extensions, package files
   (package.json, go.mod, Cargo.toml, etc.), and build configs.
3. **Commit conventions**: Read the last 20 commit messages. Detect if
   `type(scope): description` or another format is in use.
4. **Issue prefix**: Read from `.beads/config.yaml` (`issue-prefix` field).
5. **Existing agent instructions**: Check AGENTS.md and CLAUDE.md for existing
   beads-related content that needs reconciliation.
6. **Cross-cutting dimensions**: Auto-detect which label dimensions to include:
   - **Concern labels** (always): performance, security, ux, tests, accessibility
   - **Workflow labels** (always): investigate, brainstorm, collaborative
   - **Platform labels** (if cross-platform): windows, macos, linux, mobile, ios, android
   - **Browser labels** (if web project): firefox, safari, chrome

### If planned project (has docs, minimal code):

1. **Read planning docs**: scan `docs/brainstorms/`, `docs/plans/`, `docs/`, and
   root `.md` files. Extract: project description, planned directory structure,
   tech stack mentions, feature lists.
2. **Map planned structure to areas**: if docs describe a file/directory layout,
   derive candidate area labels from it.
3. **Languages and frameworks**: extract from tech stack mentions in docs.
4. **Commit conventions**: no history to detect. Default: `type(scope): description`.
5. **Issue prefix**: read from `.beads/config.yaml`.
6. **Cross-cutting dimensions**: infer from project description (e.g., "cross-platform
   CLI" → platform labels, "web app" → browser labels).

### If empty project (no code, no docs):

1. **No detection possible.** Use universal defaults:
   - Areas: `area:core` (placeholder — user will customize in Q2)
   - Commit format: `type(scope): description`
   - Languages: unknown
2. **Issue prefix**: read from `.beads/config.yaml`.
3. **Cross-cutting dimensions**: include concern and workflow labels (always useful).
   Skip platform/browser unless user adds them in Q2.

### All modes — project type detection:

Detect project type for priority calibration in Phase 3:
- **Production app**: has deployment configs (Dockerfile, vercel.json, netlify.toml,
  AWS configs, CI/CD with deploy steps).
- **Library/SDK**: has library indicators (go.mod with module path, package.json
  with `main`/`exports`/`types` fields, gem spec, setup.py with classifiers).
- **CLI tool**: has CLI entry points (cmd/ directory, bin/ scripts, console scripts).
- **Personal/internal**: no deployment configs, no library exports.
If unclear, default to "personal/internal" (most conservative priority calibration).

## Phase 2: Interview

Ask these questions **one at a time** using AskUserQuestion. Wait for each
answer before proceeding to the next question.

### Q1: Project Description

Generate three short descriptions (1-2 sentences each) that ground the project
in reality. Each description must capture: what the project IS (function, not
aspiration), what TYPE of project it is (production web app, CLI tool, library,
personal project), and who/what it serves. These descriptions help agents scope
their work — a portfolio site on S3 is fundamentally different from an internal
API serving production traffic.

Present them as options (the built-in "Other" handles freeform):

- Option 1: [Grounded description: what it is, project type, audience]
- Option 2: [Technical description: stack, architecture, deployment context]
- Option 3: [Purpose description: problem it solves, value delivered]

For planned/empty projects: base descriptions on planning docs or user input
rather than codebase analysis. Describe what the project IS (function, type,
audience), not its development stage. "A Go TUI for browsing usage data" is
stable whether the project is in planning or fully built. "Planning docs only,
no code yet" becomes stale in a week — don't include development status.

### Q2: Area Labels

Present the areas detected in Phase 1 as a structured list. Explain:
"Each bead gets exactly one area label. This is the highest-leverage
convention - it makes search, filtering, and ownership clear."

Mark the full detected list as "(Recommended)" - more granular areas reduce
ambiguity for agents assigning labels. Only recommend consolidation if you
detect 15+ areas.

Options:
- Accept these areas (Recommended) (show the detected list with descriptions)
- Consolidated (suggest fewer, broader groupings)
- Modify the list (add, remove, or rename)

### Q3: Commit Format

First determine context: "Is this an open source project where external
contributors submit PRs?"

- If **yes** (open source): Default to `type(scope): description` without bead
  references. Bead refs in commits are noise to upstream contributors.
- If **no** (personal/team): Default to `type(scope): description [prefix-xxx]`
  with optional bead reference for traceability.

**Recommendation logic:** Don't blindly recommend keeping the detected format.
Evaluate it: if the detected format lacks scope (e.g., `type: description`
with no parenthesized scope), recommend `type(scope): description` instead.
Only recommend "keep current" when the detected format already includes scope.

Every option must include a markdown preview showing the SAME example commits
rewritten into that option's format. This makes the difference between options
visually obvious. Use real examples from git log if available, or generate
realistic examples using the project's area labels as scopes. The `(scope)` is
the area of the codebase being changed — it maps to the area labels from Q2
(e.g., `feat(data):` means a feature in the data area). Options without
previews show "No preview available" which is confusing.

Options (always show all four):
- Keep current format (show it) - only mark "(Recommended)" if it already has scope
- type(scope): description (with scope, no bead ref) - recommend for open source, or when detected format lacks scope
- type(scope): description [prefix-xxx] (with scope + bead ref) - recommend for personal/team projects
- type: description [prefix-xxx] (no scope, with bead ref)

### Q4: Code Change Policy

"What's your bar for creating a bead?"

- Changelog-worthy changes + anything persisting between sessions (recommended)
- All code changes (strict - every change gets a bead)
- Only multi-session work (minimal - beads for persistence only)

### Q5: Beads + Claude Code Tasks

Before presenting the question, explain the context:

"Beads is your project's issue tracker — what work exists, why, and what
happened. Claude Code Tasks (TaskCreate, TaskUpdate, TaskList) is a built-in
checklist for tracking your steps while working, and complements beads well."

Then explain the complementary pattern:
- Beads = project memory. Cross-session, changelog-worthy, searchable history.
  Git-backed, syncs across machines and teams, works with any agent tool.
  "What work exists and why."
- Tasks = session memory. Execution steps, implementation order within a session.
  Disk-backed, survives compaction. "Where I am right now."
- The pattern: claim a bead → break it into 3-7 tasks → work the tasks → close
  the bead with a rich close_reason capturing what the tasks accomplished.

Then ask: "Do you use Claude Code Tasks alongside beads for within-session
execution?"

Options:
- Yes, I use both (generate full integration guidance in reference.md)
- No, just beads (skip detailed tasks integration guidance in reference.md)

### Q6: Close Outcome Template

Show the default template:

```
Summary: <one sentence>
Change: <what changed>
Files: <paths>
Verify: <how verified>
Risk: <if any>
Notes: <optional gotchas>
```

"This is the format agents will use when closing beads. Each field is
searchable. Want to customize it?"

Every option must include a markdown preview showing the template for that
choice. Options without previews show "No preview available" which is confusing.

Options:
- Use this template as-is (recommended) - preview: full 6-field template
- Customize fields (let them add/remove/rename) - preview: full template with note about which fields to change
- Simpler (just Summary + Change + Files) - preview: 3-field template

## Phase 3: Generate

Create `.beads/conventions/` directory if it doesn't exist. Then generate these
files using interview answers to fill the template slots marked with `[brackets]`.

### File 1: `.beads/conventions/labels.md`

```markdown
# Label Taxonomy

Use structured labels to categorize work. **Do not invent new labels** without
updating this file.

## Area

Exactly one `area:*` label is required per issue. Area labels describe ownership
of the change.

| Label | Use for |
|---|---|
[Generate one row per area from Q2, e.g.:]
| `area:cli` | CLI commands, flags, user-facing behavior |

[Include additional dimension tables based on Phase 1 detection. Use these templates:]

### Platform
| Label | Use for |
|---|---|
| `platform:windows` | Windows-specific behavior |
| `platform:macos` | macOS-specific behavior |
| `platform:linux` | Linux-specific behavior |

### Concern (cross-cutting)
| Label | Use for |
|---|---|
| `performance` | Speed, memory, optimization |
| `security` | Auth, vulnerabilities, input sanitization |
| `ux` | User experience, error messages, output quality |
| `tests` | Test failures, coverage gaps |

### Workflow
| Label | Use for |
|---|---|
| `workflow:investigate` | Unknown root cause; must close with root cause |
| `workflow:brainstorm` | Ideas, exploration, not yet implementation |
| `workflow:collaborative` | Needs human input or multi-agent coordination |

## Rules

- Labels must be assigned at creation time.
- Area labels describe ownership of the change, not symptoms.
- `workflow:investigate` must close with root cause + resolution or follow-up.
- Combine for specificity: `area:cli,platform:windows,performance`

## Examples

[Generate 3-4 `bd create` examples using this project's prefix and labels]
```

### File 2: `.beads/conventions/reference.md`

```markdown
# Beads Quality Reference

Detailed conventions for this project. AGENTS.md has the summary; this file
has the depth.

## Close Outcome Template

[Use template from Q6, default or customized]

Summary: <one sentence>
Change: <what changed>
Files: <paths>
Verify: <how verified>
Risk: <if any>
Notes: <optional gotchas>

**Why each field matters:**

| Field | Purpose |
|---|---|
| Summary | Search target — what `bd search` finds when someone hits the same problem |
| Change | Scope — what was actually modified, so future agents know the blast radius |
| Files | Navigation — where to look when verifying or extending this work |
| Verify | Confidence — how to confirm the fix still holds after other changes |
| Risk | Landmines — what could break, so future agents don't step on it |
| Notes | Gotchas — the thing that took 30 minutes to figure out, saved for the next agent |

**Formatting:** The `--reason` field is searchable and renders as markdown.
Use **blank lines** between sections (double newline). This is critical for
readability — single newlines without blank lines won't render as line breaks.

**Good example — note the literal newlines in the bash string:**

```bash
bd close [prefix]-xxx --reason="Summary: [one-sentence summary]

Change: [what changed]

Files: [paths]

Verify: [how verified]

Risk: [if any]

Notes: [optional gotchas]"
```

[Also generate one filled-in example using this project's prefix and realistic
content, showing the same bash syntax with literal newlines]

**Bad examples:**
- `"Done"` - no context, unsearchable
- `"Fixed it"` - no details for future agents
- `"Summary: x; Change: y; Files: z"` - semicolons on one line, unreadable
- Single newlines without blank lines - won't render as separate sections

## Creating Issues

Always include: `--type`, `--priority`, `--labels`, `--description`.

Priority: 0 (critical) to 4 (backlog).
Search before creating: `bd search "query"`

## Dependencies

- Only add blocking deps when work truly cannot start
- Default to NOT adding blocking deps
- Use `bd dep relate <a> <b>` for connected but parallel work
- Never use `blocks` to mean "belongs to epic" - use parent-child

## Description Structure

Structure descriptions by issue type. Without explicit section guidance, agents
default to one-line descriptions that are unsearchable and useless for future
sessions.

### Epic
Include when relevant: Context, Core Issues, Phases [NOW/SOON/LATER], Key
Decisions, Related beads

### Bug
Include when relevant: Observed behavior, Expected behavior, Steps to reproduce,
Environment

### Feature
Include when relevant: Motivation (why this matters), Acceptance criteria, Scope
boundary (what is NOT included)

### Task
Include when relevant: Context (why now, not later), Scope, Done-when

### Chore
Include when relevant: What + why now (not just "update deps" — state the reason,
e.g., "because CVE-2026-xxxx affects production")

### Examples

[Generate one complete `bd create` command per type using the project's actual
prefix and labels. Descriptions must be filled in with realistic content for this
project, not placeholder text. These are the examples agents pattern-match from.]

**Bug example:**
```
bd create --title="[realistic bug title for this project]" --type=bug --priority=1 \
  --labels=area:[detected-area],ux \
  --description="Observed: [specific behavior in this project's context]

Expected: [correct behavior]

Steps to reproduce:
1. [step using this project's actual commands/UI]
2. [step]

Environment: [relevant platform/browser/version]"
```

**Feature example:**
```
bd create --title="[realistic feature for this project]" --type=feature --priority=2 \
  --labels=area:[detected-area] \
  --description="Motivation: [why this matters for this project's users]

Acceptance criteria:
- [concrete criterion using this project's context]
- [concrete criterion]

Scope boundary: Does NOT include [explicit exclusion]"
```

**Task example:**
```
bd create --title="[realistic task for this project]" --type=task --priority=2 \
  --labels=area:[detected-area] \
  --description="Context: [why this task is needed now]

Scope: [what this covers]

Done-when: [concrete completion signal]"
```

**Chore example:**
```
bd create --title="[realistic chore for this project]" --type=chore --priority=3 \
  --labels=area:[detected-area] \
  --description="Update [specific dependency] because [specific reason — CVE,
breaking change in downstream, compatibility with new version of X].

Scope: [what's included in this update]"
```

**Epic example:**
```
bd create --title="[realistic epic for this project]" --type=epic --priority=1 \
  --labels=area:epic \
  --description="Context: [why this initiative exists]

Core Issues:
- [issue 1]
- [issue 2]

Phases:
[NOW] [immediate work]
[SOON] [next steps after NOW completes]
[LATER] [future work, may change]

Key Decisions:
- [decision with rationale]

Related: [prefix]-xxx, [prefix]-yyy"
```

## Priority Semantics

Priority determines creation-time classification. `bd ready` handles work
ordering via the dependency graph — priority is a tiebreaker when multiple
unblocked items exist.

[Generate the table matching the detected project type from Phase 1. If project
type is unclear, use the personal/internal table as default.]

**For production apps:**

| Priority | Meaning | Example |
|---|---|---|
| P0 | Site down, data loss, security incident | Database failures, auth broken, data corruption |
| P1 | Core feature broken, major regression | Checkout error, search returns nothing, API 500s |
| P2 | Meaningful improvement or non-critical bug | Visual glitch, slow page load, confusing UX |
| P3 | Minor polish, non-urgent | Tooltip wording, button alignment, minor refactor |
| P4 | Backlog, speculative | "Someday" ideas, exploratory features |

**For CLI tools:**

| Priority | Meaning | Example |
|---|---|---|
| P0 | Data loss or corruption | Corrupts config, deletes wrong data, silent failures |
| P1 | Core command broken, blocking workflow | Main command crashes, output format broken |
| P2 | Confusing behavior or error message | Misleading error, unclear help text, wrong exit code |
| P3 | Minor friction | Extra flag for common operation, verbose output |
| P4 | Backlog | Nice-to-have features, exploratory |

**For libraries/SDKs:**

| Priority | Meaning | Example |
|---|---|---|
| P0 | Breaking change in published API, data loss | Removed public method, corrupts caller state |
| P1 | Incorrect behavior in core functionality | Wrong return value, silent error swallowing |
| P2 | Missing docs, poor error messages | Unclear API, confusing type signatures |
| P3 | Minor API ergonomics | Better defaults, convenience methods |
| P4 | Backlog | Speculative features, internal refactors |

**For personal/internal projects:**

| Priority | Meaning | Example |
|---|---|---|
| P0 | Blocks all other work | Build broken, environment unusable |
| P1 | Core functionality broken | Main feature doesn't work, data issues |
| P2 | Meaningful improvement | New feature, notable bug fix |
| P3 | Minor polish | UI tweak, code cleanup |
| P4 | Backlog | Ideas, experiments |

### Scope awareness

When creating beads, consider complexity alongside priority. A self-contained
fix in one file is a different bead from a multi-file refactor, even if they
address the same problem. Smaller, focused beads are easier for agents to
complete in a single context window. If a bead's scope exceeds what you can
complete in your current session, break it into sub-beads before starting.

## Create Fields Beyond --description

`bd create` supports fields that add structured, searchable context. Use them
when the trigger condition applies. An empty `--notes ""` is worse than omitting.

| Field | Trigger | Decision Test |
|---|---|---|
| `--description` | Always | -- |
| `--notes` | Discovered something non-obvious during investigation | "Would a future agent hit this same surprise?" |
| `--design` | Chose between alternatives | "Is there a 'why not the other way?' question?" |
| `--acceptance` | Done-when isn't obvious from the title | "Could two agents disagree on whether this is complete?" |
| `--spec-id` | A planning doc is driving this work | "Is there a doc that explains why this bead exists?" |

**Examples using this project's prefix:**

```
--notes="[Realistic gotcha for this project, e.g.: Redis pool must be >=5 for
concurrent workers. Default of 3 causes intermittent failures under load.]"

--design="[Realistic design choice, e.g.: Redis pub/sub over WebSockets because
deployment has no sticky sessions.]"

--acceptance="[Realistic criteria, e.g.: Lighthouse >=90, loads in <2s on 3G
throttle, works in Safari 16+.]"

--spec-id="docs/plans/[relevant-plan-file].md"
```

## Enriching Beads During Work

Don't wait until creation and close to add context. Use `bd update` at specific
moments during implementation to build structured, searchable metadata.

| Moment during work | Action | Instead of |
|---|---|---|
| Made a design decision | `bd update <id> --design "chose X because Y"` | Commenting "decided to use X because Y" |
| Found an edge case | `bd update <id> --append-notes "edge case: Z"` | Commenting "note: also affects Z" |
| Realized "done" criteria changed | `bd update <id> --acceptance "new criteria"` | Commenting "realized we also need Y" |
| Created a related doc | `bd update <id> --spec-id "docs/..."` | Commenting "see docs/plans/foo.md" |
| Scope became clearer | `bd update <id> --add-label concern` | Not updating labels at all |

## When to Use Comments vs Fields

> **Decision test:** Would this information still be useful if all comments were
> deleted? If yes — it belongs in a structured field, not a comment.

**Comments are right for:**
- Session handoff notes: "Stopped at step 3, X works, Y needs Z"
- Questions or blockers: "Need clarification on IE11 support"
- Progress checkpoints: "Phase 1 complete, starting Phase 2"
- External references: "Related: github.com/..."
- Post-close addenda: "Found edge case after closing, see [prefix]-xxx"

**Comments are wrong for (use fields instead):**

| Instead of this comment... | ...use this update |
|---|---|
| "Decided to use Redis because no sticky sessions" | `bd update <id> --design "Redis pub/sub — no sticky sessions"` |
| "Note: config falls back to env vars in Docker" | `bd update <id> --append-notes "Config falls back to env vars when no file (Docker)"` |
| "Done when Lighthouse >= 90" | `bd update <id> --acceptance "Lighthouse >= 90, <2s on 3G"` |

## Session Discipline

- `bd show <id>` and read close_reason before starting work
- Check `bd list --status=in_progress` at session start
- Use `bd comments add <id> "..."` for session handoff notes
- Use `bd update` fields for structured metadata (design, notes, acceptance)
- Close beads before committing
- Run `bd dolt push` before ending session (if remote configured)
- Consider context window scope when claiming work — if a bead requires
  touching many files or systems, consider breaking it into sub-beads

## Beads + Claude Code Tasks

[If Q5 = yes, include full integration pattern. If no, include a brief mention.]

Beads and Tasks serve different levels of the same workflow:

| | Beads (`bd`) | Tasks (`TaskCreate`) |
|---|---|---|
| **Level** | Project-level work items | Execution-level steps |
| **Tracks** | What work exists and why | How you're implementing a specific bead |
| **Persists** | Yes (git-backed, survives everything) | Yes (disk-backed, survives compaction) |
| **Granularity** | Changelog-worthy changes | 1-3 files, completable in one focused stretch |
| **Dependencies** | Full DAG (blocks, relates, parent-child) | Simple DAG (blockedBy/blocks) |

### The integration pattern

1. **Find work:** `bd ready` → pick a bead
2. **Claim it:** `bd update <id> --status=in_progress`
3. **Break it down:** read the bead description, create 3-7 Tasks referencing
   the bead ID in each task subject
4. **Execute:** work through Tasks, mark each `completed` as you go
5. **Close the bead:** `bd close <id> --reason="..."` — the close_reason
   captures what the Tasks accomplished (handoff from execution to history)

### When to create Tasks for a bead

Create Tasks when:
- The bead has 3+ distinct implementation steps
- Complex work where tangents are likely
- You're mentally planning a sequence of changes

Skip Tasks when:
- Single-step fix (just do it)
- Exploratory work (investigate first, break down after)
- Simple enough to hold in your head

### Anti-patterns
- **Orphan Tasks:** tasks without a bead — work becomes untraceable in project history
- **Scope drift:** tasks that wander beyond the bead's scope — create a new bead instead
- **Task hoarding:** 10+ tasks means the bead is too big — split the bead first

### When to create a bead (vs just working)

Would it show up in a changelog? Will you need context in a future session?
Is there a decision worth recording? If yes to any → create a bead.

## Connecting Plans, Docs, and Commits

- In docs: `<!-- Related: [prefix]-xxx -->` near the top
- In beads: `bd comments add <id> "Doc: path/to/doc.md"`
- In commits: [commit format from Q3]
- In plans: Reference bead IDs the plan addresses

## Convention Health

Periodically audit bead quality: Are descriptions structured? Are close reasons
rich? Are labels consistent? If conventions feel stale or agents aren't following
them, re-run `/beads:onboard` to upgrade with improved guidance.
```

### File 3: AGENTS.md (beads section)

This section coordinates with `bd setup`'s marker system. The Go binary
(`cmd/bd/setup/agents.go`) writes a generic beads section wrapped in
`<!-- BEGIN BEADS INTEGRATION -->` / `<!-- END BEADS INTEGRATION -->` markers.
Onboard's output is project-specific and replaces the generic content.

**Marker strategy:**
- If markers exist in AGENTS.md: replace everything between them (inclusive).
- If no markers exist: append with markers for future `bd setup` interoperability.
- If no AGENTS.md: create the file with markers.

Always wrap the generated section in these exact markers:

```markdown
<!-- BEGIN BEADS INTEGRATION -->
## Issue Tracking

> This project uses [beads](https://github.com/steveyegge/beads) for task tracking.

**MANDATORY**: Read these files before creating, updating, or closing any issue:
1. `.beads/conventions/reference.md` — issue lifecycle, field triggers, close format
2. `.beads/conventions/labels.md` — valid label taxonomy

Do NOT create issues from memory or pattern-match from examples above.
The conventions docs are the source of truth.

[Code change policy from Q4]

### Commit Format

[Format from Q3]

### Quick Reference

| Action | Command |
|---|---|
| Find work | `bd ready` |
| Read before work | `bd show <id>` |
| Claim | `bd update <id> --claim` |
| Complete | `bd close <id> --reason="..."` |
| Session notes | `bd comments add <id> "..."` |
| Search | `bd search "query"` |
| Sync | `bd dolt push` |

### Creating Issues

Always include: `--type`, `--priority`, `--labels`, `--description`.
Use `--notes`, `--design`, `--acceptance` when trigger conditions apply
(see reference.md for triggers and decision tests).
Valid labels: `.beads/conventions/labels.md`
Description structure by type: `.beads/conventions/reference.md`

### Close Outcome Format

Use **literal newlines with blank lines** between fields in `--reason`:

[Template from Q6, shown as a bash command with literal newlines like:]
```bash
bd close <id> --reason="Summary: ...

Change: ...

Files: ...

Verify: ...

Risk: ...

Notes: ..."
```

Never use semicolons or single-line formatting. See reference.md for full guidance.

### Session Rules

- Read close_reason before working a bead to avoid re-solving
- Check for abandoned work: `bd list --status=in_progress`
- Use comments for session handoff notes (survives compaction)
- Use `bd update` fields for structured metadata during work
- Close beads before committing
- Only add blocking deps when work truly cannot start
- Don't invent labels - use `.beads/conventions/labels.md`
- Do NOT use `bd edit` - it opens $EDITOR (interactive). Use `bd update <id> --field "value"` instead

### During Work

- Use `bd update <id> --design "..."` when you make a non-obvious choice
- Use `bd update <id> --append-notes "..."` when you discover edge cases
- Use comments for session progress, questions, and handoffs
- Full lifecycle guide: `.beads/conventions/reference.md`

[If Q5 = yes:]
### Beads + Tasks

Beads for cross-session persistence. Tasks for within-session execution.

<!-- END BEADS INTEGRATION -->
```

**Post-write note:** If this replaced `bd setup`'s generic section, mention in
Phase 4 output: "Replaced the generic `bd setup` beads section with
project-specific conventions."

### File 4: CLAUDE.md (if needed)

If CLAUDE.md exists but doesn't reference AGENTS.md, add `@AGENTS.md` near
the top. If CLAUDE.md doesn't exist, create it with:

```markdown
@AGENTS.md
```

### File 5: `.beads/conventions/.onboard-state.yaml`

Write a state file for re-onboard detection. This file is NOT loaded by agents
during normal work — it's only read by the onboard command on re-run.

```yaml
# Generated by /beads:onboard — do not edit manually
version: "v3"
timestamp: "[ISO 8601 timestamp of generation]"
answers:
  project_description: "[Q1 answer text]"
  project_type: "[detected: production-app|cli-tool|library|personal-project]"
  area_labels:
    - "[label]: [description]"
  commit_format: "[Q3 selected format]"
  code_change_policy: "[Q4 answer: changelog-worthy|all-changes|multi-session-only]"
  beads_tasks_integration: "[Q5 answer: yes|no|unsure]"
  close_template: "[Q6 answer: default|custom|simple]"
  close_template_fields:
    - "[field names from selected template]"
```

## Phase 4: Validation

After generating all files:

1. Verify each file was written successfully (read them back).
2. Show summary:

| File | Purpose |
|---|---|
| `.beads/conventions/labels.md` | Label taxonomy — [N] areas + [dimensions] |
| `.beads/conventions/reference.md` | Quality reference — close template, description structure, priority, lifecycle |
| `.beads/conventions/.onboard-state.yaml` | Re-onboard state (version, answers) |
| `AGENTS.md` | Scannable enforcement layer — agents read this every session |
| `CLAUDE.md` | @AGENTS.md reference (if created/updated) |

3. Highlight the label taxonomy (list the area labels).
4. Highlight the close outcome template.
5. Surface any recommendations from Phase 0 doctor results:
   - Team mode not configured → "Consider `bd init --team` for multi-user sync"
   - Editor hooks missing → "Consider `bd setup claude` or `bd setup cursor`"
   - Other doctor warnings → list them
6. **If CLAUDE.md has substantial content (flagged in Phase 0):**
   Surface this recommendation with reasoning:

   "Your CLAUDE.md has [N lines] of project-specific content (covering [summary]).
   Beads writes conventions to AGENTS.md because it's tool-agnostic — AGENTS.md
   is the open standard that works across Claude Code, Cursor, Codex, and other
   AI coding tools. Your CLAUDE.md uses `@AGENTS.md` to pull those conventions in.

   The issue: when CLAUDE.md also has project instructions (build commands,
   architecture notes, workflow rules), agents read both sources and can get
   mixed signals when instructions overlap or conflict with the beads conventions.

   Recommendation: Move project-specific instructions into AGENTS.md (outside
   the beads markers) so there's one authoritative source. Keep CLAUDE.md
   minimal — just `@AGENTS.md` plus any personal preferences or Claude
   Code-specific settings.

   Want me to show you what's in CLAUDE.md so you can decide what to move?"

   - If user says yes: display the CLAUDE.md content and let them decide what
     to migrate. Do NOT modify or delete CLAUDE.md content without explicit approval.
   - If user declines: move on. This is a recommendation, not a gate.
7. Suggest: "Run `bd doctor` to verify full project health."
8. Remind: "These conventions are now active. Agents will follow them for every
   bead created in this project. Review and adjust as needed."
9. Offer milestone bead. Adapt the explanation based on project state:

   **If first-time onboard (no existing beads):**
   "This onboard created conventions, but since this project has no existing
   beads the guidance was based on predictions. After 30 beads of real usage,
   you'll have data to see what's actually working — are agents enriching
   beads with structured fields or just leaving comments? Is the label
   taxonomy still relevant or do you need more areas? Are close reasons
   structured or falling back to one-liners? A reminder bead surfaces this
   at the right time."

   **If re-onboard (existing beads):**
   "You just updated your conventions — a reminder bead lets you verify
   these changes are landing in practice. It'll record today's date so when
   it surfaces after 30 beads, you can compare beads created after this
   re-onboard against beads before it. Are agents following the updated
   guidance? Did the changes actually improve quality?"

   Then ask: "Want me to create a P3 reminder bead to verify these updated
   conventions are working after 30 more beads?"
   - If yes: `bd create --title="Verify re-onboard conventions are working" --type=chore
     --priority=3 --description="Re-onboarded on [today's date]. After 30+ beads,
     compare bead quality before vs after this date. Check: Are agents following
     the updated description templates? Are the new area labels the right split?
     Are close reasons richer than before? If not, re-run /beads:onboard to
     refine further." --labels=workflow:brainstorm`
   - If no: skip silently.

## Important

- Keep generated content concise and scannable.
- The close outcome template and label enforcement are the two highest-value
  outputs. Get these right.
- Don't over-engineer. The system should take less time to set up than it saves.
- Reference detailed conventions from AGENTS.md rather than duplicating inline.
- Write to AGENTS.md (tool-agnostic), not CLAUDE.md.
- One question at a time. Wait for each answer before asking the next.
