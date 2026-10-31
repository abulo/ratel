package crontab

import (
	"context"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/abulo/ratel/v3/core/logger"
	"github.com/sirupsen/logrus"
)

// bounds 描述 cron 单个字段的合法取值区间，names 保存月份、星期等别名到数值的映射。
type bounds struct {
	min, max uint
	names    map[string]uint
}

// TaskManager 管理已注册任务并在后台按计划时间触发执行。
//
// 零值不可直接使用，须通过 NewTaskManager 构造；任务注册、删除与遍历方法内部
// 加锁，支持并发调用；同一调度器重复启动不会创建多个调度 goroutine。
type TaskManager struct {
	adminTaskList map[string]Tasker
	taskLock      sync.RWMutex
	stop          chan bool
	changed       chan bool
	started       bool
	wait          sync.WaitGroup
}

// NewTaskManager 创建一个空的 TaskManager，尚未启动调度。
func NewTaskManager() *TaskManager {
	return &TaskManager{
		adminTaskList: make(map[string]Tasker),
		taskLock:      sync.RWMutex{},
		// PERF-7-39 [P2] stop/changed 缓冲为 1，配合非阻塞发送，避免调度循环退出后独立发送 goroutine 永久阻塞
		stop:    make(chan bool, 1),
		changed: make(chan bool, 1),
		started: false,
	}
}

