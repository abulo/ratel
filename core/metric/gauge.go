package metric

import (
	"github.com/abulo/ratel/core/constant"
	"github.com/prometheus/client_golang/prometheus"
)

// GaugeVecOpts ...
type GaugeVecOpts struct {
	Namespace string
	Subsystem string
	Name      string
	Help      string
	Labels    []string
}

type gaugeVec struct {
	*prometheus.GaugeVec
}

// Build ...
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

// NewGaugeVec ...
func NewGaugeVec(name string, labels []string) *gaugeVec {
	return GaugeVecOpts{
		Namespace: constant.DefaultNamespace,
		Name:      name,
		Help:      name,
		Labels:    labels,
	}.Build()
}

// Inc ...
func (gv *gaugeVec) Inc(labels ...string) {
	gv.WithLabelValues(labels...).Inc()
}

// Add ...
func (gv *gaugeVec) Add(v float64, labels ...string) {
	gv.WithLabelValues(labels...).Add(v)
}

// Set ...
func (gv *gaugeVec) Set(v float64, labels ...string) {
	gv.WithLabelValues(labels...).Set(v)
}
