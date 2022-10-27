package metric

import (
	"github.com/abulo/ratel/v3/core/constant"
	"github.com/prometheus/client_golang/prometheus"
)

// CounterVecOpts ...
type CounterVecOpts struct {
	Namespace string
	Subsystem string
	Name      string
	Help      string
	Labels    []string
}

// Build ...
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

// NewCounterVec ...
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
}

// Inc ...
func (counter *counterVec) Inc(labels ...string) {
	counter.WithLabelValues(labels...).Inc()
}

// Add ...
func (counter *counterVec) Add(v float64, labels ...string) {
	counter.WithLabelValues(labels...).Add(v)
}
