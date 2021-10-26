package proxy

import (
	"github.com/abulo/ratel/v2/store/kafka"
)

type ProxyKafka struct {
	*kafka.Clientkafka
}

//NewProxyKafka 缓存
func NewProxyKafka() *ProxyKafka {
	return &ProxyKafka{}
}

//Store 设置写库
func (proxy *ProxyKafka) Store(client *kafka.Clientkafka) {
	proxy.Clientkafka = client
}

//StoreEs 设置组
func (proxypool *ProxyPool) StoreKafka(group string, proxy *ProxyKafka) {
	proxypool.m.Store(group, proxy)
}

//LoadEs 获取分组
func (proxypool *ProxyPool) LoadKafka(group string) *kafka.Clientkafka {
	if f, ok := proxypool.m.Load(group); ok {
		return f.(*ProxyKafka).Clientkafka
	}
	return nil
}
