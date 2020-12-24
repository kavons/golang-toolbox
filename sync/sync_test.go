package sync_test

import (
	"fmt"
	"sync"
	"testing"
)

func TestSyncWaitGroup(t *testing.T) {
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)

		go func(n int) {
			defer wg.Done()
			fmt.Println(n)
		}(i)
	}

	wg.Wait()
}
