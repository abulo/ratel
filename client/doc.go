// Package client 提供按分组（group）存取 gRPC 与 etcdv3 客户端配置的代理池。
//
// Proxy 内部基于 sync.Map 维护分组到配置的映射，支持并发读写；GrpcConfig 与
// EtcdConfig 分别包装 grpc.Config、etcdv3.Config，供代理池统一存放与取回。
package client
