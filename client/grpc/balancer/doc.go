// Package balancer 提供 gRPC 客户端负载均衡器的基础框架与平滑加权轮询（SWR）均衡器。
//
// 基础部分定义 Balancer Builder、Picker Builder 抽象，跟踪子连接连通状态并提供错误 Picker；
// SWR 均衡器结合服务分组、节点与路由权重配置，在就绪的子连接之间分配流量。
package balancer
