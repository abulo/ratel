package util

import (
	"math/rand"
	"time"
)

// Do 以指数退避加随机抖动的方式重试函数 f。
//
// f 首次执行失败后，在 attempts 次重试内依次等待 sleep、2*sleep……，等待时间中
// 附加随机抖动；全部重试失败后返回最后一次错误。attempts 为 0 时只执行 f 一次。
func Do(attempts int, sleep time.Duration, f func() error) error {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if err := f(); err != nil {
		if attempts--; attempts >= 0 {
			// 加入随机抖动，避免大量客户端在同一时刻集中重试造成惊群效应。
			jitter := time.Duration(r.Int63n(int64(sleep)))
			sleep = sleep + jitter/2
			time.Sleep(sleep)
			return Do(attempts, 2*sleep, f)
		}
		return err
	}
	return nil
}
