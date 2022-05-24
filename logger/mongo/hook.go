package mongo

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/abulo/ratel/v3/logger/entry"
	"github.com/abulo/ratel/v3/logger/queue"
	"github.com/abulo/ratel/v3/stores/mongodb"
	"github.com/abulo/ratel/v3/util"
	"github.com/sirupsen/logrus"
)

var defaultOptions = options{
	maxQueues:  512,
	maxWorkers: 2,
	levels: []logrus.Level{
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	},
	out: os.Stderr,
}

type options struct {
	maxQueues  int
	maxWorkers int
	extra      map[string]interface{}
	exec       ExecCloser
	levels     []logrus.Level
	out        io.Writer
}

// SetMaxQueues 设置缓冲区的数量
func SetMaxQueues(maxQueues int) Option {
	return func(o *options) {
		o.maxQueues = maxQueues
	}
}

// SetMaxWorkers 设置工作线程数
func SetMaxWorkers(maxWorkers int) Option {
	return func(o *options) {
		o.maxWorkers = maxWorkers
	}
}

// SetExtra 设置扩展参数
func SetExtra(extra map[string]interface{}) Option {
	return func(o *options) {
		o.extra = extra
	}
}

// SetExec 设置Execer接口
func SetExec(exec ExecCloser) Option {
	return func(o *options) {
		o.exec = exec
	}
}

// SetLevels 设置可用的日志级别
func SetLevels(levels ...logrus.Level) Option {
	return func(o *options) {
		if len(levels) == 0 {
			return
		}
		o.levels = levels
	}
}

// SetOut 设置错误输出
func SetOut(out io.Writer) Option {
	return func(o *options) {
		o.out = out
	}
}

// Option 钩子参数选项
type Option func(*options)

// Default create a default mongo hook
func Default(client *mongodb.MongoDB, opts ...Option) *Hook {
	var options []Option
	options = append(options, opts...)
	options = append(options, SetExec(NewExec(client)))
	return New(options...)
}

// DefaultWithURL create a default mongo hook
func DefaultWithURL(client *mongodb.MongoDB, opts ...Option) *Hook {
	var options []Option
	options = append(options, opts...)
	options = append(options, SetExec(NewExecWithURL(client)))
	return New(options...)
}

// New 创建一个要添加到logger实例的钩子
func New(opt ...Option) *Hook {
	opts := defaultOptions
	for _, o := range opt {
		o(&opts)
	}

	if opts.exec == nil {
		// panic("Unknown Execer interface implementation")
		logrus.Info("Unknown Execer interface implementation")
	}

	q := queue.NewQueue(opts.maxQueues, opts.maxWorkers)
	q.Run()

	return &Hook{
		opts: opts,
		q:    q,
	}
}

// Hook 将日志发送到 mongo 数据库
type Hook struct {
	opts options
	q    *queue.Queue
}

// Levels 返回可用的日志记录级别
func (h *Hook) Levels() []logrus.Level {
	return h.opts.levels
}

// Fire 触发日志事件时将调用
func (h *Hook) Fire(entryLogrus *logrus.Entry) error {
	var funcVal string
	var fileVal string
	if entryLogrus.HasCaller() {
		funcVal = entryLogrus.Caller.Function
		fileVal = fmt.Sprintf("%s:%d", entryLogrus.Caller.File, entryLogrus.Caller.Line)
	}
	hostName, err := os.Hostname()
	if err != nil {
		hostName = "unknown"
	}
	level := entryLogrus.Level.String()
	newEntry := &entry.Entry{
		Host:      hostName,
		Timestamp: util.Now(),
		File:      fileVal,
		Func:      funcVal,
		Message:   entryLogrus.Message,
		Level:     strings.ToUpper(level),
		Data:      entryLogrus.Data,
	}
	h.q.Push(queue.NewJob(newEntry, func(v interface{}) {
		h.exec(v.(*entry.Entry))
	}))
	return nil
}

func (h *Hook) exec(entry *entry.Entry) {
	if extra := h.opts.extra; extra != nil {
		for k, v := range extra {
			if _, ok := entry.Data[k]; !ok {
				entry.Data[k] = v
			}
		}
	}
	err := h.opts.exec.Exec(entry)
	if err != nil && h.opts.out != nil {
		fmt.Fprintf(h.opts.out, "[Mongo-Hook] Execution error: %s", err.Error())
	}
}

// Flush 等待日志队列为空
func (h *Hook) Flush() {
	h.q.Terminate()
}
