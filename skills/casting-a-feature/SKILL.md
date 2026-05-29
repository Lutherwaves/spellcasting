---
name: casting-a-feature
description: Use when adding a complete resource to an existing magic-based service — emits migration, types, service layer, and route layer in one pass
---

# casting-a-feature

Full vertical slice for a new resource: migration → types → service → routes, plus validation and tests, in a single cast.

## Prerequisites

A magic-based service exists at the cwd (identified by `cmd/server.go` importing `github.com/tink3rlabs/magic/middlewares`). If not, abort and suggest `/cast:new`.

## Inputs

| Input | Type | Example |
|---|---|---|
| `RESOURCE` | lowercase singular | `todo` |

Derived tokens used in template substitution:
- `Resource` — capitalized singular: `Todo`
- `RESOURCES` — lowercase plural: `todos`
- `Resources` — capitalized plural: `Todos`
- `RESOURCETABLE` — snake_case plural: `todos`

## Procedure

### 1. Emit migration

Pick the next migration id:

```bash
NEXT=$(ls config/migrations/postgresql/ 2>/dev/null | sed -E 's/^([0-9]+)__.*/\1/' | sort -n | tail -1)
NEXT=$(( ${NEXT:-0} + 1 ))
NEXT_PAD=$(printf "%04d" $NEXT)
```

Copy and substitute:

```bash
cp "${CLAUDE_PLUGIN_ROOT}/templates/migration/config/migrations/postgresql/NNNN__create_RESOURCE.yaml" \
   "config/migrations/postgresql/${NEXT_PAD}__create_${RESOURCETABLE}.yaml"
sed -i.bak \
  -e "s/NNNN/${NEXT_PAD}/g" \
  -e "s/RESOURCETABLE/${RESOURCETABLE}/g" \
  -e "s/RESOURCE/${RESOURCE}/g" \
  -e "s/SERVICESCHEMA/${SCHEMA}/g" \
  "config/migrations/postgresql/${NEXT_PAD}__create_${RESOURCETABLE}.yaml"
rm "config/migrations/postgresql/${NEXT_PAD}__create_${RESOURCETABLE}.yaml.bak"
```

### 2. Emit service layer

Cascade into `casting-a-service-layer` (invoke via `Skill` tool) to produce `pkg/types/RESOURCE.go` and `pkg/features/RESOURCE/`.

### 3. Emit route layer

Cascade into `casting-a-route` to produce `pkg/routes/RESOURCE/` and append the mount line to `cmd/server.go`.

### 4. Smoke + test

```bash
go vet ./...
go test ./pkg/routes/${RESOURCE}/... ./pkg/features/${RESOURCE}/...
```

### 5. Commit

```bash
git add config/migrations/postgresql config/openapi.json pkg/ cmd/server.go
git commit -m "feat: add ${RESOURCE} resource end-to-end"
```

## Capability boundary

Stay inside `magic/storage`, `magic/middlewares`, `magic/mql`, `magic/errors`. No `blox-eng/common/*`. See `magic-capabilities` for the in-bounds list.

## Doc references

- Migration shape: see `magic/storage` migration runner — https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/storage
- Route shape: https://github.com/tink3rlabs/todo-service (`pkg/routes/<resource>/`)
- Cursor pagination + Lucene filter: `magic/mql.Parse` — https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/mql
