package task_won

import (
    "testing"
    "time"
)

func TestForTask(t *testing.T) {
    FirstTask.RunFunc()
    time.Sleep(5 * time.Second)
}
