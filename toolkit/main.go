package main

import (
	"log"

	"github.com/abulo/ratel/v3/core/env"
	"github.com/abulo/ratel/v3/toolkit/api"
	"github.com/abulo/ratel/v3/toolkit/backend"
	"github.com/abulo/ratel/v3/toolkit/build"
	"github.com/abulo/ratel/v3/toolkit/dao"
	"github.com/abulo/ratel/v3/toolkit/frontend"
	"github.com/abulo/ratel/v3/toolkit/module"
	"github.com/abulo/ratel/v3/toolkit/upgrade"
	"github.com/abulo/ratel/v3/toolkit/vue"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "toolkit",
	Short:   "toolkit: An elegant toolkit for Go microservices.",
	Long:    `toolkit: An elegant toolkit for Go microservices.`,
	Version: env.RatelVersion(),
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd:   true,
		DisableNoDescFlag:   true,
		DisableDescriptions: true,
		HiddenDefaultCmd:    true,
	},
	DisableFlagsInUseLine: true,
}

func init() {
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})
	rootCmd.AddCommand(frontend.CmdNew)
	rootCmd.AddCommand(backend.CmdNew)
	rootCmd.AddCommand(dao.CmdNew)
	rootCmd.AddCommand(module.CmdNew)
	rootCmd.AddCommand(build.CmdNew)
	rootCmd.AddCommand(upgrade.CmdNew)
	rootCmd.AddCommand(upgrade.CmdInit)
	rootCmd.AddCommand(api.CmdNew)
	rootCmd.AddCommand(vue.Vue)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
