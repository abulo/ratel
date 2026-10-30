package metric

import "github.com/prometheus/client_golang/prometheus"

// HistogramVecOpts 定义直方图向量的 Prometheus 配置项。
type HistogramVecOpts struct {
	// Namespace 为指标命名空间。
	Namespace string
	// Subsystem 为指标子系统名称。
	Subsystem string
	// Name 为指标名称。
	Name      string
	// Help 为指标说明文本，展示在 Prometheus 面板中。
	Help      string
	// Labels 为标签名列表，顺序须与上报时传入的标签值顺序一致。
	Labels    []string
	// Buckets 为直方图上界数组；为空时使用 Prometheus 默认桶划分。
	Buckets   []float64
}

type histogramVec struct {
	*prometheus.HistogramVec
}

// Build 按配置创建直方图向量并注册到 Prometheus 默认注册器；同名向量重复注册会触发 panic。
func (opts HistogramVecOpts) Build() *histogramVec {
	vec := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: opts.Namespace,
			Subsystem: opts.Subsystem,
			Name:      opts.Name,
			Help:      opts.Help,
			Buckets:   opts.Buckets,
		}, opts.Labels)
	prometheus.MustRegister(vec)
	return &histogramVec{
		HistogramVec: vec,
	}
}

// Observe 记录一个观测值 v 到对应标签组合的直方图；标签值数量与顺序须与向量定义的 Labels 一致。
func (histogram *histogramVec) Observe(v float64, labels ...string) {
	histogram.WithLabelValues(labels...).Observe(v)
}
