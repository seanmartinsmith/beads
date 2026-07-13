---
date: 2026-03-16
topic: final-architecture-position
predecessor: 2026-03-16-layered-quality-architecture-brainstorm.md
status: ready-for-plan
---

# Final Architecture Position

Successor to the layered quality architecture brainstorm. Incorporates findings
from: GitHub issue/PR analysis, Gas Town repo research, Steve Yegge interviews
(SE Daily, TWiT, O'Reilly), Dolt blog posts, community analysis, and extensive
dogfooding on tpane (85+ beads).

## The Problem (simple version)

New agents don't read convention docs. Beads has 70+ commands and rich features
but agents only see ~8 via `bd prime`. Close reasons are empty, descriptions are
thin, fields like `--acceptance` go unused, and the conventions set up via onboard
get forgotten after compaction. Users end up manually auditing bead quality.

This isn't a "quality enforcement" problem. It's a "feature discoverability"
problem. The tool is feature-rich. Agents just don't know the features exist.

## Design Philosophy

### Yegge's position: "Quality is just a choice, man"

Respect this. Never block workflows. Warn, don't error. Teach, don't gatekeep.

### Our addition: quality can't be a choice if you don't know the options

If agents are completely unaware of features and too lazy to explore, and there's
no persistent memory across sessions, "choice" is theoretical. Surfacing features
IS respecting the philosophy - it makes the choice real. We're not enforcing
anything. We're exposing what already exists and letting agents use it.

### The tool steered too far

Yegge's philosophy is sound but the implementation overshot. The project changes
so fast that docs can't keep up, and you literally have to brainstorm for hours
to figure out how to use the tool effectively. "Quality is a choice" is fine as
a principle, but if you don't surface how to use the tool without digging through
thousands of commits and issues, it's not actually a choice. The feature richness
(70+ commands) is a strength, but it becomes a weakness when 90% of it is invisible.

The fix isn't enforcement - it's exposure. Let agents discover what's useful and
they'll probably surprise us with how well they use the features. We're not saying
agents should be forced to follow rules. We're saying they should at least know
the tools exist.

### Desire Paths principle

Beads was built by watching what agents naturally try to do and making those
attempts work (100+ subcommands from this). Our teaching should align with
how agents naturally behave, not impose structure they'll fight.

## Three Workstreams

### Workstream 1: Quick Win PR (validate + close reason)

**Fix `--validate` to check `--acceptance` field.**
- Parity with `bd lint` which already does this (GH#2468)
- Small Go change in `internal/validation/template.go`
- `ValidateTemplate()` should check acceptance field, not just description text
- OR have `--validate` call `LintIssue()` directly

**Add `validation.on-close` config option.**
- Warn mode only, defaulting to OFF
- Warns on terse close reasons (<20 chars)
- Does not block any workflow unless explicitly set to error
- Mirrors existing `validation.on-create` pattern

**Also fix: `--dry-run` bypasses `--validate`** - confusing behavior, should
preview validation results too.

PR to `seanmartinsmith/beads` fork, then upstream to `steveyegge/beads`.

### Workstream 2: Prime + Template Improvements

**Improve `bd prime` output** to surface features agents need:
- Close reason template (Summary/Change/Files/Discovery/Verify/Risk)
- Field usage (`--acceptance`, `--design`, `--notes`, `--validate`)
- Description requirements by type (minimum sections)
- `validation.on-create warn` config mention
- Lifecycle commands (`bd defer`, `bd supersede`, `bd close --suggest-next`)
- Hygiene commands (`bd lint`, `bd stale`, `bd orphans`)
- `bd human` for human-in-the-loop workflow
- `bd preflight --check` before PRs
- Memory systems: `bd remember` for project knowledge

**Improve `beads-section.md` template** (full profile for non-hook agents):
- Same quality content as prime improvements
- Flows to Codex, Factory, Mux, OpenCode automatically

**Add convention quality checks to `bd doctor`**:
- `bd doctor --conventions` or `bd doctor --quality`
- Runs: lint + stale + orphans + label taxonomy drift
- Follows existing modular doctorCheck pattern

**Includes brainstorm: formulas/templates as SOP encoding.**
Research what formulas exist today, what's possible, whether pre-built formulas
(TDD, spec-driven-dev, bug-fix) belong in this PR or are a separate initiative.
This is NOT scope creep - understanding how formulas work is core to solving
the "how do agents learn development workflows" problem. If formulas can encode
SOPs, that changes what prime/templates need to teach.

PR to `seanmartinsmith/beads` fork, then upstream.

### Workstream 3: Plugin Skill (proof of concept)

**Revised `/beads:onboard` Claude Code skill:**
- Interactive project-specific convention setup (labels, priority, examples)
- Sets `bd config set validation.on-create warn`
- Generates `.beads/conventions/` files (labels.md, reference.md)
- Inline creation checklist in AGENTS.md (not just pointing to files)
- Session lifecycle teaching (plan mode + Tasks + beads + comments)
- Memory systems complementarity (bd remember vs Claude Code memories)

**Includes brainstorm: Claude Code deep-dive.**
Comprehensive audit of current Claude Code features (Tasks, plan mode, hooks,
skills, subagents, memory) and how they integrate with beads. This feeds the
plugin skill design - need to know what Claude Code offers before we can teach
agents to use both tools together.

**Includes brainstorm: Claude Code deep-dive.**
Claude Code evolves daily with new features. Need comprehensive audit of
current capabilities (Tasks, plan mode, memory, hooks, skills, subagents)
and how they complement beads. This directly feeds the plugin skill design -
can't build a good Claude Code plugin without understanding Claude Code.

**Plugin as proof of concept, not the final destination.**
The feature should eventually be accessible to ALL users, not just Claude Code.
Someone using Cursor or Codex should also be able to set up good conventions.
But plugin first because:
- Solves personal problem immediately (building for myself first)
- Acts as proof of concept with real usage data
- Faster iteration than Go CLI development
- If it works, the CLI version has evidence backing the PR

If the plugin proves the concept on real projects (tpane, portfolio, etc.),
the convention generation logic can be ported to Go CLI later. The PR pitch
writes itself: "here's bead quality before, here's after, here's the tool."

The plugin output is tool-agnostic (AGENTS.md, `.beads/conventions/`) so
anyone benefits from the generated files even if they don't use Claude Code.

## Beads Tracking (created, prefix: sms, in beads repo as contributor)

| ID | Type | P | Title | Status |
|----|------|---|-------|--------|
| sms-imn | bug | P2 | Fix --validate to check --acceptance field | ready |
| sms-q43 | feature | P2 | Add validation.on-close config | ready |
| sms-8si | feature | P2 | Improve bd prime output | blocked by sms-q3r |
| sms-n6q | feature | P2 | Improve beads-section.md template | ready |
| sms-nlw | feature | P2 | Add bd doctor --conventions | ready |
| sms-li8 | feature | P2 | Revise /beads:onboard plugin skill | blocked by sms-8si, sms-s2i |
| sms-q3r | task | P3 | Brainstorm: formulas/templates SOPs | ready |
| sms-s2i | task | P2 | Brainstorm: Claude Code deep-dive | ready |

Dependencies:
- sms-li8 (onboard skill) depends on sms-8si (prime) + sms-s2i (claude code brainstorm)
- sms-8si (prime) depends on sms-q3r (formulas brainstorm)
- sms-8si (prime) related to sms-n6q (template)

## Open Decisions

### Plugin vs CLI first?

**Current position**: Plugin first because:
- Solves personal problem immediately
- Acts as proof of concept with real usage data
- Plugin output is tool-agnostic (AGENTS.md, conventions files)
- Faster iteration cycle than Go CLI development
- If it works, CLI version has evidence backing the PR

**Alternative**: Go straight to CLI if the gap is clear enough. But we'd be
building without dogfooding evidence. Plugin-first is lower risk.

**Important nuance**: This isn't "plugin-only." People who don't use Claude Code
still deserve help setting up conventions. The eventual goal is a CLI command
anyone can run. The plugin proves the concept, the CLI makes it universal. If we
can effectively communicate WHY this is useful and prove it works, the CLI PR
has a much stronger case.

### Name for the CLI version (when we get there)?

`bd onboard` exists and was deliberately simplified. Options:
- `bd conventions init` / `bd conventions setup`
- `bd setup conventions`
- `bd quality init`
- Don't rename - enhance `bd onboard` to accept flags for convention generation
- Decide later based on what the plugin proves

### PR account

Contributions should go through `seanmartinsmith` GitHub account for public
resume building. Fork `steveyegge/beads` to `seanmartinsmith/beads`, branch
per workstream, PR upstream.

## What's NOT in Scope

- Gas Town agent role templates (those are `gt prime`, not `bd prime`)
- Federation features
- External tracker sync (Linear, GitLab, GitHub)
- Dolt migration issues
- Self-upgrade command
- Cross-project visibility (Obsidian vault approach is separate)
- New molecules/formulas (brainstorm only, not implementation)

## Evidence Base

### Dogfooding (tpane, 85+ beads)
- 24% thin descriptions, --acceptance never used, --validate unknown
- P0/P1 miscalibration, scope:* label drift (15+ undocumented)
- Inline AGENTS.md checklist > pointing to reference.md
- Audit bead at P4 never surfaced via bd ready

### Beads repo (704 beads)
- 26% terse close reasons, 77% of multi-step beads closed terse
- Quality jumped 26%→82% when templates added to instructions
- Examples > abstract rules (SkillsBench confirmed)

### Community
- Issue #2140: per-agent setup requested
- Issues #2611/#2612: 91+166 quantified failed command/flag attempts
- PRs #2638/#2639: alias fixes for most common failures

### Yegge philosophy
- "Quality is just a choice" - enable, don't enforce
- "Desire Paths" - make agent natural behavior work
- "Agent-first, human-welcome" - CATNIP for LLMs
- bd prime is SSOT, AGENTS.md is minimal pointer
- bd onboard deliberately simplified to point to prime

## Next Steps

1. **New session**: Review this doc, then write implementation plan
2. Fork beads to `seanmartinsmith/beads` if not already done
3. Start workstream 1 (smallest, fastest to prove contribution path works)
4. Run `bd ready` at session start to see what's available

## Session Recovery

All context persists across sessions via:
- **bd memories** (6 memories) - key architectural decisions and findings
- **bd ready / bd list** - tracked work with dependencies
- **This doc + predecessor** (`2026-03-16-layered-quality-architecture-brainstorm.md`)
- **Brainstorm docs** from earlier iterations in `docs/brainstorms/`
- Working in worktree at `.claude/worktrees/onboard`, branch `feature/onboard-command`
  (rebased onto v0.61.0 main)
