# mkt Teaching Angle - Agent Instruction Gap Analysis

*Generated: 2026-04-13 for bd-0il*

## Executive Summary

65% of bd's 103 subcommands are invisible to agents. The teaching layer has three structural problems:

1. **PRIME.md carries too much weight.** It's the only real teaching surface, loaded every session (~1.5k tokens). Everything not in PRIME.md is effectively unknown. But PRIME.md can't grow much without bloating every session.

2. **The beads plugin skill has depth but no progressive disclosure.** The `beads:beads` skill contains 15 resource files covering gates, molecules, chemistry patterns, workflows, etc. But agents never invoke the skill because PRIME.md doesn't tell them these capabilities exist or when to reach for the skill.

3. **No behavioral hooks.** Teaching documents tell agents what commands exist. Nothing enforces correct usage patterns or catches misuse in flight. The human/gate confusion persists because there's no hook that fires when an agent uses `bd human` thinking it blocks workflow.

The fix is a three-layer teaching architecture:

- **PRIME.md** (awareness layer): What exists, when to reach for it, 2-3 sentence pointers to deeper resources. No command syntax beyond the core 15.
- **Skills** (depth layer): On-demand detailed guidance for feature groups. Agents invoke when PRIME.md tells them the capability is relevant.
- **Hooks** (enforcement layer): Catch misuse patterns and redirect in real-time. Lightweight, focused on the highest-damage confusion points.

---

## PRIME.md Recommendations

### Must Add (Critical) - Immediate

These address confusion that is actively causing wrong decisions.

#### 1. Human vs. Gate Clarification

Replace the current "Human Decisions" section with text that disambiguates advisory flagging from blocking:

```
### Flagging for Human Input

`bd human <id>` adds a `human` label - this is ADVISORY. The issue still
appears in `bd ready`. Use this when you want to signal "a human should
look at this" without blocking workflow.

To actually BLOCK an issue until a human approves, use a gate:
  bd gate create --await human:"approval prompt" --title "Waiting for human"
Gates prevent blocked steps from appearing in `bd ready` until resolved.
See /beads skill resources for full gate documentation.
```

**Token cost:** ~80 tokens net (replaces existing ~60 token section). Negligible.

#### 2. Markdown in Fields Pointer

Add under the "Creating & Updating" section:

```
- **Rich fields**: `--description`, `--notes`, `--design` accept full markdown.
  Use checklists, headers, code blocks, and structured sections. Dense,
  structured descriptions are dramatically more searchable and recoverable
  than flat text. See /beads skill resources for field templates.
```

**Token cost:** ~40 tokens. Negligible.

#### 3. Advanced Capabilities Awareness Block

Add a new section after "Common Workflows" that points to the skill layer:

```
## Beyond the Basics

These capabilities exist but aren't needed every session. Invoke the
/beads skill for detailed guidance when the situation matches.

| Situation | Capability | Entry Point |
|-----------|-----------|-------------|
| Need to block on CI, PR, timer, or human approval | Gates | `bd gate --help` |
| Repeatable workflow (release, review, patrol) | Formulas + Molecules | `bd formula list` |
| Complex filtering beyond bd list flags | Query language | `bd query --help` |
| Parallel agents on one epic | Swarm | `bd swarm --help` |
| Agents racing to merge | Merge slot | `bd merge-slot --help` |
| Cross-project dependency | Ship | `bd ship --help` |
| Bulk create/update (5+ operations) | Batch | `bd batch --help` |
| Persist arbitrary key-value data across sessions | KV store | `bd kv --help` |
| Ephemeral work that shouldn't clutter history | Wisps | `bd mol wisp --help` |
```

**Token cost:** ~150 tokens. Worth it because it turns "unknown" features into "discoverable" features without teaching the full syntax.

### Should Add (High) - Next Release

#### 4. `bd query` One-Liner

Under "Finding Work" in the command reference:

```
- `bd query "status=open AND priority<=1 AND updated>7d"` - Compound query
  (more powerful than bd list flags)
```

**Token cost:** ~25 tokens.

#### 5. `bd batch` Mention

Under "Creating & Updating":

