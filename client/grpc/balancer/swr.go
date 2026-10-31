package balancer

import (
	"sync"

	"github.com/pkg/errors"

	"github.com/abulo/ratel/v3/core/constant"
	"github.com/abulo/ratel/v3/registry"
	"github.com/abulo/ratel/v3/server"
	"github.com/smallnest/weighted"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/grpclog"
)

const (
	// NameSmoothWeightRoundRobin 为平滑加权轮询（SWR）均衡器在 gRPC 注册中心的策略名称。
	NameSmoothWeightRoundRobin = "swr"
)

// PickerBuildInfo 构造 Picker 所需的全部信息：ReadySCs 为就绪 SubConn 到其地址的映射，
// 内嵌 Attributes 承载解析器下发的路由、提供方与消费方配置。
type PickerBuildInfo struct {
	// ReadySCs 为所有就绪状态的 SubConn 及其创建时使用的地址。
	ReadySCs map[balancer.SubConn]base.SubConnInfo
	*attributes.Attributes
}

// PickerBuilder 依据 PickerBuildInfo 构造具体 balancer.Picker 的构造器抽象。
type PickerBuilder interface {
	Build(info PickerBuildInfo) balancer.Picker
}

func init() {
	balancer.Register(
		NewBalancerBuilderV2(NameSmoothWeightRoundRobin, &swrPickerBuilder{}, base.Config{HealthCheck: true}),
	)
}

type swrPickerBuilder struct{}

// Build 基于传入的就绪连接信息构造 SWR Picker。
func (s swrPickerBuilder) Build(info PickerBuildInfo) balancer.Picker {
	return newSWRPicker(info)
}

// lockedSW 将一个 weighted.SW 桶与其互斥锁绑定；weighted.SW 内部不做并发保护，
// 每次 Next 仍需在桶级锁内串行执行，但不同路由桶之间不再共享同一把全局锁。
type lockedSW struct {
	mu   sync.Mutex
	size int
	sw   weighted.SW
}

// Add 向桶内追加候选节点并维护节点计数；仅在 Picker 构建期（parseBuildInfo）调用。
func (l *lockedSW) Add(item any, weight int) {
	l.sw.Add(item, weight)
	l.size++
}

type swrPicker struct {
	readySCs     map[balancer.SubConn]base.SubConnInfo
	global       *lockedSW
	routeBuckets map[string]*lockedSW
	*attributes.Attributes
}

func newSWRPicker(info PickerBuildInfo) *swrPicker {
	picker := &swrPicker{
		global:       &lockedSW{},
		readySCs:     info.ReadySCs,
		routeBuckets: map[string]*lockedSW{},
	}
	picker.parseBuildInfo(info)
	return picker
}

// Pick 先按方法名查找该 URI 的专用路由桶，存在则在其中加权轮询；否则使用全局桶。
// 选不到候选 SubConn 时返回错误。
func (p *swrPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	bucket := p.global
	if bs, ok := p.routeBuckets[info.FullMethodName]; ok {
		// 按 URI 将流量分组到独立的路由桶。
		bucket = bs
	}

	// PERF-5-6 [P1]：每个路由桶持有独立互斥锁，Pick 只锁住当前命中的桶，不同方法名的选择
	// 不再在一把全局锁上排队；桶在 Build 完成后节点集合不再变更，size 可在无锁下读取。
	// 单节点桶走无锁快路径：weighted.SW 对 n==1 直接返回唯一节点、不改动轮询状态。（预期收益需实测验证）
	if bucket.size == 1 {
		if sub, ok := bucket.sw.Next().(balancer.SubConn); ok {
			return balancer.PickResult{SubConn: sub}, nil
		}
		return balancer.PickResult{}, errors.New("pick failed")
	}

	bucket.mu.Lock()
	defer bucket.mu.Unlock()

	sub, ok := bucket.sw.Next().(balancer.SubConn)
	if ok {
		return balancer.PickResult{SubConn: sub}, nil
	}

	return balancer.PickResult{}, errors.New("pick failed")
}

// parseBuildInfo 汇总就绪连接与解析器属性，构建全局加权桶以及各路由 URI 的专用桶：
// 同一分组内的节点共用分组权重，节点级权重会覆盖分组权重。
func (p *swrPicker) parseBuildInfo(info PickerBuildInfo) {
	var hostedSubConns = map[string]balancer.SubConn{}
	var groupedSubConns = map[string][]balancer.SubConn{}

	for subConn, info := range info.ReadySCs {
		// PERF-5-7 [P2]：每个节点只向全局桶 Add 一次；删除循环末尾重复的第二次 Add，
		// 避免桶内节点数与权重翻倍导致后续每次 Next 的 O(n) 扫描长度被放大。
		p.global.Add(subConn, 1)
		if info.Address.Attributes != nil {
			if serviceInfo, ok := info.Address.Attributes.Value(constant.KeyServiceInfo).(server.ServiceInfo); ok {
				// TODO(owner): 按服务分组进一步精细化路由。
				group := serviceInfo.Group
				if _, ok := groupedSubConns[group]; !ok {
					groupedSubConns[group] = make([]balancer.SubConn, 0)
				}
				groupedSubConns[group] = append(groupedSubConns[group], subConn)
			}
		}
		host := info.Address.Addr
		hostedSubConns[host] = subConn
	}

	if info.Attributes == nil {
		return
	}

	providerConfig, ok := info.Attributes.Value(constant.KeyProviderConfig).(registry.ProviderConfig)
	if !ok {
		return
	}

	// PERF-5-8 [P2]：移除构建期遗留的 fmt.Printf 同步 stdout 输出，降级为 grpclog.V(2) 级别受控日志。
	if grpclog.V(2) {
		grpclog.Infoln("swr: providerConfig =>", providerConfig)
	}

	consumerConfigs, ok := info.Attributes.Value(constant.KeyConsumerConfig).(map[string]registry.ConsumerConfig)
	if !ok {
		return
	}

	if grpclog.V(2) {
		grpclog.Infoln("swr: consumerConfigs =>", consumerConfigs)
	}

	// 解析路由配置，为每个 URI 建立独立的加权轮询桶。
	routeConfigs, ok := info.Attributes.Value(constant.KeyRouteConfig).(map[string]registry.RouteConfig)
	if !ok {
		return
	}
	for _, config := range routeConfigs {
		if _, ok := p.routeBuckets[config.URI]; !ok {
			p.routeBuckets[config.URI] = &lockedSW{}
		}

		// 基于分组配置权重，同一分组下的所有节点使用相同权重。
		for group, weight := range config.Upstream.Groups {
			sConns, ok := groupedSubConns[group]
			if !ok {
				continue
			}

			for _, sConn := range sConns {
				p.routeBuckets[config.URI].Add(sConn, weight)
			}
		}

		// 基于节点 IP 配置权重；显式配置节点权重时覆盖对应分组权重。
		for node, weight := range config.Upstream.Nodes {
			if sConn, ok := hostedSubConns[node]; ok {
				p.routeBuckets[config.URI].Add(sConn, weight)
			}
		}

	}
}
