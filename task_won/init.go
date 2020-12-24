package task_won

import (
	"fmt"
	"sync"

	"github.com/astaxie/beego/logs"

	"won-utils"
)

// 个人定制部分
const (
	FirstTaskId = 1
)

var (
	Manager *utils.Scheduler
	lock    sync.Mutex
	TaskMap map[int]*Task
)

func init() {
	Manager = utils.NewScheduler()
	Manager.Start()
	RegisterTaskMap()
}

func RegisterTaskMap() {
    // 个人定制部分
	TaskMap = map[int]*Task{
		FirstTaskId: FirstTask,
	}

	for _, v := range TaskMap {
		AddTask(v)
	}
}

type Task struct {
	Id      int          `json:"id"`
	Name    string       `json:"name"`
	Spec    string       `json:"spec"`
	Params  string       `json:"params"`
	Running bool         `json:"running"`
	RunFunc func() error `json:"-"`
}

func (t *Task) Run() {
	defer func() {
		if err := recover(); err != nil {
			logs.Debug(fmt.Sprintf("task_id:%d,task_name:%s has exception:%s", t.Id, t.Name, err))
		}
	}()

	if t.Running == true {
		return
	}

	t.Running = true
	t.RunFunc()

	defer func() { t.Running = false }()
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

func GetAllRunningTasks() map[int]*Task {
	cronTasks := Manager.AllCronTasks()

	tasks := make(map[int]*Task)
	for _, c := range cronTasks {
		t := c.Task.(*Task)
		if t.Running == true {
			tasks[t.Id] = t
		}
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
func AddTask(task *Task) bool {
	lock.Lock()
	defer lock.Unlock()

	if GetTask(task.Id) != nil {
		return false
	}

	err := Manager.AddCronTask(task.Spec, task)
	if err != nil {
		logs.Error("AddTask [%d] [%s]:%s", task.Id, task.Name, err.Error())
		return false
	}

	return true
}

// remove task
func RemoveTask(id int) {
	Manager.RemoveCronTask(func(c *utils.CronTask) bool {
		if t, ok := c.Task.(*Task); ok {
			if t.Id == id {
				t.Running = false
				return true
			}
		}
		return false
	})
}

// stop all tasks
func StopAllTasks() {
	tasks := GetAllTasks()
	for _, t := range tasks {
		t.Running = false
	}
	Manager.Stop()
}

// start all tasks
func StartAllTasks() {
	Manager.ReStart()
}
