---
name: casting-a-service-layer
description: Use when adding domain logic (types + service interface + implementation) for a resource without HTTP — typically because the route is handled differently (worker, scheduled job, internal-only)
---

# casting-a-service-layer

Emits the domain layer for a resource — types DTO + service interface + storage-backed implementation + table-driven test. No HTTP wiring.

## Prerequisites

- Existing magic service at cwd (cmd/server.go imports magic middlewares).
- `pkg/features/` exists; if not, create it: `mkdir -p pkg/features pkg/types`.

## Inputs

| Input | Type | Example |
|---|---|---|
| `RESOURCE` | lowercase singular | `todo` |

Derived: `Resource`, `RESOURCES`, `Resources` as in `casting-a-feature`.

## Procedure

### 1. Copy templates

```bash
cp "${CLAUDE_PLUGIN_ROOT}/templates/service-layer/pkg/types/RESOURCE.go" \
   "pkg/types/${RESOURCE}.go"
mkdir -p "pkg/features/${RESOURCE}"
cp "${CLAUDE_PLUGIN_ROOT}/templates/service-layer/pkg/features/RESOURCE/service.go" \
   "pkg/features/${RESOURCE}/service.go"
cp "${CLAUDE_PLUGIN_ROOT}/templates/service-layer/pkg/features/RESOURCE/service_test.go" \
   "pkg/features/${RESOURCE}/service_test.go"
```

### 2. Substitute placeholders

```bash
find "pkg/types/${RESOURCE}.go" "pkg/features/${RESOURCE}" -type f -name '*.go' \
  -exec sed -i.bak \
    -e "s/RESOURCETABLE/${RESOURCETABLE}/g" \
    -e "s/Resources/${RESOURCES_CAP}/g" \
    -e "s/RESOURCES/${RESOURCES}/g" \
    -e "s/Resource/${RESOURCE_CAP}/g" \
    -e "s/RESOURCE/${RESOURCE}/g" \
    -e "s/SERVICENAME/${SERVICENAME}/g" {} \;
find pkg -name '*.go.bak' -delete
```

Note: order matters in the sed — replace the longest tokens first (`RESOURCETABLE`, `Resources`) before the shorter prefixes (`Resource`, `RESOURCE`).

### 3. Vet + test

```bash
go vet ./pkg/types/... ./pkg/features/${RESOURCE}/...
go test ./pkg/features/${RESOURCE}/...
```

### 4. Commit

```bash
git add pkg/types/${RESOURCE}.go pkg/features/${RESOURCE}/
git commit -m "feat: add ${RESOURCE} service layer"
```

## Capability boundary

Service layer uses `storage.StorageAdapter` only (no project-specific storage helpers). Lucene filters pass through to `storage.Search` — the adapter parses and evaluates the raw Lucene string. Use `mql.NewParser(input).Parse()` only when you need in-process AST evaluation (validation, transformation, inspection). See `magic-capabilities`.

## Doc references

- `storage.StorageAdapter`: https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/storage#StorageAdapter
- `storage.StorageAdapter.Search` (Lucene passthrough): https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/storage#StorageAdapter
- `mql.NewParser` (in-process AST parsing only): https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/mql#NewParser
