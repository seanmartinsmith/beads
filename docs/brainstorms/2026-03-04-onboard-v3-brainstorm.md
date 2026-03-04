---
date: 2026-03-04
topic: onboard-v3-new-project-and-bead-quality
---

# `/beads:onboard` v3 - New Projects + Bead Quality

## Context

v2 of onboard (implemented, on `feature/onboard-command` branch) handles existing projects
well: analyzes the codebase, runs a 6-question interview, generates convention files. But
testing on a brand new project (cctui) revealed two major gaps:

1. **New projects have nothing to analyze** - no directories, no git history, no config files
2. **Onboard doesn't teach bead creation quality** - it teaches labels and close templates,
   but not what makes a good bead vs a bad one

These aren't edge cases. Setting conventions at project birth is the highest-value moment,
and bead quality is the whole point of having conventions.

## What We're Building

Expand onboard to handle three project states (new empty, new with planning docs, existing)
and teach agents the *practice* of using beads well, not just the structural conventions.

## Empirical Research: Quality Correlation Analysis

Analyzed 354 beads from the portfolio project (Dec 2025 - Feb 2026) cross-referenced against
30+ commits to CLAUDE.md instruction files. Data source: `~/projects/portfolio`.

### Bead Quality by Month

| Period | CLAUDE.md state | Beads | Avg close reason | Rich close (100+) | Structured desc |
|---|---|---|---|---|---|
| Dec 2025 | 382 lines, basic beads section | 108 | 131 chars | 26% | 14% |
| Jan 2026 | Grew to 884 lines (close format, Tasks, bead IDs, epics) | 202 | 392 chars | 82% | 32% |
| Feb 2026 | Reduced to 284 lines (extracted non-beads to agent_docs) | 44 | 557 chars | 100% | 36% |

### Key Inflection Points

1. **Jan 8** - Added priority semantics, label taxonomy, creation examples. Foundation.
2. **Jan 23** - Added close reason format with double newlines + Notes field. **Close reason
   quality jumped from 26% to 82% rich.** Single biggest quality improvement.
3. **Jan 26** - Added `area:epic` label, required bead IDs in commits, bv smart commands.
   **Structured descriptions jumped from 14% to 32%.**
4. **Jan 30** - Added Session Execution with Tasks workflow. Beads-as-persistent-memory
   framing.
5. **Feb 2** - Reduction audit. Extracted non-beads content to agent_docs. All beads
   instructions kept intact. **Quality continued improving** (36% structured, 100% rich).

### What Actually Moved the Needle (ranked by impact)

1. **Close reason format with good/bad examples** - the single biggest quality jump
2. **Required bead IDs in commits** - connected beads to code, changed agent mental model
3. **Area:epic label + structured epic patterns** - gave agents patterns for rich beads
4. **Tasks-beads integration framing** - "persistent memory vs working memory" mental model
5. **Priority semantics with concrete examples** - eliminated guessing

### Verdict on Reduction Audit

The Feb reduction **did not hurt quality** - it improved. Kept all beads instructions,
moved CSS/logging/CASS stuff to separate docs. Signal-to-noise ratio improved. This
validates the layered approach: compact AGENTS.md section + detailed convention files.

## Discovery: The Bead Quality Gap

Comparing beads from the portfolio project (300+ beads over 3 months):

**Good beads have:**
- Structured description sections (Core Issues, Enhancements, Phases, Key Decisions, Context)
- Clear rationale for why the work matters
- Linked related beads and docs
- Phased plans with status markers ([NOW], [SOON], [LATER])
- Rich close_reason documenting what actually happened
- Comments tracking session progress
- Proper epic structure with parent/child relationships

**Bad beads have:**
- One-line descriptions with no context
- No linked related work
- Close reasons like "Done" or "Fixed it"
- Missing labels or wrong priority
- No session notes

## Key Design Decisions

### 1. Three-way mode detection (Phase 0)

| Mode | Signal | Behavior |
|---|---|---|
| New project (empty) | No `.beads/`, no/minimal code | Guide toward brainstorm/plan workflow, offer scaffold |
| New project (planned) | No `.beads/`, has planning docs | Auto-init, extract areas/description from docs |
| Existing project | `.beads/` exists | Current flow - analyze codebase |

**bd init integration:** Auto-run `bd init` when `.beads/` is missing. Surface what
happened: "Initialized beads with prefix `cctui`". Offer to customize prefix and
mention `--team` mode exists so users aren't unaware of options.

### 2. New project input sources (Phase 1)

For new projects, planning docs replace codebase analysis:

| Input | Existing project | New with docs | New empty |
|---|---|---|---|
| Description | AI-generated from codebase | Extract from brainstorm/plan | Manual entry |
| Area labels | Detected from directories | Extract from file structure in plan | Manual entry or defer |
| Tech stack | Detected from package files | Extract from plan | Manual entry |
| Commit format | Detected from git log | Recommend type(scope): description | Recommend type(scope): description |

**Doc detection heuristic:** Check `docs/brainstorms/`, `docs/plans/`, `docs/`, root for
`.md` files. Read them, extract project description, planned file structure, tech stack.

