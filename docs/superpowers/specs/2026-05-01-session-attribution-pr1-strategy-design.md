# session attribution: PR1 strategy + reply for #3583

> **Status**: design + audits complete; ready to post reply.
> **Scope**: how we respond on gastownhall/beads#3583, what PR1 commits to architecturally, and what we defer.
> **Companion beads**: `bd-edi` (execution), `bd-xdc` (decision record), `bt-edi` (paired bt-side hygiene followup).

## 🚦 cold-start orientation (read first if you're a new session)

You are picking up mid-handoff between strategy alignment and PR1 implementation planning. The strategy is done. The plan isn't written yet. Here's the minimum context.

**Done:**
1. Issue `#3400` filed; PRs `#3401` + `#3405` pushed on issues-table substrate; closed by sms 2026-04-28 after `#3578` exposed structural defect
2. Discord alignment with @coffeegoddd (dustin) in dolt `#beads` channel 2026-04-28 — events-table direction confirmed
3. GH issue `#3583` filed 2026-04-28 with events-table proposal, both phasing alternatives, @maphew tagged
4. Both maintainers replied 2026-05-01 (maphew 20:58 UTC; coffeegoddd 22:18 UTC)
5. This design doc captures: maintainer cognition map, three-layer scope split, three architectural fork resolutions, draft reply text, bd-edi comment text
6. **Gas City + Gas Town local audits complete 2026-05-03** — both LOW risk, two new constraints surfaced (see "audit findings" section); wasteland transitively covered (no direct bd coupling)

**Pending:**
- Reply on `#3583` not posted (draft below in "draft reply text" section, updated with audit findings)
- Final bd-edi bead comment with post-audit state (text below)
- PR1 implementation plan not written (next session, via `superpowers:writing-plans`)

**Your job (next session):**
1. Read this doc end-to-end before doing anything
2. Confirm reply has been posted (check `gh issue view 3583 --comments`); if not, post the draft from "draft reply text" section
3. Confirm post-audit bd-edi comment is added; if not, add it from the "bd-edi comment to add" section
4. Invoke `superpowers:writing-plans` to convert this design (incl. audit findings) into the executable PR1 plan. Plan must respect:
   - Migration idempotency (`ALTER TABLE ... ADD COLUMN IF NOT EXISTS session VARCHAR(255)` — Gas Town independently CREATEs `wisp_events`)
   - Derived `closed_by_session` view tolerates "closed but no event" (Gas Town reaper bypasses bd.CloseIssue)
   - bd v1.x semver — Gas Town pins `v1.0.0`; PR1 ships as `v1.1.0` (additive struct fields are minor)

**Where everything lives:**
- This doc: `docs/superpowers/specs/2026-05-01-session-attribution-pr1-strategy-design.md`
- Execution bead: `bd show bd-edi`
- Decision bead: `bd show bd-xdc`
- Paired bt-side bead: `bt-edi` (in bt project)
- GH issue (open): https://github.com/gastownhall/beads/issues/3583
- Closed reference PRs: gastownhall/beads `#3401`, `#3405`
- Triggering PR (exposed gap): gastownhall/beads `#3578`
- Discord history: cass-searchable, dolt `#beads` channel, 2026-04-28
- **This session (origin of doc):** `ec49aa31-2ff9-4027-bf73-6c5a6f49b87d` — readable via `cass transcript C:\Users\sms\.claude\projects\C--Users-sms-System-tools-beads\ec49aa31-2ff9-4027-bf73-6c5a6f49b87d.jsonl` for the full strategy-alignment dialog if you need rationale beyond what this doc captures

## context (where we are)

1. `#3400` — original issue (closed)
2. `#3401` + `#3405` — implementation PRs on issues-table substrate; maphew engaged with review; closed by sms 2026-04-28 after #3578 merge exposed the structural defect
3. dolt `#beads` Discord conversation 2026-04-28 — alignment with dustin on direction; he advised filing fresh issue for events-table redesign, looping in maphew
4. `#3583` — filed 2026-04-28 with the events-table proposal, both phasing alternatives surfaced, @maphew tagged
5. **2026-05-01** — both maintainers replied. ball in our court.

