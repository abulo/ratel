package component

import "github.com/abulo/ratel/v3/core/metric"

// Component 表示一个由应用托管的长运行组件。
//
// Start 在收到 stop 通道关闭信号前持续运行，通道关闭或运行出错时退出；
// ShouldBeLeader 返回该组件是否只允许在 Leader 节点运行。
type Component interface {
	Start(<-chan struct{}) error

	ShouldBeLeader() bool
}

var _ Component = ComponentFunc(nil)

// ComponentFunc 是普通组件函数到 Component 接口的适配器，直接复用传入函数实现 Start；
// 其 ShouldBeLeader 固定返回 false，即不作为 Leader 专属组件。
type ComponentFunc func(<-chan struct{}) error

// Start 调用底层函数并透传 stop 通道与错误返回值。
func (f ComponentFunc) Start(stop <-chan struct{}) error {
	return f(stop)
}

// ShouldBeLeader 返回 false，表示该适配器组件不在 Leader 节点上独占运行。
func (f ComponentFunc) ShouldBeLeader() bool {
	return false
}

// Manager 是组件管理器，在 Component 基础上提供向自身追加托管组件的能力。
type Manager interface {
	Component
	AddComponent(...Component) error
}

// Builder 是组件构造器，在构建前注入组件管理器与指标采集器，最终产出一个 Component。
type Builder interface {
	WithComponentManager(Manager) Builder
	WithMetrics(metric.Metrics) Builder
	Build() Component
}
