package upgrade

import (
	"fmt"

	"github.com/abulo/ratel/v3/toolkit/base"
	"github.com/spf13/cobra"
)

// CmdUpgrade represents the upgrade command.
var CmdUpgrade = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade the ratel tools",
	Long:  "Upgrade the ratel tools. Example: ratel upgrade",
	Run:   Run,
}

// Run upgrade the ratel tools.
func Run(cmd *cobra.Command, args []string) {
	err := base.GoInstall(
		"github.com/abulo/ratel/v3/toolkit@latest",
		"google.golang.org/protobuf/cmd/protoc-gen-go@latest",
		"google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest",
		"github.com/google/gnostic/cmd/protoc-gen-openapi@latest",
	)
	if err != nil {
		fmt.Println(err)
	}
}
