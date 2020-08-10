package redis

var log logger

func (config *Config) SetLogger(l logger) {
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

func (config *Config) SetTrace(t bool) {
	trace = t
}