// 各 cron 字段的取值区间，月份与星期附带英文别名映射。
var (
	seconds = bounds{0, 59, nil}
	minutes = bounds{0, 59, nil}
	hours   = bounds{0, 23, nil}
	days    = bounds{1, 31, nil}
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
	weeks = bounds{0, 6, map[string]uint{
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
	// starBit 在表达式包含通配符 * 时置位最高位，用于区分“日”与“星期”字段之间
	// 应取交集还是并集。
	starBit = 1 << 63
)

// Schedule 以位掩码形式保存 cron 表达式解析后的秒、分、时、日、月、星期六个字段，
// 每个字段的某一位为 1 表示对应取值匹配该字段。
type Schedule struct {
	Second uint64
	Minute uint64
	Hour   uint64
	Day    uint64
	Month  uint64
	Week   uint64
}

// TaskFunc 表示任务执行函数，接收上下文并在执行失败时返回错误。
type TaskFunc func(ctx context.Context) error

// Tasker 描述可被 TaskManager 调度的任务契约：任务须能够提供 cron 规则、当前状态、
// 执行入口，以及上一次、下一次执行时间和超时时间。
type Tasker interface {
	GetSpec(ctx context.Context) string
	GetStatus(ctx context.Context) string
	Run(ctx context.Context) error
	SetNext(context.Context, time.Time)
	GetNext(ctx context.Context) time.Time
	SetPrev(context.Context, time.Time)
	GetPrev(ctx context.Context) time.Time
	GetTimeout(ctx context.Context) time.Duration
}

// taskerr 记录单次任务执行失败发生的时间与错误信息。
type taskerr struct {
	t       time.Time
	errinfo string
}

// Task 是对单个定时任务的封装，保存 cron 规则、执行函数、执行时间与错误记录。
//
// Task 自身不是并发安全结构，由 TaskManager 在调度时统一访问；错误记录只保留
// ErrLimit 容量内最近的若干条。
type Task struct {
	// Taskname 为任务名称。
	Taskname string
	// Spec 为解析后的调度规则。
	Spec *Schedule
	// SpecStr 为原始 cron 表达式字符串。
	SpecStr string
	// DoFunc 为任务实际执行的函数。
	DoFunc TaskFunc
	// Prev 为上一次执行时间。
	Prev time.Time
	// Next 为下一次计划执行时间。
	Next time.Time
	// Timeout 为任务执行超时时间；为 0 表示不设置超时。
	Timeout time.Duration
	// Errlist 按 ErrLimit 容量环形保存最近的执行错误。
	Errlist []*taskerr
	// ErrLimit 为 Errlist 保留的错误条数；NewTask 初始化为 100。
	ErrLimit int
	// errCnt 累计执行失败次数，用于在 Errlist 中定位环形覆盖下标。
	errCnt int
}

// NewTask 根据任务名称、cron 表达式与执行函数创建任务，ErrLimit 默认为 100；
// opts 用于在解析 cron 规则前覆盖任务配置（如超时时间）。spec 非法时会在 SetCron 中 panic。
func NewTask(tname string, spec string, f TaskFunc, opts ...Option) *Task {
	task := &Task{
		Taskname: tname,
		DoFunc:   f,
		ErrLimit: 100,
		SpecStr:  spec,
		Errlist:  make([]*taskerr, 0),
	}

	for _, opt := range opts {
		opt.apply(task)
	}

	task.SetCron(spec)
	return task
}

// GetSpec 返回任务原始的 cron 表达式字符串。
func (t *Task) GetSpec(context.Context) string {
	return t.SpecStr
}

// GetStatus 将已记录的错误以“时间:错误信息”逐行拼接（HTML 换行）返回；无错误时返回空串。
func (t *Task) GetStatus(context.Context) string {
	var str strings.Builder
	for _, v := range t.Errlist {
		if v == nil {
			continue
		}
		str.WriteString(v.t.String() + ":" + v.errinfo + "<br>")
	}
	return str.String()
}

// Run 执行任务函数：捕获执行期间的 panic 并写入日志，避免单个任务异常波及调度器；
// DoFunc 返回错误时按环形容量记入 Errlist，并将该错误原样返回。
func (t *Task) Run(ctx context.Context) error {

	defer func() {
		if err := recover(); err != nil {
			logger.Logger.WithFields(logrus.Fields{
				"err": err,
			}).Error("recover")
		}
	}()
	err := t.DoFunc(ctx)
	if err != nil {
		index := t.errCnt % t.ErrLimit
		t.Errlist[index] = &taskerr{t: t.Next, errinfo: err.Error()}
		t.errCnt++
	}
	return err
}

// SetNext 以 now 为基准按任务规则计算下一次执行时间并写入 Next。
func (t *Task) SetNext(ctx context.Context, now time.Time) {
	t.Next = t.Spec.Next(now)
}

// GetNext 返回任务下一次计划执行时间。
func (t *Task) GetNext(context.Context) time.Time {
	return t.Next
}

// SetPrev 将任务上一次执行时间记录为 now。
func (t *Task) SetPrev(ctx context.Context, now time.Time) {
	t.Prev = now
}

// GetPrev 返回任务上一次执行时间。
func (t *Task) GetPrev(context.Context) time.Time {
	return t.Prev
}

// GetTimeout 返回任务执行超时时间；为 0 表示不设置超时。
func (t *Task) GetTimeout(context.Context) time.Duration {
	return t.Timeout
}

// Option 是任务配置选项接口，函数式选项模式由 optionFunc 实现。
type Option interface {
	apply(*Task)
}

// optionFunc 是将函数包装为 Option 的适配器。
type optionFunc func(*Task)

func (f optionFunc) apply(t *Task) {
	f(t)
}

// TimeoutOption 返回设置任务超时时间的 Option；传入 0 表示不限制超时。
func TimeoutOption(timeout time.Duration) Option {
	return optionFunc(func(t *Task) {
		t.Timeout = timeout
	})
}

// SetCron 解析 cron 表达式并更新任务的调度规则。
//
// 表达式按“秒 分 时 日 月 星期”六个字段解析，星期字段可省略；取值范围依次为
// 0-59、0-59、0-23、1-31、1-12、0-6（0 表示星期日）。字段支持 *（任意取值）、
// ,（列举）、-（区间）、/（步长），并识别 @yearly、@monthly、@weekly、@daily、
// @hourly 等预定义描述符；格式非法或描述符无法识别时会直接 panic。
func (t *Task) SetCron(spec string) {
	t.Spec = t.parse(spec)
}

// parse 解析 cron 表达式：以 @ 开头时交给 parseSpec 处理，否则按五或六个字段
// 分别生成位掩码；字段数非法时直接 panic。
func (t *Task) parse(spec string) *Schedule {
	if len(spec) > 0 && spec[0] == '@' {
		return t.parseSpec(spec)
	}
	// 按空白切分字段，顺序为：秒、分、时、日、月、星期（可选）。
	fields := strings.Fields(spec)
	if len(fields) != 5 && len(fields) != 6 {
		log.Panicf("Expected 5 or 6 fields, found %d: %s", len(fields), spec)
	}

	// 省略星期字段时等价于 *。
	if len(fields) == 5 {
		fields = append(fields, "*")
	}

	schedule := &Schedule{
		Second: getField(fields[0], seconds),
		Minute: getField(fields[1], minutes),
		Hour:   getField(fields[2], hours),
		Day:    getField(fields[3], days),
		Month:  getField(fields[4], months),
		Week:   getField(fields[5], weeks),
	}

	return schedule
}

// parseSpec 将 @yearly 等预定义描述符转换为对应的 Schedule；
// 无法识别的描述符直接 panic。
func (t *Task) parseSpec(spec string) *Schedule {
	switch spec {
	case "@yearly", "@annually":
		return &Schedule{
			Second: 1 << seconds.min,
			Minute: 1 << minutes.min,
			Hour:   1 << hours.min,
			Day:    1 << days.min,
			Month:  1 << months.min,
			Week:   all(weeks),
		}

	case "@monthly":
		return &Schedule{
			Second: 1 << seconds.min,
			Minute: 1 << minutes.min,
			Hour:   1 << hours.min,
			Day:    1 << days.min,
			Month:  all(months),
			Week:   all(weeks),
		}

	case "@weekly":
		return &Schedule{
			Second: 1 << seconds.min,
			Minute: 1 << minutes.min,
			Hour:   1 << hours.min,
			Day:    all(days),
			Month:  all(months),
			Week:   1 << weeks.min,
		}

	case "@daily", "@midnight":
		return &Schedule{
			Second: 1 << seconds.min,
			Minute: 1 << minutes.min,
			Hour:   1 << hours.min,
			Day:    all(days),
			Month:  all(months),
			Week:   all(weeks),
		}

	case "@hourly":
		return &Schedule{
			Second: 1 << seconds.min,
			Minute: 1 << minutes.min,
			Hour:   all(hours),
			Day:    all(days),
			Month:  all(months),
			Week:   all(weeks),
		}
	}
	log.Panicf("Unrecognized descriptor: %s", spec)
	return nil
}

// Next 返回 t 之后下一个满足调度规则的时间；五年内找不到匹配时间时返回零值 time.Time。
func (s *Schedule) Next(t time.Time) time.Time {
	// 从 t 的下一秒开始向后查找。
	t = t.Add(1*time.Second - time.Duration(t.Nanosecond())*time.Nanosecond)

	// added 标记是否已经向后推进过时间分量，首次推进时先把被推进分量之下的部分归零。
	added := false

	// 五年内仍未找到匹配时间即视为规则不可达，返回零值。
	yearLimit := t.Year() + 5

WRAP:
	if t.Year() > yearLimit {
		return time.Time{}
	}

	for 1<<uint(t.Month())&s.Month == 0 {
		if !added {
			added = true
			t = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
		}
		t = t.AddDate(0, 1, 0)

		if t.Month() == time.January {
			goto WRAP
		}
	}

	for !dayMatches(s, t) {
		if !added {
			added = true
			t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		}
		t = t.AddDate(0, 0, 1)

		if t.Day() == 1 {
			goto WRAP
		}
	}

	for 1<<uint(t.Hour())&s.Hour == 0 {
		if !added {
			added = true
			t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())
		}
		t = t.Add(1 * time.Hour)

		if t.Hour() == 0 {
			goto WRAP
		}
	}

	for 1<<uint(t.Minute())&s.Minute == 0 {
		if !added {
			added = true
			t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())
		}
		t = t.Add(1 * time.Minute)

		if t.Minute() == 0 {
			goto WRAP
		}
	}

	for 1<<uint(t.Second())&s.Second == 0 {
		if !added {
			added = true
			t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, t.Location())
		}
		t = t.Add(1 * time.Second)

		if t.Second() == 0 {
			goto WRAP
		}
	}

	return t
}

