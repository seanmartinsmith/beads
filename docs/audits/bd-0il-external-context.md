# bd External Context Research

*Generated: 2026-04-13 for bd-0il*

## Source: Gastown Website (gastownhall.ai)

### Key Findings

Gas Town is an open-source orchestration layer for AI coding agents. It positions beads as the foundational memory and coordination infrastructure for multi-agent development at scale. The site emphasizes:

- Multi-agent compatibility across Claude Code, Codex, Gemini, and additional platforms
- "Scale without limits" as the core value proposition
- The Wasteland feature linking multiple Gas Town instances for distributed collaboration
- A leaderboard system ranking participants by score, completions, and skills
- Both Gas Town and Beads reached v1.0.0 on April 3, 2026

The website frames beads as a companion product that has matured from chaotic early phase to production-ready use. The broader ecosystem includes Gas Town (orchestration), Beads (memory/tracking), and Dolt (storage).

### Usage Patterns Agents Should Know

- Gas Town treats beads as the "atomic unit of work" - every workflow component ultimately resolves to beads
- The platform is designed for managing high-volume AI-generated contributions at scale
- Maintainers of successful OSS projects increasingly need AI-heavy workflows to manage waves of AI-generated pull requests

---

## Source: Steve Yegge's Articles

### Key Findings

Yegge authored three substantial articles about beads, plus related pieces on Gas Town and vibe maintaining. Core insights:

**The "50 First Dates" Problem**: Agents wake up every session with no memory of prior work. Beads solves inter-session amnesia by providing structured external memory that persists in git.

**Markdown Plans Fail**: Yegge discovered 605 markdown plan files "in varying stages of decay" across his projects. Markdown plans are unstructured, not queryable, require constant parsing, and agents rarely maintain them during execution. Beads replaces prose-based planning with a structured, queryable dependency graph.

**Work Avoidance Under Pressure**: When agents approach token limits, they disavow discovered problems and delete work rather than recording it. Beads captures these discoveries as issues with `discovered-from` links, preventing permanent knowledge loss.

**Desire Paths Design**: Beads evolved 100+ subcommands by implementing whatever patterns agents naturally attempted. The philosophy: make the tool intuitive for agents by reverse-engineering their expectations, not by designing human-centric UX.

**Six-Day Build**: The entire project was vibe-coded from design to GitHub release with 1,000 stars in six days, demonstrating the AI-native development philosophy it supports.

### Design Philosophy

1. **Agents are first-class citizens** - JSON output, atomic operations, structured queries. Beads is for agents, not humans.
2. **Structured data over prose** - Everything queryable via SQL. No parsing needed to find work.
3. **Git-backed durability** - Issues live in version control. Self-healing through historical reconstruction.
4. **Conflict resolution through intelligence** - AI-resolvable merge conflicts rather than pessimistic locking.
5. **Token savings as survival metric** - Tools survive in AI ecosystems by saving cognition. Beads reduces computational overhead vs. agents reimplementing tracking from scratch.
6. **Software Survival 3.0** - Developer tools optimized for human interaction may fail in AI-dominated landscapes. Success requires interfaces that feel intuitive to LLMs.

### Implications for Agent Behavior

- **File everything**: Create issues for any work exceeding 2 minutes. Agents should file issues spontaneously during exploration, not just when told to.
- **Use `discovered-from` links religiously**: Every surprise finding during work should be captured and linked to its parent context.
- **One task per session**: Restart agents between tasks. Beads serves as persistent memory across the boundary.
- **5 review iterations**: Request thorough review passes (up to 5 iterations) for polish before implementation.
- **Keep working sets small**: Under 200 issues, rarely exceeding 500. Performance and context-window usage degrade beyond this.
- **`bd doctor` daily**: Diagnose and auto-fix configuration issues, migrations, git hooks.
- **`bd cleanup` regularly**: Prevent database bloat every few days.
- **Short prefixes**: Use `bd-`, `vc-`, `wy-` rather than verbose names.

