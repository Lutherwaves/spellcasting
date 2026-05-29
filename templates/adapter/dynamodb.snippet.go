// dynamodb.snippet.go — drop the viper keys below into config/default.yaml
// and remove this file; no Go code changes needed beyond config.
//
// config/default.yaml:
//
//   storage:
//     type: dynamodb
//     config:
//       region: us-east-1
//       table: SERVICENAME-RESOURCETABLE
//       # For local dev with DynamoDB Local:
//       endpoint: http://localhost:8000
//       # AWS credentials are picked up from the standard chain
//       # (env vars, ~/.aws/credentials, IAM role).
//
// Note: DynamoDB adapter does not use SQL-style predicates.
// If your service layer uses raw SQL in storage.Query() calls,
// convert them to filter maps or remove them when switching to DynamoDB.
//
// go.mod: no extra direct dependency; magic/storage bundles the AWS SDK.
// Run: go mod tidy

package adapter
