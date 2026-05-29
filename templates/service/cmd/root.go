package cmd

import (
	"bytes"
	"embed"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ConfigFS embed.FS

var rootCmd = &cobra.Command{
	Use:   "SERVICENAME",
	Short: "SERVICENAME service",
}

func Execute() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(serverCommand)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func initConfig() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "default"
	}
	data, err := ConfigFS.ReadFile("config/" + env + ".yaml")
	if err != nil {
		fmt.Fprintln(os.Stderr, "config read error:", err)
		os.Exit(1)
	}
	viper.SetConfigType("yaml")
	if err := viper.ReadConfig(bytes.NewReader(data)); err != nil {
		fmt.Fprintln(os.Stderr, "config parse error:", err)
		os.Exit(1)
	}
	viper.AutomaticEnv()
}
