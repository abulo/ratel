package watch

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

// GoVersion 执行 go version，返回去掉版本前缀后的本地 Go 版本字符串；命令执行失败时返回错误。
func GoVersion() (string, error) {
	var buf bytes.Buffer
	cmd := exec.Command("go", "version")
	cmd.Stdout = &buf
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(strings.TrimPrefix(buf.String(), "go version ")), nil
}

// Watch 启动热编译监视服务。
//
// 配置校验失败或获取 Go 版本失败时直接返回错误；编译过程中的错误不通过返回值
// 传递，而是写入配置的 Logger。上下文取消时结束监视，并终止现存编译进程与被编译程序。
func Watch(ctx context.Context, opt *Options) error {
	if err := opt.sanitize(); err != nil {
		return err
	}
	b := opt.newBuilder()
	env, err := GoVersion()
	if err != nil {
		return err
	}
	b.logf(LogTypeInfo, fmt.Sprintf("当前环境参数如下：%s", env))
	b.logf(LogTypeInfo, fmt.Sprintf("给程序传递了以下参数：%s", b.appArgs))
	switch {
	case len(b.exts) == 0:
		b.logf(LogTypeWarn, "将 ext 设置为空值，意味着不监视任何文件的改变！")
	case len(b.exts) > 0:
		b.logf(LogTypeInfo, fmt.Sprintf("系统将监视以下类型的文件：%s", b.exts))
	}
	b.logf(LogTypeInfo, fmt.Sprintf("输出文件为：%s", b.appName))
	return b.watch(ctx, opt.paths)
}

// watch 先主动编译一次，随后进入事件循环监视目录变更；上下文取消时清理进程后返回。
func (b *builder) watch(ctx context.Context, paths []string) error {
	go b.build() // 首次主动编译，后续变更才由事件触发。
	watcher, err := b.initWatcher(paths)
	if err != nil {
		return err
	}
	defer watcher.Close()
	var buildTime time.Time
	for {
		select {
		case <-ctx.Done():
			b.killApp()
			b.killGo()
			b.logf(LogTypeInfo, "用户取消")
			return nil
		case event := <-watcher.Events:
			if event.Op&fsnotify.Chmod == fsnotify.Chmod {
				b.logf(LogTypeIgnore, fmt.Sprintf("watcher.Events:忽略 %s 事件", event.String()))
				continue
			}
			if b.isIgnore(event.Name) {
				b.logf(LogTypeIgnore, fmt.Sprintf("watcher.Events:忽略不被监视的文件：%s", event.String()))
				continue
			}
			if time.Since(buildTime) <= b.watcherFreq {
				b.logf(LogTypeIgnore, fmt.Sprintf("watcher.Events:忽略短期内频繁修改的文件：%s", event.Name))
				continue
			}
			buildTime = time.Now()
			if event.Name == "go.mod" && b.goTidy {
				b.logf(LogTypeInfo, fmt.Sprintf("watcher.Events:%s 事件触发了 go mod tidy", event.String()))
				go b.tidy()
			} else {
				b.logf(LogTypeInfo, fmt.Sprintf("watcher.Events:%s 事件触发了编译", event.String()))
				go b.build()
			}
		case err := <-watcher.Errors:
			b.logf(LogTypeWarn, fmt.Sprintf("watcher.Errors：%s", err.Error()))
			return nil
		}
	}
}

// initWatcher 创建 fsnotify 监视器，过滤掉无可监视文件的目录后逐个加入监视；任一目录加入失败时关闭监视器并返回错误。
func (b *builder) initWatcher(paths []string) (*fsnotify.Watcher, error) {
	b.logf(LogTypeInfo, "初始化监视器...")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	paths = b.filterPaths(paths)
	ps := strings.Join(paths, ",\n")
	b.logf(LogTypeIgnore, fmt.Sprintf("以下路径或是文件将被监视：%s", ps))
	for _, path := range paths {
		if err := watcher.Add(path); err != nil {
			watcher.Close()
			return nil, err
		}
	}
	return watcher, nil
}
