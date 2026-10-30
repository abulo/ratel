// Package leastloaded 实现基于在途请求计数的 P2C 最少负载选择器。
//
// Next 随机抽取两个候选节点，选择在途请求数较少的一个；请求结束回调负责递减计数，
// 保证并发场景下负载跟踪准确。
package leastloaded
