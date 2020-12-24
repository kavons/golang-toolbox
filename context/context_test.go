package context_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type A interface {
    Do()
}

type B struct {

}

func (i *B) Do() {

}

type C struct {
    A
}

var a A = &C{}

func TestWithValue(t *testing.T) {
	ctx := context.WithValue(context.Background(), "trace_id", 123456)
	ctx = context.WithValue(ctx, "session", "won")

	traceId, ok := ctx.Value("trace_id").(int)
	if !ok {
		traceId = 654321
	}
	assert.Equal(t, traceId, 123456, "int error.")

	session, _ := ctx.Value("session").(string)
	assert.Equal(t, session, "won", "string error.")
}

type Result struct {
	r   *http.Response
	err error
}

func TestWithTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	tr := &http.Transport{}
	client := &http.Client{Transport: tr}
	c := make(chan Result, 1)
	req, err := http.NewRequest("GET", "http://www.baidu1.com", nil)
	if err != nil {
		fmt.Println("http request failed, err:", err)
		return
	}

	go func() {
		resp, err := client.Do(req)
		pack := Result{r: resp, err: err}
		c <- pack
	}()

	select {
	case <-ctx.Done(): //ctx到时，这个channel里面就会有数据。
		tr.CancelRequest(req)
		res := <-c
		fmt.Println("Timeout! err:", res.err, ctx.Err())
	case res := <-c:
		defer res.r.Body.Close()
		out, _ := ioutil.ReadAll(res.r.Body)
		fmt.Printf("Server Response: %s", out)
	}
}

func generate(ctx context.Context) <-chan int {
	dst := make(chan int)
	n := 1
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("being canceled")
				return
			case dst <- n:
				n++
			}
		}
	}()
	return dst
}

func TestWithCancel1(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	intChan := generate(ctx)
	for n := range intChan {
		fmt.Println(n)
		if n == 5 {
			break
		}
	}

	cancel()
	time.Sleep(time.Minute)
}

func watch(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println(name, "监控退出")
			return
		default:
			fmt.Println(name, "监控中")
			time.Sleep(1 * time.Second)
		}
	}
}

func TestWithCancel2(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	go watch(ctx, "1")
	go watch(ctx, "2")
	go watch(ctx, "3")

	time.Sleep(5 * time.Second)

	cancel()
	time.Sleep(5 * time.Second)
}

func TestWithDeadline(t *testing.T) {
	d := time.Now().Add(5 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()

	select {
	case <-time.After(10 * time.Second):
		fmt.Println("overslept")
	case <-ctx.Done():
		fmt.Println(ctx.Err())
	}
}
