//go:build !go1.19

package base

import (
	"fmt"
	"os"
	"os/exec"
)

// GoInstall 使用 go get -u 安装指定路径的 Go 工具，供 Go 1.19 之前的版本回退使用。
//
// 任一安装失败时返回对应错误；path 中的每个元素均视为完整的包路径。
func GoInstall(path ...string) error {
	for _, p := range path {
		fmt.Printf("go get -u %s\n", p)
		cmd := exec.Command("go", "get", "-u", p)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}
