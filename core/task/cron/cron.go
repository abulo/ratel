package cron

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/abulo/ratel/v3/core/logger"
)

// Cron 维护一组调度条目，按各自调度表在独立 goroutine 中执行关联任务；
// 可在运行期间启动、停止，并查看当前条目快照。
type Cron struct {
	entries   []*Entry
	chain     Chain
	stop      chan struct{}
	add       chan *Entry
	remove    chan EntryID
	snapshot  chan chan []Entry
	running   bool
	runningMu sync.Mutex
	location  *time.Location
	parser    ScheduleParser
	nextID    EntryID
	jobWaiter sync.WaitGroup
}

// ScheduleParser 将调度表字符串解析为 Schedule。
type ScheduleParser interface {
	Parse(spec string) (Schedule, error)
}

// Job 是提交给 Cron 执行的任务接口。
type Job interface {
	Run()
}

// Schedule 描述任务的触发周期。
type Schedule interface {
	// Next 返回晚于给定时间的下一次触发时间；初始化时以及每次任务执行后都会调用。
	Next(time.Time) time.Time
}

// EntryID 标识 Cron 实例内的一个调度条目。
type EntryID int

// Entry 由调度表与到点执行的任务组成。
type Entry struct {
	// ID 为调度器分配的条目 ID，可用于查看快照或移除条目。
	ID EntryID

	// Schedule 为任务依据的调度表。
	Schedule Schedule

	// Next 为下一次触发时间；Cron 未启动或调度表不可满足时为零值时间。
	Next time.Time

	// Prev 为上一次触发时间；从未触发时为零值时间。
	Prev time.Time

	// WrappedJob 为调度表触发时实际执行的任务，已经过包装链包装。
	WrappedJob Job

	// Job 为用户最初提交的任务，保留下来以便后续通过 Entries 取回。
	Job Job
}

// Valid 报告该条目是否为有效条目（ID 非零）。
func (e Entry) Valid() bool { return e.ID != 0 }

// byTime 用于按 Next 时间对条目排序，零值时间排在末尾。
type byTime []*Entry

// Len 返回条目数量，是 sort.Interface 的一部分。
func (s byTime) Len() int { return len(s) }

// Swap 交换两个索引处的条目，是 sort.Interface 的一部分。
func (s byTime) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Less 报告索引 i 处条目的下一次触发时间是否早于索引 j 处，是 sort.Interface 的一部分。
func (s byTime) Less(i, j int) bool {
	// 两个零值时间互不为前；零值时间恒排在其他时间之后。
	if s[i].Next.IsZero() {
		return false
	}
	if s[j].Next.IsZero() {
		return true
	}
	return s[i].Next.Before(s[j].Next)
}

