package util

import (
	"time"
)

// Do 以指数退避加随机抖动的方式重试函数 f。
//
// f 首次执行失败后，在 attempts 次重试内依次等待 sleep、2*sleep……，等待时间中
// 附加随机抖动；全部重试失败后返回最后一次错误。attempts 为 0 时只执行 f 一次。
func Do(attempts int, sleep time.Duration, f func() error) error {
	if err := f(); err != nil {
		if attempts--; attempts >= 0 {
			// 加入随机抖动，避免大量客户端在同一时刻集中重试造成惊群效应。
			// PERF-7-30 [P2] 复用进程级共享随机源生成抖动，不再每次新建 rand.Rand 与随机源（预期收益需实测验证）
			jitter := time.Duration(sharedInt63n(int64(sleep)))
			sleep = sleep + jitter/2
			time.Sleep(sleep)
			return Do(attempts, 2*sleep, f)
		}
		return err
	}
	return nil
}
