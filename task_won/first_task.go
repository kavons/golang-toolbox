package task_won

import (
    "fmt"
)

var FirstTask = &Task{
    Id:      FirstTaskId,
    Name:    "first task",
    Spec:    "* * * * * ?",
    RunFunc: FirstTaskRunFunc,
}

func FirstTaskRunFunc() error {
    fmt.Println("first task")
    return nil
}