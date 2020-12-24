package task

import (
	"sync"
	"utils/timing"

	"github.com/astaxie/beego/logs"
)

var (
	Manager *timing.Scheduler
	lock    sync.Mutex
)

func init() {
	Manager = timing.New()
	Manager.Start()
}

type Task struct {
	Id      int
	Name    string
	RunFunc func()
}

func (t *Task) Run() {
	t.RunFunc()
}

// get all tasks
func GetAllTasks() map[int]*Task {
	cronTasks := Manager.AllCronTasks()

	tasks := make(map[int]*Task)
	for _, c := range cronTasks {
		t := c.Task.(*Task)
		tasks[t.Id] = t
	}

	return tasks
}

// get one task
func GetTask(id int) *Task {
	tasks := GetAllTasks()
	if t, ok := tasks[id]; ok {
		return t
	}
	return nil
}

// add task
func AddTask(spec string, task *Task) bool {
	lock.Lock()
	defer lock.Unlock()

	if GetTask(task.Id) != nil {
		return false
	}

	err := Manager.AddCronTask(spec, task)
	if err != nil {
		logs.Error("AddTask:", err.Error())
		return false
	}

	return true
}

// remove task
func RemoveTask(id int) {
	Manager.RemoveCronTask(func(c *timing.CronTask) bool {
		if t, ok := c.Task.(*Task); ok {
			if t.Id == id {
				return true
			}
		}
		return false
	})
}

// stop all tasks
func StopAllTasks() {
	Manager.Stop()
}
