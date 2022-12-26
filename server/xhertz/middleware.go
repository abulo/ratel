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

	"github.com/abulo/ratel/v2/core/logger"
	"github.com/abulo/ratel/v2/core/metric"
	"github.com/abulo/ratel/v2/core/trace"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

func metricServerInterceptor() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		beg := time.Now()
		ctx.Next(c)
		metric.ServerHandleHistogram.Observe(time.Since(beg).Seconds(), metric.TypeHTTP, string(ctx.Request.Method())+"."+string(ctx.Request.Path()), extractAID(ctx))
		metric.ServerHandleCounter.Inc(metric.TypeHTTP, string(ctx.Request.Method())+"."+string(ctx.Request.Path()), extractAID(ctx), http.StatusText(ctx.Response.StatusCode()))
	}
}

func traceServerInterceptor() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		span, ctxNew := trace.StartSpanFromContext(
			c,
			trace.SpanHertzHttpStartName(ctx),
			trace.TagComponent("http"),
			trace.TagSpanKind("server"),
			// trace.HeaderExtractor(ctx.Request.Header.GetHeaders()),
			trace.CustomTag("http.url", string(ctx.Request.Path())),
			trace.CustomTag("http.method", string(ctx.Request.Method())),
			trace.CustomTag("peer.ipv4", ctx.ClientIP()),
		)
		c = ctxNew
		defer span.Finish()
		ctx.Next(c)
	}
}

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

func extractAID(ctx *app.RequestContext) string {
	return ctx.Request.Header.Get("AID")
}

// HTTPHeader is provided to wrap an http.Header into an HTTPHeaderCarrier.
type HTTPHeader map[string][]string

// Visit implements the HTTPHeaderCarrier interface.
func (h HTTPHeader) Visit(v func(k, v string)) {
	for k, vs := range h {
		v(k, vs[0])
	}
}

// Set sets the header entries associated with key to the single element value.
// The key will converted into lowercase as the HTTP/2 protocol requires.
func (h HTTPHeader) Set(key, value string) {
	h[strings.ToLower(key)] = []string{value}
}

// HTTPHeaderCarrier accepts a visitor to access all key value pairs in an HTTP header.
type HTTPHeaderCarrier interface {
	Visit(func(k, v string))
}

// HTTPHeaderSetter sets a key with a value into a HTTP header.
type HTTPHeaderSetter interface {
	Set(key, value string)
}

// HTTPHeaderToCGIVariable performs an CGI variable conversion.
// For example, an HTTP header key `abc-def` will result in `ABC_DEF`.
func HTTPHeaderToCGIVariable(key string) string {
	return strings.ToUpper(strings.ReplaceAll(key, "-", "_"))
}

// CGIVariableToHTTPHeader converts a CGI variable into an HTTP header key.
// For example, `ABC_DEF` will be converted to `abc-def`.
func CGIVariableToHTTPHeader(key string) string {
	return strings.ToLower(strings.ReplaceAll(key, "_", "-"))
}
