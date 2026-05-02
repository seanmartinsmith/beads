# session attribution: PR1 strategy + reply for #3583

> **Status**: design complete; paused pending Gas City compatibility audit before reply post + plan writing.
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

**Pending:**
- Reply on `#3583` not posted (draft below in "draft reply text" section)
- bd-edi bead comment not added (text below in "bd-edi comment to add" section)
- Gas City clone + compatibility audit (sms doing in parallel — see "parallel work" section)
- PR1 implementation plan not written (depends on Gas City audit outcome)

**Your job:**
1. Read this doc end-to-end before doing anything
2. Check Gas City audit findings — look for bd-edi comments dated after 2026-05-01, or `bd show bd-edi` for status
3. If audit clean (no surprising coupling): invoke `superpowers:writing-plans` to convert this design into the executable PR1 plan, then post the draft reply, then add the bd-edi comment, then proceed
4. If audit surfaced real coupling: revise the "PR1 architectural commitments" section below in light of findings, get user re-approval, *then* write the plan

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

## what's in PR1

- migration: nullable `session VARCHAR(255)` on `events` + `wisp_events`
- distributed wiring through every call site that records an event
- session resolution at cmd/bd layer: `--session > BEADS_SESSION_ID > CLAUDE_SESSION_ID > empty`
- opt-out via `core.capture-session: false` config flag (default true) — *flagged YAGNI in doc-review; defer decision pending Gas City audit, see "open questions deferred" below*
- `closed_by_session` preserved with events as source of truth; legacy field names computed at read time; `bd show --json` field shape bit-identical
- tests: every event type captures session; both tables parity; previously-gapped `bd ready --claim` covered; env precedence; opt-out (if flag stays in scope)
- docs: CLAUDE.md cross-project provenance section updated to reflect substrate change while keeping interface-stable field names

> **Migration safety note**: sms operates against a populated shared Dolt server. Plan-writing session should include a pre-merge migration probe step (run on a snapshot, confirm no Dolt-specific drift class issues; cf. `bd-y5f` historical drift incident). Adding a nullable column is the safe class of migration, but Dolt has surprised us before.

## explicit non-goals for PR1

- `identities` table
- `sessions` table
- normalized FK from events to sessions
- dropping `closed_by_session` column (separate later decision)
- `bd events <id>` read surface
- `bd list --session` filter / `bd stats --group-by=session`
- backfilling pre-existing events with session
- cross-connector session translation (Cursor/Codex/Cline → BEADS_SESSION_ID)

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

