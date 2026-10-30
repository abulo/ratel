package call

import "runtime"

// Caller 返回调用方栈帧对应的函数全名、源码文件路径与行号。
//
// skip 为需要向上跳过的栈帧数：0 表示调用 Caller 的函数所在栈帧，数值越大越向上回溯；
// 调用方通常传入自身封装层级对应的跳过层数。
func Caller(skip int) (fn string, file string, lineNo int) {
	pc, file, lineNo, _ := runtime.Caller(skip)
	fn = runtime.FuncForPC(pc).Name()
	return fn, file, lineNo
}
