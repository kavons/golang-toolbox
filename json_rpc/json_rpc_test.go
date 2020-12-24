package json_rpc_test

import (
	"errors"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"testing"
)

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

type Arithmetic struct{}

func (*Arithmetic) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (*Arithmetic) Divide(args *Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("除数不能为零")
	}

	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func TestJsonRpcServer(t *testing.T) {
	// json rpc server
	err := rpc.Register(new(Arithmetic))
	checkErr(err)

	listener, err := net.Listen("tcp", "127.0.0.1:9191")
	checkErr(err)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}

			go jsonrpc.ServeConn(conn)
		}
	}()

	forever := make(chan bool)
	<- forever
}

func TestJsonRpcClient(t *testing.T) {
	// json rpc client
	client, err := jsonrpc.Dial("tcp", ":9191")
	checkErr(err)
	defer client.Close()

	args := &Args{A: 17, B: 3}
	var multiply int
	err = client.Call("Arithmetic.Multiply", args, &multiply)
	checkErr(err)
	t.Logf("%d * %d = %d\n", args.A, args.B, multiply)

	args = &Args{A: 20, B: 6}
	quotient := new(Quotient)
	call := client.Go("Arithmetic.Divide", args, quotient, nil)
	<-call.Done
	t.Logf("%d / %d = %d, %d\n", args.A, args.B, quotient.Quo, quotient.Rem)
}