// dayMatches 判断 t 的日与星期是否同时满足规则；两个字段都使用通配符时取交集，
// 否则按 cron 惯例取并集。
func dayMatches(s *Schedule, t time.Time) bool {
	var (
		domMatch = 1<<uint(t.Day())&s.Day > 0
		dowMatch = 1<<uint(t.Weekday())&s.Week > 0
	)

	if s.Day&starBit > 0 || s.Week&starBit > 0 {
		return domMatch && dowMatch
	}
	return domMatch || dowMatch
}

// WorkerStart 启动后台调度循环，始终返回 nil。
func (t *TaskManager) WorkerStart() error {
	t.StartTask()
	return nil
}

// WorkerStop 停止后台调度循环，始终返回 nil。
func (t *TaskManager) WorkerStop() error {
	t.StopTask()
	return nil
}

// StartTask 启动后台调度 goroutine；重复调用不会重复启动。
func (t *TaskManager) StartTask() {
	t.taskLock.Lock()
	defer t.taskLock.Unlock()
	if t.started {
		return
	}
	t.started = true
	go t.run()
}

// run 是调度循环主体：计算最近一次到期时间，到期时触发任务，
// 同时响应任务变更与停止信号。
func (t *TaskManager) run() {
	now := time.Now().Local()
	// 启动时先为所有任务计算下一次执行时间。
	t.setTasksStartTime(now)

	for {
		// 仅在读锁内复制任务列表引用，排序与取最早时间都在锁外完成。
		t.taskLock.RLock()
		sortList := NewMapSorter(t.adminTaskList)
		t.taskLock.RUnlock()
		sortList.Sort()
		var effective time.Time
		// PERF-7-38 [P1] 排序器已缓存首个任务 Next，直接复用缓存值，不再重复调用 GetNext
		if len(t.adminTaskList) == 0 || sortList.nexts[0].IsZero() {
			// 暂无任务时拉长等待，期间仍可接收新增任务与停止信号。
			effective = now.AddDate(10, 0, 0)
		} else {
			effective = sortList.nexts[0]
		}

		select {
		case now = <-time.After(effective.Sub(now)):
			t.runNextTasks(sortList, effective)
			continue
		case <-t.changed:
			// 任务列表发生变更，重新计算所有任务的执行时间。
			now = time.Now().Local()
			t.setTasksStartTime(now)
			continue
		case <-t.stop:
			t.markManagerStop()
			return
		}
	}
}