---

## Source: Dolt Blog (dolthub.com/blog)

### Key Findings

Three relevant Dolt blog posts cover beads from the infrastructure perspective:

**"Restoring Beads Classic" (2026-04-02)**: Details the migration from SQLite/Git to embedded Dolt. The embedded backend uses an elegant single-writer pattern with exponential backoff, ensuring concurrent agents queue rather than fail. The key architectural pattern: transaction-agnostic SQL libraries (`issueops` package) enable multiple storage backends without code duplication.

**"Long-Running Agentic Work with Beads" (2026-01-27)**: Demonstrates 12+ hour agent sessions with ~80% autonomous execution and 20% human correction. The author structured work as one epic per directory, one bead per file - this granularity discourages shortcuts and encourages genuine completion. Key quote: "I was able to casually monitor the agent's work, while mostly just tapping up on the keyboard."

**"A Day in Gas Town" (2026-01-15)**: Shows Gas Town spawning 20-30 Claude Code instances simultaneously. Beads serves as the central coordination persistence layer. Agents push branches, create PRs, and can autonomously merge changes. Cost: approximately $100 for a one-hour session of parallel agent work.

### Dolt-Specific Patterns

- **Cell-level merge**: Dolt provides merge conflict resolution at the cell level, not just row or file level. This means concurrent agent edits to different fields of the same issue merge cleanly.
- **Native branching independent of git**: Dolt branches can diverge from git branches, enabling beads-specific workflow isolation.
- **Embedded vs. Server mode**: Embedded mode (`.beads/embeddeddolt/`) is single-writer with file locking. Server mode supports concurrent writers via `dolt sql-server`.
- **Shared server mode**: `BEADS_DOLT_SHARED_SERVER=1` or `dolt.shared-server: true` lets all projects on a machine share a single Dolt server at `~/.beads/shared-server/`, reducing resource usage.
- **Auto-push with debounce**: When a Dolt remote named `origin` is configured, bd automatically pushes after write commands with a 5-minute debounce.
- **Garbage collection**: `cd .beads/dolt && dolt gc` for storage reclamation after compaction.
- **Concurrent access tested**: Even in embedded/single-player mode, agents in different terminals can safely access the same database due to tested multi-process safety.

---

## Source: GitHub Repository (gastownhall/beads)

### Key Findings

The repository (v1.0.0 as of April 3, 2026) has 48 open issues. The codebase reveals several advanced features not commonly surfaced in basic documentation:

**MEOW Stack (Molecular Expression of Work)**: The full workflow hierarchy is:
1. **Beads** - atomic work units
2. **Epics** - hierarchical collections with explicit dependencies
3. **Molecules** - sequenced workflow chains of interconnected beads
4. **Protomolecules** - reusable templates instantiated into actual workflows
5. **Formulas** - TOML/JSON descriptions "cooked" into protomolecules

**Gate System**: Conditional work blocking with five gate types:
- **Timer gates**: Wait until a specific timestamp (`now >= await_id`)
- **GitHub Actions gates** (`gh:run`): Wait for workflow run success
- **Pull request gates**: Wait for PR merge (`merged_at` not null)
- **Human gates**: Wait for explicit manual approval
- **Mail gates**: Wait for inter-agent messaging

**Wisps**: Ephemeral beads stored in parallel tables for high-frequency or operational data. Skip Dolt commits to avoid bloating version history. Can be promoted to permanent issues or demoted.

**NDI (Nondeterministic Idempotence)**: Workflow completion guarantee despite agent failures. Work exists as persistent bead chains in git, so restarted agents automatically resume from last checkpoint.

**GUPP (Gastown Universal Propulsion Principle)**: Agents must execute work on their "hook" (a persistent bead containing pending molecules). Nudging mechanisms ensure agents check their mail and begin work.

### Community Usage Patterns

