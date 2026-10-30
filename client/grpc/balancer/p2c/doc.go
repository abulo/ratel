// Package p2c 实现基于 P2C（Power of Two Choices，随机二抽样）的最少在途负载 gRPC 均衡器。
//
// 均衡器从就绪 SubConn 中随机抽取两个候选，选择当前在途请求数较少的一个，从而把新请求引导到
// 负载更低的节点，避免单节点被持续打满。
package p2c
