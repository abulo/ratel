package hostname

import (
	"os"
	"sync"
)

var (
	// PERF-4-7 [P1]：主机名在进程生命周期内恒定，sync.Once 仅发起一次系统调用并缓存结果。
	// 语义说明：与原实现一致，首次获取失败时返回 "unknown"；失败结果同样被缓存，
	// 避免每次调用重复进入系统调用。
	initOnce sync.Once
	cached   string
)

// Hostname 返回当前主机名；获取主机名的系统调用失败时返回 "unknown"。
func Hostname() string {
	initOnce.Do(func() {
		hostName, err := os.Hostname()
		if err != nil {
			hostName = "unknown"
		}
		cached = hostName
	})
	return cached
}
