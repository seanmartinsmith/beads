# Brainstorm: Onboard Description Quality + Close Reason Improvements

Date: 2026-03-06
Status: Design approved
Related: `claude-plugin/commands/onboard.md` (lines 385-484)

## What We're Building

Improvements to the `/beads:onboard` command's description quality teaching. The
onboard currently teaches description *structure* (what sections to include per type)
but not description *quality* (why each section matters, what good vs bad looks like,
what the minimum bar is). Close reasons have the same gap - the template exists but
nothing enforces or teaches quality.

The goal: after running onboard, agents produce beads that are actionable, searchable,
and serve as semantic anchors for project history. Nothing gets "lost" as the codebase
grows because beads are written with future discoverability in mind.

## Why This Approach

### Evidence from the beads repo (704 issues)

- 97% have descriptions, 91% are >100 chars, 56% are >500 chars - creation side is
  already rich when agents have good descriptions to pattern-match from
- 26% of close reasons are terse (<20 chars). 106 issues closed with just "Closed"
  or similar. These are lost project memory.
- 43 multi-step beads found, 33 of them closed with terse reasons. Complex work with
  rich descriptions that evaporates at close time.
- Best descriptions (bd-lfak, bd-1rh) have Vision/Problem/Phases/Success Criteria.
  Worst are one-liners that don't even identify the file being changed.

### Semantic anchoring principle

Beads descriptions and close reasons function like `debug.trace()` semantic anchors
in code - stable reference points that survive refactoring and help future agents
(or humans) find the right context. The test: "If I `bd search` for this problem in
6 months, will this bead surface?"

Note: Code-level semantic traces (debug.trace patterns from portfolio) are a
complementary practice worth adopting separately. Out of scope for this PR but worth
noting as a future onboard recommendation.

## Key Decisions

### 1. Type-driven templates with quality floor

Each issue type gets a template with required and optional sections. Required sections
are the floor - every bead of that type must have them. Optional sections add depth
when applicable.

Decision tests from the existing `--notes`/`--design` guidance ("Would a future agent
need this?") serve as the WHY behind each section, not the primary teaching tool.
Agents pattern-match from examples better than they internalize principles.

### 2. Close reason minimum: Summary + Change + Files

The full 6-field template (Summary/Change/Files/Verify/Risk/Notes) is the standard
for meaningful changes. The floor is Summary + Change + Files - always structured,
always searchable. Verify/Risk/Notes added when the change is non-trivial.

"Closed" with zero context is never acceptable. Even a one-line fix gets
"Summary: [what] | Change: [what changed] | Files: [where]".

### 3. Tasks integration woven into lifecycle

Tasks aren't a separate add-on (current Q5 framing). They're part of the
description-to-close lifecycle:
- Description defines the work (what to do, why)
- Tasks decompose the work (how, in what order)
- Close reason captures the outcome (what tasks accomplished, what changed)

The close reason is richer because tasks tracked what happened during the session.

### 4. Dense, not long

The quality bar is "actionable and searchable", not a word count. A 2-sentence
description is fine if both sentences carry information. The floor is density of
useful information, not volume.

### 5. Real examples from beads repo data

Use actual beads from the repo's 704-issue database as good/bad examples rather than
placeholder text. This grounds the teaching in reality.

## Changes to onboard.md

### Description Structure section (lines 385-484)

1. **Add required vs optional framing** per type:
   - Task: Context + Scope + Done-when (required), File locations + Technical context (optional)
   - Bug: Observed + Expected + Steps (required), Environment + Impact (optional)
   - Feature: Motivation + Acceptance criteria + Scope boundary (required),
     Design decisions + Technical context (optional)
   - Epic: Why this matters + Core issues + Phases (required),
     Design decisions + Success criteria + In/out scope (optional)
   - Chore: What + Why now (required), Scope (optional)

2. **Add decision tests** as the WHY for each section:
   - "Would a future agent need this to pick up the work cold?"
   - "Would this surface in a `bd search` for the problem?"
   - "Is there a 'why not the other way?' question a future agent would ask?"

3. **Replace placeholder examples** with filled-in examples using patterns from
   the beads repo (bd-lfak style for epics, bd-1rh style for bugs, etc.)

### Close reason section (reference.md template)

1. **Add minimum viable close** guidance: Summary + Change + Files is the floor
2. **Add "when to use full template"** trigger: non-trivial changes, multi-file
   changes, changes with risk, changes that affect other systems
3. **Show gradient** using real examples from the beads repo:
   - Bad: "Closed" (bd-7yg - 1395-char description, 6-char close)
   - Minimum: "Summary: X | Change: Y | Files: Z"
   - Full: The beads_viewer example (SQLite/Dolt fallback)

### Tasks integration (Q5 / reference.md)

1. **Reframe from "do you use tasks?" to "tasks complete the lifecycle"**
2. **Show how tasks feed close reasons** - the session's task list becomes the
   skeleton of the close reason
3. **Add the decision test**: "Would breaking this into tasks help me write a
   better close reason?" If yes, create tasks.

## Open Questions

None - all resolved during brainstorm dialogue.

## Out of Scope

- Code-level semantic traces (debug.trace patterns) - separate practice, separate doc
- ADR spine documents - lower priority, separate PR
- `--json` flag pattern for AGENTS.md - small addition, can be done separately
- Line count monitoring - not actionable yet
