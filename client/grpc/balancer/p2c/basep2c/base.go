package basep2c

import "google.golang.org/grpc/balancer"

// P2c 是 P2C 负载选择器契约：实现负责维护候选集合并按自身策略选出节点，同时通过 Done 回调
// 跟踪在途请求数；上层 Picker 不感知具体负载算法。
type P2c interface {
	// Next 返回下一个被选中的候选项；返回的回调在请求结束后调用，用于更新节点负载状态。
	Next() (any, func(balancer.DoneInfo))
	// Add 向选择器中加入一个候选项。
	Add(any)
}
