package ping_pong

import (
	"fmt"
	"testing"
	"time"
)

type Balls struct {
	hits int
}

func Player(name string, table chan *Balls) {
	for {
		ball := <-table
		ball.hits++
		fmt.Println(name, ball.hits)
		time.Sleep(100 * time.Millisecond)
		table <- ball
	}
}

func TestPingPong(t *testing.T) {
	table := make(chan *Balls)
	go Player("ping", table)
	go Player("pong", table)

	table <- new(Balls)
	time.Sleep(5 * time.Second)
	<-table
}
