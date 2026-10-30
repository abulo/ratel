package cron

import "time"

// SpecSchedule 基于传统 crontab 表达式描述精度到秒的触发周期，各字段以位集合形式保存。
type SpecSchedule struct {
	// Second、Minute、Hour、Dom、Month、Dow 分别为秒、分钟、小时、日、月、星期的命中位集合。
	Second, Minute, Hour, Dom, Month, Dow uint64

	// Location 为该调度表单独指定的时区；为 time.Local 时沿用调用方时间。
	Location *time.Location
}

// bounds 描述单个字段的取值区间，以及星期/月份别名到数值的映射。
type bounds struct {
	min, max uint
	names    map[string]uint
}

// 各字段的取值区间。
var (
	seconds = bounds{0, 59, nil}
	minutes = bounds{0, 59, nil}
	hours   = bounds{0, 23, nil}
	dom     = bounds{1, 31, nil}
	months  = bounds{1, 12, map[string]uint{
		"jan": 1,
		"feb": 2,
		"mar": 3,
		"apr": 4,
		"may": 5,
		"jun": 6,
		"jul": 7,
		"aug": 8,
		"sep": 9,
		"oct": 10,
		"nov": 11,
		"dec": 12,
	}}
	dow = bounds{0, 6, map[string]uint{
		"sun": 0,
		"mon": 1,
		"tue": 2,
		"wed": 3,
		"thu": 4,
		"fri": 5,
		"sat": 6,
	}}
)

const (
	// starBit 在表达式包含星号通配时置位，用于区分日与星期字段为“或”还是“与”语义。
	starBit = 1 << 63
)

// Next 返回晚于 t 的下一次触发时间；若五年内都找不到满足调度表的时间，则返回零值时间。
func (s *SpecSchedule) Next(t time.Time) time.Time {
	// 按月、日、小时、分钟、秒逐级进位匹配；字段回绕到边界后必须回头重新校验前置字段，
	// 因此使用 WRAP 标签统一回到入口重新检查。

	// 先把给定时间转换到调度表指定的时区，并记录原时区，找到目标时间后再转回。
	// 未指定时区（time.Local）时按 t 自带时区处理。
	origLocation := t.Location()
	loc := s.Location
	if loc == time.Local {
		loc = t.Location()
	}
	if s.Location != time.Local {
		t = t.In(s.Location)
	}

	// 从 t 的下一秒开始计算。
	t = t.Add(1*time.Second - time.Duration(t.Nanosecond())*time.Nanosecond)

	// 标记是否已经为匹配而进位过某个字段。
	added := false

	// 五年内仍找不到可命中时间时视为不可满足，返回零值。
	yearLimit := t.Year() + 5

WRAP:
	if t.Year() > yearLimit {
		return time.Time{}
	}

	for 1<<uint(t.Month())&s.Month == 0 {
		// 首次进位时先把日、时分秒清零，从月初开始推进。
		if !added {
			added = true
			t = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, loc)
		}
		t = t.AddDate(0, 1, 0)

		if t.Month() == time.January {
			goto WRAP
		}
	}

	for !dayMatches(s, t) {
		if !added {
			added = true
			t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
		}
		t = t.AddDate(0, 0, 1)
		// 夏令时切换可能导致当天不存在零点，此时把时刻拉回零点附近。
		if t.Hour() != 0 {
			if t.Hour() > 12 {
				t = t.Add(time.Duration(24-t.Hour()) * time.Hour)
			} else {
				t = t.Add(time.Duration(-t.Hour()) * time.Hour)
			}
		}

		if t.Day() == 1 {
			goto WRAP
		}
	}

	for 1<<uint(t.Hour())&s.Hour == 0 {
		if !added {
			added = true
			t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, loc)
		}
		t = t.Add(1 * time.Hour)

		if t.Hour() == 0 {
			goto WRAP
		}
	}

	for 1<<uint(t.Minute())&s.Minute == 0 {
		if !added {
			added = true
			t = t.Truncate(time.Minute)
		}
		t = t.Add(1 * time.Minute)

		if t.Minute() == 0 {
			goto WRAP
		}
	}

	for 1<<uint(t.Second())&s.Second == 0 {
		if !added {
			added = true
			t = t.Truncate(time.Second)
		}
		t = t.Add(1 * time.Second)

		if t.Second() == 0 {
			goto WRAP
		}
	}

	return t.In(origLocation)
}

// dayMatches 判断给定时间是否同时满足调度表的日与星期约束；任一字段使用星号时取“与”，
// 否则日与星期取“或”。
func dayMatches(s *SpecSchedule, t time.Time) bool {
	var (
		domMatch bool = 1<<uint(t.Day())&s.Dom > 0
		dowMatch bool = 1<<uint(t.Weekday())&s.Dow > 0
	)
	if s.Dom&starBit > 0 || s.Dow&starBit > 0 {
		return domMatch && dowMatch
	}
	return domMatch || dowMatch
}
