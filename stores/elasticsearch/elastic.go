package elasticsearch

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/abulo/ratel/v3/logger"
	"github.com/abulo/ratel/v3/metric"
	"github.com/abulo/ratel/v3/trace"
	"github.com/olivere/elastic/v7"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

const MaxContentLength = 1 << 16

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
		options = append(options, elastic.SetHttpClient(ESTraceServerInterceptor(config.DisableMetric, config.DisableTrace)))
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

func ESTraceServerInterceptor(DisableMetric, DisableTrace bool) *http.Client {
	newESTracedTransport := &ESTracedTransport{}
	newESTracedTransport.Transport = &http.Transport{}
	newESTracedTransport.DisableMetric = DisableMetric
	newESTracedTransport.DisableTrace = DisableTrace
	return &http.Client{
		Transport: newESTracedTransport,
	}
}

type ESTracedTransport struct {
	*http.Transport
	DisableMetric bool
	DisableTrace  bool
}

func (t *ESTracedTransport) RoundTrip(r *http.Request) (resp *http.Response, err error) {
	start := time.Now()
	var span opentracing.Span
	if !t.DisableTrace {
		span, ctx := trace.StartSpanFromContext(
			r.Context(),
			"elastic",
			trace.CustomTag("peer.service", "elastic"),
			trace.TagSpanKind("client"),
			trace.HeaderExtractor(r.Header),
			trace.CustomTag("http.url", r.URL.Path),
			trace.CustomTag("http.method", r.Method),
		)

		defer func() {
			if err != nil {
				span.SetTag("elastic.error", err.Error())
				span.SetTag(string(ext.Error), true)
			}
			span.Finish()
		}()
		span.SetTag(string(ext.DBType), "elastic")
		span.SetTag(string(ext.DBInstance), r.URL.Host)
		span.SetTag("elastic.method", r.Method)
		span.SetTag("elastic.url", r.URL.Path)
		span.SetTag("elastic.params", r.URL.Query().Encode())
		r = r.WithContext(ctx)
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
	if !t.DisableTrace {
		span.SetTag("elastic.status_code", resp.StatusCode)
	}

	if !t.DisableMetric {
		cost := time.Since(start)
		if err != nil {
			metric.LibHandleCounter.WithLabelValues("elastic", "ERR").Inc()
		} else {
			metric.LibHandleCounter.Inc("elastic", "OK")
		}
		metric.LibHandleHistogram.WithLabelValues("elastic").Observe(cost.Seconds())
	}
	return resp, err
}
