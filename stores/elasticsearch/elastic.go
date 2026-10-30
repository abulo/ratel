package elasticsearch

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/abulo/ratel/v3/core/call"
	"github.com/abulo/ratel/v3/core/logger"
	"github.com/abulo/ratel/v3/core/metric"
	"github.com/abulo/ratel/v3/core/resource"
	globalTrace "github.com/abulo/ratel/v3/core/trace"
	"github.com/abulo/ratel/v3/util"
	"github.com/olivere/elastic/v7"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

// MaxContentLength 为允许读取进内存并参与追踪记录的请求体大小上限，单位字节，
// 超过该上限的请求体不再回读。
const MaxContentLength = 1 << 16

// Config 保存 Elasticsearch 客户端的连接地址、认证信息与功能开关。
type Config struct {
	// URL 为集群节点地址列表，至少配置一个；客户端通过 SetURL 传给底层 elastic.Client。
	URL []string
	// Username 为基础认证用户名，与 Password 同时非空时才启用认证。
	Username string
	// Password 为基础认证密码。
	Password string
	// EnableMetric 是否采集 Elasticsearch 操作指标。
	EnableMetric bool
	// EnableTrace 是否挂载链路追踪。
	EnableTrace bool
}

// Client 包装底层 elastic.Client 并持有连接配置，供上层直接调用 elastic v7 的查询 API。
//
// 实例由 NewClient 创建；Close 当前为空操作，不主动断开底层连接。
type Client struct {
	*elastic.Client
	*Config
}

// Close 释放客户端资源；当前底层 elastic.Client 不提供断开方法，该方法为空操作并返回 nil。
func (client *Client) Close() error {
	return nil
}

var clientManager = resource.NewResourceManager()

// NewClient 依据 Config 创建并按节点地址缓存复用 Elasticsearch 客户端。
//
// URL 为空、客户端初始化失败时记录日志并 panic；开启指标或追踪时为底层 HTTP 客户端
// 安装带追踪与指标能力的传输层。
func NewClient(config *Config) *Client {
	url := util.Implode(";", config.URL)
	val, err := clientManager.GetResource(url, func() (io.Closer, error) {
		client, err := getClient(config)
		if err != nil {
			return nil, err
		}
		return client, nil
	})
	if err != nil {
		logger.Logger.Panic(err)
		return nil
	}
	return val.(*Client)
}

// getClient 依据 Config 构造底层 elastic.Client 并包装为 Client；URL 为空或初始化失败时
// 记录日志并返回错误。
func getClient(config *Config) (*Client, error) {
	var options []elastic.ClientOptionFunc
	if len(config.URL) < 1 {
		logger.Logger.Panic("url not set")
		return nil, errors.New("url not set")
	}
	options = append(options, elastic.SetURL(config.URL...))
	if config.Username != "" && config.Password != "" {
		options = append(options, elastic.SetBasicAuth(config.Username, config.Password))
	}
	if config.EnableMetric || config.EnableTrace {
		options = append(options, elastic.SetHttpClient(ESTraceServerInterceptor(config.EnableMetric, config.EnableTrace, util.Implode(";", config.URL))))
	}
	options = append(options, elastic.SetSniff(false))
	client, err := elastic.NewClient(options...)
	if err != nil {
		logger.Logger.Panic(err)
		return nil, err
	}

	newClient := &Client{}
	newClient.Client = client
	newClient.Config = config
	return newClient, nil
}

// ESTraceServerInterceptor 构造带追踪与指标传输层的 http.Client；EnableMetric、EnableTrace
// 为对应开关，Addr 为节点地址，用于指标标签。
func ESTraceServerInterceptor(EnableMetric, EnableTrace bool, Addr string) *http.Client {
	newESTracedTransport := &ESTracedTransport{}
	newESTracedTransport.Transport = &http.Transport{}
	newESTracedTransport.EnableMetric = EnableMetric
	newESTracedTransport.EnableTrace = EnableTrace
	newESTracedTransport.Addr = Addr
	return &http.Client{
		Transport: newESTracedTransport,
	}
}

// ESTracedTransport 包装底层 http.Transport，在请求转发前后采集链路追踪与操作指标。
//
// 实例由 ESTraceServerInterceptor 创建；EnableMetric、EnableTrace 决定对应能力是否开启，
// Addr 作为指标中的节点标签。
type ESTracedTransport struct {
	*http.Transport
	// EnableMetric 是否采集操作指标。
	EnableMetric bool
	// EnableTrace 是否挂载链路追踪。
	EnableTrace bool
	// Addr 为节点地址，用于指标标签。
	Addr string
}

// RoundTrip 转发单个 HTTP 请求：开启追踪时创建客户端 Span 并记录方法与路径；
// 请求体未超过 MaxContentLength 时回读以便后续复用；随后交由底层 Transport 执行。
func (t *ESTracedTransport) RoundTrip(r *http.Request) (resp *http.Response, err error) {
	start := time.Now()
	if t.EnableTrace {
		tracer := globalTrace.NewTracer(trace.SpanKindClient)
		attrs := []attribute.KeyValue{
			semconv.RPCSystemKey.String("elastic"),
		}
		fn, file, line := call.Caller(6)
		attrs = append(attrs,
			semconv.CodeFunction(fn),
			semconv.CodeFilepath(file),
			semconv.CodeLineNumber(line),
		)
		md := metadata.New(nil)
		newCtx, span := tracer.Start(r.Context(), "Count", propagation.HeaderCarrier(md), trace.WithAttributes(attrs...))

		span.SetAttributes(
			attribute.String("elastic.method", r.Method),
			attribute.String("elastic.url", r.URL.Path),
		)
		defer span.End()
		r = r.WithContext(newCtx)
	}
	contentLength, _ := strconv.Atoi(r.Header.Get("Content-Length"))
	if r.Body != nil && contentLength < MaxContentLength {
		buf, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		r.Body = io.NopCloser(bytes.NewBuffer(buf))
	}
	resp, err = t.Transport.RoundTrip(r)
	if err != nil {
		return nil, err
	}
	// FIXME(owner): 指标上报条件与 EnableMetric 开关语义相反，当前在开关关闭时才上报，待修正为 t.EnableMetric。
	if !t.EnableMetric {
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
