package watch

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

type builder struct {
	// exts 为需要监视的文件扩展名。
	exts []string
	// extSet 为 exts 的集合形态，供 isIgnore 做 O(1) 后缀判定（PERF-7-40）。
	extSet map[string]struct{}
	// anyExt 标记是否监视所有文件类型。
	anyExt bool
	// excludes 为忽略文件的路径匹配模式列表。
	excludes []string
	// appName 为输出的可执行文件路径。
	appName string
	// logs 为日志输出对象。
	logs Logger
	// watcherFreq 为判定文件变更的时间间隔。
	watcherFreq time.Duration

	// appWD 为被编译程序的工作目录。
	appWD string
	// appCmd 为被编译程序的进程句柄，用于结束旧进程。
	appCmd *exec.Cmd
	// appArgs 为传给被编译程序启动命令的参数。
	appArgs []string
	// appKillMux 保护 appCmd 的互斥锁。
	appKillMux sync.Mutex

	// goTidy 标记是否自动执行 go mod tidy。
	goTidy bool
	// goCmd 为 go build 或 go mod tidy 的进程句柄。
	goCmd *exec.Cmd
	// goArgs 为传给 go build 的参数，go mod tidy 忽略该参数。
	goArgs []string
	// goKillMux 保护 goCmd 的互斥锁。
	goKillMux sync.Mutex
}

func (opt *Options) newBuilder() *builder {
	exts := opt.Exts
	if opt.anyExts {
		exts = []string{"*"}
	}
	// PERF-7-40 [P2] 构造期把监视扩展名预构建为集合，事件回调内 O(1) 判定，替代每次事件线性扫描 exts（开发期工具，预期收益需实测验证）
	extSet := make(map[string]struct{}, len(exts))
	for _, ext := range exts {
		extSet[ext] = struct{}{}
	}
	return &builder{
		exts:        exts,
		extSet:      extSet,
		anyExt:      opt.anyExts,
		excludes:    opt.Excludes,
		appName:     opt.appName,
		logs:        opt.Logger,
		watcherFreq: opt.WatcherFrequency,
		appWD:       filepath.Dir(opt.appName),
		appArgs:     opt.appArgs,
		goTidy:      opt.AutoTidy,
		goArgs:      opt.goCmdArgs,
	}
}

// logf 按指定日志类型输出一条消息。
func (b *builder) logf(typ int8, msg string) {
	b.logs.Output(typ, msg)
}

// isIgnore 判断指定路径是否应被忽略而不监视；被编译程序自身、匹配忽略模式
// 或扩展名不在监视列表中的文件返回 true。
func (b *builder) isIgnore(path string) bool {
	if b.appCmd != nil && b.appCmd.Path == path {
		return true
	}

	if b.goTidy && filepath.Base(path) == "go.mod" {
		return false
	}

	for _, p := range b.excludes {
		matched, err := filepath.Match(p, path)
		if err != nil {
			b.logf(LogTypeError, err.Error())
			return false
		}
		if matched {
			return true
		}
	}

	if b.anyExt {
		return false
	}
	// PERF-7-40 [P2] 用预构建集合按文件扩展名 O(1) 查找，替代逐后缀 HasSuffix 线性扫描
	if _, ok := b.extSet[filepath.Ext(path)]; ok {
		return false
	}

	return true
}

// tidy 中止现存编译进程后执行 go mod tidy，输出走 Go 编译器日志通道。
func (b *builder) tidy() {
	b.killGo()

	b.logf(LogTypeInfo, "执行 go mod tidy...")

	b.goCmd = exec.Command("go", "mod", "tidy")
	b.goCmd.Stderr = asWriter(LogTypeGo, b.logs)
	b.goCmd.Stdout = asWriter(LogTypeGo, b.logs)
	if err := b.goCmd.Run(); err != nil {
		b.logf(LogTypeError, fmt.Sprintf("go mod tidy 失败：%s", err.Error()))
		return
	}

	b.logf(LogTypeSuccess, "go mod tidy 完成!")
	b.goCmd = nil
}

