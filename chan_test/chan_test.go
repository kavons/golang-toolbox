package chan_test

import (
	"fmt"
	"sync"
	"testing"
	"time"
    "os"
    "os/signal"
)

type Message struct {
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Event   string `json:"event"`
	Body    string `json:"body"`
}

type Client struct {
	Id   int
	Send chan Message
}

func PrintClientMessage(c *Client) {
	for {
		select {
		case m := <-c.Send:
			fmt.Printf("%#v", m)
		}
	}
}

func TestChan(t *testing.T) {
	c1 := Client{
		Id:   1,
		Send: make(chan Message),
	}

	c2 := Client{
		Id:   2,
		Send: make(chan Message),
	}

	m := Message{
		Type:    "type",
		Channel: "channel",
		Event:   "event",
		Body:    "body",
	}

	go PrintClientMessage(&c1)
	go PrintClientMessage(&c2)

	c1.Send <- m
	c2.Send <- m

	time.Sleep(time.Second)
}

func TestCloseChan(t *testing.T) {
	data := make(chan int)
	go func() {
		data <- 1
		data <- 2
		close(data)

		fmt.Println("send over")
	}()

	for v := range data {
		fmt.Println(v)
	}

	fmt.Println("receive over")
}

func TestCloseChanAgain(t *testing.T) {
	data := make(chan int)
	go func() {
		data <- 1
		close(data)

		fmt.Println("send over")
	}()

	for {
		v, ok := <-data
		if !ok {
			fmt.Println("closed")

			v1, ok1 := <-data // close again
			fmt.Println(v1, ok1)
			break
		} else {
			fmt.Println(v)
		}
	}

	fmt.Println("receive over")
}

func TestSelectBufferedChan(t *testing.T) {
	c := make(chan int, 2)
	t.Logf("len(c) = %d", len(c))

	select {
	case c <- 10: //
		t.Logf("len(c) = %d", len(c))
	default:
		t.Log("1")
	}

	select {
	case c <- 11:
		t.Logf("len(c) = %d", len(c))
	default:
		t.Log("2") //
	}

	close(c)

	select {
	case v, ok := <-c:
		t.Logf("v - %d, ok - %v", v, ok) //
	default:
		t.Log("3")
	}

	select {
	case v, ok := <-c:
		t.Logf("v - %d, ok - %v", v, ok)
	default:
		t.Log("4") //
	}

	select {
	case v, ok := <-c:
		t.Logf("v - %d, ok - %v", v, ok)
	default:
		t.Log("5") //
	}
}

func TestSelectBlockedChan(t *testing.T) {
	c := make(chan int)
	t.Logf("len(c) = %d", len(c))

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case c <- 1:
			t.Log("insert 1")
		default:
			t.Log("closing channel")
			close(c)
		}
	}()

	select {
	case v, ok := <-c:
		t.Logf("v - %d, ok - %v", v, ok)
	default:
		t.Log("1") //
	}

	wg.Wait()
}

func TestSignalChan(t *testing.T) {
    ch := make(chan os.Signal)
    signal.Notify(ch, os.Kill)

    s := <- ch
    t.Log(s)
}
