package balancer

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/resolver"
)

// NewBalancerBuilderV2 按名称、Picker 构造器与基础配置创建一个 gRPC Balancer Builder。
func NewBalancerBuilderV2(name string, pb PickerBuilder, config base.Config) balancer.Builder {
	return &baseBuilder{
		name:            name,
		v2PickerBuilder: pb,
		config:          config,
	}
}

type baseBuilder struct {
	name            string
	v2PickerBuilder PickerBuilder
	config          base.Config
}

// Build 依据客户端连接与构建选项创建底层 baseBalancer，初始 Picker 为始终返回
// ErrNoSubConnAvailable 的错误 Picker，等待子连接状态上报后再生成可用 Picker。
func (bb *baseBuilder) Build(cc balancer.ClientConn, opt balancer.BuildOptions) balancer.Balancer {
	bal := &baseBalancer{
		cc:              cc,
		v2PickerBuilder: bb.v2PickerBuilder,

		subConns: make(map[resolver.Address]balancer.SubConn),
		scStates: make(map[balancer.SubConn]connectivity.State),
		csEvltr:  &balancer.ConnectivityStateEvaluator{},
		config:   bb.config,
	}
	bal.v2Picker = NewErrPickerV2(balancer.ErrNoSubConnAvailable)
	return bal
}

// Name 返回该 Builder 注册使用的负载均衡策略名称。
func (bb *baseBuilder) Name() string {
	return bb.name
}

// 编译期断言 baseBalancer 实现 balancer.Balancer 接口。
var _ balancer.Balancer = (*baseBalancer)(nil)

type baseBalancer struct {
	cc              balancer.ClientConn
	v2PickerBuilder PickerBuilder
	csEvltr         *balancer.ConnectivityStateEvaluator
	state           connectivity.State
	subConns        map[resolver.Address]balancer.SubConn
	scStates        map[balancer.SubConn]connectivity.State
	v2Picker        balancer.Picker
	config          base.Config
	attributes      *attributes.Attributes
}

// HandleResolvedAddrs 为旧版接口，当前实现未提供。
func (b *baseBalancer) HandleResolvedAddrs(addrs []resolver.Address, err error) {
	panic("not implemented")
}

// ResolverError 在解析器上报错误时，仅当整体处于 TransientFailure、Idle 或 Connecting 状态时，
// 将当前 Picker 替换为携带该错误的错误 Picker。
func (b *baseBalancer) ResolverError(err error) {
	switch b.state {
	case connectivity.TransientFailure, connectivity.Idle, connectivity.Connecting:
		b.v2Picker = NewErrPickerV2(err)
	}
}

// UpdateClientConnState 同步解析器上报的地址集合：为新地址创建 SubConn 并发起连接，移除已下线
// 地址对应的 SubConn，同时保存解析器携带的属性。
func (b *baseBalancer) UpdateClientConnState(s balancer.ClientConnState) error {
	// TODO(owner): s.ResolverState.Err 非 nil 时补充日志记录。
	// TODO(owner): 处理 s.ResolverState.ServiceConfig。
	if grpclog.V(2) {
		grpclog.Infoln("base.baseBalancer: got new ClientConn state: ", s)
	}
	// addrsSet 由新地址集合转换而来，用于快速判断地址是否已被解析器移除。
	addrsSet := make(map[resolver.Address]struct{})
	for _, a := range s.ResolverState.Addresses {
		addrsSet[a] = struct{}{}
		if _, ok := b.subConns[a]; !ok {
			// 每个 SubConn 只传入单个地址，这是 gRPC 推荐且向前兼容的做法。
			sc, err := b.cc.NewSubConn(
				[]resolver.Address{a},
				balancer.NewSubConnOptions{HealthCheckEnabled: b.config.HealthCheck},
			)
			if err != nil {
				grpclog.Warningf("base.baseBalancer: failed to create new SubConn: %v", err)
				continue
			}
			b.subConns[a] = sc
			b.scStates[sc] = connectivity.Idle
			sc.Connect()
		}
	}

	b.attributes = s.ResolverState.Attributes

	for a, sc := range b.subConns {
		// 地址已被解析器移除，关闭对应 SubConn；状态记录保留到 HandleSubConnStateChange 中清理。
		if _, ok := addrsSet[a]; !ok {
			b.cc.RemoveSubConn(sc)
			delete(b.subConns, a)
		}
	}
	return nil
}

