package crontab

// sched := crontab.NewScheduler()
// sched.Schedule().Every(10).Seconds().Do(something)
// sched.Schedule().Every(3).Minutes().Do(something)
// sched.Schedule().Every(4).Hours().Do(something)
// sched.Schedule().Every(2).Days().At("12:32").Do(something)
// sched.Schedule().Every(12).Weeks().Do(something)
// sched.Schedule().Every(1).Monday().Do(something)
// sched.Schedule().Every(1).Saturday().At("8:00").Do(something)
// sched.Run()

import (
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

//TimeUnit 是用于在内部处理时间单位的计数
type TimeUnit int

type Logger interface {
	Panic(args ...interface{})
	Fatal(args ...interface{})
	Error(args ...interface{})
	Warning(args ...interface{})
	Warn(args ...interface{})
	Info(args ...interface{})
	Debug(args ...interface{})
	Trace(args ...interface{})
}

var Log Logger

func SetLogger(logger Logger) {
	Log = logger
}

const (
	none = iota
	second
	minute
	hour
	day
	week
	monday
	tuesday
	wednesday
	thursday
	friday
	saturday
	sunday
)

var timeNow = func() time.Time {
	return time.Now()
}

// Job 结构处理调度和运行jobs所需的所有数据
type Job struct {
	identifier string
	scheduler  *Scheduler
	unit       TimeUnit
	frequency  int
	useAt      bool
	atHour     int
	atMinute   int
	workFunc   func()

	nextScheduledRun time.Time
}

// Every 是以既定频率填充既定Job结构的方法
func (j *Job) Every(frequencies ...int) *Job {
	l := len(frequencies)

	switch l {
	case 0:
		j.frequency = 1
	case 1:
		if frequencies[0] <= 0 {
			Log.Panic("Every expects frequency to be greater than of equal to 1")
		}
		j.frequency = frequencies[0]
	default:
		Log.Panic("Every expects 0 or 1 arguments")
	}

	return j
}

// At 方法使用提供的信息填充atHour和atMinute字段的Job结构
func (j *Job) At(t string) *Job {
	j.useAt = true
	j.atHour, _ = strconv.Atoi(strings.Split(t, ":")[0])
	j.atMinute, _ = strconv.Atoi(strings.Split(t, ":")[1])
	return j
}

// Do 执行方法
func (j *Job) Do(function func()) string {
	j.workFunc = function
	j.scheduleNextRun()
	j.scheduler.jobs = append(j.scheduler.jobs, *j)
	return j.identifier
}

func (j *Job) due() bool {
	return timeNow().After(j.nextScheduledRun)
}

func (j *Job) isAtUsedIncorrectly() bool {
	return j.useAt == true && (j.unit == second || j.unit == minute ||
		j.unit == hour || j.unit == week)
}

func (j *Job) unitNotDayOrWEEKDAY() bool {
	return j.unit == second || j.unit == minute ||
		j.unit == hour || j.unit == week

}

func (j *Job) unitNotWEEKDAY() bool {
	return j.unit == second || j.unit == minute ||
		j.unit == hour || j.unit == day ||
		j.unit == week
}

func (j *Job) scheduleNextRun() {
	if j.frequency == 1 {
		if j.isAtUsedIncorrectly() {
			Log.Panic(
				`Cannot schedule Every(1) with At()
				 when unit is not day or WEEKDAY`,
			)
		}

		if j.unitNotDayOrWEEKDAY() {
			if j.nextScheduledRun == (time.Time{}) {
				j.nextScheduledRun = timeNow()
			}

			switch j.unit {
			case second:
				j.nextScheduledRun = j.nextScheduledRun.Add(1 * time.Second)
			case minute:
				j.nextScheduledRun = j.nextScheduledRun.Add(1 * time.Minute)
			case hour:
				j.nextScheduledRun = j.nextScheduledRun.Add(1 * time.Hour)
			case week:
				// 168 hours in a week
				j.nextScheduledRun = j.nextScheduledRun.Add(168 * time.Hour)
			}
		} else {
			switch j.unit {
			case day:
				if j.nextScheduledRun == (time.Time{}) {
					now := timeNow()
					lastMidnight := time.Date(
						now.Year(),
						now.Month(),
						now.Day(),
						0, 0, 0, 0,
						time.Local,
					)
					if j.useAt == true {
						j.nextScheduledRun = lastMidnight.Add(
							time.Duration(j.atHour)*time.Hour +
								time.Duration(j.atMinute)*time.Minute,
						)
					} else {
						// If At is not specified, move the next scheduled run to next midnight
						j.nextScheduledRun = lastMidnight.Add(24 * time.Hour)
					}
				} else {
					j.nextScheduledRun = j.nextScheduledRun.Add(24 * time.Hour)
				}

			case monday:
				j.scheduleWeekday(time.Monday)
			case tuesday:
				j.scheduleWeekday(time.Tuesday)
			case wednesday:
				j.scheduleWeekday(time.Wednesday)
			case thursday:
				j.scheduleWeekday(time.Thursday)
			case friday:
				j.scheduleWeekday(time.Friday)
			case saturday:
				j.scheduleWeekday(time.Saturday)
			case sunday:
				j.scheduleWeekday(time.Sunday)
			}

		}

	} else {
		if j.isAtUsedIncorrectly() {
			Log.Panic("Cannot schedule Every(>1) with At() when unit is not day")
		}

		if j.unitNotWEEKDAY() {

			if j.unit != day {
				if j.nextScheduledRun == (time.Time{}) {
					j.nextScheduledRun = timeNow()
				}

				switch j.unit {
				case second:
					j.nextScheduledRun = j.nextScheduledRun.Add(
						time.Duration(j.frequency) * time.Second,
					)
				case minute:
					j.nextScheduledRun = j.nextScheduledRun.Add(
						time.Duration(j.frequency) * time.Minute,
					)
				case hour:
					j.nextScheduledRun = j.nextScheduledRun.Add(
						time.Duration(j.frequency) * time.Hour,
					)
				case week:
					j.nextScheduledRun = j.nextScheduledRun.Add(
						time.Duration(j.frequency*168) * time.Hour,
					)

				}
			} else {
				if j.nextScheduledRun == (time.Time{}) {
					now := timeNow()
					lastMidnight := time.Date(
						now.Year(),
						now.Month(),
						now.Day(),
						0, 0, 0, 0,
						time.Local,
					)
					if j.useAt == true {
						j.nextScheduledRun = lastMidnight.Add(
							time.Duration(j.atHour)*time.Hour +
								time.Duration(j.atMinute)*time.Minute,
						)
					} else {
						j.nextScheduledRun = lastMidnight
					}
				}

				j.nextScheduledRun = j.nextScheduledRun.Add(
					time.Duration(j.frequency*24) * time.Hour,
				)

			}

		} else {
			Log.Panic("Cannot schedule Every(>1) when unit is WEEKDAY")
		}
	}
	return
}

func (j *Job) scheduleWeekday(dayOfWeek time.Weekday) {
	if j.nextScheduledRun == (time.Time{}) {
		now := timeNow()
		lastWeekdayMidnight := time.Date(
			now.Year(),
			now.Month(),
			now.Day()-int(now.Weekday()-dayOfWeek),
			0, 0, 0, 0,
			time.Local)
		if j.useAt == true {
			j.nextScheduledRun = lastWeekdayMidnight.Add(
				time.Duration(j.atHour)*time.Hour +
					time.Duration(j.atMinute)*time.Minute,
			)
		} else {
			j.nextScheduledRun = lastWeekdayMidnight
		}
	}
	j.nextScheduledRun = j.nextScheduledRun.Add(7 * 24 * time.Hour)
}

// Second 方法用秒填充给定的作业结构
func (j *Job) Second() *Job {
	j.unit = second
	return j
}

// Seconds 方法用秒填充给定的作业结构
func (j *Job) Seconds() *Job {
	j.unit = second
	return j
}

// Minute 方法用分钟填充给定的作业结构
func (j *Job) Minute() *Job {
	j.unit = minute
	return j
}

// Minutes 方法用分钟填充给定的作业结构
func (j *Job) Minutes() *Job {
	j.unit = minute
	return j
}

// Hour 方法用小时填充给定的作业结构
func (j *Job) Hour() *Job {
	j.unit = hour
	return j
}

// Hours 方法用小时填充给定的作业结构
func (j *Job) Hours() *Job {
	j.unit = hour
	return j
}

// Day 方法用每天填充给定的作业结构
func (j *Job) Day() *Job {
	j.unit = day
	return j
}

// Days 方法用每天填充给定的作业结构
func (j *Job) Days() *Job {
	j.unit = day
	return j
}

// Week 方法用每周填充给定的作业结构
func (j *Job) Week() *Job {
	j.unit = week
	return j
}

// Weeks 方法用每周填充给定的作业结构
func (j *Job) Weeks() *Job {
	j.unit = week
	return j
}

// Monday 周一
func (j *Job) Monday() *Job {
	j.unit = monday
	return j
}

// Tuesday 周二
func (j *Job) Tuesday() *Job {
	j.unit = tuesday
	return j
}

// Wednesday 周3
func (j *Job) Wednesday() *Job {
	j.unit = wednesday
	return j
}

// Thursday 周4
func (j *Job) Thursday() *Job {
	j.unit = thursday
	return j
}

// Friday 周5
func (j *Job) Friday() *Job {
	j.unit = friday
	return j
}

// Saturday 周6
func (j *Job) Saturday() *Job {
	j.unit = saturday
	return j
}

// Sunday 周7
func (j *Job) Sunday() *Job {
	j.unit = sunday
	return j
}

// Scheduler 存储 job
type Scheduler struct {
	identifier string
	jobs       []Job
}

// NewScheduler 创建一个调度器 new Scheduler
func NewScheduler() Scheduler {
	return Scheduler{
		identifier: uuid.New().String(),
		jobs:       make([]Job, 0),
	}
}

func (s *Scheduler) activateTestMode() {
	timeNow = func() time.Time {
		return time.Date(1, 1, 1, 1, 1, 0, 0, time.Local)
	}
}

// Run method  Scheduler.
func (s *Scheduler) Run() {
	for {
		for jobIdx := range s.jobs {
			job := &s.jobs[jobIdx]
			if job.due() {
				job.scheduleNextRun()
				go job.workFunc()
			}
		}
		time.Sleep(1 * time.Second)

	}
}

// Schedule
func (s *Scheduler) Schedule() *Job {
	newJob := Job{
		identifier:       uuid.New().String(),
		scheduler:        s,
		unit:             none,
		frequency:        1,
		useAt:            false,
		atHour:           0,
		atMinute:         0,
		workFunc:         nil,
		nextScheduledRun: time.Time{}, // zero value
	}
	return &newJob
}
