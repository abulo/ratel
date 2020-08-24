package mongodb

//logger 日志接口
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

//log 日志
var _logger logger

//SetLogger 设置日志
func SetLogger(l logger) {
	_logger = l
}
