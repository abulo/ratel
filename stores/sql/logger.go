package sql

import (
	"context"
	"errors"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type WrapLogger struct {
	Logrus                *logrus.Logger
	SlowThreshold         time.Duration
	SourceField           string
	SkipErrRecordNotFound bool
	Debug                 bool
}

func NewLogger(logger *logrus.Logger) *WrapLogger {
	return &WrapLogger{
		Logrus:                logger,
		SkipErrRecordNotFound: true,
		Debug:                 true,
	}
}

func (l *WrapLogger) SetDebug(debug bool) {
	l.Debug = debug
}

func (l *WrapLogger) SetSourceField(field string) {
	l.SourceField = field
}

func (l *WrapLogger) SetSkipErrRecordNotFound(skip bool) {
	l.SkipErrRecordNotFound = skip
}

func (l *WrapLogger) SetSlowThreshold(threshold time.Duration) {
	l.SlowThreshold = threshold
}

func (l *WrapLogger) LogMode(level logger.LogLevel) logger.Interface {
	l.Logrus.SetLevel(logrus.Level(level))
	return l
}

func (l *WrapLogger) Info(ctx context.Context, s string, args ...interface{}) {
	l.Logrus.WithContext(ctx).Infof(s, args...)
}

func (l *WrapLogger) Warn(ctx context.Context, s string, args ...interface{}) {
	l.Logrus.WithContext(ctx).Warnf(s, args...)
}

func (l *WrapLogger) Error(ctx context.Context, s string, args ...interface{}) {
	l.Logrus.WithContext(ctx).Errorf(s, args...)
}

func (l *WrapLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, _ := fc()
	fields := logrus.Fields{}
	if l.SourceField != "" {
		fields[l.SourceField] = utils.FileWithLineNum()
	}
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound) && l.SkipErrRecordNotFound) {
		fields[logrus.ErrorKey] = err
		l.Logrus.WithContext(ctx).WithFields(fields).Errorf("%s [%s]", sql, elapsed)
		return
	}
	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		l.Logrus.WithContext(ctx).WithFields(fields).Warnf("%s [%s]", sql, elapsed)
		return
	}
	if l.Debug {
		l.Logrus.WithContext(ctx).WithFields(fields).Debugf("%s [%s]", sql, elapsed)
	}
}

type WrapEntry struct {
	Logrus                *logrus.Entry
	SlowThreshold         time.Duration
	SourceField           string
	SkipErrRecordNotFound bool
	Debug                 bool
}

func NewEntry(entry *logrus.Entry) *WrapEntry {
	return &WrapEntry{
		Logrus:                entry,
		SkipErrRecordNotFound: true,
		Debug:                 true,
	}
}

func (l *WrapEntry) SetDebug(debug bool) {
	l.Debug = debug
}

func (l *WrapEntry) LogMode(level logger.LogLevel) logger.Interface {
	l.Logrus.Logger.SetLevel(logrus.Level(level))
	return l
}

func (l *WrapEntry) Info(ctx context.Context, s string, args ...interface{}) {
	l.Logrus.WithContext(ctx).Infof(s, args...)
}

func (l *WrapEntry) Warn(ctx context.Context, s string, args ...interface{}) {
	l.Logrus.WithContext(ctx).Warnf(s, args...)
}

func (l *WrapEntry) Error(ctx context.Context, s string, args ...interface{}) {
	l.Logrus.WithContext(ctx).Errorf(s, args...)
}

func (l *WrapEntry) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, _ := fc()
	fields := logrus.Fields{}
	if l.SourceField != "" {
		fields[l.SourceField] = utils.FileWithLineNum()
	}
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound) && l.SkipErrRecordNotFound) {
		fields[logrus.ErrorKey] = err
		l.Logrus.WithContext(ctx).WithFields(fields).Errorf("%s [%s]", sql, elapsed)
		return
	}
	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		l.Logrus.WithContext(ctx).WithFields(fields).Warnf("%s [%s]", sql, elapsed)
		return
	}
	if l.Debug {
		l.Logrus.WithContext(ctx).WithFields(fields).Debugf("%s [%s]", sql, elapsed)
	}
}

func LoggerCaller(skip int) map[string]string {
	pc, file, lineNo, _ := runtime.Caller(skip)
	name := runtime.FuncForPC(pc).Name()
	return map[string]string{
		"file": file,
		"func": name,
		"line": cast.ToString(lineNo),
	}
}
