//go:build !go1.19

package base

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
)

// GoInstall 使用 go get -u 安装指定路径的 Go 工具，供 Go 1.19 之前的版本回退使用。
//
// 任一安装失败时返回对应错误；path 中的每个元素均视为完整的包路径。
func GoInstall(path ...string) error {
	type installResult struct {
		err error
	}
	results := make([]installResult, len(path))

	// PERF-8-14 [P2] 多工具安装受限并发（并发度 3，与 go1.19 构建版本保持一致策略）；
	// 错误按原参数顺序收集并返回首个失败，保持“任一失败即返回错误”的语义。
	sem := make(chan struct{}, 3)
	var wg sync.WaitGroup
	for i, p := range path {
		fmt.Printf("go get -u %s\n", p)
		idx := i
		wg.Add(1)
		sem <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			cmd := exec.Command("go", "get", "-u", p)
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
