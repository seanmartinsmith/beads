---
date: 2026-03-04
topic: onboard-v3-design-decisions
depends-on: 2026-03-04-onboard-v3-brainstorm.md
---

# `/beads:onboard` v3 - Design Decisions

Resolves the open questions from the [v3 brainstorm](./2026-03-04-onboard-v3-brainstorm.md).
That doc has the empirical research (354 beads analyzed, quality correlation data) and gap
identification. This doc has the concrete design decisions for implementation.

## Decision 1: Three-Way Mode Detection

**Same 6 questions across all modes. Detection sources adapt, interview doesn't fork.**

| Mode | Phase 1 Detection Source | Interview Behavior |
|---|---|---|
| Existing project | Codebase: directories, git log, config files | Auto-detected defaults presented |
| New with planning docs | Planning docs: brainstorms, plans, specs | Extracted defaults from docs |
| New empty | Nothing to detect | Sensible universal defaults |

The interview questions are the same - what changes is where the default values come from.
Phase 1 tries codebase detection first, falls back to doc extraction, falls back to
universal defaults.

**Rationale:** Can't design the perfect new-project flow without dogfooding it. This is the
simplest approach that ships the quality teaching improvements (which benefit ALL modes).
New-project-specific refinements become v4 after real usage data.

## Decision 2: Re-onboard Flow

**Adaptive detection using persisted state file.**

### State file: `.beads/conventions/.onboard-state.yaml`

Stores:
- `version`: which onboard version generated these conventions
- `timestamp`: when onboard last ran
- `answers`: all interview answers as key-value pairs (for pre-filling re-onboard)

Hidden dotfile in conventions directory. Not loaded by agents during normal work -
only read by the onboard command itself on re-run.

### Re-onboard behavior matrix

| State | Detection Signal | Behavior |
|---|---|---|
| Never onboarded | No `.beads/conventions/` | Full interview, all defaults from detection |
| Older version, unmodified | Version stamp < current, files match generated | "v3 adds X, Y, Z. Upgrade?" Extract old answers, regenerate with latest templates |
| Older version, manually customized | Version stamp < current, files differ from generated | "Your files have custom edits. Merge new guidance into existing, or start fresh?" |
| Current version | Version stamp matches | "Conventions are current. Want to review/update answers?" |

**Key insight:** Since onboard is an LLM executing instructions (not compiled code), the
re-onboard flow can read existing convention files, understand customizations, and
intelligently merge new guidance without blowing away hand-tuned content. The state file
provides the structured data; the LLM provides the adaptive intelligence.

## Decision 3: Self-Improving Scaffold

**Milestone bead offered in Phase 4 + health check guidance in reference.md.**

Two mechanisms for two time horizons:

1. **Milestone bead** (immediate trigger): After generating conventions, Phase 4 offers
   to create a bead: "Re-evaluate conventions after 30 beads." Priority 3. When the
   project hits that milestone, it surfaces in `bd ready`. Simple, no automation needed.
   User can decline.

2. **Convention health check** (ongoing trigger): A section in reference.md tells agents:
   "Periodically audit label usage, close reason quality, and description richness.
   Re-run `/beads:onboard` to upgrade conventions based on actual usage patterns."

The milestone bead is a Phase 4 offer, not an interview question. It's meta-information
about the onboard process, not a core convention decision.

## Decision 4: Priority Calibration

**Expand Q1 for project grounding. Auto-generate P0-P4 examples in reference.md.**

No new interview question. Q1 already captures the project description - expand it to also
capture project type and deployment context (production web app, CLI tool, library, personal
project). Phase 1 auto-detects what it can (deployment configs, package exports, etc.).

The priority guidance in reference.md then provides project-type-specific examples:

- Production web app: P0 = "site down, users can't access"
- CLI tool: P0 = "data loss or corruption"
- Library: P0 = "breaking change in published API"
- Personal project: P0 = "blocks all other work"

**Important framing:** Priority matters at **creation time** (assigning the right level),
not at work time (`bd ready` resolves ordering via the dependency graph). The calibration
guides "how to assign correctly" not "what to work on next."

