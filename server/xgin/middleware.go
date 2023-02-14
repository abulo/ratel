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

	"github.com/abulo/ratel/core/logger"
	"github.com/abulo/ratel/core/metric"
	"github.com/abulo/ratel/core/trace"
	"github.com/abulo/ratel/gin"
	"github.com/sirupsen/logrus"
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

func metricServerInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		beg := time.Now()
		c.Next()
		metric.ServerHandleHistogram.Observe(time.Since(beg).Seconds(), metric.TypeHTTP, c.Request.Method+"."+c.Request.URL.Path, extractAID(c))
		metric.ServerHandleCounter.Inc(metric.TypeHTTP, c.Request.Method+"."+c.Request.URL.Path, extractAID(c), http.StatusText(c.Writer.Status()))
	}
}

func traceServerInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := trace.StartSpanFromContext(
			c.Request.Context(),
			trace.SpanGinHttpStartName(c),
			trace.TagComponent("http"),
			trace.TagSpanKind("server"),
			trace.HeaderExtractor(c.Request.Header),
			trace.CustomTag("http.url", c.Request.URL.Path),
			trace.CustomTag("http.method", c.Request.Method),
			trace.CustomTag("peer.ipv4", c.ClientIP()),
		)
		c.Request = c.Request.WithContext(ctx)
		defer span.Finish()
		c.Next()
	}
}

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

// stack returns a nicely formatted stack frame, skipping skip frames.
func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
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

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastSlash := bytes.LastIndex(name, slash); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}
