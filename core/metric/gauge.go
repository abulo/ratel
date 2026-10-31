package metric

import (
	"sync"

	"github.com/abulo/ratel/v3/core/constant"
	"github.com/prometheus/client_golang/prometheus"
)

// GaugeVecOpts 定义仪表盘向量的 Prometheus 配置项。
type GaugeVecOpts struct {
	// Namespace 为指标命名空间，NewGaugeVec 使用 constant.DefaultNamespace。
	Namespace string
	// Subsystem 为指标子系统名称。
	Subsystem string
	// Name 为指标名称。
	Name      string
	// Help 为指标说明文本，展示在 Prometheus 面板中。
	Help      string
	// Labels 为标签名列表，顺序须与上报时传入的标签值顺序一致。
	Labels    []string
}

type gaugeVec struct {
	*prometheus.GaugeVec
	// PERF-3-14 [P2]：按标签组合缓存已解析的 Gauge 子指标，热点上报命中缓存后不再重复执行向量内查找；
	// 首次未命中仍走 WithLabelValues，标签数量不匹配时由 Prometheus 原样 panic，错误行为不变。
	mu       sync.RWMutex
	children map[string]prometheus.Gauge
}

// Build 按配置创建仪表盘向量并注册到 Prometheus 默认注册器；同名向量重复注册会触发 panic。
func (opts GaugeVecOpts) Build() *gaugeVec {
	vec := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: opts.Namespace,
			Subsystem: opts.Subsystem,
			Name:      opts.Name,
			Help:      opts.Help,
		}, opts.Labels)
	prometheus.MustRegister(vec)
	return &gaugeVec{
		GaugeVec: vec,
	}
}

// NewGaugeVec 以默认命名空间、给定名称与标签快速创建仪表盘向量，Help 文本与指标名称相同。
func NewGaugeVec(name string, labels []string) *gaugeVec {
	return GaugeVecOpts{
		Namespace: constant.DefaultNamespace,
		Name:      name,
		Help:      name,
		Labels:    labels,
	}.Build()
}

// getChild 返回 labels 对应的 Gauge 子指标，未命中时经 WithLabelValues 解析并回填缓存。
func (gv *gaugeVec) getChild(labels []string) prometheus.Gauge {
	key := labelKey(labels)
	gv.mu.RLock()
	child, ok := gv.children[key]
	gv.mu.RUnlock()
	if ok {
		return child
	}
	child = gv.WithLabelValues(labels...)
	gv.mu.Lock()
	if gv.children == nil {
		gv.children = make(map[string]prometheus.Gauge)
	}
	gv.children[key] = child
	gv.mu.Unlock()
	return child
}

// Inc 将对应标签组合的仪表盘加一；标签值数量与顺序须与向量定义的 Labels 一致。
func (gv *gaugeVec) Inc(labels ...string) {
	gv.getChild(labels).Inc()
}

// Add 将对应标签组合的仪表盘增加 v，v 可为负值；标签值数量与顺序须与向量定义的 Labels 一致。
func (gv *gaugeVec) Add(v float64, labels ...string) {
	gv.getChild(labels).Add(v)
}

// Set 将对应标签组合的仪表盘置为 v；标签值数量与顺序须与向量定义的 Labels 一致。
func (gv *gaugeVec) Set(v float64, labels ...string) {
	gv.getChild(labels).Set(v)
}
