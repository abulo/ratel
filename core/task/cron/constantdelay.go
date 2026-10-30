package cron

import "time"

// ConstantDelaySchedule 表示按固定时间间隔循环触发的调度表。
type ConstantDelaySchedule struct {
	// Delay 为两次触发之间的间隔。
	Delay time.Duration
}

// Every 构造一个每隔 duration 触发一次的调度表；duration 小于 1 秒时提升为 1 秒，
// 不足 1 秒的纳秒部分会被截断，因此实际触发频率不会高于每秒一次。
func Every(duration time.Duration) ConstantDelaySchedule {
	if duration < time.Second {
		duration = time.Second
	}
	return ConstantDelaySchedule{
		Delay: duration - time.Duration(duration.Nanoseconds())%time.Second,
	}
}

// Next 返回晚于 t 的下一次触发时间，并对齐到整秒。
func (schedule ConstantDelaySchedule) Next(t time.Time) time.Time {
	return t.Add(schedule.Delay - time.Duration(t.Nanosecond())*time.Nanosecond)
}
