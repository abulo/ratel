package xhertz

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/abulo/ratel/v3/core/hostname"
	"github.com/abulo/ratel/v3/core/logger"
	"github.com/abulo/ratel/v3/core/metric"
	globalTrace "github.com/abulo/ratel/v3/core/trace"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/errors"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

// metadataProvider 将 Hertz 请求头适配为链路追踪所需的元数据读写载体。
type metadataProvider struct {
	metadata map[string]string
	headers  *protocol.RequestHeader
}

func (m *metadataProvider) Get(key string) string {
	return m.headers.Get(key)
}

func (m *metadataProvider) Set(key, value string) {
	m.headers.Set(key, value)
}

func (m *metadataProvider) Keys() []string {
	out := make([]string, 0, len(m.metadata))

	m.headers.VisitAll(func(key, value []byte) {
		out = append(out, string(key))
	})

	return out
}

// metricServerInterceptor 返回 HTTP 服务指标采集中间件：请求处理结束后记录耗时直方图与请求计数器。
func metricServerInterceptor() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		beg := time.Now()
		ctx.Next(c)
		metric.ServerHandleHistogram.Observe(time.Since(beg).Seconds(), metric.TypeHTTP, string(ctx.Request.Method())+"."+string(ctx.Request.Path()), extractAID(ctx))
		metric.ServerHandleCounter.Inc(metric.TypeHTTP, string(ctx.Request.Method())+"."+string(ctx.Request.Path()), extractAID(ctx), http.StatusText(ctx.Response.StatusCode()))
	}
}

// traceServerInterceptor 返回服务端链路追踪中间件：创建服务端 Span、填充 HTTP 属性并结束 Span。
func traceServerInterceptor() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		tracer := globalTrace.NewTracer(trace.SpanKindServer)
		attrs := []attribute.KeyValue{
			semconv.RPCSystemKey.String("http"),
			semconv.HostName(hostname.Hostname()),
		}
		headers := &ctx.Request.Header
		provider := &metadataProvider{headers: headers}
		c, span := tracer.Start(c, cast.ToString(ctx.Request.Method())+" "+cast.ToString(ctx.Request.Path()), provider, trace.WithAttributes(attrs...))
		span.SetAttributes(
			semconv.HTTPTargetKey.String(cast.ToString(ctx.Request.Path())),
			semconv.HTTPMethodKey.String(cast.ToString(ctx.Request.Method())),
			semconv.HTTPUserAgentKey.String(cast.ToString(ctx.Request.Header.UserAgent())),
			semconv.NetPeerNameKey.String(cast.ToString(ctx.ClientIP())),
		)
		defer span.End()
		ctx.Next(c)
	}
}

// recoverMiddleware 返回 panic 恢复与访问日志中间件：捕获 panic 并记录错误与调用栈，
// 识别客户端断管（broken pipe、connection reset）后直接中止，其余 panic 返回 500；
// 请求耗时超过阈值时在日志中标记慢查询。
func recoverMiddleware(slowQueryThresholdInMilli int64) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {

		var beg = time.Now()
		fields := make(logrus.Fields)
		var brokenPipe bool
		defer func() {
			fields["cost"] = time.Since(beg).Seconds()
			if slowQueryThresholdInMilli > 0 {
				if cost := int64(time.Since(beg)) / 1e6; cost > slowQueryThresholdInMilli {
					fields["slow"] = cost
				}
			}
			if rec := recover(); rec != nil {
				if ne, ok := rec.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				var err = rec.(error)
				fields["stack"] = stack(3)
				fields["err"] = err.Error()
				logger.Logger.WithFields(fields).Error("access")
				if brokenPipe {
					_ = ctx.Error(err)
					ctx.Abort()
					return
				}
				ctx.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			fields["method"] = cast.ToString(ctx.Request.Method())
			fields["code"] = ctx.Response.StatusCode()
			fields["size"] = len(ctx.Response.BodyBytes())
			fields["host"] = cast.ToString(ctx.Request.Host())
			fields["path"] = cast.ToString(ctx.Request.Path())
			fields["ip"] = ctx.ClientIP()
			fields["err"] = ctx.Errors.ByType(errors.ErrorTypePrivate).String()
			logger.Logger.WithFields(fields).Info("access")
		}()
		ctx.Next(c)
	}
}

// stack 返回跳过 skip 帧后格式化的调用栈，包含文件、行号、函数名与对应源码行；源码不可读时用占位符表示。
func stack(skip int) []byte {
	buf := new(bytes.Buffer)
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// 找不到源码文件时至少输出文件、行号与程序计数器地址。
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := os.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

// function 返回程序计数器 pc 对应函数的简称：剥离包路径前缀并将包名中的中点替换为点。
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// runtime.FuncForPC 返回的全名带有包路径前缀，而文件路径已经输出在前一行；
	// 同时先按斜杠去除路径前缀，再按点去除包名，避免包名本身含点时误删。
	if lastSlash := bytes.LastIndex(name, slash); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}

// source 返回源码行切片中第 n 行去除首尾空白后的内容；行号越界时返回 "???" 占位符。
func source(lines [][]byte, n int) []byte {
	// 调用栈中行号从 1 开始计数，而切片下标从 0 开始。
	n--
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

func extractAID(ctx *app.RequestContext) string {
	return ctx.Request.Header.Get("AID")
}

// HTTPHeader 将标准库 http.Header 适配为键到单值切片的映射，供链路追踪载体使用。
type HTTPHeader map[string][]string

// Visit 遍历映射中每个键的第一个值并调用回调，实现 HTTPHeaderCarrier 接口。
func (h HTTPHeader) Visit(v func(k, v string)) {
	for k, vs := range h {
		v(k, vs[0])
	}
}

// Set 将指定键的值设置为单元素切片；键名统一转为小写，符合 HTTP/2 协议要求。
func (h HTTPHeader) Set(key, value string) {
	h[strings.ToLower(key)] = []string{value}
}

// HTTPHeaderCarrier 抽象对 HTTP 头部全部键值对的遍历访问能力。
type HTTPHeaderCarrier interface {
	Visit(func(k, v string))
}

// HTTPHeaderSetter 抽象向 HTTP 头部写入单个键值对的能力。
type HTTPHeaderSetter interface {
	Set(key, value string)
}

// HTTPHeaderToCGIVariable 将 HTTP 头部键名转换为 CGI 变量名：连字符转为下划线并整体大写，
// 例如 abc-def 转换为 ABC_DEF。
func HTTPHeaderToCGIVariable(key string) string {
	return strings.ToUpper(strings.ReplaceAll(key, "-", "_"))
}

// CGIVariableToHTTPHeader 将 CGI 变量名转换回 HTTP 头部键名：下划线转为连字符并整体小写，
// 例如 ABC_DEF 转换为 abc-def。
func CGIVariableToHTTPHeader(key string) string {
	return strings.ToLower(strings.ReplaceAll(key, "_", "-"))
}
