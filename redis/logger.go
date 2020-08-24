package redis

var _logger logger

func SetLogger(l logger) {
	_logger = l
}

type logger interface {
	Panic(args ...interface{})
	Fatal(args ...interface{})
	Error(args ...interface{})
	Warning(args ...interface{})
	Warn(args ...interface{})
	Info(args ...interface{})
	Debug(args ...interface{})
	Trace(args ...interface{})
}

var trace bool

func SetTrace(t bool) {
	trace = t
}
