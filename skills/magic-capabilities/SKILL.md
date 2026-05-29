---
name: magic-capabilities
description: Use when a cast needs to verify whether a package, type, or middleware exists in github.com/tink3rlabs/magic — the in-bounds surface for spellcasting templates
---

# magic-capabilities

The authoritative answer to "is X available in `github.com/tink3rlabs/magic`?". Pinned to **`v0.17.3`**. If your magic is newer, prefer upstream and note the drift in your cast output.

> **NOTE — reconciled against real magic v0.17.3 source on 2026-05-29.** Three major corrections vs. previous content:
> 1. `storage.NewStorageAdapter()` does not exist → use `storage.StorageAdapterFactory{}.GetInstance(type, configMap)`.
> 2. `mql.Parse(string)` does not exist → use `mql.NewParser(input).Parse()`; but service layers should pass raw Lucene to `storage.Search`, not call the parser directly.
> 3. Migration YAML has **no** `id:` field — the numeric ID is encoded in the filename prefix (`0001__create_RESOURCETABLE.yaml`). YAML fields are `description:` + `migrations: [{migrate: ..., rollback: ...}]`.

## Why this skill exists

Every cast template references only what's in this surface. If a generated file imports a package not listed here — or imports anything from a project-specific common library that wraps magic — the cast has failed. Revert and use an in-bounds equivalent.

## Package surface (v0.17.3)

### `github.com/tink3rlabs/magic/storage`

**Adapter type constants** (pass to `StorageAdapterFactory.GetInstance`):

| Constant | Value |
|---|---|
| `storage.MEMORY` | in-process map; no persistence |
| `storage.SQL` | PostgreSQL / MySQL via database/sql |
| `storage.COSMOSDB` | Azure Cosmos DB |
| `storage.DYNAMODB` | AWS DynamoDB |

**Config shape** (`config/default.yaml`):

```yaml
storage:
  type: memory          # one of the constants above (lowercase string value)
  config:               # flat map[string]string passed to GetInstance
    schema: myservice
    host: localhost
    port: "5432"
    # ... adapter-specific keys
```

**Symbols:**

| Symbol | Kind | Signature / Use | Doc |
|---|---|---|---|
| `StorageAdapterFactory` | struct | Factory; use `{}.GetInstance(...)` | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/storage#StorageAdapterFactory |
| `StorageAdapterFactory.GetInstance` | method | `GetInstance(t StorageAdapterType, config map[string]string) (StorageAdapter, error)` — constructs adapter from type constant + flat config map | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/storage#StorageAdapterFactory.GetInstance |
| `StorageAdapterType` | type | `string` alias; cast with `storage.StorageAdapterType(viper.GetString("storage.type"))` | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/storage#StorageAdapterType |
| `StorageAdapter` | interface | Common adapter contract implemented by every backend | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/storage#StorageAdapter |
| `StorageAdapter.List` | method | `List(dest any, sortKey string, filter map[string]any, limit int, cursor string, params ...map[string]any) (string, error)` — returns next cursor | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/storage#StorageAdapter |
| `StorageAdapter.Search` | method | `Search(dest any, sortKey string, query string, limit int, cursor string, params ...map[string]any) (string, error)` — accepts raw Lucene string; adapter parses and evaluates | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/storage#StorageAdapter |
| `StorageAdapter.Get` | method | `Get(dest any, filter map[string]any, params ...map[string]any) error` | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/storage#StorageAdapter |
| `StorageAdapter.Create` | method | `Create(value any, params ...map[string]any) error` | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/storage#StorageAdapter |
| `StorageAdapter.Update` | method | `Update(value any, filter map[string]any, params ...map[string]any) error` | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/storage#StorageAdapter |
| `StorageAdapter.Delete` | method | `Delete(value any, filter map[string]any, params ...map[string]any) error` | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/storage#StorageAdapter |
| `NewDatabaseMigration` | func | `NewDatabaseMigration(adapter StorageAdapter) *DatabaseMigration` — build a migration runner | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/storage#NewDatabaseMigration |
| `ConfigFs` | var (`embed.FS`) | Set by `main.go` to expose embedded `config/` to migrations | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/storage |

**`NewStorageAdapter` does not exist** — this is a common error in older docs. Always use `StorageAdapterFactory{}.GetInstance`.

### `github.com/tink3rlabs/magic/middlewares`

