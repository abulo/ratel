package consul

import (
	"fmt"
	"time"

	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc/resolver"
)

type Builder struct {
	name    string
	address string
	token   string
}

type Resolver struct {
	conn      resolver.ClientConn
	lastIndex uint64
	name      string
	address   string
	token     string
}

// NewResolver ...
func NewResolver(conf ConsulCliConf) string {
	builder := &Builder{
		name:    conf.ServiceName,
		address: conf.Address,
		token:   conf.Token,
	}
	resolver.Register(builder)

	return "consul:///" + conf.ServiceName
}

func (r *Builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {
	rl := &Resolver{conn: cc, name: r.name, address: r.address, token: r.token}
	go rl.watch()
	return rl, nil
}

func (r *Builder) Scheme() string {
	return "consul"
}

func (r *Resolver) watch() {
	config := api.DefaultConfig()
	Address = r.address
	Token = r.token
	client, err := api.NewClient(config)
	if err != nil {
		panic(err)
	}

	for {
		services, metainfo, err := client.Health().Service(r.name, "", true, &api.QueryOptions{WaitIndex: r.lastIndex})
		if err != nil {
			time.Sleep(time.Second * 5)
			continue
		}

		r.lastIndex = metainfo.LastIndex
		var newAddrs []resolver.Address
		for _, service := range services {
			addr := fmt.Sprintf("%v:%v", service.Service.Address, service.Service.Port)
			newAddrs = append(newAddrs, resolver.Address{Addr: addr})
		}

		r.conn.UpdateState(resolver.State{Addresses: newAddrs})
	}
}

func (r *Resolver) ResolveNow(opt resolver.ResolveNowOption) {
}

func (r *Resolver) Close() {
}
