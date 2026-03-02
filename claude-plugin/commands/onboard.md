---
description: Set up project-specific quality conventions for beads
argument-hint:
---

Onboard a project with beads quality conventions. Analyzes the codebase and
interactively generates project-specific rules that make agents produce
high-quality, searchable, non-redundant beads.

Run this once per project, ideally right after `bd init`.

## Phase 0: State Detection

Before starting the interview, detect the current project state.

1. Check for `.beads/` directory. **If missing, stop and suggest `bd init`.**
2. Read `.beads/config.yaml` for the issue prefix.
3. Check for `.beads/conventions/` directory:
   - If exists: this is a re-onboard. Ask: "Conventions already exist. Update them or start fresh?"
4. Check for `<!-- BEGIN BEADS INTEGRATION -->` markers in AGENTS.md:
   - If markers found: will replace content between markers (interop with `bd setup`).
   - If no markers but beads content found: will wrap replacement in markers.
   - If no AGENTS.md: will create it with markers.
5. Check if CLAUDE.md exists and whether it references `@AGENTS.md`.
6. Run `bd doctor 2>&1 || true` (capture output, don't show it yet). Note errors
   and warnings to surface as recommendations after generation.

**Do not block on doctor results.** The only hard gate is `.beads/` existence.

## Phase 1: Analyze the Codebase

Perform mechanical detection to inform the interview:

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

## Phase 2: Interview

Ask these questions **one at a time** using AskUserQuestion. Wait for each
answer before proceeding to the next question.

### Q1: Project Description

Generate three short descriptions (1-2 sentences each) based on the codebase
analysis. Present them as options (the built-in "Other" handles freeform):

- Option 1: [AI-generated description emphasizing primary function]
- Option 2: [AI-generated description emphasizing tech stack/architecture]
- Option 3: [AI-generated description emphasizing user/audience]

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

Every option must include a markdown preview showing example commits in that
format (use real examples from git log, rewritten into each format). Options
without previews show "No preview available" which is confusing.

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

### Q5: Beads + Ephemeral Tasks

"Do you use ephemeral task tools (Claude Code Tasks, Copilot Todos, etc.)
alongside beads for within-session work?"

Explain the complementary pattern regardless of answer:
- Beads = persistent across sessions, changelog-worthy, searchable history
- Tasks = ephemeral within a session, execution coordination, implementation steps

Options:
- Yes, I use both (tailor the generated docs to emphasize the integration)
- No, just beads (still mention tasks as an option in the reference doc)
- Not sure what this means (explain further, then re-ask)

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

**Formatting:** Use blank lines between sections. No semicolons.

**Good example:**
[Generate one example using this project's prefix]

**Bad examples:**
- "Done" - no context, unsearchable
- "Fixed it" - no details for future agents

## Creating Issues

Always include: `--type`, `--priority`, `--labels`, `--description`.

Priority: 0 (critical) to 4 (backlog).
Search before creating: `bd search "query"`

## Dependencies

- Only add blocking deps when work truly cannot start
- Default to NOT adding blocking deps
- Use `bd dep relate <a> <b>` for connected but parallel work
- Never use `blocks` to mean "belongs to epic" - use parent-child

## Session Discipline

- `bd show <id>` and read close_reason before starting work
- Check `bd list --status=in_progress` at session start
- Use `bd comments add <id> "..."` for session notes
- Close beads before committing
- Run `bd dolt push` before ending session (if remote configured)

## Beads + Ephemeral Tasks

[If Q5 = yes, include full integration pattern:]
Beads for persistence across sessions. Tasks for within-session execution.
Create a bead for the work item. Use Tasks to break down implementation.
Close the bead when the work ships.

When to create a bead: Would it show up in a changelog? Will you need context
in a future session? Is there a decision worth recording?

## Connecting Plans, Docs, and Commits

- In docs: `<!-- Related: [prefix]-xxx -->` near the top
- In beads: `bd comments add <id> "Doc: path/to/doc.md"`
- In commits: [commit format from Q3]
- In plans: Reference bead IDs the plan addresses
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
Valid labels: `.beads/conventions/labels.md`

### Close Outcome Format

[Template from Q6]

### Session Rules

- Read close_reason before working a bead to avoid re-solving
- Check for abandoned work: `bd list --status=in_progress`
- Use comments for session notes (survives compaction)
- Close beads before committing
- Only add blocking deps when work truly cannot start
- Don't invent labels - use `.beads/conventions/labels.md`

[If Q5 = yes:]
### Beads + Tasks

Beads for cross-session persistence. Tasks for within-session execution.

### Details

- Label taxonomy: `.beads/conventions/labels.md`
- Full reference: `.beads/conventions/reference.md`
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

## Phase 4: Validation

After generating all files:

1. Verify each file was written successfully (read them back).
2. Show summary:

| File | Purpose |
|---|---|
| `.beads/conventions/labels.md` | Label taxonomy - [N] areas + [dimensions] |
| `.beads/conventions/reference.md` | Close template, creation rules, session discipline |
| `AGENTS.md` | Scannable enforcement layer - agents read this every session |
| `CLAUDE.md` | @AGENTS.md reference (if created/updated) |

3. Highlight the label taxonomy (list the area labels).
4. Highlight the close outcome template.
5. Surface any recommendations from Phase 0 doctor results:
   - Team mode not configured → "Consider `bd init --team` for multi-user sync"
   - Editor hooks missing → "Consider `bd setup claude` or `bd setup cursor`"
   - Other doctor warnings → list them
6. Suggest: "Run `bd doctor` to verify full project health."
7. Remind: "These conventions are now active. Agents will follow them for every
   bead created in this project. Review and adjust as needed."

## Important

- Keep generated content concise and scannable.
- The close outcome template and label enforcement are the two highest-value
  outputs. Get these right.
- Don't over-engineer. The system should take less time to set up than it saves.
- Reference detailed conventions from AGENTS.md rather than duplicating inline.
- Write to AGENTS.md (tool-agnostic), not CLAUDE.md.
- One question at a time. Wait for each answer before asking the next.
