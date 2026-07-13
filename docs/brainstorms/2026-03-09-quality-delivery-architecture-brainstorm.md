# Brainstorm: Beads Quality Delivery Architecture

Date: 2026-03-09
Status: In progress (paused after Phase 2 exploration)
Predecessor: `docs/brainstorms/2026-03-06-onboard-description-quality-brainstorm.md`

## The Problem

Agents produce low-quality beads (terse descriptions, empty close reasons) when
conventions aren't in context. The conventions themselves are fine - the delivery
mechanism is broken. Same project, same day, wildly different quality depending on
whether the agent happened to have conventions loaded.

Evidence from beads repo (704 issues):
- 26% of close reasons are terse (<20 chars), mostly just "Closed"
- 77% of multi-step beads closed with terse reasons
- Descriptions are strong (91% >100 chars) when templates exist

Evidence from dogfooding (cctui + tpane, same day):
- Some beads follow conventions perfectly (structured close, labeled sections)
- Others are one-liners despite conventions being set up via onboard

## Root Cause Analysis

Two separate problems:

### 1. Persistence (conventions don't survive compaction)

- `bd prime` runs at SessionStart + PreCompact - only content guaranteed to survive
- But prime carries ZERO convention awareness (just workflow rules + CLI reference)
- AGENTS.md has convention pointers but gets compacted/summarized
- Convention files on disk are passive - no mechanism forces agents to read them

### 2. Consistency (conventions not always followed even when loaded)

- Some agents skip conventions even when they're in context
- Skill triggering is probabilistic (depends on intent matching)
- No enforcement mechanism at the tool level

## Approaches Explored

### 1. Skills (rejected as primary mechanism)

Idea: Create `beads-quality` skill that triggers when agents create/close beads.

Pros:
- Progressive disclosure (frontmatter always loaded, body on demand)
- Point-of-need injection

Cons:
- Skills are plugin-only (doesn't work for CLI-only users)
- Triggering is probabilistic, not guaranteed
- Content still gets compacted in long sessions

### 2. PreToolUse Hooks (rejected due to overhead)

Idea: Hook on Bash tool that detects `bd create`/`bd close` and injects conventions.

Pros:
- Can't be skipped (fires automatically)
- Point-of-need injection

Cons:
- Fires on EVERY Bash call (hook script spawns per call)
- Windows process spawn cost: ~100-300ms per Bash call
- Tax on all Bash usage, not just bd commands
- Claude Code-specific

### 3. CLI Output as Delivery Vehicle (explored, partially promising)

Idea: `bd create` and `bd close` include quality templates in their output.

Pros:
- Tool-agnostic (works everywhere)
- Can't be "forgotten" - enforcement IS the command
- Survives compaction by nature (output is at point of need)

Cons:
- Template in every create/close output is noisy after first time
- Timing issue: by the time `bd create` runs, description is already written

Status: User wasn't sold on showing templates every command.

### 4. Enhanced bd prime (strong candidate for persistence)

Idea: Prime learns about conventions. Includes compact quality section (~50 tokens).

```
## Quality Floor
- Bug: Observed + Expected + Steps (required)
- Feature: Motivation + Acceptance + Scope (required)
- Close: Summary + Change + Files (minimum)
- Full conventions: .beads/conventions/reference.md
```

Pros:
- Survives compaction (PreCompact hook re-injects)
- Minimal token cost (~50 extra tokens per session)
- No new mechanisms - extends existing infrastructure
- Works with any AI tool (prime outputs to any consumer)

Cons:
- Soft enforcement (agent can still ignore)
- Must be compact to avoid bloating prime output

Status: Strong candidate for solving the persistence problem.

### 5. CLI Validation (configurable strictness)

Idea: `bd create` and `bd close` validate input quality.

Default: Warning + accept (soft enforcement)
Configurable: Hard reject per project (strict mode)

Pros:
- Tool-agnostic enforcement
- Configurable per project via onboard

Cons:
- Deterministic checks are limited (length, section markers)
- LLMs are generally compliant when given instructions - validation may be overkill

Status: Accepted as safety net, not primary mechanism.

## Emerging Architecture

The decomposition that felt right before pause:

1. **Persistence layer (enhanced prime)**: bd prime includes compact quality
   conventions. Survives compaction. Costs ~50 extra tokens. Convention content
   comes from onboard-generated files.

2. **Consistency layer (TBD)**: Mechanism to ensure agents follow conventions
   even when loaded. Could be CLI validation, skill, hook, or just better
   instruction writing. This is the unsolved part.

3. **Customization layer (onboard)**: Onboard generates project-specific
   convention content. CLI/prime reads it. Each project gets tailored templates.

4. **Safety net (CLI validation)**: Minimal checks - reject "Closed" with no
   context, warn on terse close reasons. Configurable strictness per project.

## Key Design Decisions Made

1. **No hard rejection by default** - warnings, not errors. Hard reject available
   as per-project config option.

2. **Onboard remains the customization mechanism** - generates project-specific
   conventions that other systems consume.

3. **Convention files stay on disk** - they're configuration, not instructions
   agents must remember to read. Something else (prime/CLI) reads and delivers them.

4. **LLMs are compliant when instructed** - the problem is delivery, not
   enforcement. Focus on getting conventions in context, not on validation logic.

## Open Questions (resume here)

### Q1: Does fixing persistence alone fix most quality issues?
If prime keeps conventions in context (survives compaction), is that enough?
Or do we need a separate consistency mechanism? The user was mulling this when
we paused.

### Q2: What's the right consistency mechanism?
If persistence alone isn't enough, which approach for consistency:
- CLI validation (warn/reject)
- Skill at point of need
- Better instruction writing (make conventions more memorable)
- Something else entirely

### Q3: How should convention content flow from onboard to prime?
Options:
- PRIME.md override file (exists today, replaces all prime output)
- New convention-specific file that prime reads alongside its defaults
- Direct config in .beads/ that prime auto-discovers

### Q4: What about the skill idea?
The user's instinct about skills was strong - "maybe onboarding should create a
skill instead of convention files." This wasn't fully explored. Worth revisiting:
could a skill be the CONSISTENCY layer while prime handles PERSISTENCE?

## Architecture Research Summary

### Current Convention Delivery Channels

| Channel | Fires when | Survives compaction | Convention aware |
|---|---|---|---|
| bd prime | SessionStart, PreCompact | Yes (re-injected) | No |
| AGENTS.md | Session start | No (compacted) | Yes (pointers) |
| Convention files | Manual read | N/A (on disk) | Yes (content) |
| Beads skill | Intent match | No (compacted) | No |
| Task agent | Spawned | No | No |

### Gap: No channel is both compaction-proof AND convention-aware.

Prime is compaction-proof but convention-blind.
Convention files are convention-rich but have no delivery mechanism.
Closing this gap is the core of the architecture.
