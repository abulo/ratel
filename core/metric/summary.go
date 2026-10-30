package metric

import "github.com/prometheus/client_golang/prometheus"

// SummaryVecOpts 定义 Summary 向量的 Prometheus 配置项。
type SummaryVecOpts struct {
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
}

type summaryVec struct {
	*prometheus.SummaryVec
}

// Build 按配置创建 Summary 向量并注册到 Prometheus 默认注册器；同名向量重复注册会触发 panic。
func (opts SummaryVecOpts) Build() *summaryVec {
	vec := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: opts.Namespace,
			Subsystem: opts.Subsystem,
			Name:      opts.Name,
			Help:      opts.Help,
		}, opts.Labels)
	prometheus.MustRegister(vec)
	return &summaryVec{
		SummaryVec: vec,
	}
}

// Observe 记录一个观测值 v 到对应标签组合的 Summary；标签值数量与顺序须与向量定义的 Labels 一致。
func (summary *summaryVec) Observe(v float64, labels ...string) {
	summary.WithLabelValues(labels...).Observe(v)
}
