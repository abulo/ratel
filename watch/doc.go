// Package watch 提供 Go 项目热编译监视能力。
//
// 该包按 Options 配置监视指定目录下的文件变更，自动执行 go build 与可选的
// go mod tidy，编译成功后重启被编译的程序；日志通过 Logger 接口对外输出，
// 支持按扩展名、路径模式与递归目录进行过滤。
package watch
