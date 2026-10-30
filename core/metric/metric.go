package metric

import (
	"time"

	"github.com/abulo/ratel/v3/core/constant"
	"github.com/abulo/ratel/v3/core/env"
	"github.com/prometheus/client_golang/prometheus"
)

// Metrics 组合 Prometheus 的注册器与采集器，并补充批量注册采集器的能力。
type Metrics interface {
	prometheus.Registerer
	prometheus.Gatherer

	BulkRegister(...prometheus.Collector) error
}

// 调用类型标识，作为各指标 type 标签的取值。
var (
	// TypeHTTP 表示 HTTP 服务或 HTTP 调用。
	TypeHTTP = "http"
	// TypeGRPCUnary 表示 gRPC 一元调用。
	TypeGRPCUnary = "unary"
	// TypeGRPCStream 表示 gRPC 流式调用。
	TypeGRPCStream = "stream"
	// TypeRedis 表示 Redis 访问。
	TypeRedis = "redis"
	// TypeWebsocket 表示 WebSocket 连接。
	TypeWebsocket = "ws"
	// TypeMySQL 表示 MySQL 访问。
	TypeMySQL = "mysql"

	// 任务执行结果标识，作为任务类指标 code 标签的取值。
	// CodeJobSuccess 表示任务执行成功。
	CodeJobSuccess = "ok"
	// CodeJobFail 表示任务执行失败。
	CodeJobFail = "fail"
	// CodeJobReentry 表示任务重复进入。
	CodeJobReentry = "reentry"

	// 缓存访问结果标识，作为缓存类指标 code 标签的取值。
	// CodeCacheMiss 表示缓存未命中。
	CodeCacheMiss = "miss"
	// CodeCacheHit 表示缓存命中。
	CodeCacheHit = "hit"
)

// 预置的全局指标向量，进程启动时完成注册，业务代码直接按标签上报。
var (
	// ServerHandleCounter 统计服务端处理请求次数，标签为 type、method、client、code。
	ServerHandleCounter = CounterVecOpts{
		Namespace: constant.DefaultNamespace,
		Name:      "server_handle_total",
		Labels:    []string{"type", "method", "client", "code"},
	}.Build()

	// ServerHandleHistogram 统计服务端请求耗时分布，标签为 type、method、client。
	ServerHandleHistogram = HistogramVecOpts{
		Namespace: constant.DefaultNamespace,
		Name:      "server_handle_seconds",
		Labels:    []string{"type", "method", "client"},
	}.Build()

	// ClientHandleCounter 统计客户端发起调用次数，标签为 type、name、method、server、code。
	ClientHandleCounter = CounterVecOpts{
		Namespace: constant.DefaultNamespace,
		Name:      "client_handle_total",
		Labels:    []string{"type", "name", "method", "server", "code"},
	}.Build()

	// ClientHandleHistogram 统计客户端调用耗时分布，标签为 type、name、method、server。
	ClientHandleHistogram = HistogramVecOpts{
		Namespace: constant.DefaultNamespace,
		Name:      "client_handle_seconds",
		Labels:    []string{"type", "name", "method", "server"},
	}.Build()

	// JobHandleCounter 统计定时任务执行次数，标签为 type、name、code。
	JobHandleCounter = CounterVecOpts{
		Namespace: constant.DefaultNamespace,
		Name:      "job_handle_total",
		Labels:    []string{"type", "name", "code"},
	}.Build()

	// JobHandleHistogram 统计定时任务执行耗时分布，标签为 type、name。
	JobHandleHistogram = HistogramVecOpts{
		Namespace: constant.DefaultNamespace,
		Name:      "job_handle_seconds",
		Labels:    []string{"type", "name"},
	}.Build()

	// LibHandleHistogram 统计第三方库调用耗时分布，标签为 type、method、address。
	LibHandleHistogram = HistogramVecOpts{
		Namespace: constant.DefaultNamespace,
		Name:      "lib_handle_seconds",
		Labels:    []string{"type", "method", "address"},
	}.Build()

	// LibHandleCounter 统计第三方库调用次数，标签为 type、method、address、code。
	LibHandleCounter = CounterVecOpts{
		Namespace: constant.DefaultNamespace,
		Name:      "lib_handle_total",
		Labels:    []string{"type", "method", "address", "code"},
	}.Build()

	// LibHandleSummary 汇总第三方库调用的分位数统计，标签为 name、status。
	LibHandleSummary = SummaryVecOpts{
		Namespace: constant.DefaultNamespace,
		Name:      "lib_handle_stats",
		Labels:    []string{"name", "status"},
	}.Build()

	// CacheHandleCounter 统计缓存访问次数，标签为 type、name、method、code。
	CacheHandleCounter = CounterVecOpts{
		Namespace: constant.DefaultNamespace,
		Name:      "cache_handle_total",
		Labels:    []string{"type", "name", "method", "code"},
	}.Build()

	// CacheHandleHistogram 统计缓存访问耗时分布，标签为 type、name、method。
	CacheHandleHistogram = HistogramVecOpts{
		Namespace: constant.DefaultNamespace,
		Name:      "cache_handle_seconds",
		Labels:    []string{"type", "name", "method"},
	}.Build()

	// BuildInfoGauge 记录当前进程的构建信息，标签为 name、id、env、region、zone、version、go_version。
	BuildInfoGauge = GaugeVecOpts{
		Namespace: constant.DefaultNamespace,
		Name:      "build_info",
		Labels:    []string{"name", "id", "env", "region", "zone", "version", "go_version"},
	}.Build()
)

func init() {
	BuildInfoGauge.WithLabelValues(
		env.Name(),
		env.AppID(),
		env.AppMode(),
		env.AppRegion(),
		env.AppZone(),
		env.BuildVersion(),
		env.GoVersion(),
	).Set(float64(time.Now().UnixNano() / 1e6))
}