| Symbol | Kind | Use | Doc |
|---|---|---|---|
| `EnsureValidToken` | func | chi middleware: JWT validation via OIDC JWKS | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/middlewares#EnsureValidToken |
| `EnsureValidTokenConfig` | struct | `{Enabled, IssuerURL, Audience, AllowedClockSkew}` | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/middlewares#EnsureValidTokenConfig |
| `SetDefaultClaimsConfig` | func | Configure custom JWT claim keys (tenant, email, roles, groups) | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/middlewares#SetDefaultClaimsConfig |
| `ClaimsConfig` | struct | `{TenantIdKey, EmailKey, RolesKey, GroupsKey}` | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/middlewares#ClaimsConfig |
| `RequireRole` | func | chi middleware factory: 403 unless JWT `roles[]` contains one of `roles...` | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/middlewares#RequireRole |
| `TenantRequestContext` | middleware | Extracts tenant_id from JWT, puts on `context.Context` | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/middlewares#TenantRequestContext |
| `UserRequestContext` | middleware | Extracts email/groups from JWT, puts on context | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/middlewares#UserRequestContext |
| `ErrorHandler` | struct | `Wrap(func(w,r) error)` adapter that maps `magic/errors` types to HTTP status codes | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/middlewares#ErrorHandler |
| `ObservabilityWithOptions` | middleware | Per-request HTTP metrics + trace spans, with `SkipPaths`/`SkipPathPrefixes` | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/middlewares#ObservabilityWithOptions |
| `ObservabilityOptions` | struct | `{SkipPaths, SkipPathPrefixes}` | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/middlewares#ObservabilityOptions |

### `github.com/tink3rlabs/magic/observability`

| Symbol | Kind | Use | Doc |
|---|---|---|---|
| `Observer` | struct | The handle held by `cmd/server.go`, passed to middleware + `/metrics` | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/observability#Observer |
| `Init` | func | `Init(ctx, cfg) (*Observer, error)` — boots Prometheus exporter + OTLP tracer | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/observability#Init |
| `DefaultConfig` | func | Returns a `Config` with sensible defaults; service tunes fields after | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/observability#DefaultConfig |
| `Config` | struct | `{ServiceName, MetricsMode, EnableTracing, TracesOTLPEndpoint, TracesOTLPInsecure}` | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/observability#Config |
| `MetricsModePrometheus` | const | Metrics mode constant; sister value `MetricsModeDisabled` exists | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/observability |
| `Observer.MetricsHandler` | method | `http.Handler` to mount at `/metrics` | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/observability#Observer.MetricsHandler |
| `Observer.Shutdown` | method | Flush spans + exporters on graceful shutdown | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/observability#Observer.Shutdown |

### `github.com/tink3rlabs/magic/health`

| Symbol | Kind | Use | Doc |
|---|---|---|---|
| `NewHealthChecker` | func | Build a checker bound to a `StorageAdapter` | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/health#NewHealthChecker |
| `HealthChecker.Check` | method | `Check(storageEnabled bool, deps []string) error` — used in `/health/readiness` | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/health#HealthChecker.Check |

### `github.com/tink3rlabs/magic/leadership`

| Symbol | Kind | Use | Doc |
|---|---|---|---|
| `LeaderElectionProps` | struct | `{HeartbeatInterval, StorageAdapter, AdditionalProps}` | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/leadership#LeaderElectionProps |
| `NewLeaderElection` | func | Start a single-leader election against the storage adapter | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/leadership#NewLeaderElection |
| `LeaderElection.Start` | method | Begin heartbeat loop in a goroutine | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/leadership |

### `github.com/tink3rlabs/magic/mql`

| Symbol | Kind | Use | Doc |
|---|---|---|---|
| `NewParser` | func | `NewParser(input string) *Parser` — creates a parser for the given Lucene query string | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/mql#NewParser |
| `Parser.Parse` | method | `Parse() (Expr, error)` — parses and returns the AST | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/mql#Parser.Parse |
| AST node types | structs | `Query`, `Clause`, `Term`, `Range`, `Boolean` for predicate construction | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/mql |

**`mql.Parse(string)` does not exist** — use `mql.NewParser(input).Parse()`.

**Preferred pattern:** service layers pass the raw Lucene string directly to `storage.Search(...)` and let the adapter parse and evaluate it. Only call `mql.NewParser` when you need in-process AST evaluation (e.g. to validate, transform, or inspect the query before sending it to storage).

### `github.com/tink3rlabs/magic/errors`

