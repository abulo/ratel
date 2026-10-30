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
		metric.ServerHandleHistogram.Observe(time.Since(beg).Seconds(), metric.TypeHTTP, c.Request.Method+"."+c.Request.URL.Path, extractAID(c))
		metric.ServerHandleCounter.Inc(metric.TypeHTTP, c.Request.Method+"."+c.Request.URL.Path, extractAID(c), http.StatusText(c.Writer.Status()))
	}
}

// traceServerInterceptor 返回服务端链路追踪中间件：创建服务端 Span、填充 HTTP 属性并将追踪上下文注入请求。
func traceServerInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		tracer := globalTrace.NewTracer(trace.SpanKindServer)
		attrs := []attribute.KeyValue{
			semconv.RPCSystemKey.String("http"),
			semconv.HostName(hostname.Hostname()),
		}

		ctx, span := tracer.Start(c.Request.Context(), c.Request.URL.Path, propagation.HeaderCarrier(c.Request.Header), trace.WithAttributes(attrs...))
		span.SetAttributes(
			semconv.HTTPURLKey.String(c.Request.URL.String()),
			semconv.HTTPTargetKey.String(c.Request.URL.Path),
			semconv.HTTPMethodKey.String(c.Request.Method),
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
					_ = c.Error(err)
					c.Abort()
					return
				}
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			fields["method"] = c.Request.Method
			fields["code"] = c.Writer.Status()
			fields["size"] = c.Writer.Size()
			fields["host"] = c.Request.Host
			fields["path"] = c.Request.URL.Path
			fields["ip"] = c.ClientIP()
			fields["err"] = c.Errors.ByType(gin.ErrorTypePrivate).String()
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
