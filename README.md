# spellcasting

> Cast magic-based microservices.

A Claude Code plugin that scaffolds and extends Go microservices built on [`github.com/tink3rlabs/magic`](https://github.com/tink3rlabs/magic). Inspired by [`obra/superpowers`](https://github.com/obra/superpowers).

## Install

```bash
/plugin install https://github.com/Lutherwaves/spellcasting
```

Once installed, `using-sorcery` is auto-loaded on session start — your agent learns the cast catalog and capability boundary immediately.

## Casts

| Command | When to use it | Produces |
|---|---|---|
| `/cast` | Non-trivial work — let the wizard map intent to casts | A casting brief with the cast sequence |
| `/cast:new <name>` | Bootstrapping a brand-new magic-based service | Full service tree: `main.go`, `cmd/`, embedded `config/`, `Makefile`, `Dockerfile`, CI |
| `/cast:feature <name>` | Adding a complete resource end-to-end | Migration → types → service → routes, with validation + tests |
| `/cast:service <name>` | Domain logic only, no HTTP | `pkg/types/<name>.go` + `pkg/features/<name>/service.go` + test |
| `/cast:route <name>` | HTTP layer (auto-chains service layer if missing) | `pkg/routes/<name>/{routes,validation,handler,routes_test}.go` + mount |
| `/cast:tweak` | Modifying existing code (add field, add filter, add guard, swap adapter) | Targeted diffs respecting DTO-completeness |

## Capability boundary

Every cast is hard-bounded to the magic + todo-service surface. Casts never introduce project-specific common libraries, custom rate-limit/audit/cache middleware, or per-resource ACL authorizers beyond what magic ships. Anything outside magic is added by you *after* the cast finishes.

For the exhaustive in-bounds surface (every package, type, and function templates may reference), invoke the `magic-capabilities` skill.

## Magic version compatibility

| spellcasting | magic |
|---|---|
| v0.1.x | v0.17.3 |

When magic cuts a new tag and a generated service hits an API mismatch, the relevant cast skill includes doc links to `pkg.go.dev` so an agent can reconcile against the latest signatures.

## Contributing

Trunk-based, Conventional Commits, branch-protected `main`. See [`CONTRIBUTING.md`](CONTRIBUTING.md).

## License

[MIT](LICENSE) — © 2026 Lutherwaves
