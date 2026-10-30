package client

import "sync"

// Proxy 是按分组保存客户端配置的并发安全代理池，底层使用 sync.Map 存放分组与配置的映射。
//
// 零值即可使用，配置的读写由内部 sync.Map 保证并发安全；NewClientProxy 仅用于显式初始化。
type Proxy struct {
	m sync.Map
}

// NewClientProxy 创建一个空的 Proxy 代理池。
func NewClientProxy() *Proxy {
	return &Proxy{
		m: sync.Map{},
	}
}
