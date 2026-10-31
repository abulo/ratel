package xgin

import (
	"bytes"
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
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

func extractAID(ctx *gin.Context) string {
	return ctx.Request.Header.Get("AID")
}

// metricServerInterceptor 返回 HTTP 服务指标采集中间件：请求处理结束后记录耗时直方图与请求计数器。
func metricServerInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		beg := time.Now()
		c.Next()
		// PERF-6-4 [P1] method/path/aid 与标签串每请求只构造一次；路由模板替代原始 Path 控制指标标签基数，
		// 未匹配路由（FullPath 为空，如 404）回退原始 Path，保持原行为。
		method := c.Request.Method
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		label := method + "." + path
		aid := extractAID(c)
		metric.ServerHandleHistogram.Observe(time.Since(beg).Seconds(), metric.TypeHTTP, label, aid)
		metric.ServerHandleCounter.Inc(metric.TypeHTTP, label, aid, http.StatusText(c.Writer.Status()))
	}
}

// traceServerInterceptor 返回服务端链路追踪中间件：创建服务端 Span、填充 HTTP 属性并将追踪上下文注入请求。
func traceServerInterceptor() gin.HandlerFunc {
	// PERF-6-5 [P1] tracer 与主机名在装配期获取一次并闭包复用，避免每请求重复查 tracer 注册表与调用 os.Hostname()；
	// 恒定属性在工厂层构造一次，请求级属性在闭包内复用局部变量，避免重复取值与字符串分配。
	// （hostname 系统调用的实际 CPU 占比需实测验证）
	tracer := globalTrace.NewTracer(trace.SpanKindServer)
	staticAttrs := []attribute.KeyValue{
		semconv.RPCSystemKey.String("http"),
		semconv.HostName(hostname.Hostname()),
	}
	return func(c *gin.Context) {
		target := c.Request.URL.Path
		method := c.Request.Method
		ctx, span := tracer.Start(c.Request.Context(), target, propagation.HeaderCarrier(c.Request.Header), trace.WithAttributes(staticAttrs...))
		span.SetAttributes(
			semconv.HTTPURLKey.String(c.Request.URL.String()),
			semconv.HTTPTargetKey.String(target),
			semconv.HTTPMethodKey.String(method),
			semconv.HTTPUserAgentKey.String(c.Request.UserAgent()),
			semconv.NetPeerNameKey.String(c.ClientIP()),
		)
		defer span.End()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// recoverMiddleware 返回 panic 恢复与访问日志中间件：捕获 panic 并记录错误与调用栈，
// 识别客户端断管（broken pipe、connection reset）后直接中止，其余 panic 返回 500；
// 请求耗时超过阈值时在日志中标记慢查询。
func recoverMiddleware(slowQueryThresholdInMilli int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		var beg = time.Now()
		var fields logrus.Fields
		var brokenPipe bool
		defer func() {
			// PERF-6-6 [P1] 耗时只计算一次；fields 惰性初始化，成功且不写额外字段时不再分配空 map。
			elapsed := time.Since(beg)
			if slowQueryThresholdInMilli > 0 {
				if cost := elapsed.Milliseconds(); cost > slowQueryThresholdInMilli {
					if fields == nil {
						fields = make(logrus.Fields)
					}
					fields["slow"] = cost
				}
			}
			if rec := recover(); rec != nil {
				if ne, ok := rec.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						// PERF-6-6 [P1] 错误文本只小写化一次，避免两次 ToLower 分配。
						lowerErr := strings.ToLower(se.Error())
						if strings.Contains(lowerErr, "broken pipe") || strings.Contains(lowerErr, "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				err := rec.(error)
				if fields == nil {
					fields = make(logrus.Fields)
				}
				fields["cost"] = elapsed.Seconds()
				fields["stack"] = stack(3)
				fields["err"] = err.Error()
				logger.Logger.WithFields(fields).Error("access")
				if brokenPipe {
					_ = c.Error(err)
					c.Abort()
					return
				}
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			privateErrs := c.Errors.ByType(gin.ErrorTypePrivate)
			errStr := ""
			if len(privateErrs) > 0 {
				errStr = privateErrs.String()
			}
			// PERF-6-6 [P1] ClientIP 只解析一次；错误列表仅在存在私有错误时拼接。
			ip := c.ClientIP()
			if fields == nil {
				fields = make(logrus.Fields)
			}
			fields["cost"] = elapsed.Seconds()
			fields["method"] = c.Request.Method
			fields["code"] = c.Writer.Status()
			fields["size"] = c.Writer.Size()
			fields["host"] = c.Request.Host
			fields["path"] = c.Request.URL.Path
			fields["ip"] = ip
			fields["err"] = errStr
			logger.Logger.WithFields(fields).Info("access")
		}()
		c.Next()
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
