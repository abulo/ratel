package sql

import (
	"os"
	"runtime"
	"strings"

	"github.com/abulo/ratel/v3/core/metric"
	"github.com/abulo/ratel/v3/core/trace"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/spf13/cast"
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
	return func(op string, client *Client, next Handler) Handler {
		return func(scope *gorm.DB) {
			if ctx := scope.Statement.Context; ctx != nil {
				call := Caller(7)
				if parentSpan := trace.SpanFromContext(ctx); parentSpan != nil {
					parentCtx := parentSpan.Context()
					span := opentracing.StartSpan(client.DriverName, opentracing.ChildOf(parentCtx))
					ext.SpanKindRPCClient.Set(span)
					hostName, err := os.Hostname()
					if err != nil {
						hostName = "unknown"
					}
					ext.PeerHostname.Set(span, hostName)
					ext.PeerAddress.Set(span, client.Host)
					ext.DBInstance.Set(span, client.Database)
					ext.DBStatement.Set(span, client.DriverName)
					sqlRaw := scope.Dialector.Explain(scope.Statement.SQL.String(), scope.Statement.Vars...)
					span.LogFields(log.String("sql", sqlRaw))
					span.LogFields(log.Object("call", call))
					defer span.Finish()
					ctx = opentracing.ContextWithSpan(ctx, span)
					scope.Statement.Context = ctx
					next(scope)
					if scope.Error != nil {
						ext.Error.Set(span, true)
					}
					return
				}
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
