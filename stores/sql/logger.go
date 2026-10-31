package sql

import (
	"context"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

// WrapLogger 将 logrus.Logger 适配为 gorm logger.Interface，供 GORM 输出 SQL 日志。
//
// 零值不可直接使用，须由 NewLogger 构造；各 Set* 方法在运行期调整开关与阈值，
// 不保证并发修改安全。
type WrapLogger struct {
	// Logrus 为底层日志实例，日志级别与输出均由它承载。
	Logrus *logrus.Logger
	// SlowThreshold 为慢 SQL 判定阈值，超过该耗时以 Warn 级别记录；为 0 时不判定慢查询。
	SlowThreshold time.Duration
	// SourceField 为调用方源码位置在日志字段中的键名；为空时不记录源码位置。
	SourceField string
	// SkipErrRecordNotFound 开启后忽略记录不存在类错误，不作为错误日志上报。
	SkipErrRecordNotFound bool
	// Debug 是否在正常路径下输出 Debug 级别 SQL 日志。
	Debug bool
}

// NewLogger 基于给定 logrus.Logger 构造 WrapLogger，默认忽略记录不存在错误并开启调试输出。
//
// 传入 nil 会在首次记录日志时触发空指针解引用，调用方须提供有效的日志实例。
func NewLogger(logger *logrus.Logger) *WrapLogger {
	return &WrapLogger{
		Logrus:                logger,
		SkipErrRecordNotFound: true,
		Debug:                 true,
	}
}

// SetDebug 调整调试日志开关。
func (l *WrapLogger) SetDebug(debug bool) {
	l.Debug = debug
}

// SetSourceField 设置调用方源码位置在日志字段中的键名；传空字符串可关闭源码位置记录。
func (l *WrapLogger) SetSourceField(field string) {
	l.SourceField = field
}

// SetSkipErrRecordNotFound 设置是否忽略记录不存在类错误。
func (l *WrapLogger) SetSkipErrRecordNotFound(skip bool) {
	l.SkipErrRecordNotFound = skip
}

// SetSlowThreshold 设置慢 SQL 判定阈值；传 0 可关闭慢查询判定。
func (l *WrapLogger) SetSlowThreshold(threshold time.Duration) {
	l.SlowThreshold = threshold
}

// LogMode 按 GORM 日志级别调整底层 logrus 级别，并返回适配器自身以满足 logger.Interface 契约。
func (l *WrapLogger) LogMode(level logger.LogLevel) logger.Interface {
	l.Logrus.SetLevel(logrus.Level(level))
	return l
}

// Info 记录 Info 级别日志，上下文透传给底层 logrus。
func (l *WrapLogger) Info(ctx context.Context, s string, args ...any) {
	l.Logrus.WithContext(ctx).Infof(s, args...)
}

// Warn 记录 Warn 级别日志，上下文透传给底层 logrus。
func (l *WrapLogger) Warn(ctx context.Context, s string, args ...any) {
	l.Logrus.WithContext(ctx).Warnf(s, args...)
}

// Error 记录 Error 级别日志，上下文透传给底层 logrus。
func (l *WrapLogger) Error(ctx context.Context, s string, args ...any) {
	l.Logrus.WithContext(ctx).Errorf(s, args...)
}

// Trace 记录一次 SQL 执行轨迹：记录不存在错误被忽略时不报错，超过慢查询阈值记 Warn，
// 调试开关开启时常规执行记 Debug；错误在未被忽略时记 Error。
func (l *WrapLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	// PERF-2-2 [P1] 先判定本次是否会输出日志；不满足任何级别时直接返回，
	// 不构造 SQL 字符串、不做调用栈回溯、不分配 Fields 映射，日志输出内容与原逻辑保持一致。
	willLogError := Acceptable(err) != nil && !l.SkipErrRecordNotFound
	willLogSlow := l.SlowThreshold != 0 && elapsed > l.SlowThreshold
	if !willLogError && !willLogSlow && !l.Debug {
		return
	}
	sql, _ := fc()
	fields := logrus.Fields{}
	if l.SourceField != "" {
		fields[l.SourceField] = utils.FileWithLineNum()
	}
	if willLogError {
		fields[logrus.ErrorKey] = err
		l.Logrus.WithContext(ctx).WithFields(fields).Errorf("%s [%s]", sql, elapsed)
		return
	}
	if willLogSlow {
		l.Logrus.WithContext(ctx).WithFields(fields).Warnf("%s [%s]", sql, elapsed)
		return
	}
	if l.Debug {
		l.Logrus.WithContext(ctx).WithFields(fields).Debugf("%s [%s]", sql, elapsed)
	}
}

// WrapEntry 将 logrus.Entry 适配为 gorm logger.Interface，行为与 WrapLogger 一致，
// 区别在于底层持有带固定上下文字段的 Entry。
//
// 零值不可直接使用，须由 NewEntry 构造；各配置项不保证并发修改安全。
type WrapEntry struct {
	// Logrus 为底层日志条目，日志级别与输出均由它承载。
	Logrus *logrus.Entry
	// SlowThreshold 为慢 SQL 判定阈值，超过该耗时以 Warn 级别记录；为 0 时不判定慢查询。
	SlowThreshold time.Duration
	// SourceField 为调用方源码位置在日志字段中的键名；为空时不记录源码位置。
	SourceField string
	// SkipErrRecordNotFound 开启后忽略记录不存在类错误，不作为错误日志上报。
	SkipErrRecordNotFound bool
	// Debug 是否在正常路径下输出 Debug 级别 SQL 日志。
	Debug bool
}

// NewEntry 基于给定 logrus.Entry 构造 WrapEntry，默认忽略记录不存在错误并开启调试输出。
//
// 传入 nil 会在首次记录日志时触发空指针解引用，调用方须提供有效的日志条目。
func NewEntry(entry *logrus.Entry) *WrapEntry {
	return &WrapEntry{
		Logrus:                entry,
		SkipErrRecordNotFound: true,
		Debug:                 true,
	}
}

// SetDebug 调整调试日志开关。
func (l *WrapEntry) SetDebug(debug bool) {
	l.Debug = debug
}

// LogMode 按 GORM 日志级别调整底层日志级别，并返回适配器自身以满足 logger.Interface 契约。
func (l *WrapEntry) LogMode(level logger.LogLevel) logger.Interface {
	l.Logrus.Logger.SetLevel(logrus.Level(level))
	return l
}

// Info 记录 Info 级别日志，上下文透传给底层 logrus。
func (l *WrapEntry) Info(ctx context.Context, s string, args ...any) {
	l.Logrus.WithContext(ctx).Infof(s, args...)
}

// Warn 记录 Warn 级别日志，上下文透传给底层 logrus。
func (l *WrapEntry) Warn(ctx context.Context, s string, args ...any) {
	l.Logrus.WithContext(ctx).Warnf(s, args...)
}

// Error 记录 Error 级别日志，上下文透传给底层 logrus。
func (l *WrapEntry) Error(ctx context.Context, s string, args ...any) {
	l.Logrus.WithContext(ctx).Errorf(s, args...)
}

// Trace 记录一次 SQL 执行轨迹：记录不存在错误被忽略时不报错，超过慢查询阈值记 Warn，
// 调试开关开启时常规执行记 Debug；错误在未被忽略时记 Error。
func (l *WrapEntry) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	// PERF-2-2 [P1] 先判定本次是否会输出日志；不满足任何级别时直接返回，
	// 不构造 SQL 字符串、不做调用栈回溯、不分配 Fields 映射，日志输出内容与原逻辑保持一致。
	willLogError := Acceptable(err) != nil && !l.SkipErrRecordNotFound
	willLogSlow := l.SlowThreshold != 0 && elapsed > l.SlowThreshold
	if !willLogError && !willLogSlow && !l.Debug {
		return
	}
	sql, _ := fc()
	fields := logrus.Fields{}
	if l.SourceField != "" {
		fields[l.SourceField] = utils.FileWithLineNum()
	}
	if willLogError {
		fields[logrus.ErrorKey] = err
		l.Logrus.WithContext(ctx).WithFields(fields).Errorf("%s [%s]", sql, elapsed)
		return
	}
	if willLogSlow {
		l.Logrus.WithContext(ctx).WithFields(fields).Warnf("%s [%s]", sql, elapsed)
		return
	}
	if l.Debug {
		l.Logrus.WithContext(ctx).WithFields(fields).Debugf("%s [%s]", sql, elapsed)
	}
}

// LoggerCaller 返回跳过 skip 层调用栈后当前调用方的文件、函数名与行号映射，
// 用于日志或追踪信息中定位原始调用位置；skip 超出调用栈深度时返回零值字段。
func LoggerCaller(skip int) map[string]string {
	pc, file, lineNo, _ := runtime.Caller(skip)
	name := runtime.FuncForPC(pc).Name()
	return map[string]string{
		"file": file,
		"func": name,
		"line": cast.ToString(lineNo),
	}
}
