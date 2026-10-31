package call

import (
	"runtime"
	"sync"
)

// callerInfo 缓存单个 PC 对应的调用位置信息。
type callerInfo struct {
	fn   string
	file string
	line int
}

// PERF-4-8 [P1]：以 PC 为键缓存函数名/文件/行号，相同调用点第二次起跳过符号表解析。
// 说明：Go 代码地址在运行期不会被回收复用，同一 PC 的函数名、文件与行号恒定，
// 缓存结果不会失真（各调用点门控状态下的节省幅度预期收益需实测验证）。
var callerCache sync.Map // uintptr -> callerInfo

// Caller 返回调用方栈帧对应的函数全名、源码文件路径与行号。
//
// skip 为需要向上跳过的栈帧数：0 表示调用 Caller 的函数所在栈帧，数值越大越向上回溯；
// 调用方通常传入自身封装层级对应的跳过层数。
func Caller(skip int) (fn string, file string, lineNo int) {
	pc, file, lineNo, _ := runtime.Caller(skip)
	if pc == 0 {
		// 与原实现保持一致：无法解析栈帧时走原始符号查询路径。
		fn = runtime.FuncForPC(pc).Name()
		return fn, file, lineNo
	}

	if v, ok := callerCache.Load(pc); ok {
		info := v.(callerInfo)
		return info.fn, info.file, info.line
	}

	fn = runtime.FuncForPC(pc).Name()
	callerCache.Store(pc, callerInfo{fn: fn, file: file, line: lineNo})
	return fn, file, lineNo
}