**Quick wins are orthogonal to priority.** They're about complexity/scope, not urgency.
Reference.md includes a note: "When creating beads, consider scope alongside priority.
A self-contained fix in one file is a different bead from a multi-file refactor, even if
they address the same problem. Smaller, focused beads are easier for agents to complete
in a single context window."

## Decision 5: Description Templates

**Section checklists + full example `bd create` commands per type.**

SkillsBench finding: agents don't reliably generate structured content from abstract rules.
They need concrete examples to pattern-match from. Both checklists (for quick reference)
and examples (for pattern-matching) in reference.md.

### Template structure per type

| Type | Sections (include when relevant) |
|---|---|
| **Epic** | Context, Core Issues, Phases [NOW/SOON/LATER], Key Decisions, Related |
| **Bug** | Observed behavior, Expected behavior, Steps to reproduce, Environment |
| **Feature** | Motivation (why), Acceptance criteria, Scope boundary (NOT included) |
| **Task** | Context (why now), Scope, Done-when |
| **Chore** | What + why now (not just "update deps" - why specifically) |

Each type gets one complete `bd create` example using the project's actual prefix and
labels. Descriptions are filled in, not placeholders. The examples demonstrate the
section structure in practice.

**Not rigid templates.** "Include when relevant" guidance, but with enough specificity
that agents know which sections apply to which situations. The chore example specifically
addresses the laziness problem: agents write "update dependencies" when they should write
"update dependencies because CVE-2026-xxxx affects production."

## Decision 6: bd create Fields Guidance

**Concrete triggers + examples + decision tests per field.**

The "use when relevant" framing from the brainstorm is too vague - agents will skip
fields by default. Each field gets a trigger condition (not a judgment call), one filled-in
example, and a decision test.

| Field | Trigger | Decision Test |
|---|---|---|
| `--description` | Always | - |
| `--notes` | You discovered something non-obvious during investigation | "Would a future agent hit this same surprise?" |
| `--design` | You chose between alternatives | "Is there a 'why not the other way?' question someone would ask?" |
| `--acceptance` | Done-when isn't obvious from the title | "Could two agents disagree on whether this is complete?" |
| `--spec-id` | A planning doc is driving this work | "Is there a doc that explains why this bead exists?" |

Each field also gets one example using the project's prefix, showing the field populated
with realistic content (not "put details here" placeholder text).

## Decision 7: Bead Lifecycle (bd update During Work)

**Moment-triggers for when to run `bd update` + anti-patterns showing what to replace.**

The current behavior: agents stuff everything into `bd comments add`. Comments become a
wall of mixed information (design decisions, edge cases, progress notes, handoffs).

The new guidance teaches specific moments during work when `bd update` is the right tool:

| Moment | Action | Instead Of |
|---|---|---|
| Made a design decision | `bd update <id> --design "..."` | Commenting "decided to use X because Y" |
| Found an edge case | `bd update <id> --append-notes "..."` | Commenting "note: also affects Z" |
| Refined what "done" means | `bd update <id> --acceptance "..."` | Commenting "realized we also need to handle Y" |
| Created a related doc | `bd update <id> --spec-id "docs/..."` | Commenting "see docs/plans/foo.md" |
| Scope became clearer | `bd update <id> --add-label concern` | Not adding labels at all |

The pattern: **if the information describes the bead itself, it's a field update. If it
describes what happened while working on the bead, it's a comment.**

## Decision 8: Comments vs Fields Principle

**One-line decision test + concrete right/wrong examples.**

> "Would this information still be useful if all comments were deleted?"
> If yes - it belongs in a structured field, not a comment.

### When comments ARE right

| Use comments for | Why |
|---|---|
| Session handoff notes | "Stopped at step 3, X works, Y needs Z" - ephemeral session context |
| Questions or blockers | "Need clarification on IE11 support" - conversation |
| Progress checkpoints | "Phase 1 complete, starting Phase 2" - timestamped markers |
| External references | "Related: github.com/..." - links outside beads |
| Post-close addenda | "Found edge case after closing, see prefix-xxx" |

