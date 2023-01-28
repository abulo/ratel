package upgrade

import (
	"fmt"

	"github.com/abulo/ratel/v3/toolkit/base"
	"github.com/spf13/cobra"
)

// CmdUpgrade represents the upgrade command.
var CmdNew = &cobra.Command{
	Use:   "upgrade",
	Short: "升级脚手架",
	Long:  "升级脚手架命令 : toolkit upgrade",
	Run:   Run,
}

// CmdUpgrade represents the upgrade command.
var CmdInit = &cobra.Command{
	Use:   "init",
	Short: "脚手架初始化",
	Long:  "脚手架初始化 : toolkit init",
	Run:   Run,
}

// Run upgrade the ratel tools.
func Run(cmd *cobra.Command, args []string) {
	err := base.GoInstall(
		"github.com/abulo/ratel/v3/toolkit@latest",
		"google.golang.org/protobuf/cmd/protoc-gen-go@latest",
		"google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest",
		"github.com/google/gnostic/cmd/protoc-gen-openapi@latest",
		"github.com/oligot/go-mod-upgrade@latest",
		"github.com/syncore/protoc-go-inject-tag@latest",
	)
	if err != nil {
		fmt.Println(err)
	}
}