// New 创建一个 Cron 任务调度器，并按给定 Option 调整默认行为。
//
// 默认时区为 time.Local；默认解析器接受标准五段 cron 表达式；包装链默认为空。
// 时区、解析器与包装链均可通过 With* 选项覆盖。
func New(opts ...Option) *Cron {
	c := &Cron{
		entries:   nil,
		chain:     NewChain(),
		add:       make(chan *Entry),
		stop:      make(chan struct{}),
		snapshot:  make(chan chan []Entry),
		remove:    make(chan EntryID),
		running:   false,
		runningMu: sync.Mutex{},
		location:  time.Local,
		parser:    standardParser,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// FuncJob 将无参函数包装为 cron.Job。
type FuncJob func()

// Run 调用被包装的无参函数。
func (f FuncJob) Run() { f() }

// AddFunc 将函数 cmd 按调度表 spec 加入 Cron；spec 以本实例时区解析，
// 返回可用于后续移除条目的不透明 ID；表达式非法时返回解析错误。
func (c *Cron) AddFunc(spec string, cmd func()) (EntryID, error) {
	return c.AddJob(spec, FuncJob(cmd))
}

// AddJob 将任务 cmd 按调度表 spec 加入 Cron；spec 以本实例时区解析，
// 返回可用于后续移除条目的不透明 ID；表达式非法时返回解析错误。
func (c *Cron) AddJob(spec string, cmd Job) (EntryID, error) {
	schedule, err := c.parser.Parse(spec)
	if err != nil {
		return 0, err
	}
	return c.Schedule(schedule, cmd), nil
}

// Schedule 将任务 cmd 绑定到给定调度表并加入 Cron；任务会经过配置的包装链包装。
func (c *Cron) Schedule(schedule Schedule, cmd Job) EntryID {
	c.runningMu.Lock()
	defer c.runningMu.Unlock()
	c.nextID++
	entry := &Entry{
		ID:         c.nextID,
		Schedule:   schedule,
		WrappedJob: c.chain.Then(cmd),
		Job:        cmd,
	}
	if !c.running {
		c.entries = append(c.entries, entry)
	} else {
		c.add <- entry
	}
	return entry.ID
}

// Entries 返回当前调度条目的快照。
func (c *Cron) Entries() []Entry {
	c.runningMu.Lock()
	defer c.runningMu.Unlock()
	if c.running {
		replyChan := make(chan []Entry, 1)
		c.snapshot <- replyChan
		return <-replyChan
	}
	return c.entrySnapshot()
}

// Location 返回实例当前使用的时区。
func (c *Cron) Location() *time.Location {
	return c.location
}

// Entry 按 ID 返回条目快照；找不到时返回零值 Entry。
func (c *Cron) Entry(id EntryID) Entry {
	for _, entry := range c.Entries() {
		if id == entry.ID {
			return entry
		}
	}
	return Entry{}
}

// Remove 将指定条目从后续调度中移除。
func (c *Cron) Remove(id EntryID) {
	c.runningMu.Lock()
	defer c.runningMu.Unlock()
	if c.running {
		c.remove <- id
	} else {
		c.removeEntry(id)
	}
}

// Start 在独立 goroutine 中启动调度器；重复调用为无操作。
func (c *Cron) Start() {
	c.runningMu.Lock()
	defer c.runningMu.Unlock()
	if c.running {
		return
	}
	c.running = true
	go c.run()
}

// Run 在当前 goroutine 中运行调度器；重复调用为无操作。
func (c *Cron) Run() {
	c.runningMu.Lock()
	if c.running {
		c.runningMu.Unlock()
		return
	}
	c.running = true
	c.runningMu.Unlock()
	c.run()
}

func (c *Cron) run() {
	logger.Logger.Infof("start")

	// 为每个条目计算初始的下一次触发时间。
	now := c.now()
	for _, entry := range c.entries {
		entry.Next = entry.Schedule.Next(now)
		logger.Logger.Infof("schedule|now=%v, entry=%v, next=%v", now, entry.ID, entry.Next)
	}

	for {
		sort.Sort(byTime(c.entries))

		var timer *time.Timer
		if len(c.entries) == 0 || c.entries[0].Next.IsZero() {
			// 暂无条目时挂起一个超长定时器，仍需保持事件循环以响应新增条目与停止请求。
			timer = time.NewTimer(100000 * time.Hour)
		} else {
			timer = time.NewTimer(c.entries[0].Next.Sub(now))
		}

		for {
			select {
			case now = <-timer.C:
				now = now.In(c.location)

				// 执行所有触发时间不晚于 now 的条目。
				for _, e := range c.entries {
					if e.Next.After(now) || e.Next.IsZero() {
						break
					}
					c.startJob(e.WrappedJob)
					e.Prev = e.Next
					e.Next = e.Schedule.Next(now)
				}

			case newEntry := <-c.add:
				timer.Stop()
				now = c.now()
				newEntry.Next = newEntry.Schedule.Next(now)
				c.entries = append(c.entries, newEntry)
				logger.Logger.Infof("added|now=%v, entry=%v, next=%v", now, newEntry.ID, newEntry.Next)

			case replyChan := <-c.snapshot:
				replyChan <- c.entrySnapshot()
				continue

			case <-c.stop:
				timer.Stop()
				logger.Logger.Infof("stop")
				return

			case id := <-c.remove:
				timer.Stop()
				now = c.now()
				c.removeEntry(id)
				logger.Logger.Infof("removed|entry=%v", id)
			}

			break
		}
	}
}

func (c *Cron) startJob(j Job) {
	c.jobWaiter.Go(func() {
		j.Run()
	})
}

func (c *Cron) now() time.Time {
	return time.Now().In(c.location)
}

// Stop 在调度器运行时停止它；未运行时不做任何操作。
//
// 返回的 Context 在所有正在执行的任务完成后被取消，调用方可用它等待任务收尾。
func (c *Cron) Stop() context.Context {
	c.runningMu.Lock()
	defer c.runningMu.Unlock()
	if c.running {
		c.stop <- struct{}{}
		c.running = false
	}
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c.jobWaiter.Wait()
		cancel()
	}()
	return ctx
}

func (c *Cron) entrySnapshot() []Entry {
	var entries = make([]Entry, len(c.entries))
	for i, e := range c.entries {
		entries[i] = *e
	}
	return entries
}

func (c *Cron) removeEntry(id EntryID) {
	var entries []*Entry
	for _, e := range c.entries {
		if e.ID != id {
			entries = append(entries, e)
		}
	}
	c.entries = entries
}