// setTasksStartTime 以 now 为基准重新计算所有任务的下一次执行时间。
func (t *TaskManager) setTasksStartTime(now time.Time) {
	t.taskLock.Lock()
	for _, task := range t.adminTaskList {
		task.SetNext(context.Background(), now)
	}
	t.taskLock.Unlock()
}

// markManagerStop 将调度器标记为已停止。
func (t *TaskManager) markManagerStop() {
	t.taskLock.Lock()
	if t.started {
		t.started = false
	}
	t.taskLock.Unlock()
}

// runNextTasks 触发所有下一次执行时间等于 effective 的任务；
// 配置了超时的任务使用带超时的上下文运行。
func (t *TaskManager) runNextTasks(sortList *MapSorter, effective time.Time) {
	i := 0
	for _, e := range sortList.Vals {
		i++
		if e.GetNext(context.Background()) != effective {
			break
		}

		ctx := context.Background()
		t.wait.Add(1)
		if duration := e.GetTimeout(ctx); duration != 0 {
			go func(e Tasker) {
				defer t.wait.Done()
				ctx, cancelFunc := context.WithTimeout(ctx, duration)
				defer cancelFunc()
				err := e.Run(ctx)
				if err != nil {
					log.Printf("tasker.run err: %s\n", err.Error())
				}
			}(e)
		} else {
			go func(e Tasker) {
				defer t.wait.Done()
				err := e.Run(ctx)
				if err != nil {
					log.Printf("tasker.run err: %s\n", err.Error())
				}
			}(e)
		}

		e.SetPrev(context.Background(), e.GetNext(context.Background()))
		e.SetNext(context.Background(), effective)
	}
}

// StopTask 异步通知调度循环退出，不等待正在执行的任务结束。
func (t *TaskManager) StopTask() {
	// PERF-7-39 [P2] 缓冲通道 + 非阻塞发送；已有停止信号待处理时直接丢弃，不再起独立 goroutine 阻塞等待接收方
	select {
	case t.stop <- true:
	default:
	}
}