operationally, sms is hitting the `#3578`-exposed gap directly while running concurrent agent sessions. session_id capture is load-bearing for sms's work on adjacent tools. this is user-driven hardening of an existing partial implementation, exactly the lane dustin endorsed in the 4/28 Discord ("primarily focused on cleaning up correctness of existing features or hardening features for better agent usage").

## maintainer cognition map (drives reply structure)

- **dustin (coffeegoddd)** — storage/Dolt/Gas-City lane. reads code. owns Gas City compatibility risk. self-described as "not up to speed on the current state of the world there in beads" — defers bd-internal calls to maphew. binding signal: scope conservatism + Gas City impact.
- **maphew** — bd-internal architecture and contract correctness. has a Claude reviewer behind his comments (PR review style on `#3555` makes this evident). binding signal: bd-side multi-writer surface, schema parity, test coverage.
- **sms** — iterating, concurrent sessions, learning. credibility lever: the `#3578` gap isn't theoretical — actively hitting it.

three different review lanes. reply has to write to each appropriately.

## the L8 architectural framing

session attribution has three possible scope layers. confusion in `#3583`'s comments is from not naming them:

| layer | shape | constituency it unblocks |
|---|---|---|
| **1: pure attribution capture** | `events.session VARCHAR(255)` + `wisp_events.session VARCHAR(255)`. no new tables. claim mechanics unchanged. | single user + Gas City agents (multi-agent within one user) |
| **2: structured identity** | `identities` + `sessions` tables, `events.session_id_fk` | teams + Gas City (identity × session is 2D) |
| **3: new claims system** | service vs self-report, lifecycle, assignment mechanics | not what anyone proposed |

dustin's "would break Gas City" reaction is to layer 3 (his words: "a new issue claims system (essentially) would break them"). triggered by issue title's "normalized identity + sessions design" which reads layer-2-ish. **clarification opportunity, not real disagreement.**

dustin's Discord guidance — "we are not looking to extend the surface area of beads anymore" — is a constraint we can echo back. layer 1 is surface stabilization (finishing `closed_by_session`, which is already partially shipped). layer 2 is surface expansion. dustin himself sketched layer 2 in Discord then flagged "this approach is a bit more work also." he gave us permission to defer layer 2 in his own framing.

## audit findings (2026-05-03, both Gas City + Gas Town local audits)

Reports at:
- `.beads/tmp/audit-gascity-bd-coupling-2026-05-03.md`
- `.beads/tmp/audit-gastown-bd-coupling-2026-05-03.md`

**Both LOW risk — PR1 hypothesis "strictly additive" confirmed.**

### Gas City — fully clean

- Zero coupling to `events` / `wisp_events` tables (no SQL, no struct imports, no CLI surface)
- `closed_by_session` is unread anywhere (zero matches across all file types)
- No bd Go-type imports — Gas City has its own `internal/beads/` subset structs; `encoding/json` ignores unknown fields
- No `--session` flag usage today (future opt-in opportunity, not a PR1 requirement)
- Direct SQL surfaces: `internal/api/convoy_sql.go` (issues/dependencies/labels, explicit column projections, no `SELECT *`); `examples/.../reaper.sh` (wisps/mail/config writes); one `SELECT * FROM issues` in `jsonl-export.sh` writes name-keyed JSON (additive-tolerant)
- **PR1 effect on Gas City: none. No code changes required to absorb migration.**

### Gas Town — clean read paths, two new constraints PR1 must respect

