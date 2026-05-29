# spellcasting

> Cast magic-based microservices.

A multi-platform agent plugin (Claude Code / Codex CLI / Cursor) that scaffolds and extends Go microservices built **exclusively on** [`github.com/tink3rlabs/magic`](https://github.com/tink3rlabs/magic) — magic's `storage`, `middlewares`, `observability`, `health`, `mql`, `errors`, `leadership`, and `pubsub` packages are the only runtime surface every cast may reference. Inspired by [`obra/superpowers`](https://github.com/obra/superpowers).

## Install

### Claude Code

```bash
/plugin marketplace add https://github.com/Lutherwaves/spellcasting
/plugin install spellcasting@magic
/reload-plugins
```

The SessionStart hook auto-loads `using-sorcery`.

### Codex CLI

```bash
git clone https://github.com/Lutherwaves/spellcasting.git
cd spellcasting
bash scripts/install-codex.sh
```

This symlinks each skill into `~/.agents/skills/`. Codex picks them up on next launch. Tool mapping: see [`references/codex-tools.md`](references/codex-tools.md).

### Cursor

```bash
git clone https://github.com/Lutherwaves/spellcasting.git
cd /path/to/your/project
bash /path/to/spellcasting/scripts/install-cursor.sh
```

This generates `.cursor/rules/*.mdc` in your project. Trigger a rule in Cursor chat with `@<rule-name>` (e.g. `@using-sorcery`, `@divining-intent`). Tool mapping: see [`references/cursor-tools.md`](references/cursor-tools.md).

## Casts

Just say what you want — the agent picks the right cast. Following [`obra/superpowers`](https://github.com/obra/superpowers)' v5.1.0 lead, spellcasting has no slash commands; intent routes through skill descriptions and `divining-intent`.

| Skill the agent will invoke | Triggered by phrasing like… | Produces |
|---|---|---|
| `divining-intent` (the wizard) | "I want to build a service", anything ambiguous | A casting brief with the cast sequence |
| `casting-a-new-service` | "scaffold a new service", "bootstrap a magic-based microservice" | Full service tree: `main.go`, `cmd/`, embedded `config/`, `Makefile`, `Dockerfile`, CI |
| `casting-a-feature` | "add a `<resource>` end-to-end", "I need a CRUD for X" | Migration → types → service → routes, with validation + tests |
| `casting-a-service-layer` | "add a service for X, no HTTP yet" | `pkg/types/<name>.go` + `pkg/features/<name>/service.go` + test |
| `casting-a-route` | "expose `<resource>` over HTTP" (auto-chains service if missing) | `pkg/routes/<name>/{routes,validation,handler,routes_test}.go` + mount |
| `tweaking-a-cast` | "add a column", "add a Lucene filter", "guard this route with role X", "swap to dynamodb" | Targeted diffs respecting DTO-completeness |

## Capability boundary

Every cast is hard-bounded to the magic + todo-service surface. Casts never introduce project-specific common libraries, custom rate-limit/audit/cache middleware, or per-resource ACL authorizers beyond what magic ships. Anything outside magic is added by you *after* the cast finishes.

For the exhaustive in-bounds surface (every package, type, and function templates may reference), invoke the `magic-capabilities` skill.

## Magic version compatibility

| spellcasting | magic |
|---|---|
| v0.1.x | v0.17.3 |

When magic cuts a new tag and a generated service hits an API mismatch, the relevant cast skill includes doc links to `pkg.go.dev` so an agent can reconcile against the latest signatures.

## Contributing

Trunk-based, Conventional Commits, branch-protected `main` (classic protection **and** modern Rulesets — see [`scripts/setup-rulesets.sh`](scripts/setup-rulesets.sh)). Manifest versions in `package.json`, `.claude-plugin/plugin.json`, and `.claude-plugin/marketplace.json` are bumped manually before tagging via [`scripts/bump-version.sh`](scripts/bump-version.sh) — see [`CONTRIBUTING.md`](CONTRIBUTING.md) § *Cutting a release*.

## TODOs / roadmap

Operational follow-ups deferred from v0.1.x:

- [ ] **Re-enable auto-sync of plugin manifests in `release.yml`** once the repo moves under a GitHub org. GitHub's `Integration` bypass for `github-actions[bot]` is not available on user-owned repos, so the bot's push-back is rejected by the `main-protection` ruleset. Two ways to unblock: (a) move to a `lutherwaves` org and add `actor_type: "Integration", actor_id: 15368` to the ruleset bypass, or (b) add a PAT secret scoped to `contents:write` and switch the workflow's push to use it.
- [ ] **Verify Codex CLI's actual skill-discovery path.** `scripts/install-codex.sh` symlinks into `~/.agents/skills/` on the assumption Codex follows the `agentskills.io` spec. Confirm against the Codex CLI release that ships skills support and adjust if the canonical path differs.
- [ ] **Harden the SessionStart hook's JSON escaping.** `hooks/session-start` currently escapes `\`, `"`, LF, CR, TAB by hand. Replace with `jq -Rs .` to cover the full C0 control-byte range — robustness only; no active bug.
- [ ] **Extend `tests/template-smoke.sh` to render route + service-layer + migration templates** into the rendered service before `go vet`/`go build`. Today only `templates/service/` is smoke-covered.
- [ ] **Dedupe `workflow_run` release runs against rapid `main` pushes.** Today the `concurrency` block serializes but doesn't dedupe by SHA — non-issue until traffic grows.
- [ ] **Backport API drift fixes into `magic-capabilities` automatically.** When magic cuts a new tag, a script that diffs the package surface and surfaces renamed/removed symbols would prevent silent skill drift.

## License

[MIT](LICENSE) — © 2026 Lutherwaves
