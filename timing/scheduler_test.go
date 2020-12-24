package timing

import (
	"sync"
	"testing"
	"time"
)

const OneSecond = 1*time.Second + 10*time.Millisecond

func TestAddWhileRunning(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	scheduler := New()
	scheduler.Start()
	defer scheduler.Stop()

	scheduler.AddCronCmd("* * * * * ?", func() { wg.Done() })

	select {
	case <-time.After(OneSecond):
		t.FailNow()
	case <-wait(wg):
	}
}

func TestMultipleCronTasks(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(2)

	scheduler := New()
	scheduler.Start()
	defer scheduler.Stop()

	scheduler.AddCronCmd("0 0 0 1 1 ?", func() {})
	scheduler.AddCronCmd("* * * * * ?", func() { wg.Done() })
	scheduler.AddCronCmd("0 0 0 31 12 ?", func() {})
	scheduler.AddCronCmd("* * * * * ?", func() { wg.Done() })

	select {
	case <-time.After(OneSecond):
		t.FailNow()
	case <-wait(wg):
	}
}

func TestCronTaskSnapShot(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	scheduler := New()
	scheduler.AddCronCmd("@every 2s", func() { wg.Done() })
	scheduler.Start()
	defer scheduler.Stop()

	select {
	case <-time.After(OneSecond):
		scheduler.AllCronTasks()
	}

	select {
	case <-time.After(OneSecond):
		t.FailNow()
	case <-wait(wg):
	}

}

type DummyTask struct{}

func (d DummyTask) Run() {
	panic("YOLO")
}

func TestTaskPanicRecovery(t *testing.T) {
	var task DummyTask

	scheduler := New()
	scheduler.Start()
	defer scheduler.Stop()
	scheduler.AddCronTask("* * * * * ?", task)

	select {
	case <-time.After(OneSecond):
		return
	}
}

func TestEveryCronTask(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(2)

	scheduler := New()
	scheduler.Start()
	defer scheduler.Stop()

	scheduler.ScheduleCronTask(Every(time.Second), TaskWrapper(func() {
		t.Log("task running")
		wg.Done()
	}))

	select {
	case <-time.After(2 * OneSecond):
		t.Error("expected task fires 2 times")
	case <-wait(wg):
	}
}

func wait(wg *sync.WaitGroup) chan bool {
	ch := make(chan bool)
	go func() {
		wg.Wait()
		ch <- true
	}()

	return ch
}