- Deep coupling: imports `github.com/steveyegge/beads@v1.0.0` in 30+ files; 206 `bd` CLI invocations; raw SQL via `bd sql` and direct MySQL connection (gt Dolt server, port 3307)
- BUT: zero reads of `*_by_session` columns. Single occurrence (`internal/doltserver/wisps_migrate.go:369`) is a column declaration in Gas Town's own wisps DDL, never read or written
- BUT: zero reads of bd `Event.Actor` from any production code path; the `Event` fields Gas Town actually consumes are `EventType`, `NewValue`, `OldValue`, `CreatedAt`, `ID`, `IssueID`
- BUT: `Storage.CloseIssue(ctx, id, reason, actor, session)` — Gas Town already passes session as 5th arg; signature unchanged in PR1

**Constraint A: migration must be idempotent.** Gas Town independently CREATEs `wisp_events` on its own Dolt server (`internal/doltserver/wisps_migrate.go:444-457`) without a `session` column. If bd's PR1 migration uses bare `ALTER TABLE ... ADD COLUMN`, it fails when the table pre-exists from Gas Town's path. **Required**: `ALTER TABLE ... ADD COLUMN IF NOT EXISTS session VARCHAR(255)` (or column-existence check).

**Constraint B: derived `closed_by_session` view must tolerate "closed but no event".** Gas Town's reaper (`internal/reaper/reaper.go:670-678`) closes stale issues via raw `UPDATE issues SET status='closed', closed_at=NOW(), close_reason=...` — bypasses `bd.CloseIssue`, so no `closed` row lands in `events`. PR1's derived view must return empty/NULL for these issues, not error. Matches today's behavior (current `closed_by_session` column is empty for reaper-closed issues anyway), so it's not a regression — just a constraint the view definition must honor.

**Forward-looking (post-PR1, Gas Town side, not PR1-blocking):**
- Add `session` to wisp_events DDL (one-line change) so fresh installs match
- Add `session` to wisp_events INSERT column lists (data-fidelity for migrated agent beads)
- Update or delete the stale code-comment references at `internal/beads/beads.go:1545,1568` (point to non-existent `decision 009-session-events-architecture.md`)

### Decision doc `009-session-events-architecture.md` — does not exist

Confirmed across Gas Town, Gas City, and beads itself. Two stale code-comment references in Gas Town point at a doc that was never written. The decision content lives in `bd-xdc` (and the strategy spec itself). Optional follow-up: write a real `docs/adr/0003-session-events-architecture.md` in beads post-PR1 and update Gas Town's stale comments.

### Wasteland — transitively covered

`gh search code` against `gastownhall/wasteland` (federation protocol for Gas Towns):
- No `github.com/steveyegge/beads` import
- No `closed_by_session` / `wisp_events` references
- No `FROM events` / `FROM wisps` / `FROM issues` queries
- No `exec.Command("bd"...)` shellouts

Wasteland's Dolt interactions are wasteland-internal (federation transport), not bd-coupling. Any bd schema risk affecting wasteland flows through Gas Town first, which is clean.

### YAGNI opt-out flag — pulled from PR1

The `core.capture-session: false` config flag had zero use case in either Gas City or Gas Town. Audit confirms YAGNI flag was correct. **Removed from PR1 acceptance criteria** below.

---

## PR1 architectural commitments

**Fork 1: multi-writer wiring → distributed.**

- pass `session` parameter through every event writer signature
- no centralization refactor in PR1
- centralization would be feature-driven architectural change → wrong scope under "stabilize existing surface" framing
- if a natural chokepoint emerges during implementation (e.g., 5+ writers all gaining the same parameter via the same intermediate function), surface as discussion before refactoring; do not pre-decide

**Fork 2: closed_by_session disposition → derived from events at read time.**

- single source of truth: events table
- `*_by_session` field names preserved in `bd show --json` output (interface-stable)
- value computed from events at read time (SQL `VIEW` vs in-code JOIN is implementation detail; either works, settle in the plan)
- reject cache approach (denormalize-back-to-issues): reintroduces source-of-truth ambiguity we're trying to escape
- reject drop approach: forces downstream schema-drift work for bt + Gas City + others; cost-asymmetric vs. keep
- JOIN cost on indexed events table is trivial; revisit if a hot-path read need emerges later

**Fork 3: phasing → hard commit to layer 1 in PR1.**