// GracefulShutdown 停止调度循环并等待所有正在执行的任务结束，返回的通道在全部任务完成后关闭。
func (t *TaskManager) GracefulShutdown() <-chan struct{} {
	done := make(chan struct{})
	go func() {
		select {
		case t.stop <- true:
		default:
		}
		t.wait.Wait()
		close(done)
	}()
	return done
}

// AddTask 以 taskname 为键注册任务；调度器已在运行时通知调度循环重新计算执行时间。
// taskname 重复会覆盖同名旧任务。
func (t *TaskManager) AddTask(taskname string, task Tasker) {
	isChanged := false
	t.taskLock.Lock()
	task.SetNext(context.TODO(), time.Now().Local())
	t.adminTaskList[taskname] = task
	if t.started {
		isChanged = true
	}
	t.taskLock.Unlock()

	if isChanged {
		// PERF-7-39 [P2] 缓冲通道 + 非阻塞合并发送；已有变更信号待处理时丢弃本次通知，不再起阻塞 goroutine
		select {
		case t.changed <- true:
		default:
		}
	}
}

// Len 返回当前已注册任务的数量。
func (t *TaskManager) Len() int {
	return len(t.adminTaskList)
}

// DeleteTask 按名称移除任务；调度器已在运行时通知调度循环重新计算执行时间。
func (t *TaskManager) DeleteTask(taskname string) {
	isChanged := false

	t.taskLock.Lock()
	delete(t.adminTaskList, taskname)
	if t.started {
		isChanged = true
	}
	t.taskLock.Unlock()

	if isChanged {
		// PERF-7-39 [P2] 缓冲通道 + 非阻塞合并发送；已有变更信号待处理时丢弃本次通知，不再起阻塞 goroutine
		select {
		case t.changed <- true:
		default:
		}
	}
}

// ClearTask 移除全部已注册任务；调度器已在运行时通知调度循环重新计算执行时间。
func (t *TaskManager) ClearTask() {
	isChanged := false

	t.taskLock.Lock()
	t.adminTaskList = make(map[string]Tasker)
	if t.started {
		isChanged = true
	}
	t.taskLock.Unlock()

	if isChanged {
		// PERF-7-39 [P2] 缓冲通道 + 非阻塞合并发送；已有变更信号待处理时丢弃本次通知，不再起阻塞 goroutine
		select {
		case t.changed <- true:
		default:
		}
	}
}

// MapSorter 将任务 map 拆为键、值两个平行切片，并实现 sort.Interface，
// 按任务下一次执行时间升序排列。
type MapSorter struct {
	// Keys 为任务名称切片，与 Vals 相同下标对齐。
	Keys []string
	// Vals 为任务实例切片，与 Keys 相同下标对齐。
	Vals []Tasker
	// nexts 为与 Vals 对齐的 Next 时间戳缓存，避免 Less 内反复调用 GetNext（PERF-7-38）。
	nexts []time.Time
}

// NewMapSorter 由任务 map 构造对应的 MapSorter，键与值在相同下标对齐。
func NewMapSorter(m map[string]Tasker) *MapSorter {
	ms := &MapSorter{
		Keys:  make([]string, 0, len(m)),
		Vals:  make([]Tasker, 0, len(m)),
		nexts: make([]time.Time, 0, len(m)),
	}
	for k, v := range m {
		ms.Keys = append(ms.Keys, k)
		ms.Vals = append(ms.Vals, v)
		// PERF-7-38 [P1] 构造时一次性缓存每个任务的 Next 时间戳，供排序比较复用
		ms.nexts = append(ms.nexts, v.GetNext(context.Background()))
	}
	return ms
}

// Sort 就地按下一次执行时间升序排列任务。
func (ms *MapSorter) Sort() {
	sort.Sort(ms)
}

// Len 返回待排序任务数量，实现 sort.Interface。
func (ms *MapSorter) Len() int { return len(ms.Keys) }

