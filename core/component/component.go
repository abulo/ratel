package component

import "github.com/abulo/ratel/core/metric"

type Component interface {
	// Start blocks until the channel is closed or an error occurs.
	// The component will stop running when the channel is closed.
	Start(<-chan struct{}) error

	ShouldBeLeader() bool
}

var _ Component = ComponentFunc(nil)

type ComponentFunc func(<-chan struct{}) error

func (f ComponentFunc) Start(stop <-chan struct{}) error {
	return f(stop)
}

func (f ComponentFunc) ShouldBeLeader() bool {
	return false
}

// Component manager, aggregate multiple components to one
type Manager interface {
	Component
	AddComponent(...Component) error
}

// Component builder, build component with injecting govern plugin
type Builder interface {
	WithComponentManager(Manager) Builder
	WithMetrics(metric.Metrics) Builder
	Build() Component
}
