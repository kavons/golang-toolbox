package crontask_test

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/robfig/cron"
)

type TaskEntry func()

type Task struct {
	Name         string
	SchedulePlan string
	Entry        TaskEntry
}

type Scheduler struct {
	Cron  *cron.Cron
	Tasks []*Task
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		Cron: cron.New(),
	}
}

func (s *Scheduler) Start() {
	s.Cron.Start()
}

func (s *Scheduler) AddTask(task *Task) {
	s.Tasks = append(s.Tasks, task)
}

func (s *Scheduler) DoTasks() {
	for _, task := range s.Tasks {
		err := s.Cron.AddFunc(task.SchedulePlan, task.Entry)
		if err != nil {
			log.Fatal(fmt.Sprintf("Scheduling task %s failed", task.Name))
		}
	}
}

func TestSchedulerStart(t *testing.T) {
	s := NewScheduler()

	s.AddTask(&Task{
		Name:         "First Task",
		SchedulePlan: "*/5 * * * * *",
		Entry: func() {
			fmt.Println("First Task")
		},
	})

	s.AddTask(&Task{
		Name:         "Second Task",
		SchedulePlan: "*/8 * * * * *",
		Entry: func() {
			fmt.Println("Second Task")
		},
	})

	s.DoTasks()
	s.Start()
	time.Sleep(30 * time.Second)
}
