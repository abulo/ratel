package sql

import (
	"strings"

	"github.com/abulo/ratel/v3/core/call"
	"github.com/abulo/ratel/v3/core/hostname"
	"github.com/abulo/ratel/v3/core/metric"
	globalTrace "github.com/abulo/ratel/v3/core/trace"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
	"gorm.io/gorm"
)

type (
	// Handler 为 GORM 回调中被拦截包装的下一层处理器。
	Handler func(*gorm.DB)
	// Interceptor 以装饰器形式包装 GORM 回调处理器：op 为回调阶段名，client 为连接配置，
	// next 为原处理器，返回包装后的处理器。
	Interceptor func(op string, client *Client, next Handler) Handler
)

// processor 抽象 GORM 回调注册中心，便于对 Create、Query 等阶段统一取出现有处理器并替换。
type processor interface {
	Get(name string) func(*gorm.DB)
	Replace(name string, handler func(*gorm.DB)) error
}

// getStatement 将错误信息收敛为短状态标签，避免错误文案整体进入指标导致标签基数过高。
func getStatement(err string) string {
	if !strings.HasPrefix(err, "Errord") {
		return "Unknown"
	}
	slice := strings.Split(err, ":")
	if len(slice) < 2 {
		return "Unknown"
	}
	return slice[0]
}

// RegisterInterceptor 将给定拦截器依次包装到 GORM 的 Create、Query、Delete、Update、
// Row、Raw 六个回调阶段；替换阶段处理器失败时直接 panic，属于初始化期配置错误。
//
// 各拦截器按传入顺序叠加，后注册的在外层包装。
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

// MetricInterceptor 返回采集 SQL 操作计数的拦截器：执行原处理器后按是否出错上报
// 成功或错误计数，错误状态收敛为短标签以控制指标标签基数。
func MetricInterceptor() Interceptor {
	return func(op string, client *Client, next Handler) Handler {
		return func(scope *gorm.DB) {
			next(scope)
			if scope.Error != nil {
				metric.LibHandleCounter.WithLabelValues(client.DriverName, client.Database+"."+scope.Name(), client.Host, getStatement(scope.Error.Error())).Inc()
			} else {
				metric.LibHandleCounter.WithLabelValues(client.DriverName, client.Database+"."+scope.Name(), client.Host, "OK").Inc()
			}
		}
	}
}

// TraceInterceptor 返回挂载客户端链路追踪的拦截器：仅当请求上下文非空时开启 Span，
// 在 Span 上记录数据库名、主机、用户名与 SQL 语句，并将错误写入 Span 状态；
// 上下文为空时直接透传给原处理器。
func TraceInterceptor() Interceptor {
	tracer := globalTrace.NewTracer(trace.SpanKindClient)
	attrs := []attribute.KeyValue{
		semconv.RPCSystemKey.String("gorm"),
		semconv.HostName(hostname.Hostname()),
	}
	return func(op string, client *Client, next Handler) Handler {
		return func(scope *gorm.DB) {
			if ctx := scope.Statement.Context; ctx != nil {
				stmt := scope.Statement
				next(scope)
				sqlRaw := scope.Dialector.Explain(stmt.SQL.String(), stmt.Vars...)
				fn, file, line := call.Caller(4)
				attrs = append(attrs,
					semconv.CodeFunction(fn),
					semconv.CodeFilepath(file),
					semconv.CodeLineNumber(line),
				)
				md, ok := metadata.FromIncomingContext(ctx)
				if ok {
					md = md.Copy()
				} else {
					md = metadata.MD{}
				}
				_, span := tracer.Start(ctx, op, globalTrace.MetadataReaderWriter(md), trace.WithAttributes(attrs...))
				span.SetAttributes(
					semconv.DBNameKey.String(client.DriverName),
					semconv.DBConnectionStringKey.String(client.Host),
					semconv.DBUserKey.String(client.Username),
					semconv.DBStatementKey.String(sqlRaw),
				)
				defer span.End()
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
