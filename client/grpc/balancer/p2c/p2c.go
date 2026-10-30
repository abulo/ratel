package p2c

import (
	"github.com/abulo/ratel/v3/client/grpc/balancer/p2c/basep2c"
	"github.com/abulo/ratel/v3/client/grpc/balancer/p2c/leastloaded"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/grpclog"
)

const (
	// Name 为 P2C 最少负载均衡器在 gRPC 注册中心的策略名称。
	Name = "p2c_least_loaded"
)

// newBuilder 创建带健康检查开关的 P2C 均衡器 Builder。
func newBuilder() balancer.Builder {
	return base.NewBalancerBuilder(Name, &p2cPickerBuilder{}, base.Config{HealthCheck: true})
}

func init() {
	balancer.Register(newBuilder())
}

type p2cPickerBuilder struct{}

// Build 在无就绪连接时返回错误 Picker；否则初始化最少在途负载选择器，装入全部就绪 SubConn。
func (*p2cPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	grpclog.Infof("p2cPickerBuilder: newPicker called with readySCs: %v", info.ReadySCs)
	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}

	var p2c = leastloaded.New()

	for sc := range info.ReadySCs {
		p2c.Add(sc)
	}

	rp := &p2cPicker{
		p2c: p2c,
	}
	return rp
}

type p2cPicker struct {
	p2c basep2c.P2c
}

// Pick 通过 P2C 选择器取出在途负载最低的 SubConn；无可用候选时返回 ErrNoSubConnAvailable。
// 返回的 Done 回调会在请求结束后递减对应节点的在途计数。
func (p *p2cPicker) Pick(opts balancer.PickInfo) (balancer.PickResult, error) {
	item, done := p.p2c.Next()
	if item == nil {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}

	return balancer.PickResult{SubConn: item.(balancer.SubConn), Done: done}, nil
}
