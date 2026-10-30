package es

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/abulo/ratel/v3/core/logger/entry"
	"github.com/abulo/ratel/v3/core/logger/queue"
	"github.com/abulo/ratel/v3/stores/elasticsearch"
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
	extra      map[string]any
	exec       ExecCloser
	levels     []logrus.Level
	out        io.Writer
}

// SetMaxQueues 设置异步队列允许缓冲的最大日志条目数。
func SetMaxQueues(maxQueues int) Option {
	return func(o *options) {
		o.maxQueues = maxQueues
	}
}

// SetMaxWorkers 设置后台消费日志条目数的工作协程数量。
func SetMaxWorkers(maxWorkers int) Option {
	return func(o *options) {
		o.maxWorkers = maxWorkers
	}
}

// SetExtra 为每条日志附加固定的扩展字段；条目已存在的同名字段不会被覆盖。
func SetExtra(extra map[string]any) Option {
	return func(o *options) {
		o.extra = extra
	}
}

// SetExec 自定义日志条目的后端写入器；未设置时由 Default 使用内部默认实现。
func SetExec(exec ExecCloser) Option {
	return func(o *options) {
		o.exec = exec
	}
}

// SetLevels 指定钩子需要处理的日志级别集合；传入空级别列表时保留默认级别集合。
func SetLevels(levels ...logrus.Level) Option {
	return func(o *options) {
		if len(levels) == 0 {
			return
		}
		o.levels = levels
	}
}

// SetOut 指定写入失败时错误信息的输出目标；缺省输出到标准错误。
func SetOut(out io.Writer) Option {
	return func(o *options) {
		o.out = out
	}
}

// Option 是 Hook 的配置选项函数，以函数式选项模式在 New 中修改内部配置。
type Option func(*options)

// Default 创建默认的 Elasticsearch 日志钩子，使用 NewExec 将条目写入指定索引。
//
// 索引创建失败时构造函数直接 panic；调用方可通过 opts 覆盖队列容量、级别等配置。
func Default(client *elasticsearch.Client, index string, opts ...Option) *Hook {
	var options []Option
	options = append(options, opts...)
	options = append(options, SetExec(NewExec(client, index)))

	if err := CreateIndex(client); err != nil {
		panic(err)
	}
	return New(options...)
}

// DefaultWithURL 创建默认的 Elasticsearch 日志钩子，使用 NewExecWithURL 将条目写入指定索引。
//
// 索引创建失败时构造函数直接 panic。
func DefaultWithURL(client *elasticsearch.Client, index string, opts ...Option) *Hook {
	var options []Option
	options = append(options, opts...)
	options = append(options, SetExec(NewExecWithURL(client, index)))
	if err := CreateIndex(client); err != nil {
		panic(err)
	}
	return New(options...)
}

// New 根据给定选项构造日志钩子并启动后台异步队列。
//
// 未提供写入器时仅记录一条提示日志，条目将无法写入后端。
func New(opt ...Option) *Hook {
	opts := defaultOptions
	for _, o := range opt {
		o(&opts)
	}

	if opts.exec == nil {
		logrus.Info("Unknown Execer interface implementation")
	}

	q := queue.NewQueue(opts.maxQueues, opts.maxWorkers)
	q.Run()

	return &Hook{
		opts: opts,
		q:    q,
	}
}

// Hook 是挂载到 logrus 的 Elasticsearch 日志钩子，按配置级别把日志异步提交到后台队列。
//
// Hook 须由 New 或 Default 创建，不要直接构造零值；队列中的日志条目在 Flush 时统一等待处理完成。
type Hook struct {
	opts options
	q    *queue.Queue
}

// Levels 返回钩子启用的日志级别集合，供 logrus 判断是否触发 Fire。
func (h *Hook) Levels() []logrus.Level {
	return h.opts.levels
}

// Fire 在日志级别命中时被 logrus 调用，将 logrus 条目转换为 entry.Entry 后投入异步队列。
//
// 转换与入队流程不阻断业务日志，Fire 始终返回 nil；写入后端的失败在后台消费阶段记录到错误输出。
func (h *Hook) Fire(entryLogrus *logrus.Entry) error {
	var funcVal string
	var fileVal string
	if entryLogrus.Caller != nil {
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
	h.q.Push(queue.NewJob(newEntry, func(v any) {
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
		fmt.Fprintf(h.opts.out, "[es-Hook] Execution error: %s", err.Error())
	}
}

// Flush 等待后台队列中的日志条目处理完毕并关闭工作协程。
func (h *Hook) Flush() {
	h.q.Terminate()
}