quick re-anchor before specifics: the Discord conversation that shaped
this issue (dolt #beads, 4/28) was explicitly framed around dustin's
"we are not looking to extend the surface area of beads anymore"
guidance. PR1 here is about stabilizing an existing surface, not
extending one — closed_by_session has been on the issues table since
b362b3682 / decision 009-session-events-architecture.md, but only one
of three lifecycle events captures session today. that's what #3578
made painful (bd ready --claim silently un-attributed) and what this
proposal is trying to finish, not introduce.

@maphew — appreciate the detailed read, and thanks for the prior
review on #3401/#3405. you're right that "single wire-up at
RecordFullEventInTable" was sloppy framing on my part. concrete
answers:

- multi-writer surface: distributed wiring. pass `session` through
  every event writer signature; no centralization refactor in PR1.
  centralizing event insertion would be a feature-driven architectural
  change, which contradicts the "stabilize existing surface" framing.
  if a natural chokepoint emerges during implementation i'll surface
  it for discussion before refactoring rather than pre-deciding.
- wisp_events parity: caught, adding to acceptance criteria. migration
  covers both tables, schema parity tests for both, every event writer
  family updates regardless of target table.
- closed_by_session disposition: lean keep-as-view (computed from
  events at read time) over cache (denormalize back to issues). cache
  reintroduces the source-of-truth ambiguity we're trying to escape;
  the view approach keeps events as the single source. JOIN cost on an
  indexed events table is trivial. fine to revisit if a hot-path read
  need shows up later, but PR1 ships with events as the source of
  truth and the legacy field names preserved by read-time
  computation.

@coffeegoddd — re-reading my own issue title, i think the "normalized
identity + sessions design" phrasing landed wrong. the actual PR1
proposal is much smaller than that suggests:

- add nullable `session VARCHAR(255)` to events and wisp_events
- no new tables
- no claim-mechanics changes (your "new issue claims system" concern —
  claim semantics from #3578 are unchanged; we're only adding capture
  of WHO claimed in WHICH session, parallel to the existing actor
  capture)

the identities + sessions tables you sketched in Discord (4/28
4:47 PM) is real and likely the right layer-2 destination, but you
flagged it there as "this approach is a bit more work" — i agree, and
i think it benefits from your direct involvement on identity lifecycle
semantics + Gas City sign-off rather than rolling into PR1. happy to
file a separate issue for that under your lead when the timing's
right.

on Gas City compatibility specifically — PR1 is strictly additive:
column additions on events/wisp_events, no schema removals,
closed_by_session preserved as JOIN-view so the issues-table read
shape stays bit-identical from outside. existing Gas City queries
against issues or events tables continue to work unchanged. happy to
share the migration SQL + view definition before merge for your
review against Gas City's coupling.

if PR1 scope as outlined works for both of you i'll open a clean
feature branch and start implementation. if either of you want me to
refine specifics (multi-writer call sites, view shape, migration
safety details) before code lands, happy to iterate here first. layer
2 stays a separate conversation under your lead, dustin, when the
time's right.
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

## parallel work in progress (sms-side, not this session)

sms is cloning Gas City to audit its actual coupling against bd's events/issues schema. Goals:

1. **Verify "additive only" claim is true.** Run the proposed migration locally against a Gas City instance; confirm no Gas City reads break. Specifically test:
   - any Gas City SQL that does `SELECT *` against `events`, `wisp_events`, or `issues`
   - any positional column reads (rare but possible)
   - any code paths that reflect on schema columns dynamically
2. **Surface real coupling we don't know about.** Anything Gas City depends on that we haven't accounted for becomes a new acceptance criterion or non-goal.
3. **Dogfood our own builds.** Empirical evidence for the reply: "ran the migration locally against Gas City and these queries still work" is much stronger than "trust us, additive."
4. **Side benefit.** sms has been meaning to try Gas City anyway — this side quest pays double.

Output: a comment on `bd-edi` with audit findings, dated after 2026-05-01. The next session reads that comment to know if PR1 commitments hold, need revision, or are blocked.

## next-session entry point

When a fresh session opens to write the PR1 plan:

1. **Read this doc end-to-end.** Especially the cold-start orientation at top.
2. **Check Gas City audit status.** `bd show bd-edi` — look for comments dated after 2026-05-01 mentioning Gas City. If none yet, ask sms whether to proceed without (post reply with "additive promise + invite review" framing) or wait.
3. **If audit clean:**
   - Post the draft reply text from this doc (verbatim or with sms's edits) to `#3583`
   - Add the bd-edi comment text from this doc to `bd-edi`
   - Invoke `superpowers:writing-plans` with this design doc as input. Produce the executable PR1 implementation plan: migration files, canonical call-site enumeration (do the audit maphew's claim implies), view definition, test coverage, documentation updates.
4. **If audit surfaced real coupling:**
   - Pause. Do not post reply. Do not write plan.
   - Revise this doc's "PR1 architectural commitments" section in light of findings.
   - Get sms re-approval.
   - Then proceed to step 3.

## open questions deferred to next session

- **Opt-out config flag (`core.capture-session: false`)**: doc-review flagged as YAGNI. Pulling it tightens the proposal and avoids contradicting our "no surface expansion" framing. Kept in scope pending Gas City audit — if GC audit surfaces a real config-flag use case, it stays. If not, pull it from PR1 scope and the bd-edi acceptance criteria. **Default action absent new info: pull it.**
- **SQL `VIEW` vs in-code JOIN** for `closed_by_session` derivation: implementation detail, settled in plan.
- **Centralization of event insertion**: deferred unless implementation reveals a natural chokepoint per the threshold above.

## acceptance criteria (this design)

- [x] design doc committed to repo
- [x] user reviews + approves design
- [x] cold-start handoff blocks added (top + parallel work + next-session entry point)
- [ ] bd-edi comment captured for durability (after this section, executed in this session)
- [ ] reply text posted to `#3583` (deferred to next session, post-GC audit)
- [ ] PR1 implementation plan written via `superpowers:writing-plans` (deferred to next session)