### When comments are WRONG

| Instead of this comment... | ...use this field |
|---|---|
| "Decided to use Redis because no sticky sessions" | `bd update <id> --design "Redis pub/sub over WebSockets - deployment has no sticky sessions"` |
| "Note: config loader falls back to env vars in Docker" | `bd update <id> --append-notes "Config loader falls back to env vars when no file exists (intentional for Docker)"` |
| "This is done when Lighthouse score >= 90" | `bd update <id> --acceptance "Lighthouse score >= 90, loads in <2s on 3G throttle"` |

## Decision 9: SkillsBench-Informed Design Principles

Applied throughout all generated content:

1. **Focused sections > comprehensive docs.** Each section of reference.md is internally
   focused on one topic. No encyclopedic walls of text.

2. **Expert-curated examples > abstract rules.** Agents can't reliably generate structured
   content from rules alone. Every guidance section includes a filled-in example.

3. **Specificity > interpretation.** Trigger conditions are concrete ("you chose between
   alternatives"), not judgmental ("when it makes sense"). Agents default to skipping
   vague guidance.

4. **Compact always-loaded + detailed on-demand.** AGENTS.md section stays 60-80 lines
   (loaded every session). Reference.md has the depth (loaded when creating/updating beads).

## Decision 10: Context Window Awareness

**Brief guidance in reference.md, not a full orchestration system.**

Reference.md includes a note in the session discipline section:

> "Consider context window scope when claiming work. A self-contained bug fix in one file
> is completable in a single session. A refactor touching 15 files may need to be broken
> into sub-beads or may benefit from worktrees/subagents. If a bead's scope exceeds what
> you can complete in your current context window, break it down before starting."

This plants the seed without trying to solve the larger orchestration problem, which is
out of scope for onboard.

## Interview Changes Summary

**v2 → v3: Still 6 questions, but richer.**

| Question | v2 | v3 Change |
|---|---|---|
| Q1: Project description | 3 AI-generated descriptions | Expanded to include project grounding (type, deployment context, audience). Feeds priority calibration. |
| Q2: Area labels | Detected areas + confirm | No change |
| Q3: Commit format | 4 format options | No change |
| Q4: Code change policy | 3 options | No change |
| Q5: Beads + Tasks | 3 options | No change |
| Q6: Close outcome template | Default template + customize | No change |

**Phase 0 changes:**
- Three-way mode detection (empty, planned, existing)
- Auto-`bd init` when `.beads/` missing (with prefix customization offer)
- Read `.onboard-state.yaml` for re-onboard detection

**Phase 1 changes:**
- Fall back to planning docs when no codebase to analyze
- Fall back to universal defaults when no docs either
- Auto-detect project type for priority calibration

**Phase 3 changes (generated content):**
- reference.md gets 5 new sections: description templates, priority calibration,
  bd create fields, bead lifecycle, comments vs fields
- AGENTS.md section gets brief pointers to new reference.md sections
- `.onboard-state.yaml` written with version + answers

**Phase 4 changes:**
- Milestone bead offered: "Create a reminder to revisit conventions after 30 beads?"
- Convention health check guidance included in reference.md

## What's NOT in v3

- Progressive disclosure (quick vs advanced setup) - unnecessary complexity, skip it
- New interview questions - keep at 6, enrich the output instead
- Full orchestration/parallelization guidance - out of scope, plant seeds only
- Cross-agent portability - already tool-agnostic output, no additional work needed
- Time estimates in beads - agents don't exist in time, measure in complexity/scope

## Open Questions (for dogfooding)

1. Will three-way mode detection work well for truly empty repos, or will the universal
   defaults feel too generic? Needs real usage to evaluate.
2. Will the expanded reference.md be too long? Start thorough, whittle down based on
   which sections agents actually follow vs ignore.
3. Will the `.onboard-state.yaml` approach work smoothly for re-onboard, or will edge
   cases (corrupted state, manually deleted files) need special handling?
4. Does the milestone bead actually get worked when it surfaces in `bd ready`, or do
   users ignore/close it without re-running onboard?
