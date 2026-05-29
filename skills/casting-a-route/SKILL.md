---
name: casting-a-route
description: Use when adding the HTTP layer for a resource that already has (or will get) a service layer — chains casting-a-service-layer first if the service is missing
---

# casting-a-route

Emits `pkg/routes/<resource>/` with five handlers (List, Get, Create, Update, Delete), validation, table-driven test, and appends the mount line to `cmd/server.go`. Auto-chains the service layer if missing.

## Prerequisites

- Existing magic service at cwd.
- A service layer for `RESOURCE` exists at `pkg/features/${RESOURCE}/service.go`.

If the service layer is missing AND the user did not pass `--custom`:

1. Invoke `casting-a-service-layer` via the `Skill` tool, passing the same `RESOURCE` argument.
2. After it completes, resume from step 1 of this skill's procedure.

If `--custom` was passed (user said e.g. "just the routes, I'll wire it myself"), skip the chain and proceed — leave the service-interface import unresolved and instruct the user to wire it.

## Inputs

`RESOURCE` (lowercase singular). Same token derivations as `casting-a-feature`.

## Procedure

### 1. Copy templates

```bash
mkdir -p "pkg/routes/${RESOURCE}"
cp "${CLAUDE_PLUGIN_ROOT}/templates/route/pkg/routes/RESOURCE/"*.go \
   "pkg/routes/${RESOURCE}/"
```

### 2. Substitute placeholders

```bash
find "pkg/routes/${RESOURCE}" -type f -name '*.go' \
  -exec sed -i.bak \
    -e "s/RESOURCETABLE/${RESOURCETABLE}/g" \
    -e "s/Resources/${RESOURCES_CAP}/g" \
    -e "s/RESOURCES/${RESOURCES}/g" \
    -e "s/Resource/${RESOURCE_CAP}/g" \
    -e "s/RESOURCE/${RESOURCE}/g" \
    -e "s/SERVICENAME/${SERVICENAME}/g" {} \;
find pkg/routes -name '*.go.bak' -delete
```

### 3. Mount in cmd/server.go

Find the block `router.Route("/", func(r chi.Router) {` and append, after the existing `r.Mount` lines (or as the first one if none exist):

```go
import_line='	"SERVICENAME/pkg/features/RESOURCE"'
import_line2='	RESOURCEroutes "SERVICENAME/pkg/routes/RESOURCE"'
mount_line='		r.Mount("/RESOURCES", RESOURCEroutes.NewResourceRouter(RESOURCE.NewService(storageAdapter)).Router)'
```

The agent appends these. Use `goimports` after the edit to fix import ordering, or `gofmt` if `goimports` isn't available.

### 4. Vet + test

```bash
go vet ./pkg/routes/${RESOURCE}/...
go test ./pkg/routes/${RESOURCE}/...
```

### 5. Commit

```bash
git add pkg/routes/${RESOURCE}/ cmd/server.go
git commit -m "feat: add ${RESOURCE} route layer"
```

## Capability boundary

Routes use `chi`, `chi/render`, `magic/middlewares.ErrorHandler`, `magic/errors`, and (for guarded subroutes only) `magic/middlewares.RequireRole`. No `blox-eng/common/*`, no gate `Guard`. See `magic-capabilities`.

## Doc references

- `chi.Router`: https://pkg.go.dev/github.com/go-chi/chi/v5
- `middlewares.ErrorHandler`: https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/middlewares#ErrorHandler
- `middlewares.RequireRole`: https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/middlewares#RequireRole
- Reference route shape: https://github.com/tink3rlabs/todo-service (`pkg/routes/`)
