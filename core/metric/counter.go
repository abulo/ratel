package metric

import (
	"sync"

	"github.com/abulo/ratel/v3/core/constant"
	"github.com/prometheus/client_golang/prometheus"
)

// CounterVecOpts 定义计数器向量的 Prometheus 配置项。
type CounterVecOpts struct {
	// Namespace 为指标命名空间，NewCounterVec 使用 constant.DefaultNamespace。
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

// Build 按配置创建计数器向量并注册到 Prometheus 默认注册器；同名向量重复注册会触发 panic。
func (opts CounterVecOpts) Build() *counterVec {
	vec := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: opts.Namespace,
			Subsystem: opts.Subsystem,
			Name:      opts.Name,
			Help:      opts.Help,
		}, opts.Labels)
	prometheus.MustRegister(vec)
	return &counterVec{
		CounterVec: vec,
	}
}

// NewCounterVec 以默认命名空间、给定名称与标签快速创建计数器向量，Help 文本与指标名称相同。
func NewCounterVec(name string, labels []string) *counterVec {
	return CounterVecOpts{
		Namespace: constant.DefaultNamespace,
		Name:      name,
		Help:      name,
		Labels:    labels,
	}.Build()
}

type counterVec struct {
	*prometheus.CounterVec
	// PERF-3-14 [P2]：按标签组合缓存已解析的 Counter 子指标，热点上报命中缓存后不再重复执行向量内查找；
	// 首次未命中仍走 WithLabelValues，标签数量不匹配时由 Prometheus 原样 panic，错误行为不变。
	mu       sync.RWMutex
	children map[string]prometheus.Counter
}

// getChild 返回 labels 对应的 Counter 子指标，未命中时经 WithLabelValues 解析并回填缓存。
func (counter *counterVec) getChild(labels []string) prometheus.Counter {
	key := labelKey(labels)
	counter.mu.RLock()
	child, ok := counter.children[key]
	counter.mu.RUnlock()
	if ok {
		return child
	}
	child = counter.WithLabelValues(labels...)
	counter.mu.Lock()
	if counter.children == nil {
		counter.children = make(map[string]prometheus.Counter)
	}
	counter.children[key] = child
	counter.mu.Unlock()
	return child
}

// Inc 将对应标签组合的计数器加一；标签值数量与顺序须与向量定义的 Labels 一致。
func (counter *counterVec) Inc(labels ...string) {
	counter.getChild(labels).Inc()
}

// Add 将对应标签组合的计数器增加 v；计数器只允许单调递增，调用方不得传入负值。
func (counter *counterVec) Add(v float64, labels ...string) {
	counter.getChild(labels).Add(v)
}