```
- `bd batch` - Run multiple write operations in one transaction (faster for 5+ ops)
```

**Token cost:** ~15 tokens.

#### 6. `bd bootstrap` in Error Handling

Add to the error handling table:

```
| `database not found` (fresh clone) | `bd bootstrap` (safe) or `bd init` (destructive) |
```

**Token cost:** ~15 tokens, replaces existing row.

### Consider (Medium) - Future

#### 7. Integration Command Awareness

A single line under a "Sync" heading:

```
- External tracker sync: `bd github`, `bd jira`, `bd linear`, `bd gitlab`,
  `bd ado`, `bd notion` (bidirectional, configure via `bd config`)
```

**Token cost:** ~25 tokens. Only add if users start working in environments with external trackers.

#### 8. `bd kv` in Memory Section

```
- `bd kv set <key> <value>` / `bd kv get <key>` - Persist arbitrary data
  across sessions (feature flags, environment markers)
```

**Token cost:** ~20 tokens.

---

## New Skills Recommendations

The existing `beads:beads` skill with its 15 resource files is already comprehensive. The problem isn't missing content - it's that agents never reach for it. The PRIME.md awareness block (recommendation #3 above) is the primary fix.

However, two new dedicated skills would serve situations where the general beads skill is too broad:

### Skill 1: `/beads-gates` (New)

**Rationale:** The gate system is the single highest-confusion area. Agents conflate `bd human` (advisory) with blocking. A dedicated skill that agents can invoke when they need to block on external conditions would solve this without requiring them to navigate the full beads skill.

**Trigger phrases:** "block until human approves", "wait for CI", "wait for PR merge", "gate", "approval gate", "timer gate", "block workflow"

**Content sketch:**

```
# Gates - Blocking on External Conditions

## When to Use Gates (not bd human)

bd human = advisory label. Issue still in bd ready.
bd gate = actual blocker. Issue hidden from bd ready until condition met.

## Decision: bd human vs bd gate

"Do I need to PREVENT an agent from picking this up?"
  YES -> bd gate create --await human:"prompt"
  NO  -> bd human <id> (just flag it)

## Gate Types
[Table of 5 gate types with create syntax]

## Common Patterns
- Pre-deploy approval gate
- Wait for CI before merging
- Timer gate for propagation delay
- Cross-project bead gate

## Monitoring
bd gate list / bd gate eval / bd gate approve
```

**Token budget:** ~400 words in SKILL.md, referencing the existing `ASYNC_GATES.md` resource.

### Skill 2: `/beads-workflows` (New)

**Rationale:** The formula/molecule/wisp system is powerful but impenetrable from PRIME.md alone. Agents need a skill that walks them through "I have a repeatable workflow" to "Here's how to template it."

**Trigger phrases:** "create a workflow template", "repeatable workflow", "formula", "molecule", "proto", "wisp", "ephemeral issue", "patrol cycle", "release process template"

**Content sketch:**

```
# Workflow Templates - Formulas, Molecules, Wisps

## The 30-Second Version

formula (YAML) --cook--> proto (template epic) --pour/wisp--> molecule/wisp (real work)

## Decision: Do I Need This?

"Am I doing the same multi-step workflow more than twice?"
  YES -> Encode as formula, pour instances.
  NO  -> Just create issues normally.

"Does this work need permanent record?"
  YES -> bd mol pour (persistent molecule)
  NO  -> bd mol wisp (ephemeral, burn or squash when done)

## Creating Your First Formula
[Minimal YAML example -> bd cook -> bd mol pour with variable substitution]

## Common Patterns
- Release workflow
- Code review checklist (wisp, squash findings)
- Patrol cycle (wisp, burn when clean)
- Feature with rollback (conditional bond)

## Lifecycle
[Proto creation, variable substitution, pour vs wisp, squash vs burn vs promote]
```

**Token budget:** ~500 words in SKILL.md, referencing existing `MOLECULES.md` and `CHEMISTRY_PATTERNS.md` resources.

### Skill 3: `/beads-swarm` (Consider, not immediate)

**Rationale:** Multi-agent parallel work is a growing use case. Swarm + merge-slot are both unknown. But the user base for this is smaller than gates or workflows.

**Defer until:** Multi-agent scenarios become more common in the user's workflow.

---

## beads-onboard Enhancements

### v5 Targets

beads-onboard runs once per project. It's the right place for project-setup decisions, not ongoing workflow guidance. v5 should:

#### 1. Gate Awareness in Phase 2 Interview

Add a new question (Q7):

> "Does your project use approval gates or CI-blocking workflows?"
>
> - **Yes, human approvals** - Generate gate examples with human type
> - **Yes, CI gates** - Generate gate examples with gh:run type
> - **No** - Skip gate setup
> - **Not sure yet** - Add a P3 reminder bead to revisit gates after 10 beads

If yes, generate example gate commands in `.beads/conventions/reference.md` and add the human-vs-gate disambiguation to the AGENTS.md beads section.

#### 2. Markdown Field Templates in reference.md

Extend the "Creating Issues" section to include markdown templates by issue type:

```markdown
### Bug Description Template (markdown)
\`\`\`
## Observed Behavior
[What happens]

## Expected Behavior
[What should happen]

## Reproduction
- [ ] Step 1
- [ ] Step 2

## Environment
- OS:
- bd version:
\`\`\`
```

This teaches agents to use markdown structure in fields at project-setup time.

#### 3. Query Language Mention in reference.md

Add a "Power Queries" section showing 3-4 `bd query` examples relevant to the project's label taxonomy:

```
bd query "status=open AND label=area:cli AND priority<=1"
bd query "type=bug AND updated<14d"
bd query "assignee='' AND status=open"  # Unassigned work
```

#### 4. Integration Detection in Phase 1

If `.github/` exists, ask about GitHub Issues sync. If `linear.app` or `jira` references found in configs, ask about those integrations. Generate `bd config set` commands in AGENTS.md for detected integrations.

---

## Hook Recommendations

### Hook 1: Human/Gate Confusion Guard (Critical, Immediate)

**Type:** PostToolUse on Bash
**Trigger:** Command matches `bd human <id>` (not `bd human list` or bare `bd human`)
**Behavior:** Inject a brief advisory after the command runs:

```
Note: `bd human` is advisory only - it adds a label but does NOT block
the issue from bd ready. If you need to actually block workflow until a
human responds, use `bd gate create --await human:"prompt" --timeout 24h`.
```

**Why hook, not just docs:** This fires in the moment of confusion. The agent just ran `bd human <id>` - if they intended blocking behavior, the hook catches it immediately rather than relying on them having read and remembered the docs.

**Token cost per trigger:** ~50 tokens. Only fires when the specific pattern matches, not every command.

### Hook 2: Rich Description Nudge (High, Next Release)

**Type:** PostToolUse on Bash
**Trigger:** Command matches `bd create` and `--description` value is a single sentence (no newlines, no markdown markers)
**Behavior:** Gentle nudge:

```
Tip: Descriptions support full markdown. For bugs, use a checklist
(observed/expected/repro). For features, use headers for context/approach/
acceptance criteria. Structured descriptions are dramatically more
searchable.
```

**Why hook:** Agents default to flat text because nothing pushes them otherwise. A one-time nudge per session (use a session flag to avoid repeating) changes the default behavior.

### Hook 3: Batch Optimization Hint (Medium, Future)

**Type:** PostToolUse on Bash
**Trigger:** Agent has run 3+ `bd create` or `bd update` commands in the current session
**Behavior:** One-time suggestion:

```
Tip: For bulk operations, `bd batch` runs multiple writes in one
transaction - faster and atomic. Pipe commands via stdin.
```

**Why hook:** Agents creating many issues in sequence don't know they're paying per-command Dolt overhead. This is a performance teaching moment.

### Hook 4: Bootstrap Safety Net (Medium, Future)

**Type:** PreToolUse on Bash
**Trigger:** Command matches `bd init` in a directory that already has `.beads/`
**Behavior:** Warning before execution:

```
Warning: This directory already has .beads/. `bd init` may be destructive.
Did you mean `bd bootstrap` (non-destructive recovery)?
```

---

## Feature-by-Feature Analysis

### bd human / bd gate (The Critical Confusion)

**Gap Severity:** Critical
**Teaching Vehicle:** PRIME.md rewrite (awareness) + PostToolUse hook (enforcement) + /beads-gates skill (depth)
**Instruction Draft:** See PRIME.md recommendation #1 and Hook #1 above. The three layers work together: PRIME.md plants the correct mental model, the hook catches misuse in real-time, the skill provides full gate documentation on demand.
**Priority:** Immediate. This is the #1 source of incorrect agent behavior.

### Markdown in Fields

**Gap Severity:** Critical
**Teaching Vehicle:** PRIME.md pointer (awareness) + PostToolUse hook on `bd create` (enforcement) + beads-onboard v5 templates (project setup)
**Instruction Draft:** See PRIME.md recommendation #2, Hook #2, and beads-onboard enhancement #2. The goal is shifting agent defaults from "flat sentence" to "structured markdown" for descriptions, notes, and design fields.
**Priority:** Immediate. Every issue created with flat text is a missed opportunity for searchability and context recovery.

### bd gate (Full Gate System)

**Gap Severity:** Critical
**Teaching Vehicle:** /beads-gates dedicated skill + PRIME.md awareness table entry
**Instruction Draft:** See New Skills recommendation #1. The skill provides the decision framework (when gates vs. when not), all 5 gate types, creation/monitoring/closing patterns, and the gates-vs-issues comparison table.
**Priority:** Immediate (skill creation). Gates are the actual coordination primitive for async workflows.

### bd formula / bd mol / bd cook (Workflow Templating)

**Gap Severity:** High
**Teaching Vehicle:** /beads-workflows dedicated skill + PRIME.md awareness table entry
**Instruction Draft:** See New Skills recommendation #2. The existing `MOLECULES.md` and `CHEMISTRY_PATTERNS.md` resources in the beads plugin are comprehensive. The new skill provides the entry ramp: "when do I need this?" decision tree and a minimal first-formula walkthrough. The existing resources handle depth.
**Priority:** Phase 2 (after gate confusion is fixed). Formulas are powerful but the user needs gates working correctly first since formulas create gates.

### bd query (Query Language)

**Gap Severity:** High
**Teaching Vehicle:** PRIME.md one-liner (awareness) + beads-onboard v5 (project examples) + existing beads skill resources
**Instruction Draft:** See PRIME.md recommendation #4 and beads-onboard enhancement #3. A single example in PRIME.md shows the syntax pattern; project-specific examples in conventions/reference.md make it actionable.
**Priority:** Phase 2. High leverage (agents currently overuse `bd list` + jq) but not causing wrong decisions, just inefficiency.

### bd swarm (Parallel Work Coordination)

**Gap Severity:** High
**Teaching Vehicle:** PRIME.md awareness table entry + future /beads-swarm skill
**Instruction Draft:** One line in the awareness table: "Parallel agents on one epic -> Swarm -> `bd swarm --help`". Detailed skill deferred until multi-agent usage grows.
**Priority:** Phase 2 for awareness, Phase 3 for dedicated skill.

### bd merge-slot (Serialized Conflict Resolution)

**Gap Severity:** High
**Teaching Vehicle:** PRIME.md awareness table entry + mentioned in /beads-swarm skill when created
**Instruction Draft:** One line: "Agents racing to merge -> Merge slot -> `bd merge-slot --help`". Naturally pairs with swarm documentation.
**Priority:** Phase 2 for awareness, Phase 3 for depth.

### bd ship (Cross-Project Capability Publishing)

**Gap Severity:** High
**Teaching Vehicle:** PRIME.md awareness table entry + existing beads skill `MOLECULES.md` resource (already documents ship)
**Instruction Draft:** One line: "Cross-project dependency -> Ship -> `bd ship --help`". The `MOLECULES.md` resource already has a "Cross-Project Dependencies" section with full documentation. The problem is purely discoverability.
**Priority:** Phase 2. Only relevant when user works across multiple beads projects.

### bd set-state / bd state (Operational State Management)

**Gap Severity:** High
**Teaching Vehicle:** PRIME.md awareness table entry + /beads-workflows skill
**Instruction Draft:** One line in awareness table. Detailed in the workflows skill under "Operational Patterns" (patrol cycles, health monitoring). State management is a building block for long-running agent workflows.
**Priority:** Phase 2.

### bd batch (Bulk Operations)

**Gap Severity:** Medium-High
**Teaching Vehicle:** PRIME.md one-liner + PostToolUse hook (after 3+ sequential writes)
**Instruction Draft:** See PRIME.md recommendation #5 and Hook #3. Teaching happens at the moment agents would benefit from it.
**Priority:** Phase 2.

### bd kv (Cross-Session Persistence)

**Gap Severity:** Medium-High
**Teaching Vehicle:** PRIME.md awareness table entry + optional PRIME.md memory section addition
**Instruction Draft:** See PRIME.md recommendations #3 (table) and #8 (memory section). KV is distinct from `bd remember` - it's for operational state (feature flags, environment markers), not project learnings.
**Priority:** Phase 2.

### bd graph (Visualization)

**Gap Severity:** Medium
**Teaching Vehicle:** Existing beads skill resources (already partially covered in DEPENDENCIES.md)
**Instruction Draft:** Not worth PRIME.md tokens. When an agent needs visualization, `bd graph --help` is self-documenting. Could mention in beads-onboard v5 reference.md under "Useful Commands."
**Priority:** Phase 3.

### bd find-duplicates / bd duplicates

**Gap Severity:** Medium
**Teaching Vehicle:** Existing beads skill resources + beads-onboard v5 reference.md
**Instruction Draft:** Add to the "Project Hygiene" section of conventions/reference.md: "Before creating: `bd search` for text match. Periodically: `bd find-duplicates` for semantic similarity. On detection: `bd duplicate <id> --of <canonical>`."
**Priority:** Phase 3.

### bd bootstrap

**Gap Severity:** Medium
**Teaching Vehicle:** PRIME.md error handling table replacement + Hook #4
**Instruction Draft:** See PRIME.md recommendation #6 and Hook #4. The hook catches the dangerous case (running `bd init` on existing DB); the error table teaches the safe alternative.
**Priority:** Phase 2.

### bd comments / bd comment

**Gap Severity:** Medium
**Teaching Vehicle:** beads-onboard v5 reference.md
**Instruction Draft:** The existing beads-onboard reference.md already has a "When to Use Comments vs Fields" section. Just needs the actual commands added: `bd comment <id> "text"` and `bd comments <id>`.
**Priority:** Phase 3.

### bd children

**Gap Severity:** Low-Medium
**Teaching Vehicle:** Not worth dedicated teaching. `bd children <id>` is discoverable from `--help`.
**Priority:** Phase 3 or never.

### bd epic (close-eligible)

**Gap Severity:** Medium
**Teaching Vehicle:** beads-onboard v5 reference.md under "Lifecycle" section
**Instruction Draft:** "When all children of an epic are closed: `bd epic close-eligible` to auto-close eligible epics."
**Priority:** Phase 3.

### bd promote (wisp to bead)

**Gap Severity:** Medium
**Teaching Vehicle:** /beads-workflows skill, under wisp lifecycle
**Instruction Draft:** Already partially covered in existing `CHEMISTRY_PATTERNS.md` resource. The /beads-workflows skill should include this in the wisp decision tree: "Wisp turned out important? `bd promote <wisp-id>` to make it permanent."
**Priority:** Phase 2 (bundled with workflows skill).

### bd reopen

**Gap Severity:** Low
**Teaching Vehicle:** Not worth dedicated teaching. Agents can use `bd update --status open` which works. `bd reopen` is cleaner but not critical.
**Priority:** Phase 3 or never.

### bd undefer

**Gap Severity:** Low-Medium
**Teaching Vehicle:** PRIME.md could add `bd undefer <id>` next to the existing `bd defer` line
**Instruction Draft:** One-word addition. Currently PRIME.md teaches `bd defer` but not the reverse.
**Priority:** Phase 2 (trivial addition).

### Integration Commands (GitHub, Jira, Linear, GitLab, ADO, Notion)

**Gap Severity:** Medium (varies by user environment)
**Teaching Vehicle:** beads-onboard v5 auto-detection (Phase 1) + PRIME.md awareness line
**Instruction Draft:** See beads-onboard enhancement #4 and PRIME.md recommendation #7. Integration teaching should be driven by detected project environment, not blanket.
**Priority:** Phase 3 unless user actively uses external trackers.

### bd worktree

**Gap Severity:** Medium
**Teaching Vehicle:** Existing beads skill `WORKTREES.md` resource (already exists)
**Instruction Draft:** Discoverability problem only. The resource exists. Add to PRIME.md awareness table: "Parallel development dirs -> Worktree -> `bd worktree --help`".
**Priority:** Phase 3.

### bd vc / bd branch / bd diff

**Gap Severity:** Low
**Teaching Vehicle:** Not worth dedicated teaching for most users. Dolt version control is admin-level functionality.
**Priority:** Phase 3 or never.

### bd todo

**Gap Severity:** Low
**Teaching Vehicle:** Not worth teaching. Agents use `bd create` which is more explicit. `bd todo` is a human convenience shorthand.
**Priority:** Never.

### bd q (Quick Capture)

**Gap Severity:** Low
**Teaching Vehicle:** Not worth teaching. Agents use `bd create --json` and parse the output. `bd q` saves a few characters.
**Priority:** Never.

### bd count

**Gap Severity:** Low
**Teaching Vehicle:** Not worth teaching. Agents use `bd list` and count. `bd count` is marginally more efficient.
**Priority:** Phase 3 or never.

### bd history

**Gap Severity:** Low-Medium
**Teaching Vehicle:** Could mention in the beads skill's `RESUMABILITY.md` resource
**Instruction Draft:** "To see how an issue evolved over time: `bd history <id>` shows all Dolt commits where the issue was modified."
**Priority:** Phase 3.

### bd audit

**Gap Severity:** Low (niche use case)
**Teaching Vehicle:** Existing beads:audit skill
**Instruction Draft:** Already has a dedicated skill. Usage is niche (dataset generation, compliance).
**Priority:** Phase 3 or never.

### Maintenance Commands (bd gc, bd compact, bd admin, bd flatten, bd purge)

**Gap Severity:** N/A for agents
**Teaching Vehicle:** None. These are human-initiated maintenance operations.
**Priority:** Never teach to agents. These should be behind `bd human` flags or manual invocation.

### bd sql

**Gap Severity:** N/A for agents
**Teaching Vehicle:** None. Raw SQL access is dangerous in agent hands.
**Priority:** Never teach to agents.

### bd create-form / bd completion

**Gap Severity:** N/A for agents
**Teaching Vehicle:** None. Human convenience features.
**Priority:** Never.

---

## Implementation Roadmap

### Phase 1: Fix Critical Confusions (Immediate)

**Scope:** 3 PRIME.md changes, 1 new hook, 1 new skill

1. **PRIME.md: Human/Gate disambiguation** (recommendation #1)
   - Rewrite the "Human Decisions" section
   - Add gate awareness with correct mental model
   - Token delta: ~+20 tokens net

2. **PRIME.md: Markdown in fields pointer** (recommendation #2)
   - Add one line to "Creating & Updating"
   - Token delta: ~+40 tokens

3. **PRIME.md: Advanced capabilities awareness table** (recommendation #3)
   - New section after "Common Workflows" pointing to skill layer
   - Token delta: ~+150 tokens

4. **Hook: Human/Gate confusion guard** (hook #1)
   - PostToolUse on Bash, fires on `bd human <id>` pattern
   - Immediate correction at point of confusion

5. **Skill: /beads-gates** (new skill #1)
   - Dedicated gate skill with decision framework
   - References existing `ASYNC_GATES.md` resource for depth
   - Solves the "I need to block on something" use case

**Total PRIME.md token increase:** ~210 tokens (~14% growth on ~1.5k base). Acceptable because it converts 8 features from "unknown" to "discoverable."

**Effort:** Low sequential depth. PRIME.md changes and hook can be done in parallel with skill creation. 2 agent-slots, 1 sequential step.

### Phase 2: Surface Power Features (Next Release)

**Scope:** 4 PRIME.md additions, 1 new skill, 2 hooks, beads-onboard v5

1. **Skill: /beads-workflows** (new skill #2)
   - Formula/molecule/wisp lifecycle skill
   - Entry ramp to the existing MOLECULES.md and CHEMISTRY_PATTERNS.md resources

2. **PRIME.md additions:**
   - `bd query` one-liner (recommendation #4)
   - `bd batch` mention (recommendation #5)
   - `bd bootstrap` in error table (recommendation #6)
   - `bd undefer` next to `bd defer`

3. **Hooks:**
   - Rich description nudge (hook #2)
   - Batch optimization hint (hook #3)

4. **beads-onboard v5:**
   - Gate awareness question (Q7)
   - Markdown field templates
   - Query language examples
   - Integration auto-detection

**Total PRIME.md token increase:** ~75 tokens beyond Phase 1. Cumulative ~285 tokens growth, still well within budget.

**Effort:** Medium. beads-onboard v5 has the most work (updating 3 reference files + adding interview question). 3 agent-slots, 2 sequential steps (PRIME.md/hooks parallel with skill, then beads-onboard depends on skill being done).

### Phase 3: Complete Coverage (Long-tail)

**Scope:** Incremental additions, no new skills

1. **beads-onboard v5 reference.md additions:**
   - `bd find-duplicates` / `bd duplicate` in hygiene section
   - `bd comment` / `bd comments` commands in the existing comments-vs-fields section
   - `bd epic close-eligible` in lifecycle section
   - `bd graph` in useful commands

2. **PRIME.md (only if user demand emerges):**
   - Integration commands awareness line
   - `bd kv` in memory section
   - `bd worktree` in awareness table

3. **Existing beads skill resource updates:**
   - Add `bd history` to RESUMABILITY.md
   - Ensure WORKFLOWS.md mentions `bd set-state`/`bd state` for operational patterns
   - Update INTEGRATION_PATTERNS.md with current integration commands

**Effort:** Low. All incremental additions to existing files. 1 agent-slot, 1 step per addition.

---

## Token Budget Summary

| Layer | Current | After Phase 1 | After Phase 2 | After Phase 3 |
|-------|---------|---------------|---------------|---------------|
| PRIME.md | ~1,500 tokens | ~1,710 tokens | ~1,785 tokens | ~1,835 tokens |
| /beads-gates skill | 0 | ~400 words | ~400 words | ~400 words |
| /beads-workflows skill | 0 | 0 | ~500 words | ~500 words |
| beads-onboard (runs once) | v4 | v4 | v5 | v5 |
| Hooks (per-trigger) | 0 | ~50 tokens | ~120 tokens | ~120 tokens |

PRIME.md stays under 2k tokens through all phases. The strategy deliberately pushes depth into skills (loaded on demand) and hooks (loaded only on trigger), keeping the every-session cost low.

---

## Key Design Decisions

**Why not just expand PRIME.md?** PRIME.md loads every session, every turn in some modes. At 1.5k tokens it's already substantial. Adding full documentation for gates, molecules, query language, etc. would push it to 4-5k tokens - a permanent tax on every interaction. Skills load on demand; they're free when not needed.

**Why not just update the existing beads skill resources?** The resources are already good. The problem is discoverability - agents don't know to invoke the skill for capabilities they don't know exist. The awareness table in PRIME.md is the bridge: it tells agents "this exists, invoke /beads for details."

**Why hooks for human/gate confusion?** Documentation alone hasn't fixed this. The confusion persists because agents read PRIME.md at session start and then forget the nuance when they're deep in work. A hook fires at the exact moment of confusion, providing correction when it matters most.

**Why two new skills instead of expanding beads:beads?** The general beads skill is already the size of a small manual (15 resource files, ~150KB). Two focused skills with clear trigger phrases ("I need to block on something" -> /beads-gates, "I need a workflow template" -> /beads-workflows) are more likely to be invoked than hoping agents navigate the general skill to the right resource.