// build 中止现存编译进程后执行 go build；编译成功后重启被编译程序，失败仅记录日志。
func (b *builder) build() {
	b.killGo()

	b.logf(LogTypeInfo, "编译代码...")

	b.goCmd = exec.Command("go", b.goArgs...)
	b.goCmd.Stderr = asWriter(LogTypeGo, b.logs)
	b.goCmd.Stdout = asWriter(LogTypeGo, b.logs)
	if err := b.goCmd.Run(); err != nil {
		b.logf(LogTypeError, fmt.Sprintf("编译失败：%s", err.Error()))
		return
	}

	b.logf(LogTypeSuccess, "编译成功!")
	b.goCmd = nil

	b.restartApp()
}

// killGo 结束仍在运行的 go 编译进程并等待其退出。
func (b *builder) killGo() {
	b.goKillMux.Lock()
	defer b.goKillMux.Unlock()

	if b.goCmd != nil && b.goCmd.Process != nil {
		b.logf(LogTypeInfo, "中止旧的编译进程")
		if err := b.goCmd.Process.Kill(); err != nil {
			b.logf(LogTypeError, fmt.Sprintf("中止旧的编译进程失败：%s", err.Error()))
		}
		if err := b.goCmd.Wait(); err != nil {
			b.logf(LogTypeError, fmt.Sprintf("被中止编译进程非正常退出：%s", err.Error()))
		}
		b.logf(LogTypeSuccess, "旧的编译进程被终止!")
		b.goCmd = nil
	}
}

// restartApp 先结束旧进程，再以配置的工作目录启动新进程；启动过程中的异常被 recover 后仅记录日志。
func (b *builder) restartApp() {
	defer func() {
		if err := recover(); err != nil {
			b.logf(LogTypeError, fmt.Sprintf("重启失败：%v", err))
		}
	}()

	b.killApp()

	b.logf(LogTypeInfo, fmt.Sprintf("启动新进程：%s", b.appName))
	b.appCmd = exec.Command(b.appName, b.appArgs...)
	b.appCmd.Dir = b.appWD
	b.appCmd.Stderr = asWriter(LogTypeApp, b.logs)
	b.appCmd.Stdout = asWriter(LogTypeApp, b.logs)
	if err := b.appCmd.Start(); err != nil {
		b.logf(LogTypeError, fmt.Sprintf("启动进程时出错：%s", err))
	}
}

// killApp 结束仍在运行的被编译程序并等待其退出。
func (b *builder) killApp() {
	b.appKillMux.Lock()
	defer b.appKillMux.Unlock()

	if b.appCmd != nil && b.appCmd.Process != nil {
		b.logf(LogTypeInfo, fmt.Sprintf("中止旧进程：%s", b.appName))
		if err := b.appCmd.Process.Kill(); err != nil {
			b.logf(LogTypeError, fmt.Sprintf("中止旧进程失败：%s", err.Error()))
		}
		if err := b.appCmd.Wait(); err != nil {
			b.logf(LogTypeError, fmt.Sprintf("被中止进程非正常退出：%s", err.Error()))
		}
		b.logf(LogTypeSuccess, "旧进程被终止!")
		b.appCmd = nil
	}
}

// filterPaths 剔除目录下不存在任何需监视文件的目录。
func (b *builder) filterPaths(paths []string) []string {
	ret := make([]string, 0, len(paths))

	for _, dir := range paths {
		fs, err := os.ReadDir(dir)
		if err != nil {
			b.logf(LogTypeError, err.Error())
			continue
		}

		ignore := true
		for _, f := range fs {
			if f.IsDir() {
				continue
			}
			if !b.isIgnore(f.Name()) {
				ignore = false
				break
			}
		}
		if !ignore {
			ret = append(ret, dir)
		}
	}

	return ret
}
