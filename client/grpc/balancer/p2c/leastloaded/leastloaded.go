package leastloaded

import (
	randv2 "math/rand/v2"
	"sync/atomic"

	"github.com/abulo/ratel/v3/client/grpc/balancer/p2c/basep2c"
	"google.golang.org/grpc/balancer"
)

type leastLoadedNode struct {
	item     any
	inflight int64
	// done 在 Add 时绑定当前节点创建一次，Next 直接复用，避免每条 RPC 现场新建闭包。
	done func(balancer.DoneInfo)
}

// noopDone 为无候选节点时返回的包级共享空回调，避免每次新建闭包。
var noopDone = func(balancer.DoneInfo) {}

type leastLoaded struct {
	items []*leastLoadedNode
}

// New 创建最少在途负载选择器；随机数由 math/rand/v2 顶层并发安全随机源提供。
func New() basep2c.P2c {
	return &leastLoaded{
		items: make([]*leastLoadedNode, 0),
	}
}

// Add 向选择器加入一个候选节点，初始在途计数为 0。
func (p *leastLoaded) Add(item any) {
	node := &leastLoadedNode{item: item}
	// PERF-5-10 [P2]：done 回调在 Add 时绑定节点创建一次，Next 直接返回该预存回调，
	// 消除每条 RPC 现场构造闭包导致的堆逃逸分配；DoneInfo 透传语义保持不变。
	node.done = func(balancer.DoneInfo) {
		atomic.AddInt64(&node.inflight, -1)
	}
	p.items = append(p.items, node)
}

// Next 选出下一个候选节点：无候选时返回 nil 与共享空回调；仅一个候选时直接选中；多个候选时随机抽取
// 两个并选择在途请求数较少者。选中后在途计数加 1，返回的 Done 回调在请求结束后将计数减回。
func (p *leastLoaded) Next() (any, func(balancer.DoneInfo)) {
	var sc, backsc *leastLoadedNode

	switch len(p.items) {
	case 0:
		return nil, noopDone
	case 1:
		sc = p.items[0]
	default:
		// PERF-5-9 [P1]：改用 math/rand/v2 顶层并发安全随机源（内部基于 per-P 快随机，无需调用方持锁），
		// 移除原 per-instance rand.Rand 与共享互斥锁，消除每次 Pick 在同一把锁上的排队竞争；
		// 两个不同下标的生成与去重逻辑保持不变。（预期收益需实测验证）
		a := randv2.IntN(len(p.items))
		b := randv2.IntN(len(p.items) - 1)

		if b >= a {
			b = b + 1
		}
		sc, backsc = p.items[a], p.items[b]

		// PERF-5-11 [P2]：在途计数比较改用 atomic.LoadInt64 读取，与 atomic.AddInt64 的写入处于
		// 同一原子域，消除裸读构成的数据竞争；读取的仍是最新在途数，节点选择语义不变。
		if atomic.LoadInt64(&sc.inflight) > atomic.LoadInt64(&backsc.inflight) {
			sc = backsc
		}
	}

	atomic.AddInt64(&sc.inflight, 1)

	return sc.item, sc.done
}
