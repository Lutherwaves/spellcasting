# Template Provenance

Each template file in this directory mirrors a canonical upstream source. The table below records those sources, pinned to the versions current at template-creation time.

| Template file | Upstream source |
|---|---|
| `route/pkg/routes/RESOURCE/routes.go` | https://github.com/tink3rlabs/todo-service/blob/main/pkg/routes/todo.go (router wiring pattern) |
| `route/pkg/routes/RESOURCE/handler.go` | https://github.com/tink3rlabs/todo-service/blob/main/pkg/routes/todo.go (handler method shape) |
| `route/pkg/routes/RESOURCE/validation.go` | Derived from todo-service inline `createSchema` / `replaceSchema` JSON Schema maps |
| `route/pkg/routes/RESOURCE/routes_test.go` | Original; uses magic in-memory adapter (`storage.MEMORY`) via `StorageAdapterFactory` |
| `service-layer/pkg/types/RESOURCE.go` | https://github.com/tink3rlabs/todo-service/blob/main/pkg/types/todo.go — extended with `tenant_id`, `created_at`, `updated_at` |
| `service-layer/pkg/features/RESOURCE/service.go` | https://github.com/tink3rlabs/todo-service/blob/main/pkg/features/todo/todoService.go — adapted to 5-method CRUD interface with context |
| `service-layer/pkg/features/RESOURCE/service_test.go` | Original; mirrors todoService_test.go patterns with in-memory adapter |
| `migration/config/migrations/postgresql/0001__create_RESOURCETABLE.yaml` | https://github.com/tink3rlabs/todo-service/blob/main/config/migrations/postgresql/01__base.yaml — extended with `tenant_id`, timestamps, index |
| `adapter/postgres.snippet.go` | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/storage — `SQL` adapter + `POSTGRESQL` provider |
| `adapter/mysql.snippet.go` | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/storage — `SQL` adapter + `MYSQL` provider |
| `adapter/dynamodb.snippet.go` | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/storage — `DYNAMODB` adapter |
| `adapter/cosmos.snippet.go` | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/storage — `COSMOSDB` adapter |
| `adapter/memory.snippet.go` | https://pkg.go.dev/github.com/tink3rlabs/magic@v0.17.3/storage — `MEMORY` adapter |

**Magic version pinned:** `github.com/tink3rlabs/magic@v0.17.3`
**todo-service ref pinned:** `main` branch (commit at template-creation time)

## Drift protocol

When magic releases a new minor or major version:

1. Read the magic changelog or release notes at https://github.com/tink3rlabs/magic/releases.
2. Re-fetch the pkg.go.dev pages for `storage`, `mql`, `errors`, and `middlewares` at the new version tag.
3. Diff the `StorageAdapter` interface methods against what the service and route templates call. Pay special attention to:
   - Method renames (`List` / `Search` / `Get` / `Create` / `Update` / `Delete`)
   - Signature changes (new required params, context moves, return type changes)
   - New `ContextualStorageAdapter` promotions (if the context-aware variant becomes the default)
4. Re-fetch `https://raw.githubusercontent.com/tink3rlabs/todo-service/main/pkg/features/todo/todoService.go` and compare against `service-layer/pkg/features/RESOURCE/service.go`.
5. Re-fetch `https://raw.githubusercontent.com/tink3rlabs/todo-service/main/pkg/routes/todo.go` and compare against `route/pkg/routes/RESOURCE/handler.go` and `routes.go`.
6. Update the pinned version string in this file and in `magic-capabilities/SKILL.md` (see that file's "Version update" section).
7. Commit the updated templates with message: `chore: sync templates to magic@vX.Y.Z`.

**Rule:** If any method signature in a template does not match the installed version of magic, the generated service will not compile. Template drift is a build failure, not a runtime failure — catch it at codegen time.
