package logger

import (
	"github.com/sirupsen/logrus"
)

// Logger 是进程内共享的默认 logrus 日志器，包加载时即由 logrus.New 完成初始化。
// LoggerEntry 是绑定在 Logger 上的默认日志条目，便于通过 WithField 等方式附加字段后写日志。
var (
	Logger      *logrus.Logger = logrus.New()
	LoggerEntry *logrus.Entry  = logrus.NewEntry(Logger)
)