**Empty repos:** Don't just scaffold beads - suggest the workflow. Beads aren't useful if
you don't know what you're building. Recommend: brainstorm first, plan second, then onboard.
Onboard can help track those activities, but shouldn't replace them.

### 3. Bead creation quality teaching

**Currently teaches:**
- Label taxonomy (area, platform, concern, workflow)
- Close outcome template (Summary, Change, Files, Verify, Risk, Notes)
- Commit format
- Code change policy
- Beads + tasks integration

**Should also teach (validated by portfolio correlation data):**
- Description structure templates by issue type (not just epics)
- Priority calibration with project-specific examples
- bd create fields beyond --description (--notes, --design, --acceptance, --spec-id)
- Bead lifecycle: using bd update to enrich beads during work, not just at creation/close
- Comments vs fields: when to use each (structured metadata vs conversational context)
- Dependency rules - blocking vs relate vs parent-child
- Session discipline - read before work, check abandoned work, search before creating
- Good vs bad bead examples using the project's actual prefix
- When to use epics and how to structure them

### 4. Description structure by issue type

Structure guidance should apply to ALL types, not just epics. Type-appropriate templates:

| Type | Key sections | Why |
|---|---|---|
| **Epic** | Context, Core Issues, Phases [NOW/SOON/LATER], Key Decisions, Related | Planning artifacts, need the most structure |
| **Bug** | Observed behavior, Expected behavior, Steps to reproduce, Environment | Reproducibility - agents skip "expected" making fixes ambiguous |
| **Feature** | Motivation (why), Acceptance criteria, Scope boundary (NOT included) | Scope control - prevents gold-plating or under-delivery |
| **Task** | Context (why now), Scope, Done-when | Clear completion signal |
| **Chore** | What + why now (not just "update deps" but "because X vulnerability") | Chores get lazy descriptions because they seem routine |

These are "include when relevant" guidance, not rigid templates. The data shows agents
don't structure descriptions unless explicitly told which sections to use.

### 5. Priority calibration

What P0-P4 means varies by project. Onboard should generate project-specific examples:

- Production web app: P0 = "site down, users can't access", P2 = "visual glitch on one page"
- CLI tool: P0 = "data loss or corruption", P2 = "confusing error message"
- Personal project: P0 = "blocks all other work", P2 = "annoying but workaround exists"
- Library/SDK: P0 = "breaking change in published API", P2 = "missing docs for edge case"

**Implementation:** Either a new interview question ("What type of project is this?") or
auto-detect from project signals (has deployment config = production app, is a library =
look for go.mod/package.json with exports, etc.) and generate calibrated examples.

### 6. bd create fields (in scope for v3)

`bd create` supports fields beyond `--description` that agents should use when relevant.
The key framing: more fields = more searchable context, but only when meaningful. An
empty `--notes ""` or `--design "standard approach"` is worse than omitting it.

| Field | Use when | Skip when |
|---|---|---|
| `--description` | Always. The "what and why" | Never skip |
| `--notes` | Gotchas, edge cases, things future agents should know | Simple fixes with no surprises |
| `--design` | You made a non-obvious architectural choice | Implementation is straightforward |
| `--acceptance` | "Done" criteria aren't obvious from the title | "Fix typo" doesn't need acceptance criteria |
| `--spec-id` | A plan/spec doc is driving this work | Ad hoc work with no spec |

### 7. Bead lifecycle: bd update during work + comments

Currently agents stuff everything into `bd comments add` during the work lifecycle.
Comments are being overloaded for things that have dedicated fields. This creates
noise and makes beads harder to parse - session notes, design decisions, and progress
updates all mixed together in a comment thread.

**bd update fields for mid-work enrichment:**

| Field | Use for | Instead of |
|---|---|---|
| `bd update <id> --design "..."` | Design decisions made during implementation | Commenting "decided to use X approach because..." |
| `bd update <id> --acceptance "..."` | Refined acceptance criteria discovered during work | Commenting "realized we also need to handle Y" |
| `bd update <id> --append-notes "..."` | Accumulated gotchas, edge cases, caveats | Commenting "note: this also affects Z" |
| `bd update <id> --spec-id "docs/..."` | Linking to spec/plan docs created during work | Commenting "see docs/plans/foo.md" |
| `bd update <id> --add-label concern` | Adding labels as scope becomes clearer | Not doing it at all (labels stuck at creation) |

**When comments ARE the right tool:**

| Use comments for | Why |
|---|---|
| Session handoff notes | "Stopped at step 3, X is working, Y still needs Z" - ephemeral session context |
| Questions or blockers | "Need clarification on whether we support IE11" - conversation, not metadata |
| External references | "Related discussion: github.com/..." - links to things outside beads |
| Progress checkpoints | "Phase 1 complete, starting Phase 2" - timestamped progress markers |
| Post-close addenda | "Found an additional edge case after closing, see portfolio-xyz" |

