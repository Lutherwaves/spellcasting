// mysql.snippet.go — drop the viper keys below into config/default.yaml
// and remove this file; no Go code changes needed beyond config.
//
// config/default.yaml:
//
//   storage:
//     type: sql
//     config:
//       provider: mysql
//       host: localhost
//       port: "3306"
//       dbname: SERVICENAME
//       schema: SERVICESCHEMA
//       user: root
//       password: ""
//       max_open_conns: "25"
//       max_idle_conns: "5"
//       conn_max_lifetime: "5m"
//
// go.mod: no extra direct dependency; magic/storage embeds the mysql driver.
// Run: go mod tidy

package adapter
