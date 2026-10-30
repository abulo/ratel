package leastloaded

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/abulo/ratel/v3/client/grpc/balancer/p2c/basep2c"
	"google.golang.org/grpc/balancer"
)

type leastLoadedNode struct {
	item     any
	inflight int64
}

type leastLoaded struct {
	items []*leastLoadedNode
	mu    sync.Mutex
	rand  *rand.Rand
}

// New 创建最少在途负载选择器，内部随机源以当前时间为种子。
func New() basep2c.P2c {
	return &leastLoaded{
		items: make([]*leastLoadedNode, 0),
		rand:  rand.New(rand.NewSource(time.Now().Unix())),
	}
}

// Add 向选择器加入一个候选节点，初始在途计数为 0。
func (p *leastLoaded) Add(item any) {
	p.items = append(p.items, &leastLoadedNode{item: item})
}

// Next 选出下一个候选节点：无候选时返回 nil 与空回调；仅一个候选时直接选中；多个候选时随机抽取
// 两个并选择在途请求数较少者。选中后在途计数加 1，返回的 Done 回调在请求结束后将计数减回。
func (p *leastLoaded) Next() (any, func(balancer.DoneInfo)) {
	var sc, backsc *leastLoadedNode

	switch len(p.items) {
	case 0:
		return nil, func(balancer.DoneInfo) {}
	case 1:
		sc = p.items[0]
	default:
		// rand.Rand 非并发安全，取随机数期间需要加锁。
		p.mu.Lock()
		a := p.rand.Intn(len(p.items))
		b := p.rand.Intn(len(p.items) - 1)
		p.mu.Unlock()

		if b >= a {
			b = b + 1
		}
		sc, backsc = p.items[a], p.items[b]

		// 比较在途数，选择负载较低的节点。
		if sc.inflight > backsc.inflight {
			sc, _ = backsc, sc
		}
	}

	atomic.AddInt64(&sc.inflight, 1)

	return sc.item, func(balancer.DoneInfo) {
		atomic.AddInt64(&sc.inflight, -1)
	}
}
