package redis

var log logger

func (configs *Configs) SetLogger(l logger) {
	log = l
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

func (configs *Configs) SetTrace(t bool) {
	trace = t
}
