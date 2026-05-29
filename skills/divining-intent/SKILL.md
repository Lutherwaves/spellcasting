---
name: divining-intent
description: Use when the user asks to scaffold or extend a magic-based service and the right cast sequence is ambiguous — runs a question tree and emits a cast plan
---

# divining-intent

The interactive wizard. Asks one question at a time, builds a brief, and ends with a recommended cast sequence the user can run. Constrained to magic + todo-service capabilities (see `magic-capabilities`).

## When to invoke

- The user said something like "I want to build a new service" without specifying details
- Multiple casts could apply and you can't tell which
- The user's description leaves multi-tenancy, soft-delete, auth, or storage choice ambiguous

## How to ask

Ask **one question at a time**. Prefer `AskUserQuestion` for multiple choice; fall back to plain text for free-form answers. Don't batch — the user's later answers may depend on yours.

## Question tree

### Q1. New or extending?

A) New service from scratch → continue with **new-service questions** (Q2 + Q3–Q10)
B) Extending an existing one → continue with **extension questions** (Q2 + Q3'–Q6')

### Q2. Domain in one sentence?

Free text. Used in commit messages, README, and the casting brief.

### New-service path

**Q3. Multi-tenancy?**
- A) Single-tenant (skip `middlewares.TenantRequestContext`)
- B) Multi-tenant (wire `TenantRequestContext` + tenant-scoped storage queries; user supplies RLS policy)

**Q4. Auth?**
- A) None (disable `EnsureValidToken`)
- B) JWT only (`middlewares.EnsureValidToken` enabled)
- C) JWT + role guards (per-subroute `middlewares.RequireRole`)

**Q5. Storage adapter?**
- A) postgres
- B) mysql
- C) dynamodb
- D) cosmos
- E) memory (dev-only, no persistence)

**Q6. Soft delete?**
- A) Yes — `deleted_at` column + filter helper on List/Get
- B) No — hard `DELETE`

**Q7. List endpoint shape?**
- A) Cursor pagination only (`?limit=10&next=...`)
- B) + Lucene `?filter=` via `mql.Parse`
- C) + sort param

**Q8. Observability?**
- A) Off
- B) Prometheus only (`/metrics`)
- C) Prometheus + OTLP traces (recommended for production)

**Q9. Leader-elected background work?**
- A) Yes (`leadership.NewLeaderElection`)
- B) No

**Q10. Pub/sub?**
- A) Yes (`magic/pubsub`)
- B) No

### Extension path

**Q3'. Full vertical slice (migration + types + service + route)?**
- A) Yes → `casting-a-feature`
- B) No → continue to Q4'

**Q4'. Service layer only (types + features pkg, no HTTP)?**
- A) Yes → `casting-a-service-layer`
- B) No → continue to Q5'

**Q5'. Route layer only?**
- A) Yes → `casting-a-route` (auto-chains `casting-a-service-layer` if missing)
- B) No → continue to Q6'

**Q6'. Modifying existing code (add field, filter, guard, soft-delete, swap adapter)?**
- A) Yes → `tweaking-a-cast`
- B) No → ask the user what they actually want; this skill doesn't fit

## Output format

When the answers are collected, emit a markdown brief in this exact shape:

```markdown
## Casting Brief
<one-sentence summary built from the domain answer + key choices>

## Answers
- Domain: <Q2>
- Multi-tenancy: <Q3 A or B>
- Auth: <Q4 A | B | C>
- Storage: <Q5 A | B | C | D | E>
- Soft delete: <Q6>
- List shape: <Q7>
- Observability: <Q8>
- Leader election: <Q9>
- Pub/sub: <Q10>

## Cast Sequence
1. `casting-a-new-service` for `<service-name>` (multi-tenancy=<...>, auth=<...>, storage=<...>, observability=<...>, leader-election=<...>, pubsub=<...>)
2. `casting-a-feature` for `<first-resource>`
3. <additional casts as needed>
```

For the extension path, the brief shows only the single chosen cast and any auto-chained prerequisites.

## After emitting the brief

Tell the user: "Ready to run this sequence? (yes / edit answers / cancel)"

- **yes** → invoke the first cast in the sequence via the `Skill` tool, passing the relevant answers as context. After it returns, proceed to the next cast.
- **edit answers** → re-ask only the question the user wants to change, then re-emit the brief.
- **cancel** → stop. Don't run anything.

## Capability boundary

Every cast you orchestrate is hard-bounded to the magic + todo-service surface. If the user's domain requires something not in `magic-capabilities` (e.g. their own rate-limit middleware, an internal "common" helper package), tell them: "That's outside spellcasting's v0.1 scope — you'll need to add it after the cast finishes." Do not invent imports.

Cross-reference: `magic-capabilities` for the exhaustive in-bounds surface.
