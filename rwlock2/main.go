package main

import (
	"fmt"
	"sync"
	"time"
)

var gRWLock *sync.RWMutex

var gVar int

func init() {
	gRWLock = new(sync.RWMutex)
	gVar = 1
}

func main() {
	var wg sync.WaitGroup
	go Write(4, &wg)
	wg.Add(1)
	go Read(1, &wg)
	wg.Add(1)
	go Read(2, &wg)
	wg.Add(1)
	go Read(3, &wg)
	wg.Add(1)

	wg.Wait()
}

func Read(id int, wg *sync.WaitGroup) {
	//fmt.Printf("Read Coroutine: %d start\n", id)
	//defer fmt.Printf("Read Coroutine: %d end\n", id)
	gRWLock.RLock()
	fmt.Printf("%d - gVar %d\n", id, gVar)
	time.Sleep(time.Second)
	gRWLock.RUnlock()

	wg.Done()

}

func Write(id int, wg *sync.WaitGroup) {
	//fmt.Printf("Write Coroutine: %d start\n", id)
	//defer fmt.Printf("Write Coroutine: %d end\n", id)
	gRWLock.Lock()
	gVar = gVar + 100
	fmt.Printf("%d - gVar %d\n", id, gVar)
	time.Sleep(time.Second)
	gRWLock.Unlock()
	wg.Done()

}