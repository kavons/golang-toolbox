package master_worker_job

import (
	"fmt"
	"testing"
	"time"
)

type ExampleJob struct {
	num int
}

func (j *ExampleJob) Do() {
	fmt.Println(j.num)
}

func TestExample(t *testing.T) {
	master := NewMaster(4, 8)
	master.Run()

	for i := 0; i < 20; i++ {
		master.AddJob(&ExampleJob{i})
	}

	time.Sleep(time.Second)
}
