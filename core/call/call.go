package call

import "runtime"

func Caller(skip int) (fn string, file string, lineNo int) {
	pc, file, lineNo, _ := runtime.Caller(skip)
	fn = runtime.FuncForPC(pc).Name()
	return fn, file, lineNo
}
