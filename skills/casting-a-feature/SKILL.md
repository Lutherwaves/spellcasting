---
name: casting-a-feature
description: Use when adding a complete resource to an existing magic-based service — emits migration, types, service layer, and route layer in one pass
---

# casting-a-feature

Full vertical slice for a new resource: migration → types → service → routes, plus validation and tests, in a single cast.

## Prerequisites

A magic-based service exists at the cwd (identified by `cmd/server.go` importing `github.com/tink3rlabs/magic/middlewares`). If not, abort and suggest `casting-a-new-service`.

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

Copy and substitute. The migration ID comes **only from the filename prefix** — the YAML body has no `id:` field. The template's real fields are `description:` and `migrations: [{migrate: ..., rollback: ...}]`. The `NNNN` token in the template filename is a placeholder; only the filename matters for ordering.

```bash
cp "${CLAUDE_PLUGIN_ROOT}/templates/migration/config/migrations/postgresql/NNNN__create_RESOURCE.yaml" \
   "config/migrations/postgresql/${NEXT_PAD}__create_${RESOURCETABLE}.yaml"
sed -i.bak \
  -e "s/RESOURCETABLE/${RESOURCETABLE}/g" \
  -e "s/RESOURCE/${RESOURCE}/g" \
  -e "s/SERVICESCHEMA/${SCHEMA}/g" \
  "config/migrations/postgresql/${NEXT_PAD}__create_${RESOURCETABLE}.yaml"
rm "config/migrations/postgresql/${NEXT_PAD}__create_${RESOURCETABLE}.yaml.bak"
```

Note: the `s/NNNN/${NEXT_PAD}/g` substitution is intentionally omitted — `NNNN` does not appear in the YAML body, only in the source template's filename. The destination filename already has the correct prefix.

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

Stay inside `magic/storage`, `magic/middlewares`, `magic/mql`, `magic/errors`. No project-specific common libraries. See `magic-capabilities` for the in-bounds list.

## Doc references

- Migration shape: see `magic/storage` migration runner — https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/storage#NewDatabaseMigration
- Migration YAML schema: `description:` + `migrations: [{migrate: ..., rollback: ...}]`; no `id:` field — ID is the filename prefix only
- Route shape: https://github.com/tink3rlabs/todo-service (`pkg/routes/<resource>/`)
- Cursor pagination via `storage.Search` (Lucene passthrough): https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/storage#StorageAdapter
- In-process Lucene parsing (rarely needed): `mql.NewParser(input).Parse()` — https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/mql#NewParser