- explicitly defer layer 2 (identities + sessions tables) to dustin's design surface
- frame using his own "this approach is a bit more work" Discord language
- do not keep layer 2 alive in this reply (kept layer 2 alive in `#3583` body and that triggered dustin's misread; correcting now)

**wisp_events parity (maphew's catch).**

- migration covers both `events` and `wisp_events` tables
- schema parity tests on both (column exists, capture behavior matches across tables)
- every call site that records an event updates regardless of target table

> Note: maphew's claim about multiple direct event-insert paths in current main (create/close helpers, labels, comments, transaction comments, rename, promote/demote, bulk ops) is the empirical source for this commitment. We trust it pending our own audit during plan-writing in the next session — that audit produces the canonical call-site list for the plan.

**Fork 4: migration idempotency (audit-derived).**

- migration uses `ALTER TABLE ... ADD COLUMN IF NOT EXISTS session VARCHAR(255)` (or equivalent column-existence pre-check)
- reason: Gas Town independently CREATEs `wisp_events` on its own Dolt server before bd's migration runs against that instance; bare `ADD COLUMN` would fail
- migration must be re-runnable across mixed substrate states
- acceptance test: run migration twice on a Dolt instance pre-seeded with Gas-Town-style `wisp_events` (no `session` column) — second run is a no-op, not an error

**Fork 5: derived view tolerates "closed but no event" (audit-derived).**

- bd's derived `closed_by_session` lookup must return empty/NULL for issues with `closed_at IS NOT NULL` but no corresponding `closed` event row
- reason: Gas Town's reaper (`internal/reaper/reaper.go:670-678`) closes stale issues via raw `UPDATE issues SET status='closed', closed_at=NOW()` — bypasses `bd.CloseIssue`, no `closed` event written
- pre-PR1 behavior: `closed_by_session` column is empty for reaper-closed issues today, so this is preserving existing semantics, not introducing new behavior
- acceptance test: insert issue with `status='closed', closed_at=NOW()` directly (no event); `bd show --json` returns `closed_by_session: ""` or null, no error

## what's in PR1

- migration: nullable `session VARCHAR(255)` on `events` + `wisp_events`, **idempotent** (`IF NOT EXISTS`)
- distributed wiring through every call site that records an event
- session resolution at cmd/bd layer: `--session > BEADS_SESSION_ID > CLAUDE_SESSION_ID > empty`
- `closed_by_session` preserved with events as source of truth; legacy field names computed at read time; `bd show --json` field shape bit-identical
- derived view tolerates "closed but no event" (returns empty/NULL, doesn't error)
- tests:
  - every event type captures session
  - both `events` and `wisp_events` parity
  - previously-gapped `bd ready --claim` path captures session
  - env precedence: `--session > BEADS_SESSION_ID > CLAUDE_SESSION_ID > empty`
  - migration re-run on Gas-Town-style `wisp_events` (no `session` column) is a no-op
  - issue with `closed_at` but no closed event yields empty `closed_by_session` (no error)
- docs: CLAUDE.md cross-project provenance section updated to reflect substrate change while keeping interface-stable field names
- semver: ships as `v1.1.0` (additive struct fields are minor); Gas Town pins `v1.0.0`, picks up cleanly via `go get -u`

> **Migration safety note**: sms operates against a populated shared Dolt server. Plan-writing session should include a pre-merge migration probe step (run on a snapshot, confirm no Dolt-specific drift class issues; cf. `bd-y5f` historical drift incident). Adding a nullable column with `IF NOT EXISTS` is the safe class of migration, but Dolt has surprised us before.

## explicit non-goals for PR1

- `identities` table
- `sessions` table
- normalized FK from events to sessions
- dropping `closed_by_session` column (separate later decision)
- `bd events <id>` read surface
- `bd list --session` filter / `bd stats --group-by=session`
- backfilling pre-existing events with session
- cross-connector session translation (Cursor/Codex/Cline → BEADS_SESSION_ID)
- ~~opt-out flag (`core.capture-session: false`)~~ — pulled per audit findings (no use case in either Gas City or Gas Town); env-absence already satisfies opt-out
- Gas Town `wisp_events` DDL update / INSERT-list updates — separate Gas-Town-side PR after PR1 merges
- ADR file `docs/adr/0003-session-events-architecture.md` — optional follow-up to formalize the decision content (currently in `bd-xdc`)

each is its own concern in its own PR or issue.

## reply structure (ready to post)

four moves, in order:

1. **opener** — re-anchor in Discord framing. PR1 stabilizes existing surface; not extending bd; consistent with dustin's stated scope.
2. **for maphew** — concede multi-writer surface (commit to distributed). concede wisp_events parity. concede view-not-cache. ack his prior `#3401`/`#3405` engagement.
3. **for dustin** — clarify layer split using his own "more work" framing. promise additive-only schema. invite Gas City migration review. defer layer 2 to his design surface.
4. **convergent close** — invite scope refinement on PR1 specifics here before code lands. layer 2 stays a separate conversation under his lead.

## draft reply text for #3583

```
thanks both for the detailed reads.

quick re-anchor before specifics: PR1 is about stabilizing an existing
surface, not extending one — closed_by_session has been on the issues
table since b362b3682, but only one of three lifecycle events captures
session today. #3578 (bd ready --claim silently un-attributed) helped
me identify closed_by_session was on the issues table and my Discord
conversation with @coffeegoddd helped me shape the scope. but, this
proposal is trying to correct existing architecture, not introduce any
new features.

before getting into specifics: i ran local audits on Gas City and Gas
Town to verify the "strictly additive" claim against real downstream
code, plus a code-search probe on wasteland (federation protocol).
both Gas City and Gas Town came back LOW risk; wasteland has no
direct bd coupling. detailed findings + new constraints below.

@maphew — thanks for the detailed read + the prior review on
#3401/#3405. "single wire-up at RecordFullEventInTable" was sloppy
framing on my part; the multi-writer surface is real. concrete
answers:

- multi-writer surface: distributed wiring. `session` threaded
  through every call site that records an event; no centralization
  refactor in PR1 — that's a feature-driven architectural change
  that contradicts the "stabilize existing surface" framing. if a
  natural chokepoint emerges during implementation (5+ writers
  gaining the same parameter via the same intermediate function),
  i'll surface it for discussion before refactoring.
- wisp_events parity: in acceptance criteria. migration covers both
  tables, schema parity tests on both, every writer family updates
  regardless of target table.
- closed_by_session disposition: keep-as-view over cache. cache
  reintroduces the source-of-truth ambiguity we're trying to
  escape; view keeps events as the single source, JOIN cost trivial
  on an indexed events table. bd show --json output bit-identical.

@coffeegoddd — the "normalized identity + sessions design" phrasing
in the issue title oversold the proposal. PR1 is much smaller:

- add nullable `session VARCHAR(255)` to events and wisp_events
- no new tables
- no claim-mechanics changes — claim semantics from #3578 are
  unchanged; we're only adding capture of WHO claimed in WHICH
  session, parallel to the existing actor capture

the identities + sessions tables you sketched in Discord (4/28 4:47
PM) is the right layer-2 destination, but you flagged it there as
"this approach is a bit more work." that benefits from your direct
involvement on identity lifecycle semantics + Gas City alignment
rather than rolling into PR1. happy to file a separate issue under
your lead when the timing's right.

Gas City audit (gastownhall/gascity@HEAD, local static audit):

- zero direct coupling to events/wisp_events (no SQL, no struct
  imports, no CLI surface)
- closed_by_session unread across all file types
- no bd Go-type imports — Gas City uses its own internal/beads/
  subset structs; encoding/json ignores unknown fields, so additive
  bd output flows through unchanged
- direct SQL touches issues/dependencies/labels in
  `internal/api/convoy_sql.go` (explicit column projections, no
  SELECT *) and wisps/mail/config in `examples/.../reaper.sh` —
  all unaffected by column additions on events/wisp_events

PR1 effect on Gas City: none.

Gas Town audit (gastownhall/gastown@HEAD, local static audit):

deeper coupling than Gas City — imports the beads SDK in 30+ files,
206 bd CLI invocations, raw SQL via `bd sql` and a direct MySQL
connection to its own Dolt server. read paths are clean (no
`*_by_session` reads, no `Event.Actor` reads from production code),
but the audit surfaced two real constraints PR1 must respect:

1. **migration must be idempotent.** Gas Town independently CREATEs
   wisp_events on its own Dolt server (port 3307) without a
   `session` column (`internal/doltserver/wisps_migrate.go:444-457`).
   PR1's migration uses `ALTER TABLE … ADD COLUMN IF NOT EXISTS
   session VARCHAR(255)` so it's a no-op when the column already
   exists.
2. **derived closed_by_session tolerates "closed but no event".**
   Gas Town's reaper (`internal/reaper/reaper.go:670-678`) closes
   stale issues via raw `UPDATE issues SET status='closed',
   closed_at=NOW()`, bypassing bd.CloseIssue — no closed event
   written. the derived view returns empty/NULL for these, matching
   today's behavior (the column is empty for reaper-closed issues
   currently anyway). not a regression, just a constraint the view
   has to honor.

both are in PR1's acceptance criteria with explicit tests.

post-PR1 forward-looking (Gas Town side, separate PR, not
PR1-blocking): add `session` to Gas Town's wisp_events DDL and
INSERT column lists in `wisps_migrate.go` so fresh installs match
and agent-bead → wisp migration preserves session attribution. one-
or two-line changes; happy to open that as a Gas Town PR if useful.

`Storage.CloseIssue` signature: Gas Town already passes session as
the 5th positional arg today. PR1 doesn't change the signature —
only what bd does internally with the value (writes events.session
instead of issues.closed_by_session). Gas Town behavior unchanged.

semver: PR1 ships as v1.1.0 (additive struct fields are minor under
Go module semantics). Gas Town pins v1.0.0 in go.mod, picks up
cleanly via `go get -u`.

happy to share migration SQL, view definition, and full audit
reports if useful. otherwise: if PR1 scope works for both of you,
i'll open a clean feature branch and start implementation. layer 2
stays a separate conversation under your lead, dustin, when the
timing's right.
```

tone: lowercase casual, first person singular, contractions natural, conversational. matches sms's existing voice on `#3509` and other PRs.

## bd-edi comment to add (durable session capture)

```
## Maintainer replies on #3583 (2026-05-01)

Both replies in 2026-05-01:
- @maphew (20:58 UTC): events-table direction +1; simple-first +1.
  Concrete engineering signals — multi-writer surface (not single
  wire-up at RecordFullEventInTable), wisp_events parity, lean
  keep-as-compat-cache/view on closed_by_session.
- @coffeegoddd (22:18 UTC): Gas City compatibility primary concern.
  Cautious on identities/sessions normalized model — wants identity
  + session lifecycle questions worked out before committing.
  Reading layer-2 in our pitch (issue title's "normalized identity
  + sessions design" likely triggered it).

## Decision: hard commit to layer 1 / simple-first in PR1

Three forks resolved this session. Full rationale in design doc:
`docs/superpowers/specs/2026-05-01-session-attribution-pr1-strategy-design.md`

1. Multi-writer wiring → distributed. Pass `session` through every
   call site that records an event; no centralization refactor in PR1.
2. closed_by_session disposition → events as source of truth, legacy
   field names computed at read time. Reject cache (source-of-truth
   ambiguity) and drop (asymmetric downstream cost).
3. Phasing → hard layer 1 commit. Layer 2 (identities + sessions
   tables) deferred to dustin's design surface, separate issue, his
   timing. Use his Discord "this approach is a bit more work" framing
   to defer cleanly.

wisp_events parity (maphew's catch) added to acceptance criteria.
Migration + tests cover both events and wisp_events.

## Stable as of this session

- Layer 1 = PR1 scope.
- closed_by_session stays (events as source of truth). Bead's earlier
  DROP lean reversed on cost-asymmetry grounds.
- wisp_events parity required.
- Distributed wiring; centralization deferred until organic
  justification appears (threshold: 5+ writers gaining the same
  parameter via the same intermediate function).

## Open items

- **Opt-out flag (`core.capture-session: false`)**: doc-review
  flagged YAGNI. Default action absent new info: pull from PR1
  scope. Gas City audit may surface real use case.
- **Gas City compatibility audit**: sms cloning + auditing locally
  in parallel. Audit findings land as a follow-up comment here.
- **PR1 implementation plan**: written next session via
  `superpowers:writing-plans`, after Gas City audit.

## Origin session

Strategy was aligned in CC session
`ec49aa31-2ff9-4027-bf73-6c5a6f49b87d` on 2026-05-01. Full transcript
readable via `cass transcript` if rationale beyond the design doc is
needed.
```

## audit history (executed 2026-05-03)

sms cloned `gastownhall/gascity` and `gastownhall/gastown` and ran static coupling audits in fresh CC sessions in each workspace. Reports at `.beads/tmp/audit-{gascity,gastown}-bd-coupling-2026-05-03.md`. Wasteland (`gastownhall/wasteland`) probed via `gh search code` — no direct bd coupling, transitively covered via Gas Town.

Audit prompts (for reference / future reuse): `.beads/tmp/audit-prompt-{gascity,gastown}.md`.

Outcome: both LOW risk. Two new constraints surfaced (migration idempotency, derived-view tolerance for reaper-bypass) — folded into Forks 4 and 5 above.

## next-session entry point

When a fresh session opens to write the PR1 plan:

1. **Read this doc end-to-end** — especially the cold-start orientation at top and the audit findings section.
2. **Confirm reply has been posted** on `#3583` (check `gh issue view 3583 --comments`); if not, post the draft from "draft reply text" section below.
3. **Confirm the post-audit bd-edi comment is added** (check `bd show bd-edi` for a 2026-05-03 comment); if not, add it from the "bd-edi comment to add" section.
4. **Invoke `superpowers:writing-plans`** with this design doc as input. Plan must respect:
   - All five forks above (1–5)
   - Migration idempotency (`IF NOT EXISTS`) — Fork 4
   - Derived-view tolerance for closed-no-event — Fork 5
   - bd v1.x semver: ship as `v1.1.0`
   - Pre-merge migration probe step (cf. `bd-y5f` drift class)
5. The audit findings section is canonical for "what we know about downstream coupling." Don't re-audit; reference it.

## settled questions (resolved this session)

- **Opt-out config flag (`core.capture-session: false`)**: pulled from PR1 scope. Audit found zero use case in Gas City or Gas Town; env-absence already satisfies opt-out semantics.
- **SQL `VIEW` vs in-code JOIN** for `closed_by_session` derivation: implementation detail, settle in plan. Both work; either is fine. Plan-writer's call.
- **Centralization of event insertion**: deferred unless implementation reveals 5+ writers gaining the same parameter via the same intermediate function.
- **Decision doc `009-session-events-architecture.md`**: confirmed never written. Decision content lives in `bd-xdc`. Optional follow-up: write a real `docs/adr/0003-session-events-architecture.md` post-PR1.

## acceptance criteria (this design)

- [x] design doc committed to repo
- [x] user reviews + approves design (rounds 1 + 2 incl. audit findings)
- [x] cold-start handoff blocks added (top + audit history + next-session entry point)
- [x] Gas City + Gas Town audits complete; wasteland transitively covered
- [x] bd-edi comment captured for durability (added 2026-05-01; updated 2026-05-03 with post-audit state)
- [ ] reply text posted to `#3583` (this session — see "draft reply text" below)
- [ ] PR1 implementation plan written via `superpowers:writing-plans` (next session)
