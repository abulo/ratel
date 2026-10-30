package upgrade

import (
	"fmt"

	"github.com/abulo/ratel/v3/toolkit/base"
	"github.com/spf13/cobra"
)

// CmdNew 是 toolkit upgrade 子命令，将 toolkit 本体与 protoc 代码生成工具升级到最新版本。
var CmdNew = &cobra.Command{
	Use:   "upgrade",
	Short: "升级脚手架",
	Long:  "升级脚手架命令 : toolkit upgrade",
	Run:   Run,
}

// CmdInit 是 toolkit init 子命令，执行与 CmdNew 相同的工具安装流程，用于脚手架首次初始化。
var CmdInit = &cobra.Command{
	Use:   "init",
	Short: "脚手架初始化",
	Long:  "脚手架初始化 : toolkit init",
	Run:   Run,
}

// Run 通过 go install 安装最新版本的 toolkit、protoc-gen-go、protoc-gen-go-grpc、
// protoc-gen-openapi、go-mod-upgrade 与 protoc-go-inject-tag；安装失败时打印错误，
// 不向调用方返回。
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