// regeneratePicker 基于当前子连接快照重建 Picker：整体处于 TransientFailure 时使用错误 Picker，
// 其余情况将所有 Ready 状态的 SubConn 交给 PickerBuilder 构造可用 Picker。
func (b *baseBalancer) regeneratePicker(err error) {
	if b.state == connectivity.TransientFailure {
		if err != nil {
			b.v2Picker = NewErrPickerV2(err)
		} else {
			// 上一次子连接迁移并非进入 TransientFailure，但整体聚合状态已是 TransientFailure，
			// 说明已没有可用地址。
			b.v2Picker = NewErrPickerV2(errors.New("resolver returned no addresses"))
		}
		return
	}
	readySCs := make(map[balancer.SubConn]base.SubConnInfo)

	// 从全部 SubConn 中筛出 Ready 状态的连接。
	for addr, sc := range b.subConns {
		if st, ok := b.scStates[sc]; ok && st == connectivity.Ready {
			readySCs[sc] = base.SubConnInfo{Address: addr}
		}
	}
	b.v2Picker = b.v2PickerBuilder.Build(
		PickerBuildInfo{
			ReadySCs:   readySCs,
			Attributes: b.attributes,
		},
	)
}

// HandleSubConnStateChange 为旧版接口，当前实现未提供。
func (b *baseBalancer) HandleSubConnStateChange(sc balancer.SubConn, s connectivity.State) {
	panic("not implemented")
}

// UpdateSubConnState 跟踪单个 SubConn 的状态迁移：Idle 时主动发起连接，Shutdown 时清理状态记录；
// 在 Ready 与非 Ready、TransientFailure 与非 TransientFailure 之间切换时重建 Picker，
// 并将最新整体状态与 Picker 上报给 ClientConn。
func (b *baseBalancer) UpdateSubConnState(sc balancer.SubConn, state balancer.SubConnState) {
	s := state.ConnectivityState
	if grpclog.V(2) {
		grpclog.Infof("base.baseBalancer: handle SubConn state change: %p, %v", sc, s)
	}
	oldS, ok := b.scStates[sc]
	if !ok {
		if grpclog.V(2) {
			grpclog.Infof("base.baseBalancer: got state changes for an unknown SubConn: %p, %v", sc, s)
		}
		return
	}
	b.scStates[sc] = s
	switch s {
	case connectivity.Idle:
		sc.Connect()
	case connectivity.Shutdown:
		// 解析器移除地址时只调用了 RemoveSubConn，这里补删对应的状态记录。
		delete(b.scStates, sc)
	}

	oldAggrState := b.state
	b.state = b.csEvltr.RecordTransition(oldS, s)

	// 以下任一情形需要重建 Picker：Ready 与非 Ready 切换，或整体状态在
	// TransientFailure 与非 TransientFailure 之间切换。
	if (s == connectivity.Ready) != (oldS == connectivity.Ready) ||
		(b.state == connectivity.TransientFailure) != (oldAggrState == connectivity.TransientFailure) {
		b.regeneratePicker(state.ConnectionError)
	}

	b.cc.UpdateState(balancer.State{ConnectivityState: b.state, Picker: b.v2Picker})
}

// Close 不执行清理：基础均衡器没有需要释放的内部状态，也无需在此主动移除 SubConn。
func (b *baseBalancer) Close() {
}

// ExitIdle 将处于 Idle 状态的 SubConn 全部重新发起连接。
func (b *baseBalancer) ExitIdle() {
	for _, sc := range b.subConns {
		if state, ok := b.scStates[sc]; ok && state == connectivity.Idle {
			sc.Connect()
		}
	}
}

// NewErrPickerV2 创建一个每次 Pick 都返回指定错误的错误 Picker。
func NewErrPickerV2(err error) balancer.Picker {
	return &errPickerV2{err: err}
}

type errPickerV2 struct {
	// err 为 Pick 始终返回的错误。
	err error
}

// Pick 直接返回构造时携带的错误，不选择任何 SubConn。
func (p *errPickerV2) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	return balancer.PickResult{}, p.err
}