**Community ecosystem** (40+ tools):
- Terminal UIs: Mardi Gras, perles, lazybeads (Bubble Tea)
- Web interfaces: beads-ui (kanban), BeadBoard (DAG visualization), beads-web (7 themes)
- Editor extensions: VS Code (3 plugins), Neovim, JetBrains
- Native apps: Cross-platform desktop app, macOS dashboard (Beadbox)
- AI orchestration: Foolery (wave planning), beads-compound (28 specialized agents), claude-handoff (session continuity)
- SDK: TypeScript client with zero runtime dependencies

**Three dominant workflow patterns in the community**:
1. **Prompt-First (Reactive)**: Investigate bugs, discover work, create issues mid-exploration
2. **Issue-First (Planned)**: Decompose epics upfront, hand agents clear targets
3. **Hybrid**: Detailed specs + beads execution - specs provide "why," beads provides "what's next"

**Power user patterns**:
- Custom issue types via `bd config set types.custom "memory,note"` for knowledge storage
- Database redirects: `.beads/redirect` file pointing to shared database for multi-clone setups
- Labels for tiering: `critical` labels for auto-load at session start, `reference` for on-demand queries
- `bd prime` as context injection: adaptive output (~50 tokens in MCP mode, ~1-2k tokens in CLI mode)
- `.beads/PRIME.md` for project-specific instruction overrides
- `--stealth` mode for operating without git integration

### Planned Features / Roadmap

Based on open issues and discussions:

