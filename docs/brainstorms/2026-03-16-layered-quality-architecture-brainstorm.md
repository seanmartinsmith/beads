---
date: 2026-03-16
topic: layered-quality-architecture
predecessor: 2026-03-09-quality-delivery-architecture-brainstorm.md
status: ready-for-plan
---

# Layered Quality Architecture for Beads

## Context

Brainstorm session spanning: onboard v3 code review, tpane dogfooding (85+ beads),
comprehensive CLI audit (70+ commands), template/ecosystem research, community issue
analysis (#2140, #2611, #2612), and industry best practices research.

This supersedes the March 9 "quality delivery architecture" brainstorm which paused
with open questions about persistence vs consistency. Those questions are now answered.

## The Core Problem

Beads CLI is feature-rich (70+ commands, molecules, gates, federation, query language,
graph visualization, lifecycle management) but agents only see ~8 commands via `bd prime`.
The result: 90% of features are invisible, agents produce low-quality beads, and users
manually audit what the tool could enforce.

Evidence:
- tpane dogfooding: 24% of open beads had thin descriptions missing required structure
- 26% of close reasons across 704 beads were terse (<20 chars), mostly just "Closed"
- `bd setup claude` template literally shows `--reason "Completed"` as an example
- 91 failed `bd comment` attempts, 166 failed `--comment` attempts (issue #2611/#2612)
- `--acceptance`, `--design`, `--notes`, `--validate` exist but are never mentioned
- `bd lint`, `bd stale`, `bd orphans`, `bd find-duplicates` exist but are invisible

## The Architecture: Three Layers

### Layer 1: Beads Core (Go CLI - upstream PR, benefits all users)

**Prime output improvements** (immediate effect via hooks, survives compaction):
- Close reason template with structured format
- Field usage teaching (`--acceptance`, `--design`, `--notes`, `--validate`)
- Description requirements by type (bug/feature/task/epic/chore)
- Lifecycle commands (`bd defer`, `bd supersede`, `bd duplicate`, `bd reopen`)
- Hygiene commands (`bd lint`, `bd stale`, `bd orphans`, `bd find-duplicates`)
- `bd close --suggest-next` / `--claim-next`
- Enrichment during work (`bd update --design`, `--append-notes`, `--acceptance`)
- Fix the `--reason "Completed"` anti-pattern

**Template improvements** (flows to all 12 integrations via existing system):
- Same quality content in `beads-section.md` (full profile)
- Minimal profile points to prime (already correct pattern)
- Universal template for Cursor/Windsurf/Cody/Kilo Code
- Hash staleness detection handles refresh automatically

**`bd health` command** (unified quality audit):
- Runs: lint + stale + orphans + find-duplicates + dep cycles + label check
- Single report output
- Convention drift detection (labels in use but not in taxonomy)
- Replaces the need for manual 6-command audit runs

### Layer 2: Project Onboard (CLI command - tool-agnostic)

**`bd onboard` as a Go CLI command** (any LLM can run via bash):
- Interactive mode (stdin/stdout prompts, like `npm init`)
- Flag mode for non-interactive use / LLM invocation:
  ```
  bd onboard --project-type=cli \
    --areas="cli,ui,storage,config" \
    --commit-format="type(scope): description" \
    --close-template=default
  ```
- Generates: `.beads/conventions/labels.md`, `.beads/conventions/reference.md`
- Updates AGENTS.md within markers (using existing render.go system)
- Writes `.beads/conventions/.onboard-state.yaml` for re-onboard
- Re-onboard detection and adaptive upgrade

**What it customizes (project-specific):**
- Area labels for this specific project
- Project description / project type detection
- Priority calibration (what P0 means for THIS project type)
- Commit format choice (open source vs personal/team)
- Project-specific examples using real prefix and labels

### Layer 3: Claude Code Integration (plugin-specific)

**`/beads:onboard` skill** (wraps Layer 2):
- Uses AskUserQuestion for interview UX
- Collects answers, passes as flags to `bd onboard` CLI
- Adds Tasks + beads integration teaching
- Thin wrapper - all logic in Layer 2

**Plugin enhancements:**
- Updated beads plugin hooks if needed
- Session workflow (plan mode → beads → tasks → close)

## Key Design Decisions

### D1: Close Reason Format (research-informed)

Current 6-field template (Summary/Change/Files/Verify/Risk/Notes) is good but
research says the highest-value missing field is **Learnings/Discovery** - what
was learned that wasn't known at start. This is institutional knowledge gold.

Revised format:
```
Summary: <one sentence - what was done>
Change: <what changed - scope/blast radius>
Files: <paths touched>
Verify: <how to confirm it still works>
Discovery: <what was learned that wasn't known at start>
Risk: <if any - what could break>
```

"Discovery" replaces "Notes" (vague) with a specific prompt that captures the
institutional knowledge. "Risk" moves to optional/last since it's not always
applicable.

Minimum floor remains: Summary + Change + Files (always required).
Discovery/Verify/Risk added for non-trivial changes.

Resolution categories already handled by separate close paths (`bd duplicate`,
`bd supersede`, `bd defer`) - no new enum needed.

### D2: `--design` Field as Mini-ADR

Research on Architecture Decision Records says format is:
Context → Decision → Consequences.

Teach `--design` as: "Why X over Y? Because Z. Tradeoff: W."

Example: `--design="Chose SQLite over PostgreSQL - single-binary deployment
requirement. Tradeoff: no concurrent writers, acceptable for single-user CLI."`

### D3: Spec-Driven Development Mapping

GitHub Spec Kit validates: Specify → Plan → Tasks → Implement.

Beads field mapping:
- `--description` = Specify (motivation, problem, scope)
- `--design` = Plan (approach, decisions, tradeoffs)
- `--acceptance` = Success criteria (testable done-when)
- `--spec-id` = Link to spec doc
- Claude Code Tasks = Implementation steps
- `--reason` = Outcome documentation

This isn't a new feature - it's teaching agents to use existing fields in a
workflow pattern. The fields already exist. The workflow already works. Agents
just don't know to connect them.

### D4: Session Handoff via Plan Mode + Beads

Claude Code plan mode (shift+tab) creates a plan doc in `.claude/` that doesn't
persist in the session but provides context recovery. Combined with beads:

1. Working on a bead → enriching with `bd update --design`, `--append-notes`
2. Need to hand off → `bd comments add <id> "session handoff: ..."` for
   cross-session context
3. Use plan mode for within-session context management (survives compaction)
4. Close bead with structured reason that captures outcome

Beads comments for cross-session handoff. Plan mode for within-session context.
Tasks for within-session execution steps. Three complementary mechanisms.

### D5: `bd health` Replaces Audit Bead

The tpane dogfooding proved the "create a P4 audit bead after 30 beads" approach
doesn't work - P4 never surfaces via `bd ready`. Instead:

`bd health` as a unified command that runs:
- `bd lint` (missing sections by type)
- `bd stale` (forgotten beads)
- `bd orphans` (referenced in commits, still open)
- `bd find-duplicates` (overlapping work)
- `bd dep cycles` (circular dependencies)
- Label taxonomy check (labels in use but not documented)
- Convention drift detection

This could also be surfaced in `bd doctor` output or as a periodic hook.

### D6: Memory System Complementarity (beads + Claude Code)

Current beads plugin instructions say "do NOT use MEMORY.md files" and push all
persistent knowledge into `bd remember`. This is too aggressive - it conflates
two different memory systems that serve different purposes.

**Beads memories** (`bd remember`):
- Project knowledge that ANY agent/tool needs
- Stored in Dolt, syncs via `bd dolt push/pull`
- Tool-agnostic (Cursor, Codex, Gemini, Claude all get them via `bd prime`)
- Shared across team members
- Examples: "always run tests with -race flag", "auth uses JWT not sessions"

**Claude Code memories** (`~/.claude/projects/`):
- Personal preferences and behavioral corrections
- Claude Code specific, don't sync
- Types: user profile, feedback, project context, references
- Examples: "I prefer terse responses", "don't create MEMORY.md files"

**The distinction**: beads memories are about THE PROJECT (anyone working on it
needs these facts). Claude memories are about THE PERSON (how this specific user
wants to interact with Claude).

**At-scale consideration**: when running many parallel agents executing tasks,
personal preference memories may not be useful - agents just need project facts
and execution context. But for interactive sessions (brainstorming, planning,
code review), personal preferences matter a lot. The memory system needs to
serve both modes.

**What needs to change**:
- Layer 1 (prime output): clarify when to use `bd remember` (project knowledge)
  vs when knowledge belongs in tool-specific memory (personal preferences)
- Layer 3 (Claude Code): explain the complementary relationship, stop telling
  agents to avoid Claude memories entirely
- The beads plugin SessionStart hook should not override Claude's memory system,
  it should complement it

### D7: Molecules/Formulas as Development SOP Templates

Beads already has molecules (workflow instances) and formulas (workflow templates).
These are invisible to most users but are exactly the right mechanism for encoding
development SOPs:

- TDD formula: red → green → refactor → verify
- Bug-fix formula: reproduce → diagnose → fix → verify → close
- Feature formula: spec → plan → implement → test → close
- Spec-driven formula: specify → plan → tasks → implement → review

This is out of scope for this PR but worth noting as the natural evolution. The
infrastructure exists. What's missing is pre-built formulas and teaching.

### D8: Validation and Lint System (Current State + Gaps)

**`bd create --validate` is implemented and functional.** It checks for required
heading text in the description via case-insensitive substring matching:

| Type | Required headings in description |
|---|---|
| bug | "Steps to Reproduce", "Acceptance Criteria" |
| feature | "Acceptance Criteria" |
| task | "Acceptance Criteria" |
| epic | "Success Criteria" |
| decision | "Decision", "Rationale", "Alternatives Considered" |
| chore | (none) |

Source: `internal/types/types.go:RequiredSections()` and
`internal/validation/template.go:ValidateTemplate()`

**`bd lint` (LintIssue) is smarter than `--validate`:**
- Also checks the `AcceptanceCriteria` field - if `--acceptance` is populated,
  it satisfies the "Acceptance Criteria" requirement WITHOUT needing the heading
  in the description text (GH#2468)
- Source: `internal/validation/template.go:LintIssue()`
- This means `bd lint` correctly handles "use `--acceptance` as its own field"
  but `bd create --validate` does NOT

**Config-driven auto-validation exists:**
- `bd config set validation.on-create error` - auto-validate every create
- `bd config set validation.on-create warn` - warn but don't block
- Currently unset (off) by default
- Onboard should set this to `warn` as part of project setup

**Known gaps:**

1. **`--validate` doesn't check `--acceptance` field** - only checks description
   text. `bd lint` has this fix (GH#2468) but `--validate` doesn't. Either
   `--validate` needs the same `LintIssue`-style field checking, or we document
   that `bd lint` is the proper quality gate and `--validate` is the creation
   floor.

2. **`--dry-run` bypasses `--validate`** - `bd create --dry-run --validate`
   shows what would be created but doesn't run validation. Confusing - you'd
   expect dry-run to preview validation results too.

3. **Close reason has zero validation** - No `--validate` on `bd close`. No
   `validation.on-close` config. `bd close --reason "Closed"` succeeds silently.
   This is the biggest quality gap: 26% of close reasons are terse (<20 chars).

4. **Section requirements are minimal** - Only checks for specific headings
   (Acceptance Criteria, Steps to Reproduce, etc.). Doesn't check for the richer
   structure onboard teaches (Motivation, Scope boundary, etc.). The CLI enforces
   the floor; conventions teach the ceiling. This is probably correct - the floor
   should be low enough that it doesn't block legitimate workflows.

5. **No label validation** - `bd create --validate` doesn't check that an
   `area:*` label is included. Label taxonomy drift (tpane: 15+ undocumented
   labels) has no enforcement point.

**Implications for our work:**

- Layer 1: Teach agents about `validation.on-create warn` config in prime/templates
- Layer 2: `bd onboard` should run `bd config set validation.on-create warn`
- Layer 1 (Go PR): Consider adding `--acceptance` field checking to `--validate`
  (parity with `bd lint`)
- Layer 1 (Go PR): Consider adding `validation.on-close warn` config for close
  reason minimum floor checking
- Layer 1: Surface `bd lint` in prime output as periodic quality gate

### D9: `bd human` for Multi-Agent Human-in-the-Loop

`bd human` is a built-in system for surfacing beads that need human input:
- `bd human list` - shows all beads with 'human' label
- `bd human respond <id>` - add comment and close
- `bd human dismiss <id>` - permanent dismiss
- `bd human stats` - summary statistics

This is directly relevant to the multi-agent parallel workflow. When an agent
hits a decision that needs user input, it should create a bead with the `human`
label. The user runs `bd human list` to see everything waiting across the project.

**Dogfooding finding**: The user has been reinventing this with custom labels and
"blocking brainstorm" beads in their portfolio project. The `bd human` system
should have been surfaced from the beginning.

**What needs to change**:
- Layer 1 (prime/templates): Surface `bd human` in workflow guidance
- Layer 2 (onboard): Include `human` in the label taxonomy explanation
- Teaching: when an agent needs a decision, use `human` label, not custom labels

### D10: PRIME.md Override as Convention Injection Point

`bd prime` has a 3-tier override mechanism:
1. `.beads/PRIME.md` (project-local, clone-specific)
2. `.beads/PRIME.md` at redirected location (shared across clones)
3. `~/.config/beads/PRIME.md` (global user config)
4. Fall back to auto-generated content

If `.beads/PRIME.md` exists, it REPLACES the default prime output entirely.
`bd prime --export` dumps the default content for customization.

**Implication**: onboard could generate a `.beads/PRIME.md` that includes the
project's creation checklist, close reason template, and field usage guidance.
This survives compaction (re-injected at PreCompact), is project-specific, and
doesn't require any Go code changes.

**Concern**: must be mindful of context bloat. Prime output is injected at every
SessionStart and PreCompact. Keep it compact. The current default is ~80 lines.
A PRIME.md override should stay under ~120 lines to avoid eating context window.

### D11: `bd preflight` Before PRs

`bd preflight --check` runs automated pre-PR checks: tests, lint, formatting,
version mismatches, stale JSONL, nix vendorHash. Invisible to agents currently.

Should be surfaced in prime/templates as part of the "Landing the Plane" protocol.

### D12: Template System Safety

Research confirms modifications are safe:
- Template body content changes only affect target integrations
- Hash staleness detection handles refresh
- Each integration is isolated (12 tools, 4 patterns)
- Prime output changes take effect immediately via hooks
- No breaking changes required for existing setups

Profile system (full vs minimal) is well-designed:
- Full: Codex, Factory, Mux, OpenCode (no hooks, need everything inline)
- Minimal: Claude, Gemini (hooks deliver prime, AGENTS.md is a pointer)
- File-based: Cursor, Windsurf, Cody, Kilo Code (separate template)
- Multi-file: Aider, Junie (custom content)

## Community Alignment

- Issue #2140: Per-agent setup - exactly what we're building
- Issue #2611: Command discoverability (91 failed attempts quantified)
- Issue #2612: Flag discoverability (166 failed attempts quantified)
- PR #2638/#2639: Alias fixes for common failures
- No community requests for close reason templates specifically, but the
  discoverability data proves the teaching gap

## What's NOT in Scope

- New molecules/formulas (D6 noted as future work)
- Cross-tool portability of the skill itself (Layer 3 is Claude Code only)
- External tracker sync improvements (Linear, GitLab, GitHub)
- Dolt migration pain points (#2573, #2489)
- Self-upgrade command (#949)
- Federation features

## Dogfooding Data Points

### tpane project (85+ beads, March 2026)
- 24% of open beads had thin descriptions missing required structure
- 15+ beads used undocumented `scope:*` labels (taxonomy drift)
- P0 brainstorm bead, P1 design beads (priority miscalibration)
- 9 descriptions upgraded manually (should have been caught at creation)
- `--acceptance` field not used at all (stuffed into description)
- `--validate` flag unknown to agents
- Audit bead (P4) never surfaced - sat 10 days and 50+ beads past trigger

### beads repo (704 beads, Dec 2025 - March 2026)
- 26% of close reasons terse (<20 chars)
- 77% of multi-step beads closed with terse reasons
- Quality jumped from 26% to 82% "rich" when templates added to instructions
- Structured descriptions went from 14% to 36% with guidance
- Reduction audit (884 → 284 line instructions) did NOT hurt quality

### Portfolio research (354 beads, analyzed)
- Close reason format was single biggest quality improvement
- Examples > abstract rules (SkillsBench confirmed)
- Compact always-loaded + detailed on-demand is the right pattern

## Open Questions for Planning

1. How much of `bd onboard` is Go code vs delegating to the LLM?
   - Pure Go: fully deterministic, works everywhere, but rigid
   - LLM-assisted: flexible, adaptive, but requires an LLM
   - Hybrid: Go for file generation, LLM for interview/analysis

2. Should `bd health` be a new command or integrated into `bd doctor`?
   - New command: clear separation of concerns
   - Doctor integration: one place for all health checks
   - Could be `bd doctor --quality` or `bd doctor --conventions`

3. How do we test this without shipping broken things?
   - Run on tpane, beads itself, and at least one more project
   - Compare bead quality before/after template changes
   - Verify `bd onboard` generates correct conventions
   - Verify no regressions in existing integrations

4. PR strategy: single PR or split?
   - Single: complete story, easier to review holistically
   - Split: Layer 1 first (templates), then Layer 2 (onboard CLI)
   - Recommendation: discuss with maintainer

5. Should `--validate` be fixed to check `--acceptance` field (parity with
   `bd lint`)? Or is the current behavior intentional (description-only check)?
   This affects whether we teach agents to use `--validate` or `bd lint` as
   the quality gate.

6. Should we propose `validation.on-close` config for close reason minimum
   floor? This is the biggest quality gap but also the most likely to disrupt
   existing workflows (agents closing 20 beads in parallel don't want to be
   blocked). Warn mode (not error) is probably the right default.

7. PRIME.md override strategy: should onboard generate a `.beads/PRIME.md`
   that augments/replaces default prime output? Pro: project-specific conventions
   survive compaction automatically. Con: replaces ALL default prime content,
   so onboard would need to include the essential workflow bits too. Could get
   stale if beads updates its default prime format.

8. How should `bd human` be integrated into the workflow teaching? It's a
   label-based system (`human` label) with dedicated commands. Need to decide
   if it becomes part of the standard label taxonomy or stays optional.

## Next Step

Write implementation plan: `docs/plans/2026-03-16-layered-quality-architecture-plan.md`
