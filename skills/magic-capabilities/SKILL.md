---
name: magic-capabilities
description: Use when a cast needs to verify whether a package, type, or middleware exists in github.com/tink3rlabs/magic — the in-bounds surface for spellcasting templates
---

# magic-capabilities

The authoritative answer to "is X available in `github.com/tink3rlabs/magic`?". Pinned to **`v0.17.3`**. If your magic is newer, prefer upstream and note the drift in your cast output.

## Why this skill exists

Every cast template references only what's in this surface. If a generated file imports a package not listed here — or imports anything from `github.com/blox-eng/common/*` or another project-specific helper — the cast has failed. Revert and use an in-bounds equivalent.

## Package surface (v0.17.3)

### `github.com/tink3rlabs/magic/storage`

| Symbol | Kind | Use | Doc |
|---|---|---|---|
| `StorageAdapter` | interface | Common adapter contract (Get/Put/Delete/Query) implemented by every backend | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/storage#StorageAdapter |
| `NewStorageAdapter` | func | Construct an adapter from viper config (`storage.type` selects backend) | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/storage |
| `NewDatabaseMigration` | func | Build a migration runner against a `StorageAdapter` | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/storage |
| `ConfigFs` | var (`embed.FS`) | Set by `main.go` to expose embedded `config/` to migrations | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/storage |

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
| `Parse` | func | Parse a Lucene query string into an AST; service layer feeds AST to storage adapter | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/mql#Parse |
| AST node types | structs | `Query`, `Clause`, `Term`, `Range`, `Boolean` for predicate construction | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/mql |

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

## todo-service shape (reference layout for `/cast:new`)

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
- Any custom storage bootstrap wrapper that hides `storage.NewStorageAdapter()` / `storage.NewDatabaseMigration(...).Migrate()`
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
| Storage bootstrap + migrations | `storage.NewStorageAdapter()` + `storage.NewDatabaseMigration(adapter).Migrate()` |
| Role-based authorization | `middlewares.RequireRole` |
| Per-resource ACL | not in v0.1; user adds after the cast |

## Drift protocol

When magic cuts a new tag:

1. Bump the `Pinned to **v0.17.3**` line at the top of this skill body.
2. Re-fetch each `pkg.go.dev` URL; any 404 → that symbol was renamed or removed; reconcile with upstream.
3. Add any new exported types/funcs that templates would reference.
4. If a magic constructor name changed (e.g. `storage.NewStorageAdapter` → something else), update `templates/service/cmd/server.go` to match.

v0.1 is manually maintained. v0.2+ may auto-generate from `go doc -all`.
