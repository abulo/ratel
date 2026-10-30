package watch

import (
	"io"
	"strings"

	"github.com/issue9/term/v3/colors"
)

// 日志类型标识，用于区分日志的重要程度与来源。
const (
	// LogTypeSuccess 表示操作成功类信息。
	LogTypeSuccess int8 = iota
	// LogTypeInfo 表示普通提示信息。
	LogTypeInfo
	// LogTypeWarn 表示需要关注的告警信息。
	LogTypeWarn
	// LogTypeError 表示错误信息。
	LogTypeError
	// LogTypeIgnore 表示信息量大、优先级低的忽略类内容，默认不在控制台展示。
	LogTypeIgnore
	// LogTypeApp 为被编译程序输出的信息。
	LogTypeApp
	// LogTypeGo 为 Go 编译器或 go 命令输出的信息。
	LogTypeGo
)

type (
	// Logger 是热编译过程中的日志接收抽象，实现负责把不同类型的日志路由到对应输出通道。
	Logger interface {
		// Output 输出一条日志；t 为日志类型，标识重要程度或来源。
		Output(t int8, message string)
	}

	loggerWriter struct {
		t int8
		w Logger
	}

	consoleLogger struct {
		showIgnore bool
		writers    map[int8]*consoleWriter
	}

	consoleWriter struct {
		out    io.Writer
		color  colors.Color
		prefix string
	}
)

func (w *loggerWriter) Write(bs []byte) (int, error) {
	w.w.Output(w.t, string(bs))
	return len(bs), nil
}

func asWriter(t int8, w Logger) io.Writer { return &loggerWriter{t: t, w: w} }

func (c *consoleLogger) Output(t int8, msg string) {
	if !c.showIgnore && t == LogTypeIgnore {
		return
	}

	w := c.writers[t]
	colors.Fprint(w.out, colors.Normal, w.color, colors.Default, w.prefix)
	msg = strings.TrimRight(msg, "\n")
	colors.Fprintln(w.out, colors.Normal, colors.Default, colors.Default, msg)
}

// NewConsoleLogger 创建将日志输出到控制台的 Logger。
//
// showIgnore 为 false 时忽略 LogTypeIgnore 类型日志；错误类日志写入 err，
// 其余类型写入 out。
func NewConsoleLogger(showIgnore bool, err, out io.Writer) Logger {
	newCW := func(out io.Writer, color colors.Color, prefix string) *consoleWriter {
		return &consoleWriter{
			out:    out,
			color:  color,
			prefix: prefix,
		}
	}

	return &consoleLogger{
		showIgnore: showIgnore,
		writers: map[int8]*consoleWriter{
			LogTypeSuccess: newCW(out, colors.Green, "[SUCC] "),
			LogTypeInfo:    newCW(out, colors.Blue, "[INFO] "),
			LogTypeWarn:    newCW(err, colors.Magenta, "[WARN] "),
			LogTypeError:   newCW(err, colors.Red, "[ERRO] "),
			LogTypeIgnore:  newCW(out, colors.Default, "[IGNO] "),
			LogTypeApp:     newCW(out, colors.Default, "[APP] "),
			LogTypeGo:      newCW(out, colors.Default, "[GO] "),
		},
	}
}
