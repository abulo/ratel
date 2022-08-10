package elasticsearch

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/abulo/ratel/v3/logger"
	"github.com/abulo/ratel/v3/metric"
	"github.com/abulo/ratel/v3/trace"
	"github.com/abulo/ratel/v3/util"
	"github.com/olivere/elastic/v7"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

// MaxContentLength ...
const MaxContentLength = 1 << 16

// Config ...
type Config struct {
	URL           []string
	Username      string //账号 root
	Password      string //密码
	DisableMetric bool   // 关闭指标采集
	DisableTrace  bool   // 关闭链路追踪
}

// Client --
type Client struct {
	*elastic.Client
	*Config
}

// NewClient ...
func NewClient(config *Config) *Client {
	var options []elastic.ClientOptionFunc
	if len(config.URL) < 1 {
		logger.Logger.Panic("url not set")
	}
	options = append(options, elastic.SetURL(config.URL...))
	if config.Username != "" && config.Password != "" {
		options = append(options, elastic.SetBasicAuth(config.Username, config.Password))
	}
	if !config.DisableTrace || !config.DisableTrace {
		options = append(options, elastic.SetHttpClient(ESTraceServerInterceptor(config.DisableMetric, config.DisableTrace, util.Implode(";", config.URL))))
	}
	options = append(options, elastic.SetSniff(false))
	client, err := elastic.NewClient(options...)
	if err != nil {
		logger.Logger.Panic(err)
	}
	newClient := &Client{}
	newClient.Client = client
	newClient.Config = config
	return newClient
}

// ESTraceServerInterceptor ...
func ESTraceServerInterceptor(DisableMetric, DisableTrace bool, Addr string) *http.Client {
	newESTracedTransport := &ESTracedTransport{}
	newESTracedTransport.Transport = &http.Transport{}
	newESTracedTransport.DisableMetric = DisableMetric
	newESTracedTransport.DisableTrace = DisableTrace
	newESTracedTransport.Addr = Addr
	return &http.Client{
		Transport: newESTracedTransport,
	}
}

// ESTracedTransport ...
type ESTracedTransport struct {
	*http.Transport
	DisableMetric bool
	DisableTrace  bool
	Addr          string
}

// RoundTrip ...
func (t *ESTracedTransport) RoundTrip(r *http.Request) (resp *http.Response, err error) {
	start := time.Now()
	var span opentracing.Span
	if !t.DisableTrace {
		if parentSpan := trace.SpanFromContext(r.Context()); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span = opentracing.StartSpan("elastic", opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			hostName, err := os.Hostname()
			if err != nil {
				hostName = "unknown"
			}
			ext.PeerHostname.Set(span, hostName)
			defer func() {
				if err != nil {
					span.SetTag("elastic.error", err.Error())
					span.SetTag(string(ext.Error), true)
				}
				span.Finish()
			}()
			ctx := opentracing.ContextWithSpan(r.Context(), span)
			span.SetTag(string(ext.DBType), "elastic")
			span.SetTag(string(ext.DBInstance), r.URL.Host)
			span.SetTag("elastic.method", r.Method)
			span.SetTag("elastic.url", r.URL.Path)
			span.SetTag("elastic.params", r.URL.Query().Encode())
			r = r.WithContext(ctx)
		}
	}
	contentLength, _ := strconv.Atoi(r.Header.Get("Content-Length"))
	if r.Body != nil && contentLength < MaxContentLength {
		buf, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		if !t.DisableTrace {
			span.SetTag(string(ext.DBStatement), string(buf))
			span.LogFields(log.String("params", string(buf)))
		}
		r.Body = io.NopCloser(bytes.NewBuffer(buf))
	}
	resp, err = t.Transport.RoundTrip(r)
	if err != nil {
		return nil, err
	}
	if !t.DisableMetric {
		cost := time.Since(start)
		if err != nil {
			metric.LibHandleCounter.WithLabelValues("elastic", "elastic", t.Addr, "ERR").Inc()
		} else {
			metric.LibHandleCounter.Inc("elastic", "elastic", t.Addr, "OK")
		}
		metric.LibHandleHistogram.WithLabelValues("elastic", "elastic", t.Addr).Observe(cost.Seconds())
	}
	return resp, err
}
