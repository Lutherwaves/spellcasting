---
name: casting-a-new-service
description: Use when bootstrapping a brand-new microservice on github.com/tink3rlabs/magic — clones the service template, applies wizard answers, runs a smoke build
---

# casting-a-new-service

Produces a complete, building Go service from `templates/service/` with placeholders replaced and wizard answers applied. The cast ends with `go vet ./... && go build ./...` passing in the target directory.

## Prerequisites

- Target directory is empty, OR a freshly `git init`-ed directory with no Go files.
- If neither, **abort** — tell the user: "Target directory `<path>` already contains files; either choose an empty path or use `/cast:tweak` to modify what's there."

## Inputs

Collected by `divining-intent` if invoked through `/cast`; otherwise ask directly:

| Input | Type | Default |
|---|---|---|
| `SERVICENAME` | kebab-case identifier | (required) |
| `SERVICEPORT` | string (numeric port) | `8080` |
| `SERVICESCHEMA` | snake_case identifier | SERVICENAME with `-` → `_` |
| Storage adapter | `postgres` / `mysql` / `dynamodb` / `cosmos` / `memory` | `memory` |
| Auth | `none` / `jwt` / `jwt+roles` | `jwt` |
| Multi-tenancy | `yes` / `no` | `yes` |
| Observability | `off` / `prom` / `prom+otlp` | `prom+otlp` |
| Leader election | `yes` / `no` | `no` |
| Pub/sub | `yes` / `no` | `no` |

## Procedure

### 1. Clone the template

```bash
cp -r "${CLAUDE_PLUGIN_ROOT}/templates/service/." ./
mv go.mod.tmpl go.mod
mv README.md.tmpl README.md
```

### 2. Substitute placeholders

```bash
find . -type f \( -name '*.go' -o -name '*.yaml' -o -name '*.yml' -o -name 'Makefile' -o -name 'Dockerfile' -o -name 'go.mod' -o -name 'README.md' \) -exec sed -i.bak \
  -e "s/SERVICENAME/${NAME}/g" \
  -e "s/SERVICEPORT/${PORT}/g" \
  -e "s/SERVICESCHEMA/${SCHEMA}/g" {} \;
find . -name '*.bak' -delete
```

### 3. Apply wizard-conditional edits

In `cmd/server.go`:

| Wizard answer | Edit |
|---|---|
| `multi-tenancy: no` | Remove the line `r.Use(middlewares.TenantRequestContext)` |
| `auth: none` | Remove `r.Use(authMiddleware)`, the whole `EnsureValidTokenConfig` block, and `middlewares.SetDefaultClaimsConfig(...)`. Drop the `auth.*` keys from `config/default.yaml` |
| `auth: jwt+roles` | No code change here — the role guards land per-route via `/cast:feature` or `/cast:route` |
| `observability: off` | Replace the `middlewares.ObservabilityWithOptions(...)` call with chi's `middleware.Logger`, remove the `/metrics` handler line, and delete `cmd/observability.go`. Drop the `observability.*` keys from `config/default.yaml` |
| `observability: prom` | Set `observability.tracing.enabled: false` in `config/default.yaml`; keep prom + `/metrics` |
| `observability: prom+otlp` | Default — leave as-is |
| `leader-election: yes` | Insert this in `runServer` before `initRoutes(...)`: `leadership.NewLeaderElection(leadership.LeaderElectionProps{HeartbeatInterval: viper.GetDuration("leadership.heartbeat"), StorageAdapter: storageAdapter, AdditionalProps: map[string]any{}}).Start()`. Add `leadership.heartbeat: 5s` to `config/default.yaml`. Add `"github.com/tink3rlabs/magic/leadership"` to imports |
| `pubsub: yes` | Add `"github.com/tink3rlabs/magic/pubsub"` to imports; add stub initialization in `runServer`: `// TODO: wire publisher/subscriber using magic/pubsub here` (user finishes wiring) |

Storage adapter choice is applied via `config/default.yaml` — `storage.type: <choice>`. For non-`memory` choices, the wizard should have also asked the user for connection details; populate the relevant `storage.<type>.{...}` keys.

### 4. Tidy + smoke-build

```bash
go mod tidy
go vet ./...
go build ./...
```

If any of these fail because the magic API at the pinned version differs from the template:
1. Read the upstream URL from `magic-capabilities` for the symbol that errored.
2. Reconcile the template edit against the upstream signature.
3. Re-run vet + build.
4. If still failing, report DONE_WITH_CONCERNS with the specific reconciliation needed.

### 5. Commit

```bash
git init -q  # only if not already initialized
git add .
git commit -m "feat: bootstrap ${NAME} service via spellcasting"
```

## Verification

```bash
go run . server &
sleep 2
curl -s localhost:${PORT}/health/liveness
# expect: {"status":"ok"}
kill %1
```

## Capability boundary

Out-of-bounds for this cast — if any of these appear in the generated tree, you've drifted; revert the offending file and use the in-bounds equivalent from `magic-capabilities`:

- Any project-specific "common" Go module that wraps magic
- Any custom auth, rate-limit, audit, or cache middleware not provided by magic
- Any per-resource ACL authorizer beyond `middlewares.RequireRole`

In-bounds surface: see `magic-capabilities`.

## Doc references (pinned to magic v0.17.3)

- `storage.NewStorageAdapter`, `storage.NewDatabaseMigration`: https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/storage
- `middlewares.ObservabilityWithOptions`, `EnsureValidToken`, `TenantRequestContext`, `UserRequestContext`, `RequireRole`, `ErrorHandler`: https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/middlewares
- `observability.Init`, `Observer`, `Config`: https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/observability
- `health.NewHealthChecker`: https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/health
- `leadership.NewLeaderElection`: https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/leadership
- Reference shape: https://github.com/tink3rlabs/todo-service/blob/main/cmd/server.go
