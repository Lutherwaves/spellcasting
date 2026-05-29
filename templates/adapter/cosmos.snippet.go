// cosmos.snippet.go — drop the viper keys below into config/default.yaml
// and remove this file; no Go code changes needed beyond config.
//
// config/default.yaml:
//
//   storage:
//     type: cosmosdb
//     config:
//       endpoint: https://ACCOUNTNAME.documents.azure.com:443/
//       key: ""           # primary or secondary key
//       database: SERVICENAME
//       container: RESOURCETABLE
//       partition_key: /tenant_id
//
// Note: Cosmos DB adapter does not use SQL-style predicates.
// If your service layer uses raw SQL in storage.Query() calls,
// convert them to filter maps or remove them when switching to CosmosDB.
//
// go.mod: no extra direct dependency; magic/storage bundles the Azure SDK.
// Run: go mod tidy

package adapter
