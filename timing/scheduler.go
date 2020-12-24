package timing

import (
	"log"
	"runtime"
	"sort"
	"time"
)

// 调度员 - singular
type Scheduler struct {
	Running      bool                    // 运行状态标识
	CronTasks    []*CronTask             // 所有的计划任务
	AddChan      chan *CronTask          // 添加任务通道
	RemoveChan   chan RemoveCronTaskFunc // 删除任务通道
	StopChan     chan struct{}           // 暂停所有任务通道
	SnapShotChan chan []*CronTask        // 任务快照通道
}

// 计划任务 - plural
type CronTask struct {
	schedule Schedule  // 时间安排
	Task     Task      // 任务
	Prev     time.Time // 最近一次的执行时刻
	Next     time.Time // 下一次的执行时刻
}

// 时间安排
type Schedule interface {
	Next(time.Time) time.Time
}

// 任务
type Task interface {
	Run()
}

// 删除任务
type RemoveCronTaskFunc func(*CronTask) bool

// 1. 任务排序
type CronTaskSorter []*CronTask

func (s CronTaskSorter) Len() int {
	return len(s)
}

func (s CronTaskSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s CronTaskSorter) Less(i, j int) bool {
	if s[i].Next.IsZero() {
		return false
	}

	if s[j].Next.IsZero() {
		return true
	}

	return s[i].Next.Before(s[j].Next)
}

// 2. 调度者创建及其行为
func New() *Scheduler {
	return &Scheduler{
		Running:      false,
		AddChan:      make(chan *CronTask),
		RemoveChan:   make(chan RemoveCronTaskFunc),
		StopChan:     make(chan struct{}),
		SnapShotChan: make(chan []*CronTask),
	}
}

func (s *Scheduler) Start() {
	if s.Running {
		return
	}
	s.Running = true
	go s.loop()
}

// blocking Running method
func (s *Scheduler) Run() {
	if s.Running {
		return
	}
	s.Running = true
	s.loop()
}

func (s *Scheduler) Stop() {
	if s.Running == false {
		return
	}

	s.StopChan <- struct{}{}
}

type TaskWrapper func()

func (w TaskWrapper) Run() { w() }

func (s *Scheduler) AddCronCmd(spec string, cmd func()) error {
	return s.AddCronTask(spec, TaskWrapper(cmd))
}

func (s *Scheduler) AddCronTask(spec string, task Task) error {
	schedule, err := Parse(spec)
	if err != nil {
		return err
	}

	s.ScheduleCronTask(schedule, task)
	return nil
}

func (s *Scheduler) ScheduleCronTask(schedule Schedule, task Task) {
	cronTask := &CronTask{
		schedule: schedule,
		Task:     task,
	}

	if s.Running == false {
		s.CronTasks = append(s.CronTasks, cronTask)
		return
	}

	s.AddChan <- cronTask
	return
}

func (s *Scheduler) RemoveCronTask(f RemoveCronTaskFunc) {
	s.RemoveChan <- f
}

func (s *Scheduler) AllCronTasks() []*CronTask {
	if s.Running {
		s.SnapShotChan <- nil
		l := <-s.SnapShotChan
		return l
	}

	return s.CronTaskSnapShot()
}

func (s *Scheduler) CronTaskSnapShot() []*CronTask {
	cronTasks := make([]*CronTask, 0)
	for _, c := range s.CronTasks {
		cronTasks = append(cronTasks, &CronTask{
			schedule: c.schedule,
			Task:     c.Task,
			Next:     c.Next,
			Prev:     c.Prev,
		})
	}
	return cronTasks
}

func (s *Scheduler) KickOffTask(t Task) {
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("task panic %v\n%s", err, buf)
		}
	}()
	t.Run()
}

// 3. scheduler engine loop
func (s *Scheduler) loop() {
	now := time.Now().Local()
	for _, cronTask := range s.CronTasks {
		cronTask.Next = cronTask.schedule.Next(now)
	}

	for {
		sort.Sort(CronTaskSorter(s.CronTasks))

		var effective time.Time
		if len(s.CronTasks) == 0 || s.CronTasks[0].Next.IsZero() {
			effective = now.AddDate(10, 0, 0)
		} else {
			effective = s.CronTasks[0].Next
		}

		select {
		// 按照时间顺序运行任务
		case now = <-time.After(effective.Sub(now)):
			for _, c := range s.CronTasks {
				if c.Next.After(now) || c.Next.IsZero() {
					break
				}

				go s.KickOffTask(c.Task)
				c.Prev = c.Next
				c.Next = c.schedule.Next(now)
			}
			continue

		// 动态添加任务
		case newCronTask := <-s.AddChan:
			newCronTask.Next = newCronTask.schedule.Next(now)
			s.CronTasks = append(s.CronTasks, newCronTask)

		// 动态删除任务
		case remove := <-s.RemoveChan:
			newCronTasks := make([]*CronTask, 0)
			for _, c := range s.CronTasks {
				if remove(c) {
					continue
				}

				newCronTasks = append(newCronTasks, c)
			}

			s.CronTasks = newCronTasks

		case <-s.SnapShotChan:
			s.SnapShotChan <- s.CronTaskSnapShot()

		// 动态停止所有任务
		case <-s.StopChan:
			return
		}

		now = time.Now().Local()
	}
}