**The principle:** Structured fields are for searchable, parseable metadata that describes
the bead itself. Comments are for timestamped human/agent communication about the bead.
If you'd want to find it via `bd search` or filter by it, it's a field. If it's
conversational context, it's a comment.

**This is a significant workflow improvement.** Teaching agents to use `bd update` fields
during work (not just at creation and close) means beads accumulate structured context
throughout their lifecycle instead of having a sparse description and a wall of comments.

### 8. Context efficiency

The portfolio data validates the layered approach:
- 884-line CLAUDE.md → 284 lines (67% reduction) with NO quality degradation
- Non-beads content moved to separate docs agents load on demand
- All beads instructions stayed in the main file

**Proposed structure:**
- AGENTS.md section: ~60-80 lines, scannable rules and quick reference (always loaded)
- `.beads/conventions/reference.md`: full depth including description templates, dependency
  rules, session discipline, priority calibration, good/bad examples
- `.beads/conventions/labels.md`: label taxonomy (loaded on demand)

### 9. Cross-agent compatibility

AGENTS.md is the universal standard (Linux Foundation, 60k+ repos). Every major coding
agent reads it. The onboard output is already tool-agnostic. The interactive process is
Claude Code-specific but could port to Codex/Amp via `.agents/skills/` format later.

Not in scope for this iteration.

### 10. Code change policy

Remains a user preference (Q4 in the interview), not a default recommendation.

## Recommendations

Based on the correlation analysis, in priority order:

1. **Description structure templates in reference.md** - type-specific section guidance
   with examples. This addresses the biggest quality gap (14% → 36% structured when
   agents had patterns). Expand beyond epics to bugs, features, tasks, chores.

2. **Priority calibration in the interview** - add a question or sub-question that
   generates project-specific priority examples. "What does critical mean for this
   project?" with auto-detected defaults.

3. **Example beads in conventions** - 3-4 complete `bd create` commands showing good beads
   for each type. The portfolio data shows examples were the most effective teaching tool.

4. **Keep AGENTS.md compact with pointers** - 60-80 lines max. Put depth in
   `.beads/conventions/reference.md`. The reduction audit proved this works.

5. **Bead lifecycle guidance in reference.md** - teach agents to use `bd update` fields
   (--design, --acceptance, --append-notes) during work, not just at creation and close.
   Define when comments vs fields are appropriate. This is the largest underutilized
   gap in current workflows - beads accumulate walls of comments instead of structured,
   searchable metadata.

6. **bd create fields guidance** - teach "use when relevant" framing for --notes,
   --design, --acceptance, --spec-id. Not rigid forms, but awareness that these fields
   exist and when they add value.

7. **Dependency rules in reference.md** - blocking vs relate vs parent-child. Commonly
   misused by agents, was explicitly addressed in the portfolio's evolved instructions.

8. **New project flow** - three-way mode detection, auto-init, planning doc extraction.
   Important for adoption but lower priority than quality teaching (existing projects
   benefit from quality improvements too).

## Open Questions

1. **Interview expansion** - v2 has 6 questions. Adding priority calibration + project
   type detection could make it 7-8. Is that too many? Or can we fold it into existing
   questions?
2. **Reference.md size** - the portfolio's beads-reference.md is 280 lines. How much of
   that should onboard generate? Could start smaller (150 lines) and let users expand.
3. **bv integration** - the portfolio heavily uses bv for smart triage. Should onboard
   mention bv? It's a separate tool and not everyone has it. Probably mention it as
   optional in reference.md without making it a dependency.
4. **Validation of description templates** - will agents actually follow "include these
   sections" guidance, or do we need to show them filled-in examples? The portfolio data
   suggests examples work better than rules alone.
5. **Progressive disclosure** - should onboard offer a "quick setup" (labels + close
   template only, current v2) and an "advanced setup" (full quality teaching)?

## Portfolio Project Reference

The portfolio project at `~/projects/portfolio` has 354 beads (300 closed, 51 open)
evolved over 3 months (Dec 2025 - Feb 2026).

**Key files:**
- `CLAUDE.md` (284 lines) - battle-tested instruction set, ~30 commits of evolution
- `docs/agent_docs/beads-reference.md` (280 lines) - full beads workflow reference
- `docs/agent_docs/labels.md` (97 lines) - label taxonomy
- `.beads/issues.jsonl` (354 beads) - the actual bead data

**Key commits in instruction evolution:**
- `fd19beb` (Dec 9) - first beads instructions added
- `b2184d9` (Jan 8) - major expansion: priorities, labels, commands, examples
- `b81f26b` / `6782d1a` / `478af93` (Jan 23) - close reason format (biggest quality jump)
- `d433fb6` (Jan 26) - area:epic, bead IDs in commits
- `30ae4e1` (Jan 30) - Tasks workflow integration
- `6b26c05` / `5fc60be` (Feb 2) - reduction audit (extracted non-beads to agent_docs)

## Next Steps

1. Design the expanded reference.md template with quality teaching
2. Add description structure templates by issue type
3. Add priority calibration to the interview
4. Update Phase 0/1 for new project flow
5. Implementation plan via `/ce:plan`