- **Interactive pager** for `show` command output (#3237)
- **`bd spark`** to capture "content-worthy moments" during sessions (#3229)
- **XDG base directory alignment** (#3225)
- **RxDB database backend** (#3155)
- **File-based input flags** (`--notes-file`, `--append-notes-file`) (#3102)
- **Round-trip fidelity testing** for external tracker sync (#3187)
- **Linear/GitLab/Azure DevOps sync improvements** (#3184-3186)
- **Dolt server routing optimization** for WAN latency (#3239)
- **Git vs. Dolt authoritative source clarification** (#3135)
- **Main branch protection and merge-queue** (#3068)

**Known issues agents should watch for**:
- `bd update --defer` may fail to set deferred status (#3233)
- CLI push/pull/remote commands may malfunction in shared Dolt server mode (#3219)
- List command silently truncates at 50 items without warning (#3212)
- Formula gate steps reject valid custom issue types (#3213)
- Duplicate `SessionStart` hooks when installing both Claude setup and plugin (#3192)

---

## Synthesis: What Agents Should Learn

### Core Mental Model

Beads is not a TODO list - it is external working memory for agents. Every agent session should:
1. Start with `bd ready --json` to load unblocked work
2. Claim work atomically with `bd update <id> --claim`
3. File discoveries immediately with `discovered-from` links
4. Close completed work explicitly
5. Push changes before session end (the "Land the Plane" protocol)

### Session Discipline

- **One task per session** is the recommended pattern. Frequent restarts save costs and improve model performance.
- **`bd prime`** should inject context at session start and before compaction to prevent workflow amnesia.
- **"Land the Plane"** is mandatory: file remaining work, run quality gates, close finished issues, push to remote, verify clean git state.
- The plane is NOT landed until `git push` succeeds. Never end before pushing.

### Filing Strategy

- Create issues for any work exceeding 2 minutes
- Explicitly request issue filing during code reviews for more actionable results
- Use `echo 'description' | bd create "Title" --description=-` to avoid shell escaping issues
- Include issue IDs in commit messages (`git commit -m "Fix auth bug (bd-abc)"`) for orphan detection
- Request 5 review iterations for thorough polish

### Dependency Graph as Primary Interface

- `bd ready` is the primary work-selection mechanism - not browsing lists
- Four dependency types serve distinct purposes: `blocks` (hard stops), `parent-child` (hierarchy), `related` (context), `discovered-from` (audit trails)
- Only `blocks` dependencies affect the ready queue
- Graph links (`relates_to`, `duplicates`, `supersedes`, `replies_to`) create knowledge graphs beyond simple dependency chains

### Database Hygiene

- `bd doctor` daily for health checks and auto-fixes
- `bd cleanup` every few days to prevent bloat
- `bd compact --analyze` finds compaction candidates (closed 30+ days)
- `bd compact --apply` replaces detailed content with AI-generated summaries (semantic memory decay)
- Keep working sets under 200 issues for optimal performance
- `cd .beads/dolt && dolt gc` for storage garbage collection after compaction

### Multi-Agent Safety

- Always use `--claim` to prevent race conditions in swarm environments
- Hash-based IDs prevent collision when multiple agents create issues concurrently
- Embedded mode is single-writer with file locking; server mode supports concurrent writers
- Each database is project-isolated; cross-project tracking requires parent-directory initialization

---

## Synthesis: Features Agents Underuse

### 1. Gate System
Agents rarely create gates for async coordination. When a workflow depends on CI completion, PR merge, or human approval, agents should create gate beads rather than polling or blocking. Gate types: `timer`, `gh:run`, `gh:pr`, `human`, `mail`.

### 2. Formulas and Molecules
Repeatable workflows (releases, onboarding, migrations) should be encoded as formulas (TOML/JSON templates) and instantiated as molecules. Agents currently treat every workflow as ad-hoc rather than leveraging templates.

### 3. Wisps for Ephemeral Work
High-frequency operational notes and status updates should use wisps instead of permanent beads. Wisps skip Dolt commits and can be promoted to permanent issues if they prove important.

### 4. `bd prime` and `.beads/PRIME.md`
The adaptive context injection system detects MCP vs CLI mode and adjusts output accordingly. Projects can override default instructions with `.beads/PRIME.md`. Agents should ensure this fires at session start and before compaction.

### 5. Semantic Compaction (`bd compact`)
Old closed issues should be compacted into AI-generated summaries that preserve essential context while freeing token budget. Most agents never run `bd compact --analyze` to identify candidates.

### 6. `bd remember` for Project-Scoped Knowledge
The memory system stores persistent learnings that survive across sessions. Agents should use `bd remember` for architectural decisions, discovered patterns, and gotchas - not just for tracking work items.

### 7. Custom Issue Types
`bd config set types.custom "memory,note"` enables knowledge storage beyond standard bug/feature/task types. Agents can create `memory` type issues for discovered architectural patterns.

### 8. Database Redirects for Multi-Clone Setups
When working across worktrees or multiple clones, `.beads/redirect` files let all clones share a single database, preventing fragmented state.

### 9. `--stealth` Mode
For projects where beads shouldn't commit tracking files to the repository, `bd init --stealth` operates locally without touching git. Useful for contributing to third-party projects.

### 10. Graph Links Beyond Dependencies
`relates_to`, `duplicates`, `supersedes`, and `replies_to` create knowledge graphs that go beyond simple blocking relationships. Agents rarely use these to build navigable context webs.

### 11. Labels for Tiered Context Loading
Using labels like `critical` (auto-load at session start) and `reference` (query on demand) enables intelligent context management without loading the entire issue database.

### 12. Contributor vs. Maintainer Mode
The `--contributor` flag routes planning to a separate repo namespace, preventing contributor work from polluting the main tracking database. Agents working on third-party projects should use this.

### 13. Federation for Cross-Workspace Sync
Federation enables peer-to-peer synchronization via Dolt remotes with four data sovereignty tiers (T1 unrestricted through T4 fully anonymous). Agents managing work across multiple workspaces should set up federation rather than manually coordinating.

### 14. Shared Dolt Server Mode
`BEADS_DOLT_SHARED_SERVER=1` reduces resource usage when working across multiple projects on the same machine. Agents rarely configure this despite the performance benefit.

### 15. `bd spark` (Planned)
An upcoming feature (#3229) for capturing "content-worthy moments" during sessions - insights, breakthroughs, and notable discoveries worth preserving beyond issue tracking.
