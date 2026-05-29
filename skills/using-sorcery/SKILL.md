---
name: using-sorcery
description: Use when a request involves creating, scaffolding, or extending a microservice built on github.com/tink3rlabs/magic — establishes which cast to invoke and the magic-only capability boundary
---

<SUBAGENT-STOP>
If you were dispatched as a subagent to execute a specific task, skip this skill.
</SUBAGENT-STOP>

<EXTREMELY-IMPORTANT>
If the user's request touches a magic-based microservice — scaffolding, adding a feature, adding a route, swapping storage, wiring observability — you MUST invoke a spellcasting skill BEFORE writing code or asking implementation questions.

IF A CAST APPLIES, YOU DO NOT HAVE A CHOICE. INVOKE IT.
</EXTREMELY-IMPORTANT>

## Instruction Priority

1. **User instructions** (CLAUDE.md, direct requests) — highest
2. **Spellcasting skills** — override default behavior for magic-based work
3. **Superpowers skills** — brainstorming, TDD, debugging still compose
4. **Default system prompt** — lowest

## Capability Boundary

Every cast and every template is hard-bounded to:

**In-bounds:** `github.com/tink3rlabs/magic/{storage, middlewares, observability, health, leadership, mql, errors, logger, pubsub}` · `github.com/tink3rlabs/todo-service` source layout · standard ecosystem (`chi`, `chi/middleware`, `chi/cors`, `chi/render`, `cobra`, `viper`, `slog`).

**Out-of-bounds — NEVER quote, import, or name in templates or generated code:** any project-specific common library (e.g. an org's shared middlewares, attachments, storage helpers, encryption wrappers, config initializers) · any per-resource ACL authorizer beyond what magic ships · any custom rate-limit, audit, cache, or tenant-override middleware not in magic.

If the user's project has such helpers, they're applied *after* the cast completes — never inside the cast.

Exhaustive surface: invoke the `magic-capabilities` skill.

## Cast Catalog

The user invokes casts in natural language ("scaffold a new service", "add a todos resource", "drop a soft-delete on invoices"). You pick the matching skill from this catalog and invoke it via the `Skill` tool. When intent is ambiguous, default to `divining-intent`.

| Skill | Use when | Produces |
|---|---|---|
| `divining-intent` | Non-trivial work; intent is ambiguous; the user said something like "I want to build…" without details | A markdown brief with the cast sequence and answers |
| `casting-a-new-service` | Bootstrapping a brand-new magic-based service | Full service tree: main.go, cmd/, embedded config/, Makefile, Dockerfile, CI |
| `casting-a-feature` | Adding a complete resource end-to-end | Migration → types → service → route, with validation + tests |
| `casting-a-service-layer` | Domain logic only, no HTTP yet | `pkg/types/<name>.go` + `pkg/features/<name>/service.go` + test |
| `casting-a-route` | HTTP layer; auto-chains service layer if missing | `pkg/routes/<name>/{routes.go, validation.go, handler.go, routes_test.go}` + mount |
| `tweaking-a-cast` | Modifying existing code (add field, add filter, add guard, swap adapter) | Targeted diffs respecting DTO-completeness |

## Decision Flow

```dot
digraph cast_flow {
  "Magic-based service work?" [shape=diamond];
  "New service?" [shape=diamond];
  "Full vertical slice?" [shape=diamond];
  "Service layer only?" [shape=diamond];
  "Modifying existing?" [shape=diamond];
  "Respond normally" [shape=box];
  "divining-intent (wizard)" [shape=box];
  "casting-a-new-service" [shape=box];
  "casting-a-feature" [shape=box];
  "casting-a-service-layer" [shape=box];
  "casting-a-route" [shape=box];
  "tweaking-a-cast" [shape=box];

  "Magic-based service work?" -> "Respond normally" [label="no"];
  "Magic-based service work?" -> "New service?" [label="yes"];
  "New service?" -> "casting-a-new-service" [label="yes"];
  "New service?" -> "Modifying existing?" [label="no"];
  "Modifying existing?" -> "tweaking-a-cast" [label="yes"];
  "Modifying existing?" -> "Full vertical slice?" [label="no"];
  "Full vertical slice?" -> "casting-a-feature" [label="yes"];
  "Full vertical slice?" -> "Service layer only?" [label="no"];
  "Service layer only?" -> "casting-a-service-layer" [label="yes"];
  "Service layer only?" -> "casting-a-route" [label="no, route-only"];

  "Magic-based service work?" -> "divining-intent (wizard)" [label="ambiguous / non-trivial"];
}
```

## Red Flags — Stop and Invoke a Cast

| Thought | Reality |
|---|---|
| "User just asked for a route, skip the wizard" | A route pulls types, validation, often service + migration. Invoke `casting-a-route` — it auto-chains. |
| "I'll copy from another service in the user's monorepo since it's right there" | Other services may use project-specific helpers. Out-of-bounds for templates. Use only `magic` + `todo-service` shapes. |
| "I'll just use the user's project rate-limit middleware" | Not in magic. Forbidden inside a cast — if the user wants it, they add it after the cast completes. |
| "Tiny tweak, just edit the file" | DTO-completeness: DB struct → Create DTO → Update DTO → validation → service → migration. Use `tweaking-a-cast`. |
| "Wizard is overkill for an obvious service" | Multi-tenancy, soft-delete, role guards aren't derivable from the description. Wizard exists for that. |
| "Scaffold first, tests after" | Casts ship `*_test.go` in the same pass. No "after". |
| "I'll write a custom auth middleware" | `middlewares.EnsureValidToken` + `RequireRole` are the only options in scope. |
| "Just add the field to the struct, the rest can wait" | Five-step DTO-completeness rule. All steps in one pass. |

## Composition with Superpowers

Spellcasting runs alongside `obra/superpowers`. `superpowers:brainstorming` still applies before broader product decisions; `superpowers:test-driven-development` discipline applies to all generated code; `superpowers:subagent-driven-development` is the right execution mode for a multi-cast sequence emitted by `divining-intent`.

## Skill Index

- `using-sorcery` — this skill, the entrypoint
- `magic-capabilities` — exhaustive magic surface reference
- `divining-intent` — interactive wizard
- `casting-a-new-service` — bootstrap a new service
- `casting-a-feature` — full vertical slice
- `casting-a-service-layer` — service layer only
- `casting-a-route` — route layer, auto-chains service if missing
- `tweaking-a-cast` — modify existing code