| Symbol | Kind | Use | Doc |
|---|---|---|---|
| `BadRequest` | struct | 400; `{Message}`; returned from handlers, mapped by `ErrorHandler` | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/errors#BadRequest |
| `Unauthorized` | struct | 401 | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/errors#Unauthorized |
| `Forbidden` | struct | 403 | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/errors#Forbidden |
| `NotFound` | struct | 404 | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/errors#NotFound |
| `Conflict` | struct | 409 | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/errors#Conflict |
| `ServiceUnavailable` | struct | 503 | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/errors#ServiceUnavailable |
| `Internal` | struct | 500 fallback | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/errors#Internal |

### `github.com/tink3rlabs/magic/logger`

| Symbol | Kind | Use | Doc |
|---|---|---|---|
| slog handler helpers | funcs | Configure structured logging with trace correlation | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/logger |

### `github.com/tink3rlabs/magic/pubsub`

| Symbol | Kind | Use | Doc |
|---|---|---|---|
| `Publisher`, `Subscriber` | interfaces | Cross-adapter pub/sub contract | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/pubsub |
| adapter constructors | funcs | Backend-specific constructors selected via viper config | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/pubsub |

## todo-service shape (reference layout for `casting-a-new-service`)

| File | Role | URL |
|---|---|---|
| `main.go` | `//go:embed config` + `cmd.Execute()` entry | https://github.com/tink3rlabs/todo-service/blob/main/main.go |
| `cmd/root.go` | cobra root, viper config bind | https://github.com/tink3rlabs/todo-service/blob/main/cmd/root.go |
| `cmd/server.go` | chi router, middleware stack, `/health`, `/metrics`, graceful shutdown | https://github.com/tink3rlabs/todo-service/blob/main/cmd/server.go |
| `cmd/observability.go` | `setupObservability(ctx, name)` | https://github.com/tink3rlabs/todo-service/blob/main/cmd/observability.go |
| `Makefile` | run/test/vet/lint/build targets | https://github.com/tink3rlabs/todo-service/blob/main/Makefile |
| `Dockerfile` | distroless static, nonroot user | https://github.com/tink3rlabs/todo-service/blob/main/Dockerfile |
| `.goreleaser.yaml` | binary builds for tagged releases | https://github.com/tink3rlabs/todo-service/blob/main/.goreleaser.yaml |
| `.github/workflows/ci.yml` | Conventional Commits + golangci + tests | https://github.com/tink3rlabs/todo-service/blob/main/.github/workflows/ci.yml |
| `.github/workflows/release.yml` | go-semantic-release on main | https://github.com/tink3rlabs/todo-service/blob/main/.github/workflows/release.yml |

## Hard boundary — out-of-bounds for every cast

If a generated file imports or names any of the following categories, the cast has failed:

- Any project-specific "common" library (request loggers, rate limiters, audit middleware, cache middleware, default-tenant overrides, admin overrides) that lives outside `github.com/tink3rlabs/magic`
- Any custom storage bootstrap wrapper that hides `storage.StorageAdapterFactory{}.GetInstance(...)` / `storage.NewDatabaseMigration(...).Migrate()`
- Any per-resource ACL authorizer beyond `middlewares.RequireRole`
- Any encryption-at-rest wrapper not provided by magic

In-bounds equivalents for common needs:

| You might want | In-bounds option |
|---|---|
| Structured request logging | chi's `middleware.Logger`, or magic's slog handler |
| Rate limiting | not in v0.1; user adds after the cast |
| Audit logging | not in v0.1; user adds after the cast |
| Response caching | not in v0.1; user adds after the cast |
| Tenant context on every request | `middlewares.TenantRequestContext` |
| Storage bootstrap + migrations | `storage.StorageAdapterFactory{}.GetInstance(type, configMap)` + `storage.NewDatabaseMigration(adapter).Migrate()` |
| Role-based authorization | `middlewares.RequireRole` |
| Per-resource ACL | not in v0.1; user adds after the cast |

## Drift protocol

When magic cuts a new tag:

1. Bump the `Pinned to **v0.17.3**` line at the top of this skill body.
2. Re-fetch each `pkg.go.dev` URL; any 404 → that symbol was renamed or removed; reconcile with upstream.
3. Add any new exported types/funcs that templates would reference.
4. If a magic constructor name changed (e.g. `storage.StorageAdapterFactory{}.GetInstance` → something else), update `templates/service/cmd/server.go` to match.

v0.1 is manually maintained. v0.2+ may auto-generate from `go doc -all`.
