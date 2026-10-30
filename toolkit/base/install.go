//go:build go1.19

package base

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// GoInstall 使用 go install 命令安装指定路径的 Go 工具。
//
// 路径未包含 @ 版本后缀时自动追加 @latest；任一安装失败时返回对应错误。
func GoInstall(path ...string) error {
	for _, p := range path {
		if !strings.Contains(p, "@") {
			p += "@latest"
		}
		fmt.Printf("go install %s\n", p)
		cmd := exec.Command("go", "install", p)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}
