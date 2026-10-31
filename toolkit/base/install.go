//go:build go1.19

package base

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

// GoInstall 使用 go install 命令安装指定路径的 Go 工具。
//
// 路径未包含 @ 版本后缀时自动追加 @latest；任一安装失败时返回对应错误。
func GoInstall(path ...string) error {
	type installResult struct {
		err error
	}
	results := make([]installResult, len(path))

	// PERF-8-14 [P2] 多工具安装受限并发（并发度 3）：命令行按参数顺序先串行打印，
	// 安装过程并发发起；错误按原参数顺序收集并返回首个失败，与原实现“任一失败即返回错误”的语义一致。
	sem := make(chan struct{}, 3)
	var wg sync.WaitGroup
	for i, raw := range path {
		p := raw
		if !strings.Contains(p, "@") {
			p += "@latest"
		}
		fmt.Printf("go install %s\n", p)
		idx := i
		wg.Add(1)
		sem <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			cmd := exec.Command("go", "install", p)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			results[idx] = installResult{err: cmd.Run()}
		}()
	}
	wg.Wait()
	for _, res := range results {
		if res.err != nil {
			return res.err
		}
	}
	return nil
}
