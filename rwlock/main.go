package main

import (
    "sync"
    "time"
)

var m *sync.RWMutex

func main() {
    m = new(sync.RWMutex)

    // 写的时候啥也不能干
    go write(1)
    go read(2)
    go write(3)

    time.Sleep(3*time.Second)
}

func read(i int) {
    //println(i,"read start")

    m.RLock()
    defer m.RUnlock()

    println(i,"reading")
    time.Sleep(1*time.Second)

    //println(i,"read over")
}

func write(i int) {
    //println(i,"write start")

    m.Lock()
    defer m.Unlock()

    println(i,"writing")
    time.Sleep(1*time.Second)

    //println(i,"write over")
}