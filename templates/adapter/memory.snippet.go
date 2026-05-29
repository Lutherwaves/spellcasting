// memory.snippet.go — drop the viper keys below into config/default.yaml
// and remove this file; no Go code changes needed beyond config.
//
// The memory adapter is backed by an in-process SQLite database.
// It is suitable for local development and tests; it does NOT persist across restarts.
//
// config/default.yaml:
//
//   storage:
//     type: memory
//     config: {}
//
// No connection parameters are required. The adapter creates a fresh SQLite
// instance each time the process starts.
//
// go.mod: no extra direct dependency; magic/storage bundles modernc.org/sqlite.
// Run: go mod tidy

package adapter
