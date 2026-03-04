<img src="https://r2cdn.perplexity.ai/pplx-full-logo-primary-dark%402x.png" style="height:64px;margin-right:32px"/>

# how do people currently use beads priority system steve yegge are people letting agents decipher this, is there guidance in beads docs, whats the consensus on this and has anyone actually developed methodology against it. is there a sop for other task a management or just swe in general that provide a reference guide? what about across different coding languages and projects

Beads is Steve Yegge's open-source, git-native issue tracker designed specifically as a "cognitive upgrade for coding agents" — it solves the "50 First Dates" problem where agents wake up with zero memory of prior sessions. The priority system is P1–P4, and the whole point is that `bd ready` handles priority resolution automatically so neither you nor the agent needs to manually decipher ordering.[^1][^2]

## How the Priority System Actually Works

Beads uses four priority levels in practice:[^3]

- **P1** — Highest, actively blocking or release-critical work
- **P2** — Important, next in line after P1s clear
- **P3** — Nice-to-have, not urgent
- **P4** — Backlog (Yegge says use sparingly; the tool is optimized for near-term work, not a 500-item someday list)

The `bd ready` command does the heavy lifting — it queries the dependency graph and returns only *unblocked*, prioritized tasks at that moment. This means the agent doesn't need to reason about ordering; the DAG computes it structurally. Mixing P4 backlog noise in will degrade `bd ready` results and waste context tokens.[^2][^4]

## Are People Letting Agents Decipher Priority?

The consensus is **no — and that's by design**. The whole point of the graph-based model is that explicit dependency edges (`blocks`, `parent-child`, `related`, `discovered-from`) replace prose-based inference. Agents should not be parsing intent from markdown — they run `bd ready --json` and get back structured, machine-consumable prioritized work. Letting an agent "figure out" priority from flat files is exactly the failure mode Beads was built to escape.[^4][^2][^3]

However, a real-world gap exists: **Claude and other agents don't proactively use Beads on their own**. You still have to prompt "check bd ready" at session start. The tool provides memory infrastructure; you trigger its use.[^2]

## Official Docs and Guidance

The repo has two key instruction files that serve as the de facto SOP:[^5][^6]

- **`AGENTS.md`** — General agent instructions for any coding agent
- **`AGENT_INSTRUCTIONS.md`** — Detailed operational instructions for agents working on Beads development itself

The community-accepted onboarding pattern for Claude Code is:

```bash
bd init
bd setup claude
echo "Use 'bd' for task tracking. Run 'bd ready --json' to find work." >> CLAUDE.md
```


## Established Community Methodologies

Three patterns have emerged as the community consensus:[^3][^2]


| Pattern | Use Case | How It Works |
| :-- | :-- | :-- |
| **Prompt-First** | Ad-hoc bugs, unknown scope | Describe problem → AI investigates → creates issues as it discovers them |
| **Issue-First** | Planned features, tracked sprints | Create epic + tasks first → hand off to agent with a clear target |
| **Hybrid (Specs + Beads)** | Complex features needing architectural intent | Write markdown spec → agent converts to Beads epics/tasks → execute via `bd ready` |

The hybrid is gaining the most traction for serious projects — the spec provides "the why," Beads provides "what's next". The "Land the Plane" protocol (tests → file discovered work → close issues → sync → push) has also become a standard session-end SOP.[^2][^3]

## Cross-Language and Multi-Project Use

Beads is intentionally language-agnostic — it's just a CLI + SQLite + JSONL that lives in `.beads/` alongside any codebase. The community has already built integrations in Go, Rust, TypeScript, Python, Lua, and Swift. Notable cross-language tooling includes:[^1][^2]

- **`nvim-beads`** (Lua/Neovim), **`vscode-beads`** (TypeScript), **`beads.el`** (Emacs) for editor integrations
- **Rust port** by @doodlestein that simplifies agentic workflows[^7]
- **`beads-orchestration`** for Claude Code — an orchestrator that delegates to language/tech-specific supervisor agents on isolated git branches[^2]

For multi-project use, the recommended pattern is git worktrees — each agent gets its own worktree + branch, Beads syncs via git, and agents work independently then merge normally. Yegge's own upcoming **Gas Town** orchestration layer adds "convoys" (task bundles per agent) and "rigs" (workspace configs) to formalize this further.[^2]

## The Honest Gaps vs. Jira/Linear

Yegge is explicit: *"Beads isn't a planning tool, a PRD generator, or Jira. It's orchestration for what you're working on today and this week."*  If your team lives in Jira or Linear as the source of truth, Beads has two-way Jira sync scripts but isn't meant to replace org-wide project management. The recommended split: use Beads for active-sprint work, keep distant backlog in your existing tracker, and promote items to a Beads epic when they move to "now."[^8][^3]
<span style="display:none">[^10][^11][^12][^13][^14][^15][^9]</span>

<div align="center">⁂</div>

[^1]: https://www.linkedin.com/posts/steveyegge_github-steveyeggebeads-beads-a-memory-activity-7383408928665042944-tkcj

[^2]: https://ianbull.com/posts/beads/

[^3]: https://betterstack.com/community/guides/ai/beads-issue-tracker-ai-agents/

[^4]: https://dev.to/koustubh/building-apps-with-ai-how-beads-changed-my-development-workflow-2p7

[^5]: https://github.com/steveyegge/beads/blob/main/AGENT_INSTRUCTIONS.md

[^6]: https://github.com/steveyegge/beads/blob/main/AGENTS.md

[^7]: https://www.linkedin.com/posts/steveyegge_github-steveyeggebeads-beads-a-memory-activity-7418745364104622080-a4II

[^8]: https://paddo.dev/blog/beads-memory-for-coding-agents/

[^9]: https://github.com/steveyegge/beads/discussions/1044

[^10]: https://x.com/Steve_Yegge/status/1993919120170209498

[^11]: https://www.youtube.com/watch?v=EsFa7W-FYdM

[^12]: https://www.youtube.com/watch?v=s96O9oWI_tI

[^13]: https://dev.to/ruarfff/a-gaggle-of-agents-5f9

[^14]: https://www.linkedin.com/posts/ghislainbourgin_beyond-instructions-how-beads-lets-ai-agents-activity-7420416284951056384-I-wU

[^15]: https://virtuslab.com/blog/ai/beads-give-ai-memory/