// Less 判断下标 i 的任务是否应排在 j 之前；未设置下次执行时间的任务排在最后。
func (ms *MapSorter) Less(i, j int) bool {
	// PERF-7-38 [P1] 直接比较缓存的 Next 时间戳，替代原实现内最多 4 次 GetNext 接口调用
	ni, nj := ms.nexts[i], ms.nexts[j]
	if ni.IsZero() {
		return false
	}
	if nj.IsZero() {
		return true
	}
	return ni.Before(nj)
}

// Swap 交换两个下标的任务及其名称，实现 sort.Interface。
func (ms *MapSorter) Swap(i, j int) {
	ms.Vals[i], ms.Vals[j] = ms.Vals[j], ms.Vals[i]
	ms.Keys[i], ms.Keys[j] = ms.Keys[j], ms.Keys[i]
	ms.nexts[i], ms.nexts[j] = ms.nexts[j], ms.nexts[i]
}

// getField 解析逗号分隔的多个区间表达式，合并为该字段的位掩码。
func getField(field string, r bounds) uint64 {
	var bits uint64
	ranges := strings.FieldsFunc(field, func(r rune) bool { return r == ',' })
	for _, expr := range ranges {
		bits |= getRange(expr, r)
	}
	return bits
}

// getRange 解析单个 cron 子表达式（数字、区间 number-number、步长 number/step 或 *），
// 返回对应取值的位掩码；超出区间或格式非法时直接 panic。
func getRange(expr string, r bounds) uint64 {
	var (
		start, end, step uint
		rangeAndStep     = strings.Split(expr, "/")
		lowAndHigh       = strings.Split(rangeAndStep[0], "-")
		singleDigit      = len(lowAndHigh) == 1
	)

	var extrastar uint64
	if lowAndHigh[0] == "*" || lowAndHigh[0] == "?" {
		start = r.min
		end = r.max
		extrastar = starBit
	} else {
		start = parseIntOrName(lowAndHigh[0], r.names)
		switch len(lowAndHigh) {
		case 1:
			end = start
		case 2:
			end = parseIntOrName(lowAndHigh[1], r.names)
		default:
			log.Panicf("Too many hyphens: %s", expr)
		}
	}

	switch len(rangeAndStep) {
	case 1:
		step = 1
	case 2:
		step = mustParseInt(rangeAndStep[1])

		// 单数字搭配步长时等价于 N-max/step。
		if singleDigit {
			end = r.max
		}
	default:
		log.Panicf("Too many slashes: %s", expr)
	}

	if start < r.min {
		log.Panicf("Beginning of range (%d) below minimum (%d): %s", start, r.min, expr)
	}
	if end > r.max {
		log.Panicf("End of range (%d) above maximum (%d): %s", end, r.max, expr)
	}
	if start > end {
		log.Panicf("Beginning of range (%d) beyond end of range (%d): %s", start, end, expr)
	}

	return getBits(start, end, step) | extrastar
}

// parseIntOrName 先按别名表（月份、星期）查找 expr，未命中时按整数解析。
func parseIntOrName(expr string, names map[string]uint) uint {
	if names != nil {
		if namedInt, ok := names[strings.ToLower(expr)]; ok {
			return namedInt
		}
	}
	return mustParseInt(expr)
}

// mustParseInt 将 expr 解析为非负整数；解析失败或出现负数时直接 panic。
func mustParseInt(expr string) uint {
	num, err := strconv.Atoi(expr)
	if err != nil {
		log.Panicf("Failed to parse int from %s: %s", expr, err)
	}
	if num < 0 {
		log.Panicf("Negative number (%d) not allowed: %s", num, expr)
	}

	return uint(num)
}

// getBits 置位区间 [min, max] 内按 step 步长取值对应的所有位。
func getBits(min, max, step uint) uint64 {
	var bits uint64

	// 步长为 1 时直接用位运算生成掩码，避免循环。
	if step == 1 {
		return ^(math.MaxUint64 << (max + 1)) & (math.MaxUint64 << min)
	}

	for i := min; i <= max; i += step {
		bits |= 1 << i
	}
	return bits
}

// all 返回给定区间的全部取值位，并附加通配符标记位。
func all(r bounds) uint64 {
	return getBits(r.min, r.max, 1) | starBit
}
