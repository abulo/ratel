package entry

import (
	"runtime"
	"strconv"

	"github.com/sirupsen/logrus"
)

// levelNames 预存 logrus 各等级的大写名称，避免每条日志重复执行 Level.String() 与 strings.ToUpper。
// 取值与 logrus.Level.String() 的大写结果严格一一对应：WarnLevel 为 WARNING，不得简写为 WARN。
var levelNames = map[logrus.Level]string{
	logrus.PanicLevel: "PANIC",
	logrus.FatalLevel: "FATAL",
	logrus.ErrorLevel: "ERROR",
	logrus.WarnLevel:  "WARNING",
	logrus.InfoLevel:  "INFO",
	logrus.DebugLevel: "DEBUG",
	logrus.TraceLevel: "TRACE",
}

// LevelUpper 返回日志等级的大写名称；未识别等级返回 UNKNOWN，与原
// strings.ToUpper(level.String()) 的兜底行为保持一致。
func LevelUpper(level logrus.Level) string {
	if name, ok := levelNames[level]; ok {
		return name
	}
	return "UNKNOWN"
}

// CallerLocation 从调用栈帧提取函数名与 “文件:行号”；caller 为 nil 时两个返回值均为空字符串。
// 行号经 strconv.AppendInt 写入临时字节切片，避免 fmt.Sprintf 的格式解析与反射开销。
func CallerLocation(caller *runtime.Frame) (funcVal, fileVal string) {
	if caller == nil {
		return "", ""
	}
	funcVal = caller.Function
	var buf []byte
	buf = append(buf, caller.File...)
	buf = append(buf, ':')
	buf = strconv.AppendInt(buf, int64(caller.Line), 10)
	return funcVal, string(buf)
}
