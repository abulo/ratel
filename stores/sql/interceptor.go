package sql

import (
	"runtime"
	"strings"

	"github.com/abulo/ratel/v3/core/call"
	"github.com/abulo/ratel/v3/core/hostname"
	"github.com/abulo/ratel/v3/core/metric"
	globalTrace "github.com/abulo/ratel/v3/core/trace"
	"github.com/spf13/cast"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
	"gorm.io/gorm"
)

type (
	Handler     func(*gorm.DB)
	Interceptor func(op string, client *Client, next Handler) Handler
)

type processor interface {
	Get(name string) func(*gorm.DB)
	Replace(name string, handler func(*gorm.DB)) error
}

// 收敛status，避免prometheus日志太多
func getStatement(err string) string {
	if !strings.HasPrefix(err, "Errord") {
		return "Unknown"
	}
	slice := strings.Split(err, ":")
	if len(slice) < 2 {
		return "Unknown"
	}
	// 收敛错误
	return slice[0]
}

func RegisterInterceptor(db *gorm.DB, options *Client, interceptors ...Interceptor) {
	var processors = []struct {
		Name      string
		Processor processor
	}{
		{"gorm:create", db.Callback().Create()},
		{"gorm:query", db.Callback().Query()},
		{"gorm:delete", db.Callback().Delete()},
		{"gorm:update", db.Callback().Update()},
		{"gorm:row", db.Callback().Row()},
		{"gorm:raw", db.Callback().Raw()},
	}

	for _, interceptor := range interceptors {
		for _, processor := range processors {
			handler := processor.Processor.Get(processor.Name)
			handler = interceptor(processor.Name, options, handler)
			if err := processor.Processor.Replace(processor.Name, handler); err != nil {
				panic(err)
			}
		}
	}
}

func MetricInterceptor() Interceptor {
	return func(op string, client *Client, next Handler) Handler {
		return func(scope *gorm.DB) {
			// beg := time.Now()
			next(scope)
			// cost := time.Since(beg)
			// error metric
			if scope.Error != nil {
				metric.LibHandleCounter.WithLabelValues(client.DriverName, client.Database+"."+scope.Name(), client.Host, getStatement(scope.Error.Error())).Inc()
			} else {
				metric.LibHandleCounter.WithLabelValues(client.DriverName, client.Database+"."+scope.Name(), client.Host, "OK").Inc()
			}
		}
	}
}

func TraceInterceptor() Interceptor {
	tracer := globalTrace.NewTracer(trace.SpanKindClient)
	attrs := []attribute.KeyValue{
		semconv.RPCSystemKey.String("gorm"),
	}
	return func(op string, client *Client, next Handler) Handler {
		return func(scope *gorm.DB) {
			if ctx := scope.Statement.Context; ctx != nil {
				sqlRaw := scope.Dialector.Explain(scope.Statement.SQL.String(), scope.Statement.Vars...)
				fn, file, line := call.Caller(7)
				attrs = append(attrs,
					semconv.HostName(hostname.Hostname()),
					semconv.CodeFunction(fn),
					semconv.CodeFilepath(file),
					semconv.CodeLineNumber(line),
				)
				md := metadata.New(nil)
				_, span := tracer.Start(ctx, op, propagation.HeaderCarrier(md), trace.WithAttributes(attrs...))
				span.SetAttributes(
					semconv.DBNameKey.String(client.DriverName),
					semconv.DBConnectionStringKey.String(client.Host),
					semconv.DBUserKey.String(client.Username),
					semconv.DBStatementKey.String(sqlRaw),
				)
				defer span.End()
				next(scope)
				if scope.Error != nil {
					span.RecordError(scope.Error)
					span.SetStatus(codes.Error, scope.Error.Error())
				}
				return
			}
			next(scope)
		}
	}
}

func Caller(skip int) map[string]string {
	pc, file, lineNo, _ := runtime.Caller(skip)
	name := runtime.FuncForPC(pc).Name()
	return map[string]string{
		"path": file + ":" + cast.ToString(lineNo),
		"func": name,
	}
}
