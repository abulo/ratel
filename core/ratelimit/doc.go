// Package ratelimit 提供令牌桶（token bucket）限流实现，以及基于令牌桶的限流 Reader 与 Writer。
//
// 桶支持按固定间隔、固定速率与每次补充量子三种方式填充令牌；取令牌操作分为阻塞等待、
// 限时等待与非阻塞立即取用三类。桶内部加锁，可被多个 goroutine 并发使用。
package ratelimit
