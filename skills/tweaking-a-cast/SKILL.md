---
name: tweaking-a-cast
description: Use when modifying an existing resource on a magic-based service — add a column, add a filter, add a role guard, add soft-delete, swap a storage adapter — enforces DTO-completeness on every field change
---

# tweaking-a-cast

Targeted modifications to a resource that already exists in a magic-based service. Routes by the user's intent. Every field-touching change goes through the **DTO-completeness rule**.

## Prerequisites

- Existing magic service at cwd.
- The file being modified imports `github.com/tink3rlabs/magic/*`. If not, abort with: "Target file doesn't look like a magic-based service — out of scope."

## The DTO-completeness rule (mandatory)

When you add a column to a database struct, you MUST complete ALL of these in the same change:

1. Add the field to the DB struct (in `pkg/types/<resource>.go`).
2. Add to the **Create DTO** (`CreateResource`) if user-settable on creation.
3. Add to the **Update DTO** (`UpdateResource`) if user-modifiable.
4. Add to the validation schema in `pkg/routes/<resource>/validation.go`.
5. Wire into the service layer (`Create` and `Update` functions in `pkg/features/<resource>/service.go`).
6. Add a migration: `config/migrations/postgresql/NNNN__add_COLUMN_to_RESOURCETABLE.yaml` with `ALTER TABLE ... ADD COLUMN IF NOT EXISTS ...`.

For **date-only** fields (OpenAPI `format: date`), use a `DateOnly`-equivalent type — NOT `time.Time`, which marshals to RFC3339 with a time component.

If any step is skipped, the field exists in the DB but is unreachable from the API. That's the failure mode this rule prevents.

## Sub-procedures

### Add a column

Run the six DTO-completeness steps above. Smoke:

```bash
go vet ./...
go test ./pkg/features/${RESOURCE}/... ./pkg/routes/${RESOURCE}/...
```

Commit:

```bash
git add pkg/ config/migrations
git commit -m "feat: add ${COLUMN} to ${RESOURCE}"
```

### Add a Lucene filter clause

Edit `pkg/features/${RESOURCE}/service.go` — the `List` method. Pass the user-supplied Lucene string directly to `storage.Search(...)` — the adapter parses and evaluates it. Do **not** call `mql.NewParser` in the service layer; only use it if you need in-process AST evaluation (validation, transformation, etc.). Add a test case to `service_test.go` exercising the new clause. Commit:

```bash
git commit -am "feat: support ${CLAUSE} filter on ${RESOURCE} list"
```

### Add a role guard to a subroute

Edit `pkg/routes/${RESOURCE}/routes.go`. Wrap the protected subset:

```go
r.Group(func(r chi.Router) {
    r.Use(middlewares.RequireRole("${ROLE}"))
    // mount or define the guarded handlers here
})
```

Add a test asserting `403 Forbidden` for a request whose JWT lacks `${ROLE}`. Commit:

```bash
git commit -am "feat: require ${ROLE} for ${RESOURCE} ${OP}"
```

### Add soft-delete

1. Apply DTO-completeness for `deleted_at *time.Time` (DB struct only — not in Create/Update DTOs, not user-settable).
2. In `pkg/features/${RESOURCE}/service.go`:
   - `List` and `Get`: add `WHERE deleted_at IS NULL` to the storage query.
   - `Delete`: replace the hard delete with `UPDATE ... SET deleted_at = now() WHERE id = $1`.
3. Migration: `ALTER TABLE ... ADD COLUMN IF NOT EXISTS deleted_at timestamptz`.
4. Add a test asserting `List` excludes soft-deleted rows.

Commit: `feat: add soft-delete to ${RESOURCE}`.

### Swap storage adapter

1. Read current adapter from `viper.GetString("storage.type")` in `config/default.yaml`.
2. Copy the new adapter snippet from `templates/adapter/${NEW}.snippet.go` (if helper code is needed) into `cmd/` or `pkg/features/<resource>/`.
3. Update the viper config keys — use the real config shape: `storage.type: ${NEW}` (one of `memory`, `sql`, `cosmosdb`, `dynamodb`) plus connection details as a flat map under `storage.config.*` (e.g. `storage.config.host`, `storage.config.port`, `storage.config.schema`). **Do not** use the old `storage.<type>.host` nested style — `GetInstance` takes `viper.GetStringMapString("storage.config")`.
4. `go mod tidy && go vet ./... && go build ./...`.
5. If the new adapter is `dynamodb` or `cosmos`, double-check that the storage queries don't use SQL-specific syntax — `magic/storage` adapter contracts abstract this, but custom predicates may need adaptation.

Commit: `feat: swap storage adapter to ${NEW}`.

## Bail conditions

- The target file imports nothing from `magic/*` → abort, tell the user.
- Multiple resources need the same change → run this cast per resource (don't batch).
- The user asks for a change outside these five sub-procedures → tell them: "That's outside `tweaking-a-cast`'s sub-procedure set; either describe it as one of {add column, add filter, add guard, add soft-delete, swap adapter} or implement it directly."

## Capability boundary

Tweaks stay inside the in-bounds surface. If a tweak would require importing a project-specific common library or middleware, that's a sign the modification is out of scope for spellcasting — point the user at `magic-capabilities` to find an in-bounds equivalent or punt the change for them to make by hand.

## Doc references

- `storage.StorageAdapter.Search` (Lucene passthrough, canonical): https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/storage#StorageAdapter
- `mql.NewParser` (in-process AST parsing only): https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/mql#NewParser
- `storage.StorageAdapterFactory.GetInstance` (adapter construction): https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/storage#StorageAdapterFactory.GetInstance
- `middlewares.RequireRole`: https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/middlewares#RequireRole
- Migration shape (reference): https://github.com/tink3rlabs/todo-service/tree/main/config/migrations/postgresql
