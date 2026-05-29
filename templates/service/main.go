package main

import (
	"embed"

	"SERVICENAME/cmd"

	"github.com/tink3rlabs/magic/storage"
)

//go:embed config
var configFS embed.FS

func main() {
	storage.ConfigFs = configFS
	cmd.ConfigFS = configFS
	cmd.Execute()
}
