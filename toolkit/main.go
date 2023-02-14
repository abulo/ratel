package main

import (
	"log"

	"github.com/abulo/ratel/core/env"
	"github.com/abulo/ratel/toolkit/api"
	"github.com/abulo/ratel/toolkit/build"
	"github.com/abulo/ratel/toolkit/dao"
	"github.com/abulo/ratel/toolkit/module"
	"github.com/abulo/ratel/toolkit/project"
	"github.com/abulo/ratel/toolkit/upgrade"
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
	rootCmd.AddCommand(project.CmdNew)
	rootCmd.AddCommand(dao.CmdNew)
	rootCmd.AddCommand(module.CmdNew)
	rootCmd.AddCommand(build.CmdNew)
	rootCmd.AddCommand(upgrade.CmdNew)
	rootCmd.AddCommand(upgrade.CmdInit)
	rootCmd.AddCommand(api.CmdNew)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
